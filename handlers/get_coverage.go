package handlers

import (
	"net/http"

	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/mapper"
	"github.com/ONSdigital/dp-net/v2/handlers"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

func GetCoverage(rc RenderClient, fc FilterClient, p PluralizeClient) http.HandlerFunc {
	return handlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		getCoverage(w, req, rc, fc, p, lang, accessToken, collectionID)
	})
}

func getCoverage(w http.ResponseWriter, req *http.Request, rc RenderClient, fc FilterClient, p PluralizeClient, lang, accessToken, collectionID string) {
	ctx := req.Context()
	vars := mux.Vars(req)
	filterID := vars["filterID"]

	geography, _, err := fc.GetDimension(ctx, accessToken, "", collectionID, filterID, "geography")
	if err != nil {
		log.Error(ctx, "failed to get geography dimension", err, log.Data{"filter_id": filterID})
		setStatusCode(req, w, err)
		return
	}

	geog := p.Plural(geography.Label)

	basePage := rc.NewBasePageModel()
	m := mapper.CreateGetCoverage(req, basePage, lang, filterID, geog)
	rc.BuildPage(w, m, "coverage")
}
