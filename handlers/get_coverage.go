package handlers

import (
	"net/http"

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

	geography, _, err := fc.GetDimension(ctx, accessToken, "", collectionID, filterID, "geography")
	if err != nil {
		log.Error(ctx, "failed to get geography dimension", err, log.Data{"filter_id": filterID})
		setStatusCode(req, w, err)
		return
	}

	basePage := rc.NewBasePageModel()
	m := mapper.CreateGetCoverage(req, basePage, lang, filterID, geography.Label)
	rc.BuildPage(w, m, "coverage")
}
