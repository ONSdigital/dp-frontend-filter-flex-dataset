package handlers

import (
	"net/http"

	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/config"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/mapper"
	"github.com/ONSdigital/dp-net/v2/handlers"
	"github.com/gorilla/mux"
)

// DimensionsSelector Handler
func DimensionsSelector(cfg config.Config, rc RenderClient) http.HandlerFunc {
	return handlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		dimensionsSelector(w, req, rc, lang)
	})
}

func dimensionsSelector(w http.ResponseWriter, req *http.Request, rc RenderClient, lang string) {
	vars := mux.Vars(req)

	// TODO: Get name from endpoint
	dimension := vars["name"]
	basePage := rc.NewBasePageModel()
	m := mapper.CreateSelector(req, basePage, dimension, lang)

	rc.BuildPage(w, m, "selector")
}
