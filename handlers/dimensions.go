package handlers

import (
	"net/http"
	"strconv"

	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/mapper"
	"github.com/ONSdigital/dp-net/v2/handlers"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// DimensionsSelector Handler
func DimensionsSelector(rc RenderClient, fc FilterClient, pc PopulationClient) http.HandlerFunc {
	return handlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		dimensionsSelector(w, req, rc, fc, pc, collectionID, accessToken, lang)
	})
}

func dimensionsSelector(w http.ResponseWriter, req *http.Request, rc RenderClient, fc FilterClient, pc PopulationClient, collectionID, accessToken, lang string) {
	ctx := req.Context()
	vars := mux.Vars(req)
	filterID := vars["filterID"]
	dimensionName := vars["name"]

	logData := log.Data{
		"filter_id": filterID,
	}

	currentFilter, _, err := fc.GetJobState(ctx, accessToken, "", "", collectionID, filterID)
	if err != nil {
		log.Error(ctx, "failed to get job state", err, logData)
		setStatusCode(req, w, err)
		return
	}

	filterDimension, _, err := fc.GetDimension(ctx, accessToken, "", collectionID, filterID, dimensionName)
	if err != nil {
		log.Error(ctx, "failed to find dimension in filter", err, logData)
		setStatusCode(req, w, err)
		return
	}

	basePage := rc.NewBasePageModel()

	if !isAreaType(filterDimension) {
		selector := mapper.CreateSelector(req, basePage, filterDimension.Name, lang, filterID)
		rc.BuildPage(w, selector, "selector")
		return
	}

	// The total_count is the only field required
	opts, _, err := fc.GetDimensionOptions(ctx, accessToken, "", collectionID, filterID, dimensionName, &filter.QueryParams{Limit: 0})
	if err != nil {
		log.Error(ctx, "failed to get options for dimension", err, logData)
		setStatusCode(req, w, err)
		return
	}

	hasOpts := opts.TotalCount > 0

	areaTypes, err := pc.GetAreaTypes(ctx, population.GetAreaTypesInput{
		AuthTokens: population.AuthTokens{
			UserAuthToken: accessToken,
		},
		PaginationParams: population.PaginationParams{
			Limit: 1000,
		},
		PopulationType: currentFilter.PopulationType,
	})
	if err != nil {
		log.Error(ctx, "failed to get population area types", err, logData)
		setStatusCode(req, w, err)
		return
	}

	isValidationError, _ := strconv.ParseBool(req.URL.Query().Get("error"))
	selector := mapper.CreateAreaTypeSelector(req, basePage, lang, filterID, areaTypes.AreaTypes, filterDimension, "LADCD", isValidationError, hasOpts)
	rc.BuildPage(w, selector, "selector")
}

// isAreaType determines if the current dimension is an area type
func isAreaType(dimension filter.Dimension) bool {
	if dimension.IsAreaType == nil {
		return false
	}

	return *dimension.IsAreaType
}
