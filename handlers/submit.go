package handlers

import (
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-net/v2/handlers"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// Submit filter outputs handler
func Submit(fc FilterClient) http.HandlerFunc {
	return handlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		submit(w, req, accessToken, collectionID, fc)
	})
}

func submit(w http.ResponseWriter, req *http.Request, accessToken, collectionID string, fc FilterClient) {
	vars := mux.Vars(req)
	filterID := vars["filterID"]
	ctx := req.Context()

	filter, _, err := fc.GetJobState(ctx, accessToken, "", "", collectionID, filterID)
	if err != nil {
		log.Error(ctx, "failed to get job state", err, log.Data{"filter_id": filterID})
		setStatusCode(req, w, err)
		return
	}

	// TODO: this endpoint is changing, update to specific submit method
	mdl, _, err := fc.UpdateFlexBlueprint(ctx, accessToken, "", "", collectionID, filter, true, filter.PopulationType, "")
	if err != nil {
		log.Error(ctx, "failed to submit filter blueprint", err, log.Data{"filter_id": filterID})
		setStatusCode(req, w, err)
		return
	}

	dsID := filter.DatasetID
	ed := filter.Edition
	v := filter.Version
	foID := mdl.Links.FilterOutputs.ID
	http.Redirect(w, req, fmt.Sprintf("/datasets/%s/editions/%s/versions/%s/filter-outputs/%s", dsID, ed, v, foID), http.StatusFound)
}
