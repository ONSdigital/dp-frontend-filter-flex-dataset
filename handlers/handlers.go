package handlers

import (
	"net/http"
	"strconv"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
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
	filterInput := &filter.GetFilterInput{
		FilterID: filterID,
		AuthHeaders: filter.AuthHeaders{
			UserAuthToken: accessToken,
			CollectionID:  collectionID,
		},
	}

	filterJob, err := fc.GetFilter(ctx, *filterInput)
	if err != nil {
		log.Error(ctx, "failed to get filter", err, log.Data{"filter_id": filterID})
		setStatusCode(req, w, err)
		return
	}

	dims, _, err := fc.GetDimensions(ctx, accessToken, "", collectionID, filterID, &filter.QueryParams{Limit: 500})
	if err != nil {
		log.Error(ctx, "failed to get dimensions", err, log.Data{"filter_id": filterID})
		setStatusCode(req, w, err)
		return
	}

	for i, dim := range dims.Items {
		// Needed to determine whether demension is_area_type
		filterDimension, _, err := fc.GetDimension(ctx, accessToken, "", collectionID, filterJob.FilterID, dim.Name)
		if err != nil {
			log.Error(ctx, "failed to get dimension", err, log.Data{"dimension_name": dim.Name})
			setStatusCode(req, w, err)
			return
		}
		dims.Items[i].IsAreaType = filterDimension.IsAreaType
		if len(dim.Options) == 0 {
			q := dataset.QueryParams{Offset: 0, Limit: 1000}
			opts, err := dc.GetOptions(ctx, accessToken, "", collectionID, filterJob.Dataset.DatasetID, filterJob.Dataset.Edition, strconv.Itoa(filterJob.Dataset.Version), dim.Name, &q)
			if err != nil {
				log.Error(ctx, "failed to get options for dimension", err, log.Data{"dimension_name": dim.Name})
				setStatusCode(req, w, err)
				return
			}
			for _, opt := range opts.Items {
				dims.Items[i].Options = append(dims.Items[i].Options, opt.Label)
			}
		}
	}

	basePage := rc.NewBasePageModel()
	showAll := req.URL.Query()["showAll"]
	path := req.URL.Path
	m := mapper.CreateFilterFlexOverview(req, basePage, lang, path, showAll, *filterJob, dims)
	rc.BuildPage(w, m, "overview")
}
