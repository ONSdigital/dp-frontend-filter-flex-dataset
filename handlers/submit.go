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

	filter, eTag, err := fc.GetJobState(ctx, accessToken, "", "", collectionID, filterID)
	if err != nil {
		log.Error(ctx, "failed to get job state", err, log.Data{"filter_id": filterID})
		setStatusCode(req, w, err)
		return
	}

	// etag is coming back empty??

	mdl, _, err := fc.UpdateFlexBlueprint(ctx, accessToken, "", "", collectionID, filter, true, filter.PopulationType, eTag)
	if err != nil {
		log.Error(ctx, "failed to submit filter blueprint", err, log.Data{"filter_id": filterID})
		setStatusCode(req, w, err)
		return
	}

	filterOutputID := mdl.Links.FilterOutputs.ID

	// redirect to dataset controller
	http.Redirect(w, req, fmt.Sprintf("/filter-outputs/%s", filterOutputID), http.StatusFound)
}
