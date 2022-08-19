package handlers

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/mapper"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/model"
	"github.com/ONSdigital/dp-net/v2/handlers"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// GetCoverage handler
func GetCoverage(rc RenderClient, fc FilterClient, pc PopulationClient) http.HandlerFunc {
	return handlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		getCoverage(w, req, rc, fc, pc, lang, accessToken, collectionID)
	})
}

func getCoverage(w http.ResponseWriter, req *http.Request, rc RenderClient, fc FilterClient, pc PopulationClient, lang, accessToken, collectionID string) {
	ctx := req.Context()
	vars := mux.Vars(req)
	filterID := vars["filterID"]
	q := req.URL.Query().Get("q")
	isSearch := strings.Contains(req.URL.RawQuery, "q=")

	filterJob, err := fc.GetFilter(ctx, filter.GetFilterInput{
		FilterID: filterID,
		AuthHeaders: filter.AuthHeaders{
			UserAuthToken: accessToken,
			CollectionID:  collectionID,
		},
	})
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

	var geogLabel, geogID, dimension string
	for _, dim := range filterDims.Items {
		// Needed to determine whether dimension is_area_type
		// Only one dimension will be is_area_type=true
		filterDimension, _, err := fc.GetDimension(ctx, accessToken, "", collectionID, filterID, dim.Name)
		if err != nil {
			log.Error(ctx, "failed to get dimension", err, log.Data{"dimension_name": dim.Name})
			setStatusCode(req, w, err)
			return
		}
		if *filterDimension.IsAreaType {
			geogLabel = filterDimension.Label
			geogID = filterDimension.ID
			dimension = filterDimension.Name
		}
	}

	opts, _, err := fc.GetDimensionOptions(ctx, accessToken, "", collectionID, filterID, dimension, &filter.QueryParams{})
	if err != nil {
		log.Error(ctx, "failed to get dimension options", err, log.Data{"dimension_name": dimension})
		setStatusCode(req, w, err)
		return
	}

	options := []model.SelectableElement{}
	for _, opt := range opts.Items {
		var option model.SelectableElement
		// TODO: Temporary fix until GetArea endpoint is created
		areas, err := pc.GetAreas(ctx, population.GetAreasInput{
			UserAuthToken: accessToken,
			DatasetID:     filterJob.PopulationType,
			AreaTypeID:    geogID,
			Text:          opt.Option,
		})
		if err != nil {
			log.Error(ctx, "failed to get area", err, log.Data{
				"area": geogID,
				"ID":   opt.Option,
			})
			setStatusCode(req, w, err)
			return
		}
		option.Value = opt.Option
		// needed to ensure label matches the ID
		for _, area := range areas.Areas {
			if area.ID == option.Value {
				option.Text = area.Label
				break
			}
		}
		options = append(options, option)
	}

	areas := population.GetAreasResponse{}
	if isSearch && q != "" {
		areas, err = pc.GetAreas(ctx, population.GetAreasInput{
			UserAuthToken: accessToken,
			DatasetID:     filterJob.PopulationType,
			AreaTypeID:    geogID,
			Text:          url.QueryEscape(strings.TrimSpace(q)),
		})
		if err != nil {
			log.Error(ctx, "failed to get areas", err, log.Data{"area": geogID})
			setStatusCode(req, w, err)
			return
		}
	}

	basePage := rc.NewBasePageModel()
	m := mapper.CreateGetCoverage(req, basePage, lang, filterID, geogLabel, q, dimension, areas, options, isSearch)
	rc.BuildPage(w, m, "coverage")
}
