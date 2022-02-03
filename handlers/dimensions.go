package handlers

import (
	"net/http"

	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/config"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/mapper"
	"github.com/gorilla/mux"
)

// DimensionsSelector Handler
func DimensionsSelector(cfg config.Config, rc RenderClient) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		dimensionsSelector(w, req, rc, cfg)
	}
}

func dimensionsSelector(w http.ResponseWriter, req *http.Request, rc RenderClient, cfg config.Config) {
	ctx := req.Context()
	vars := mux.Vars(req)

	// TODO: Get name from endpoint
	dimension := vars["name"]
	basePage := rc.NewBasePageModel()
	m := mapper.CreateSelector(ctx, req, basePage, cfg, dimension)

	rc.BuildPage(w, m, "selector")
}
