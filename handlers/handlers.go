package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/dimension"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/mapper"
	"github.com/ONSdigital/dp-net/v2/handlers"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

func setStatusCode(req *http.Request, w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if err, ok := err.(ClientError); ok {
		status = err.Code()
	}
	log.Error(req.Context(), "setting-response-status", err)
	w.WriteHeader(status)
}

// FilterFlexOverview Handler
func FilterFlexOverview(rc RenderClient, fc FilterClient, dc DatasetClient, dimsc DimensionClient) http.HandlerFunc {
	return handlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		filterFlexOverview(w, req, rc, fc, dc, dimsc, accessToken, collectionID, lang)
	})
}

func filterFlexOverview(w http.ResponseWriter, req *http.Request, rc RenderClient, fc FilterClient, dc DatasetClient, dimsc DimensionClient, accessToken, collectionID, lang string) {
	ctx := req.Context()
	vars := mux.Vars(req)
	filterID := vars["filterID"]
	filterInput := &filter.GetFilterInput{
		FilterID: filterID,
		AuthHeaders: filter.AuthHeaders{
			UserAuthToken: accessToken,
			CollectionID:  collectionID,
		},
	}

	filterJob, err := fc.GetFilter(ctx, *filterInput)
	if err != nil {
		log.Error(ctx, "failed to get filter", err, log.Data{"filter_id": filterID})
		setStatusCode(req, w, err)
		return
	}

	filterDims, _, err := fc.GetDimensions(ctx, accessToken, "", collectionID, filterID, &filter.QueryParams{Limit: 500})
	if err != nil {
		log.Error(ctx, "failed to get dimensions", err, log.Data{"filter_id": filterID})
		setStatusCode(req, w, err)
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

	getAreaOptions := func(dim filter.Dimension) ([]string, error) {
		areas, err := dimsc.GetAreas(ctx, dimension.GetAreasInput{
			UserAuthToken: accessToken,
			DatasetID:     filterJob.PopulationType,
			AreaTypeID:    dim.ID,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get dimension areas: %w", err)
		}

		var options []string
		for _, area := range areas.Areas {
			options = append(options, area.Label)
		}

		return options, nil
	}

	getOptions := func(dim filter.Dimension) ([]string, error) {
		if dim.IsAreaType != nil && *dim.IsAreaType {
			return getAreaOptions(dim)
		}

		return getDimensionOptions(dim)
	}

	for i, dim := range filterDims.Items {
		// Needed to determine whether dimension is_area_type
		filterDimension, _, err := fc.GetDimension(ctx, accessToken, "", collectionID, filterJob.FilterID, dim.Name)
		if err != nil {
			log.Error(ctx, "failed to get dimension", err, log.Data{"dimension_name": dim.Name})
			setStatusCode(req, w, err)
			return
		}
		dim.IsAreaType = filterDimension.IsAreaType

		if len(dim.Options) != 0 {
			continue
		}

		options, err := getOptions(dim)
		if err != nil {
			log.Error(ctx, "failed to get options for dimension", err, log.Data{"dimension_name": dim.Name})
			setStatusCode(req, w, err)
			return
		}

		dim.Options = append(filterDims.Items[i].Options, options...)
		filterDims.Items[i] = dim
	}

	datasetDims, err := dc.GetVersionDimensions(ctx, accessToken, "", collectionID, filterJob.Dataset.DatasetID, filterJob.Dataset.Edition, fmt.Sprint(filterJob.Dataset.Version))
	if err != nil {
		log.Error(ctx, "failed to get versions dimensions", err, log.Data{
			"dataset": filterJob.Dataset.DatasetID,
			"edition": filterJob.Dataset.Edition,
			"version": fmt.Sprint(filterJob.Dataset.Version),
		})
		setStatusCode(req, w, err)
		return
	}

	basePage := rc.NewBasePageModel()
	showAll := req.URL.Query()["showAll"]
	path := req.URL.Path
	m := mapper.CreateFilterFlexOverview(req, basePage, lang, path, showAll, *filterJob, filterDims, datasetDims)
	rc.BuildPage(w, m, "overview")
}
