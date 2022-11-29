package handlers

import (
	"context"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
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
