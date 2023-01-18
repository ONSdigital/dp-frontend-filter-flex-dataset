package handlers

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"sync"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/mapper"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/model"
	"github.com/ONSdigital/dp-net/v2/handlers"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// FilterFlexOverview Handler
func (f *FilterFlex) FilterFlexOverview() http.HandlerFunc {
	return handlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		filterFlexOverview(w, req, f, accessToken, collectionID, lang)
	})
}

func filterFlexOverview(w http.ResponseWriter, req *http.Request, f *FilterFlex, accessToken, collectionID, lang string) {
	ctx := req.Context()
	vars := mux.Vars(req)
	filterID := vars["filterID"]

	var filterDims filter.Dimensions
	var datasetDims dataset.VersionDimensions
	var filterJob *filter.GetFilterResponse
	var eb zebedee.EmergencyBanner
	var fErr, dErr, fdsErr, imErr, zErr error
	var isMultivariate bool
	var serviceMsg string

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		eb, serviceMsg, zErr = getZebContent(ctx, f.ZebedeeClient, accessToken, collectionID, lang)
	}()

	go func() {
		defer wg.Done()
		filterInput := &filter.GetFilterInput{
			FilterID: filterID,
			AuthHeaders: filter.AuthHeaders{
				UserAuthToken: accessToken,
				CollectionID:  collectionID,
			},
		}
		filterJob, fErr = f.FilterClient.GetFilter(ctx, *filterInput)
		if fErr != nil {
			log.Error(ctx, "failed to get filter", fErr, log.Data{"filter_id": filterID})
			setStatusCode(req, w, fErr)
			return
		}

		if f.EnableMultivariate {
			isMultivariate, imErr = isMultivariateDataset(ctx, f.DatasetClient, accessToken, collectionID, filterJob.Dataset.DatasetID)
			if imErr != nil {
				log.Error(ctx, "failed to determine if dataset type is multivariate", imErr, log.Data{
					"filter_id": filterID,
				})
				setStatusCode(req, w, imErr)
				return
			}
		}

		datasetDims, dErr = f.DatasetClient.GetVersionDimensions(ctx, accessToken, "", collectionID, filterJob.Dataset.DatasetID, filterJob.Dataset.Edition, fmt.Sprint(filterJob.Dataset.Version))
		if dErr != nil {
			log.Error(ctx, "failed to get versions dimensions", dErr, log.Data{
				"dataset": filterJob.Dataset.DatasetID,
				"edition": filterJob.Dataset.Edition,
				"version": fmt.Sprint(filterJob.Dataset.Version),
			})
			setStatusCode(req, w, dErr)
			return
		}
	}()

	go func() {
		defer wg.Done()
		filterDims, _, fdsErr = f.FilterClient.GetDimensions(ctx, accessToken, "", collectionID, filterID, &filter.QueryParams{Limit: 500})
		if fdsErr != nil {
			log.Error(ctx, "failed to get dimensions", fdsErr, log.Data{"filter_id": filterID})
			setStatusCode(req, w, fdsErr)
			return
		}
	}()

	wg.Wait()

	// log zebedee error but don't set a server error
	if zErr != nil {
		log.Error(ctx, "unable to get homepage content", zErr, log.Data{"homepage_content": zErr})
	}
	if fErr != nil {
		log.Error(ctx, "failed to get filter", fErr, log.Data{"filter_id": filterID})
		setStatusCode(req, w, fErr)
		return
	}
	if fdsErr != nil {
		log.Error(ctx, "failed to get dimensions", fdsErr, log.Data{"filter_id": filterID})
		setStatusCode(req, w, fdsErr)
		return
	}
	if imErr != nil {
		log.Error(ctx, "failed to determine if dataset type is multivariate", imErr, log.Data{
			"filter_id": filterID,
		})
		setStatusCode(req, w, imErr)
		return
	}
	if dErr != nil {
		log.Error(ctx, "failed to get versions dimensions", dErr, log.Data{
			"dataset": filterJob.Dataset.DatasetID,
			"edition": filterJob.Dataset.Edition,
			"version": fmt.Sprint(filterJob.Dataset.Version),
		})
		setStatusCode(req, w, dErr)
		return
	}

	getDimensionOptions := func(dim filter.Dimension) ([]string, int, error) {
		q := dataset.QueryParams{Offset: 0, Limit: 1000}

		opts, err := f.DatasetClient.GetOptions(ctx, accessToken, "", collectionID, filterJob.Dataset.DatasetID, filterJob.Dataset.Edition, strconv.Itoa(filterJob.Dataset.Version), dim.Name, &q)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get options for dimension: %w", err)
		}

		var options []string
		for _, opt := range sortOptionsByCode(opts.Items) {
			options = append(options, opt.Label)
		}

		return options, opts.TotalCount, nil
	}

	var hasNoAreaOptions bool
	getAreaOptions := func(dim filter.Dimension) ([]string, int, error) {
		q := filter.QueryParams{Offset: 0, Limit: 500}
		opts, _, err := f.FilterClient.GetDimensionOptions(ctx, accessToken, "", collectionID, filterID, dim.Name, &q)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get options for dimension: %w", err)
		}

		options := []string{}
		if opts.TotalCount == 0 {
			areas, err := f.PopulationClient.GetAreas(ctx, population.GetAreasInput{
				AuthTokens: population.AuthTokens{
					UserAuthToken: accessToken,
				},
				PaginationParams: population.PaginationParams{
					Limit: 1000,
				},
				PopulationType: filterJob.PopulationType,
				AreaTypeID:     dim.ID,
			})
			if err != nil {
				return nil, 0, fmt.Errorf("failed to get dimension areas: %w", err)
			}

			hasNoAreaOptions = true
			return options, areas.TotalCount, nil
		}

		var wg sync.WaitGroup
		totalCount := opts.TotalCount
		optsIDs := []string{}
		for _, opt := range opts.Items {
			wg.Add(1)
			go func(opt filter.DimensionOption) {
				defer wg.Done()
				optsIDs = append(optsIDs, opt.Option)
				var areaTypeID string
				if dim.FilterByParent != "" {
					areaTypeID = dim.FilterByParent
				} else {
					areaTypeID = dim.ID
				}

				area, err := f.PopulationClient.GetArea(ctx, population.GetAreaInput{
					AuthTokens: population.AuthTokens{
						UserAuthToken: accessToken,
					},
					PopulationType: filterJob.PopulationType,
					AreaType:       areaTypeID,
					Area:           opt.Option,
				})
				if err != nil {
					log.Error(ctx, "failed to get area", err, log.Data{
						"area": dim.ID,
						"ID":   opt.Option,
					})
					setStatusCode(req, w, err)
					return
				}

				options = append(options, area.Area.Label)
			}(opt)
		}
		wg.Wait()

		// TODO: pc.GetParentAreaCount is causing issues in production
		// if dim.FilterByParent != "" {
		// 	count, err := pc.GetParentAreaCount(ctx, population.GetParentAreaCountInput{
		// 		AuthTokens: population.AuthTokens{
		// 			UserAuthToken: accessToken,
		// 		},
		// 		PopulationType:   filterJob.PopulationType,
		// 		AreaTypeID:       dim.ID,
		// 		ParentAreaTypeID: dim.FilterByParent,
		// 		Areas:            optsIDs,
		// 		SVarID:           supVar,
		// 	})
		// 	if err != nil {
		// 		log.Error(ctx, "failed to get parent area count", err, log.Data{
		// 			"population_type":     filterJob.PopulationType,
		// 			"area_type_id":        dim.ID,
		// 			"parent_area_type_id": dim.FilterByParent,
		// 			"areas":               optsIDs,
		// 		})
		// 		return nil, 0, err
		// 	}
		// 	totalCount = count
		// }

		return options, totalCount, nil
	}

	getOptions := func(dim filter.Dimension) ([]string, int, error) {
		if dim.IsAreaType != nil && *dim.IsAreaType {
			return getAreaOptions(dim)
		}

		return getDimensionOptions(dim)
	}

	var fDims []model.FilterDimension
	for i := len(filterDims.Items) - 1; i >= 0; i-- {
		// Needed to determine whether dimension is_area_type
		filterDimension, _, err := f.FilterClient.GetDimension(ctx, accessToken, "", collectionID, filterID, filterDims.Items[i].Name)
		if err != nil {
			log.Error(ctx, "failed to get dimension", err, log.Data{"dimension_name": filterDims.Items[i].Name})
			setStatusCode(req, w, err)
			return
		}
		filterDims.Items[i].IsAreaType = filterDimension.IsAreaType
		// TODO: pc.GetParentAreaCount is causing production issues
		// if !*filterDims.Items[i].IsAreaType {
		// 	supVar = filterDims.Items[i].ID
		// }
		filterDims.Items[i].FilterByParent = filterDimension.FilterByParent

		options, count, err := getOptions(filterDims.Items[i])
		if err != nil {
			log.Error(ctx, "failed to get options for dimension", err, log.Data{"dimension_name": filterDims.Items[i].Name})
			setStatusCode(req, w, err)
			return
		}
		filterDims.Items[i].Options = options
		fDims = append(fDims, model.FilterDimension{
			Dimension:    filterDims.Items[i],
			OptionsCount: count,
		})
	}

	basePage := f.Render.NewBasePageModel()
	showAll := req.URL.Query()["showAll"]
	path := req.URL.Path
	m := mapper.CreateFilterFlexOverview(req, basePage, lang, path, showAll, *filterJob, fDims, datasetDims, hasNoAreaOptions, isMultivariate, eb, serviceMsg)
	f.Render.BuildPage(w, m, "overview")
}

// sorts options by code - numerically if possible, with negatives listed last
func sortOptionsByCode(items []dataset.Option) []dataset.Option {
	sorted := []dataset.Option{}
	sorted = append(sorted, items...)

	doNumericSort := func(items []dataset.Option) bool {
		for _, item := range items {
			_, err := strconv.Atoi(item.Option)
			if err != nil {
				return false
			}
		}
		return true
	}

	if doNumericSort(items) {
		sort.Slice(sorted, func(i, j int) bool {
			left, _ := strconv.Atoi(sorted[i].Option)
			right, _ := strconv.Atoi(sorted[j].Option)
			if left*right < 0 {
				return right < 0
			} else {
				return left*left < right*right
			}
		})
	} else {
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Option < sorted[j].Option
		})
	}
	return sorted
}
