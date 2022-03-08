package handlers

import (
	"net/http"

	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/config"
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
func FilterFlexOverview(cfg config.Config, rc RenderClient, fc FilterClient) http.HandlerFunc {
	return handlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		filterFlexOverview(w, req, rc, fc, accessToken, collectionID, lang)
	})
}

func filterFlexOverview(w http.ResponseWriter, req *http.Request, rc RenderClient, fc FilterClient, accessToken, collectionID, lang string) {
	ctx := req.Context()
	vars := mux.Vars(req)
	filterID := vars["filterID"]

	dims, _, err := fc.GetDimensions(ctx, accessToken, "", collectionID, filterID, nil)
	if err != nil {
		log.Error(ctx, "failed to get dimensions", err, log.Data{"filter_id": filterID})
		setStatusCode(req, w, err)
		return
	}

	basePage := rc.NewBasePageModel()
	showAll := req.URL.Query()["showAll"]
	m := mapper.CreateFilterFlexOverview(req, basePage, lang, dims, showAll)
	rc.BuildPage(w, m, "overview")
}
