package handlers

import (
	"net/http"
	"strconv"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
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
func FilterFlexOverview(rc RenderClient, fc FilterClient, dc DatasetClient) http.HandlerFunc {
	return handlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		filterFlexOverview(w, req, rc, fc, dc, accessToken, collectionID, lang)
	})
}

func filterFlexOverview(w http.ResponseWriter, req *http.Request, rc RenderClient, fc FilterClient, dc DatasetClient, accessToken, collectionID, lang string) {
	ctx := req.Context()
	vars := mux.Vars(req)
	filterID := vars["filterID"]

	filterJob, _, err := fc.GetJobState(ctx, accessToken, "", "", collectionID, filterID)
	if err != nil {
		log.Error(ctx, "failed to get job state", err, log.Data{"filter_id": filterID})
		setStatusCode(req, w, err)
		return
	}

	for i, dim := range filterJob.Dimensions {
		if len(dim.Options) == 0 {
			q := dataset.QueryParams{Offset: 0, Limit: 1000}
			opts, err := dc.GetOptions(ctx, accessToken, "", collectionID, filterJob.Dataset.DatasetID, filterJob.Dataset.Edition, strconv.Itoa(filterJob.Dataset.Version), dim.Name, &q)
			if err != nil {
				log.Error(ctx, "failed to get options for dimension", err, log.Data{"dimension_name": dim.Name})
				setStatusCode(req, w, err)
				return
			}
			for _, opt := range opts.Items {
				filterJob.Dimensions[i].Options = append(filterJob.Dimensions[i].Options, opt.Label)
			}
		}
	}

	basePage := rc.NewBasePageModel()
	showAll := req.URL.Query()["showAll"]
	path := req.URL.Path
	m := mapper.CreateFilterFlexOverview(req, basePage, lang, path, showAll, filterJob)
	rc.BuildPage(w, m, "overview")
}
