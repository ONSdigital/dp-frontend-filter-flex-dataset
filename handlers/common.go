package handlers

import (
	"context"

	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"

	"github.com/ONSdigital/log.go/v2/log"
)

type FormAction int

const (
	Unknown FormAction = iota
	CoverageAll
	Add
	Delete
	Search
	Continue
	ParentCoverageSearch
	CoverageDefault = "default"
	NameSearch      = "name-search"
	ParentSearch    = "parent-search"
)

func getReleaseDate(ctx context.Context, dc DatasetClient, userAuthToken, collectionID, datasetID, edition, versionID string) (string, error) {
	var vErr error
	var version, initialVersion dataset.Version

	version, vErr = dc.GetVersion(ctx, userAuthToken, "", "", collectionID, datasetID, edition, versionID)
	if vErr != nil {
		return "", vErr
	}

	if version.Version != 1 {
		initialVersion, vErr = dc.GetVersion(ctx, userAuthToken, "", "", collectionID, datasetID, edition, "1")
		if vErr != nil {
			return "", vErr
		}
		return initialVersion.ReleaseDate, nil
	} else {
		return version.ReleaseDate, nil
	}
}

func setStatusCode(req *http.Request, w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if err, ok := err.(ClientError); ok {
		status = err.Code()
	}
	log.Error(req.Context(), "setting-response-status", err)
	w.WriteHeader(status)
}
