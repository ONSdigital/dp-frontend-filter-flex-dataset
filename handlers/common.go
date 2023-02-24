package handlers

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"

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

// getZebContent is a helper function that returns the homepage content required to map the emergency banner and service message
func getZebContent(ctx context.Context, zc ZebedeeClient, userAuthToken, collectionID, lang string) (zebedee.EmergencyBanner, string, error) {
	hpc, err := zc.GetHomepageContent(ctx, userAuthToken, collectionID, lang, "/")
	return hpc.EmergencyBanner, hpc.ServiceMessage, err
}

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

// isMultivariateDataset determines whether the given filter record is based on a multivariate dataset type
func isMultivariateDataset(ctx context.Context, dc DatasetClient, accessToken, collectionID, did string) (bool, error) {
	d, err := dc.Get(ctx, accessToken, "", collectionID, did)
	if err != nil {
		return false, fmt.Errorf("failed to get dataset: %w", err)
	}

	if strings.Contains(d.Type, "multivariate") {
		return true, nil
	}
	return false, nil
}

func setStatusCode(req *http.Request, w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if err, ok := err.(ClientError); ok {
		status = err.Code()
	}
	log.Error(req.Context(), "setting-response-status", err)
	w.WriteHeader(status)
}

// getBlockedAreaCount is a helper function that does the required sorting and checks before making the api request
func (f *FilterFlex) getBlockedAreaCount(ctx context.Context, accessToken, populationType, areaTypeID, parent string, dimensionIds, areaOptions []string) (*cantabular.GetBlockedAreaCountResult, error) {
	sort.Slice(dimensionIds, func(i, j int) bool {
		return dimensionIds[i] == areaTypeID || dimensionIds[i] == parent
	})

	if parent != "" {
		areaTypeID = parent
	}

	// set default coverage
	if len(areaOptions) == 0 {
		areaOptions = []string{"K04000001"}
		areaTypeID = "nat"
	}
	sdc, err := f.PopulationClient.GetBlockedAreaCount(ctx, population.GetBlockedAreaCountInput{
		AuthTokens: population.AuthTokens{
			UserAuthToken: accessToken,
		},
		PopulationType: populationType,
		Variables:      dimensionIds,
		Filter: population.Filter{
			Codes:    areaOptions,
			Variable: areaTypeID,
		},
	})
	return sdc, err
}
