package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/helpers"

	"github.com/ONSdigital/dp-net/v2/handlers"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// Submit filter outputs handler
func (f *FilterFlex) Submit() http.HandlerFunc {
	return handlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		submit(w, req, accessToken, collectionID, f.FilterClient)
	})
}

func submit(w http.ResponseWriter, req *http.Request, accessToken, collectionID string, fc FilterClient) {
	vars := mux.Vars(req)
	filterID := vars["filterID"]
	ctx := req.Context()

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

	filterRequest := &filter.SubmitFilterRequest{
		FilterID:       filterJob.FilterID,
		PopulationType: filterJob.PopulationType,
	}
	resp, _, err := fc.SubmitFilter(ctx, accessToken, "", "", filterJob.ETag, *filterRequest)
	if err != nil {
		log.Error(ctx, "failed to submit filter", err, log.Data{"submit_filter_request": filterRequest})
		setStatusCode(req, w, err)
		return
	}

	dataset := filterJob.Dataset
	dsID := dataset.DatasetID
	ed := dataset.Edition
	v := strconv.Itoa(dataset.Version)
	foID := resp.FilterOutputID

	isCustom := helpers.IsBoolPtr(filterJob.Custom)
	if isCustom {
		http.Redirect(w, req, fmt.Sprintf("/datasets/create/filter-outputs/%s#get-data", foID), http.StatusFound)
	} else {
		http.Redirect(w, req, fmt.Sprintf("/datasets/%s/editions/%s/versions/%s/filter-outputs/%s#get-data", dsID, ed, v, foID), http.StatusFound)
	}
}
