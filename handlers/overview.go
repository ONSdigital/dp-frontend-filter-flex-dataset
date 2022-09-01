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
	var fErr, dErr, fdsErr, fdErr error

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

	getDimensionOptions := func(dim filter.Dimension) ([]string, error) {
		q := dataset.QueryParams{Offset: 0, Limit: 1000}

		opts, err := dc.GetOptions(ctx, accessToken, "", collectionID, filterJob.Dataset.DatasetID, filterJob.Dataset.Edition, strconv.Itoa(filterJob.Dataset.Version), dim.Name, &q)
		if err != nil {
			return nil, fmt.Errorf("failed to get options for dimension: %w", err)
		}

		var options []string
		for _, opt := range opts.Items {
			options = append(options, opt.Label)
		}

		return options, nil
	}

	var hasNoAreaOptions bool
	getAreaOptions := func(dim filter.Dimension) ([]string, error) {
		opts, _, err := fc.GetDimensionOptions(ctx, accessToken, "", collectionID, filterID, dim.Name, &filter.QueryParams{})
		if err != nil {
			return nil, fmt.Errorf("failed to get options for dimension: %w", err)
		}

		options := []string{}
		if opts.TotalCount == 0 {
			areas, err := pc.GetAreas(ctx, population.GetAreasInput{
				UserAuthToken: accessToken,
				DatasetID:     filterJob.PopulationType,
				AreaTypeID:    dim.ID,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to get dimension areas: %w", err)
			}

			for _, area := range areas.Areas {
				options = append(options, area.Label)
			}

			hasNoAreaOptions = true
			return options, nil
		}

		var wg sync.WaitGroup
		for _, opt := range opts.Items {
			wg.Add(1)

			go func(opt filter.DimensionOption) {
				defer wg.Done()
				var areaTypeID string
				if dim.FilterByParent != "" {
					areaTypeID = dim.FilterByParent
				} else {
					areaTypeID = dim.ID
				}
				// TODO: Temporary fix until GetArea endpoint is created
				areas, err := pc.GetAreas(ctx, population.GetAreasInput{
					UserAuthToken: accessToken,
					DatasetID:     filterJob.PopulationType,
					AreaTypeID:    areaTypeID,
					Text:          opt.Option,
				})
				if err != nil {
					log.Error(ctx, "failed to get area", err, log.Data{
						"area": dim.ID,
						"ID":   opt.Option,
					})
					setStatusCode(req, w, err)
					return
				}

				for _, area := range areas.Areas {
					if area.ID == opt.Option {
						options = append(options, area.Label)
						break
					}
				}
			}(opt)
		}
		wg.Wait()

		return options, nil
	}

	getOptions := func(dim filter.Dimension) ([]string, error) {
		if dim.IsAreaType != nil && *dim.IsAreaType {
			return getAreaOptions(dim)
		}

		return getDimensionOptions(dim)
	}

	for i, dim := range filterDims.Items {
		wg.Add(1)
		go func(dim filter.Dimension, i int) {
			defer wg.Done()
			var filterDimension filter.Dimension
			// Needed to determine whether dimension is_area_type
			filterDimension, _, fdErr = fc.GetDimension(ctx, accessToken, "", collectionID, filterID, dim.Name)
			if fdErr != nil {
				log.Error(ctx, "failed to get dimension", fdErr, log.Data{"dimension_name": dim.Name})
				setStatusCode(req, w, fdErr)
				return
			}
			dim.IsAreaType = filterDimension.IsAreaType
			dim.FilterByParent = filterDimension.FilterByParent

			options, err := getOptions(dim)
			if err != nil {
				log.Error(ctx, "failed to get options for dimension", err, log.Data{"dimension_name": dim.Name})
				setStatusCode(req, w, err)
				return
			}

			dim.Options = append(filterDims.Items[i].Options, options...)
			filterDims.Items[i] = dim
		}(dim, i)
	}
	wg.Wait()

	if fdErr != nil {
		log.Error(ctx, "failed to get dimension", fdErr, log.Data{"filter_id": filterID})
		setStatusCode(req, w, fdErr)
		return
	}

	basePage := rc.NewBasePageModel()
	showAll := req.URL.Query()["showAll"]
	path := req.URL.Path
	m := mapper.CreateFilterFlexOverview(req, basePage, lang, path, showAll, *filterJob, filterDims, datasetDims, hasNoAreaOptions)
	rc.BuildPage(w, m, "overview")
}
