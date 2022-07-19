package handlers

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/mapper"
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

	var geogLabel, geogID string
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
		}
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
	m := mapper.CreateGetCoverage(req, basePage, lang, filterID, geogLabel, q, areas, isSearch)
	rc.BuildPage(w, m, "coverage")
}
