package handlers

import (
	"net/http"

	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/config"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/mapper"
	"github.com/ONSdigital/log.go/v2/log"
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
func FilterFlexOverview(cfg config.Config, rc RenderClient) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		filterFlexOverview(w, req, rc, cfg)
	}
}

func filterFlexOverview(w http.ResponseWriter, req *http.Request, rc RenderClient, cfg config.Config) {
	ctx := req.Context()
	basePage := rc.NewBasePageModel()
	m := mapper.CreateFilterFlexOverview(ctx, basePage, cfg)

	rc.BuildPage(w, m, "overview")
}
