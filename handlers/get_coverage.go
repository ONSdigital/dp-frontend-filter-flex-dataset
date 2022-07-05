package handlers

import (
	"net/http"
	"sort"

	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/mapper"
	"github.com/ONSdigital/dp-net/v2/handlers"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

func GetCoverage(rc RenderClient, fc FilterClient) http.HandlerFunc {
	return handlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		getCoverage(w, req, rc, fc, lang, accessToken, collectionID)
	})
}

func getCoverage(w http.ResponseWriter, req *http.Request, rc RenderClient, fc FilterClient, lang, accessToken, collectionID string) {
	ctx := req.Context()
	vars := mux.Vars(req)
	filterID := vars["filterID"]

	filterDims, _, err := fc.GetDimensions(ctx, accessToken, "", collectionID, filterID, &filter.QueryParams{Limit: 500})
	if err != nil {
		log.Error(ctx, "failed to get dimensions", err, log.Data{"filter_id": filterID})
		setStatusCode(req, w, err)
		return
	}

	for i, dim := range filterDims.Items {
		// Needed to determine whether dimension is_area_type
		filterDimension, _, err := fc.GetDimension(ctx, accessToken, "", collectionID, filterID, dim.Name)
		if err != nil {
			log.Error(ctx, "failed to get dimension", err, log.Data{"dimension_name": dim.Name})
			setStatusCode(req, w, err)
			return
		}
		dim.IsAreaType = filterDimension.IsAreaType
		filterDims.Items[i] = dim
	}

	// Only one dimension with `is_area_type=true'
	sort.Search(len(filterDims.Items), func(i int) bool {
		return *filterDims.Items[i].IsAreaType == true
	})

	basePage := rc.NewBasePageModel()
	m := mapper.CreateGetCoverage(req, basePage, lang, filterID, filterDims.Items[0].Label)
	rc.BuildPage(w, m, "coverage")
}
