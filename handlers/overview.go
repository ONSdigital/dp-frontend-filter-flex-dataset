package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/mapper"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/model"
	"github.com/ONSdigital/dp-net/v2/handlers"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// FilterFlexOverview Handler
func FilterFlexOverview(rc RenderClient, fc FilterClient, dc DatasetClient, pc PopulationClient) http.HandlerFunc {
	return handlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		filterFlexOverview(w, req, rc, fc, dc, pc, accessToken, collectionID, lang)
	})
}

func filterFlexOverview(w http.ResponseWriter, req *http.Request, rc RenderClient, fc FilterClient, dc DatasetClient, pc PopulationClient, accessToken, collectionID, lang string) {
	ctx := req.Context()
	vars := mux.Vars(req)
	filterID := vars["filterID"]

	var filterDims filter.Dimensions
	var datasetDims dataset.VersionDimensions
	var filterJob *filter.GetFilterResponse
	var fErr, dErr, fdsErr error

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		filterInput := &filter.GetFilterInput{
			FilterID: filterID,
			AuthHeaders: filter.AuthHeaders{
				UserAuthToken: accessToken,
				CollectionID:  collectionID,
			},
		}
		filterJob, fErr = fc.GetFilter(ctx, *filterInput)
		if fErr != nil {
			log.Error(ctx, "failed to get filter", fErr, log.Data{"filter_id": filterID})
			setStatusCode(req, w, fErr)
			return
		}

		datasetDims, dErr = dc.GetVersionDimensions(ctx, accessToken, "", collectionID, filterJob.Dataset.DatasetID, filterJob.Dataset.Edition, fmt.Sprint(filterJob.Dataset.Version))
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
		filterDims, _, fdsErr = fc.GetDimensions(ctx, accessToken, "", collectionID, filterID, &filter.QueryParams{Limit: 500})
		if fdsErr != nil {
			log.Error(ctx, "failed to get dimensions", fdsErr, log.Data{"filter_id": filterID})
			setStatusCode(req, w, fdsErr)
			return
		}
	}()

	wg.Wait()

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

		opts, err := dc.GetOptions(ctx, accessToken, "", collectionID, filterJob.Dataset.DatasetID, filterJob.Dataset.Edition, strconv.Itoa(filterJob.Dataset.Version), dim.Name, &q)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get options for dimension: %w", err)
		}

		var options []string
		for _, opt := range opts.Items {
			options = append(options, opt.Label)
		}

		return options, opts.TotalCount, nil
	}

	var hasNoAreaOptions bool
	getAreaOptions := func(dim filter.Dimension) ([]string, int, error) {
		q := filter.QueryParams{Offset: 0, Limit: 500}
		opts, _, err := fc.GetDimensionOptions(ctx, accessToken, "", collectionID, filterID, dim.Name, &q)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get options for dimension: %w", err)
		}

		options := []string{}
		if opts.TotalCount == 0 {
			areas, err := pc.GetAreas(ctx, population.GetAreasInput{
				UserAuthToken: accessToken,
				DatasetID:     filterJob.PopulationType,
				AreaTypeID:    dim.ID,
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

				area, err := pc.GetArea(ctx, population.GetAreaInput{
					UserAuthToken:  accessToken,
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

		if dim.FilterByParent != "" {
			count, err := pc.GetParentAreaCount(ctx, population.GetParentAreaCountInput{
				UserAuthToken:    accessToken,
				DatasetID:        filterJob.PopulationType,
				AreaTypeID:       dim.ID,
				ParentAreaTypeID: dim.FilterByParent,
				Areas:            optsIDs,
			})
			if err != nil {
				log.Error(ctx, "failed to get parent area count", err, log.Data{
					"dataset_id":          filterJob.PopulationType,
					"area_type_id":        dim.ID,
					"parent_area_type_id": dim.FilterByParent,
					"areas":               optsIDs,
				})
				return nil, 0, err
			}
			totalCount = count
		}

		return options, totalCount, nil
	}

	getOptions := func(dim filter.Dimension) ([]string, int, error) {
		if dim.IsAreaType != nil && *dim.IsAreaType {
			return getAreaOptions(dim)
		}

		return getDimensionOptions(dim)
	}

	var fDims []model.FilterDimension
	for _, dim := range filterDims.Items {
		// Needed to determine whether dimension is_area_type
		filterDimension, _, err := fc.GetDimension(ctx, accessToken, "", collectionID, filterID, dim.Name)
		if err != nil {
			log.Error(ctx, "failed to get dimension", err, log.Data{"dimension_name": dim.Name})
			setStatusCode(req, w, err)
			return
		}
		dim.IsAreaType = filterDimension.IsAreaType
		dim.FilterByParent = filterDimension.FilterByParent

		options, count, err := getOptions(dim)
		if err != nil {
			log.Error(ctx, "failed to get options for dimension", err, log.Data{"dimension_name": dim.Name})
			setStatusCode(req, w, err)
			return
		}
		dim.Options = options
		fDims = append(fDims, model.FilterDimension{
			Dimension:    dim,
			OptionsCount: count,
		})
	}

	basePage := rc.NewBasePageModel()
	showAll := req.URL.Query()["showAll"]
	path := req.URL.Path
	m := mapper.CreateFilterFlexOverview(req, basePage, lang, path, showAll, *filterJob, fDims, datasetDims, hasNoAreaOptions)
	rc.BuildPage(w, m, "overview")
}
