package handlers

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"sync"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/helpers"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/mapper"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/model"
	"github.com/ONSdigital/dp-net/v3/handlers"
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
	var dimDescriptions population.GetDimensionsResponse
	var dimCategories population.GetDimensionCategoriesResponse
	var filterJob *filter.GetFilterResponse
	var eb zebedee.EmergencyBanner
	var pop population.GetPopulationTypeResponse
	var sdc *cantabular.GetBlockedAreaCountResult
	var fErr, dErr, fdsErr, imErr, zErr, sErr, dcErr, pErr error
	var isMultivariate bool
	var serviceMsg, areaTypeID, parent string
	var dimIds, nonAreaIds, areaOpts []string

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
	}()

	go func() {
		defer wg.Done()
		filterDims, _, fdsErr = f.FilterClient.GetDimensions(ctx, accessToken, "", collectionID, filterID, &filter.QueryParams{Limit: 500})
		if fdsErr != nil {
			log.Error(ctx, "failed to get dimensions", fdsErr, log.Data{"filter_id": filterID})
			setStatusCode(req, w, fdsErr)
			return
		}

		for _, dim := range filterDims.Items {
			dimIds = append(dimIds, dim.ID)
			if !helpers.IsBoolPtr(dim.IsAreaType) {
				nonAreaIds = append(nonAreaIds, dim.ID)
			}
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

	wg.Add(3)
	go func() {
		defer wg.Done()
		pop, pErr = f.PopulationClient.GetPopulationType(ctx, population.GetPopulationTypeInput{
			PopulationType: filterJob.PopulationType,
			AuthTokens: population.AuthTokens{
				UserAuthToken: accessToken,
			},
		})
	}()

	go func() {
		defer wg.Done()
		if len(nonAreaIds) > 0 {
			dimCategories, dcErr = f.PopulationClient.GetDimensionCategories(ctx, population.GetDimensionCategoryInput{
				AuthTokens: population.AuthTokens{
					UserAuthToken: accessToken,
				},
				PaginationParams: population.PaginationParams{
					Limit:  1000,
					Offset: 0,
				},
				PopulationType: filterJob.PopulationType,
				Dimensions:     nonAreaIds,
			})
		} else {
			dimCategories = population.GetDimensionCategoriesResponse{}
		}
	}()

	go func() {
		defer wg.Done()
		dimDescriptions, dErr = f.PopulationClient.GetDimensionsDescription(ctx, population.GetDimensionsDescriptionInput{
			AuthTokens: population.AuthTokens{
				UserAuthToken: accessToken,
			},
			PopulationType: filterJob.PopulationType,
			DimensionIDs:   dimIds,
		})
		if dErr != nil {
			log.Error(ctx, "failed to get dimension descriptions", dErr, log.Data{
				"population_type": filterJob.PopulationType,
				"dimension_ids":   dimIds,
			})
			setStatusCode(req, w, dErr)
			return
		}
	}()

	wg.Wait()

	if pErr != nil {
		log.Error(ctx, "failed to get population type", pErr, log.Data{
			"filter_id":       filterID,
			"population_type": filterJob.PopulationType,
		})
		setStatusCode(req, w, pErr)
		return
	}

	if dErr != nil {
		log.Error(ctx, "failed to get dimension descriptions", dErr, log.Data{
			"population_type": filterJob.PopulationType,
			"dimension_ids":   dimIds,
		})
		setStatusCode(req, w, dErr)
		return
	}

	if dcErr != nil {
		log.Error(ctx, "failed to get dimension categories", dErr, log.Data{
			"population_type": filterJob.PopulationType,
			"dimension_ids":   dimIds,
		})
		setStatusCode(req, w, dcErr)
		return
	}
	dimensionCategoriesMap := mapDimensionCategories(dimCategories)

	getDimensionOptions := func(dim filter.Dimension) ([]string, int, error) {
		dimensionCategory := dimensionCategoriesMap[dim.ID]

		var options []string
		for _, opt := range sortCategoriesByID(dimensionCategory.Categories) {
			options = append(options, opt.Label)
		}

		return options, len(options), nil
	}

	getDimensionCategorisations := func(populationType string, dimension string) (int, error) {
		cats, err := f.PopulationClient.GetCategorisations(ctx, population.GetCategorisationsInput{
			AuthTokens: population.AuthTokens{
				UserAuthToken: accessToken,
			},
			PaginationParams: population.PaginationParams{
				Limit: 1000,
			},
			PopulationType: populationType,
			Dimension:      dimension,
		})
		return cats.PaginationResponse.TotalCount, err
	}

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

		areaOpts = optsIDs

		return options, totalCount, nil
	}

	getOptions := func(dim filter.Dimension) ([]string, int, error) {
		if helpers.IsBoolPtr(dim.IsAreaType) {
			areaTypeID = dim.ID
			parent = dim.FilterByParent
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
		filterDims.Items[i].FilterByParent = filterDimension.FilterByParent

		options, count, err := getOptions(filterDims.Items[i])
		if err != nil {
			log.Error(ctx, "failed to get options for dimension", err, log.Data{"dimension_name": filterDims.Items[i].Name})
			setStatusCode(req, w, err)
			return
		}

		categorisationCount := 0
		if !isAreaType(filterDimension) {
			categorisationCount, _ = getDimensionCategorisations(filterJob.PopulationType, filterDimension.Name)
		}

		filterDims.Items[i].Options = options
		fDims = append(fDims, model.FilterDimension{
			Dimension:           filterDims.Items[i],
			OptionsCount:        count,
			CategorisationCount: categorisationCount,
		})
	}

	if isMultivariate {
		sdc, sErr = f.getBlockedAreaCount(ctx, accessToken, filterJob.PopulationType, areaTypeID, parent, dimIds, areaOpts)
		if sErr != nil {
			log.Error(ctx, "failed to get blocked area count", sErr, log.Data{
				"population_type": filterJob.PopulationType,
				"variables":       dimIds,
				"area_codes":      areaOpts,
				"area_type_id":    areaTypeID,
			})
			setStatusCode(req, w, sErr)
			return
		}
	} else {
		sdc = &cantabular.GetBlockedAreaCountResult{}
	}

	basePage := f.Render.NewBasePageModel()
	m := mapper.NewMapper(req, basePage, eb, lang, serviceMsg, filterID)
	overview := m.CreateFilterFlexOverview(*filterJob, fDims, dimDescriptions, pop, *sdc, isMultivariate)
	f.Render.BuildPage(w, overview, "overview")
}

func mapDimensionCategories(dimCategories population.GetDimensionCategoriesResponse) map[string]population.DimensionCategory {
	dimensionCategoryMap := make(map[string]population.DimensionCategory)
	for _, dimensionCategory := range dimCategories.Categories {
		dimensionCategoryMap[dimensionCategory.Id] = dimensionCategory
	}
	return dimensionCategoryMap
}

// sorts population.DimensionCategoryItems - numerically if possible, with negatives listed last
func sortCategoriesByID(items []population.DimensionCategoryItem) []population.DimensionCategoryItem {
	sorted := []population.DimensionCategoryItem{}
	sorted = append(sorted, items...)

	doNumericSort := func(items []population.DimensionCategoryItem) bool {
		for _, item := range items {
			_, err := strconv.Atoi(item.ID)
			if err != nil {
				return false
			}
		}
		return true
	}

	if doNumericSort(items) {
		sort.Slice(sorted, func(i, j int) bool {
			left, _ := strconv.Atoi(sorted[i].ID)
			right, _ := strconv.Atoi(sorted[j].ID)
			if left*right < 0 {
				return right < 0
			} else {
				return left*left < right*right
			}
		})
	} else {
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].ID < sorted[j].ID
		})
	}
	return sorted
}
