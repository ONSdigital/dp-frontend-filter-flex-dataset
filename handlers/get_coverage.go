package handlers

import (
	"net/http"

	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/mapper"
	"github.com/ONSdigital/dp-net/v2/handlers"
	"github.com/gorilla/mux"
)

func GetCoverage(rc RenderClient) http.HandlerFunc {
	return handlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		getCoverage(w, req, rc, lang, accessToken, collectionID)
	})
}

func getCoverage(w http.ResponseWriter, req *http.Request, rc RenderClient, lang, accessToken, collectionID string) {
	vars := mux.Vars(req)
	filterID := vars["filterID"]
	basePage := rc.NewBasePageModel()
	m := mapper.CreateGetCoverage(req, basePage, lang, filterID)
	rc.BuildPage(w, m, "coverage")
}
