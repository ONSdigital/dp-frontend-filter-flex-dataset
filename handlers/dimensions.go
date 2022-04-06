package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/mapper"
	"github.com/ONSdigital/dp-net/v2/handlers"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// DimensionsSelector Handler
func DimensionsSelector(rc RenderClient, fc FilterClient, dimsc DimensionClient) http.HandlerFunc {
	return handlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		dimensionsSelector(w, req, rc, fc, dimsc, collectionID, accessToken, lang)
	})
}

func dimensionsSelector(w http.ResponseWriter, req *http.Request, rc RenderClient, fc FilterClient, dimsc DimensionClient, collectionID, accessToken, lang string) {
	ctx := req.Context()
	vars := mux.Vars(req)
	filterID := vars["filterID"]
	nameParam := vars["name"]

	logData := log.Data{
		"filter_id": filterID,
	}

	currentFilter, _, err := fc.GetJobState(ctx, accessToken, "", "", collectionID, filterID)
	if err != nil {
		log.Error(ctx, "failed to get job state", err, logData)
		setStatusCode(req, w, err)
		return
	}

	filterDimension, err := findDimension(nameParam, currentFilter.Dimensions)
	if err != nil {
		log.Error(ctx, "failed to find dimension in filter", err, logData)
		setStatusCode(req, w, err)
		return
	}

	basePage := rc.NewBasePageModel()

	if !isAreaType(filterDimension) {
		selector := mapper.CreateSelector(req, basePage, filterDimension.Name, lang)
		rc.BuildPage(w, selector, "selector")
		return
	}

	areaTypes, err := dimsc.GetAreaTypes(ctx, accessToken, "", currentFilter.PopulationType)
	if err != nil {
		log.Error(ctx, "failed to get geography dimensions", err, logData)
		setStatusCode(req, w, err)
		return
	}

	selector := mapper.CreateAreaTypeSelector(req, basePage, lang, areaTypes.AreaTypes, nameParam)
	rc.BuildPage(w, selector, "selector")
}

// isAreaType determines if the current dimension is an area type
func isAreaType(dimension *filter.ModelDimension) bool {
	if dimension == nil {
		return false
	}

	if dimension.IsAreaType == nil {
		return false
	}

	return *dimension.IsAreaType
}

// findDimension attempts to find a dimension based on the dimension value in a URL param.
func findDimension(selectionParam string, dimensions []filter.ModelDimension) (*filter.ModelDimension, error) {
	// Params are matched by name, and therefore may contain spaces/escaped punctuation.
	selection, err := url.QueryUnescape(selectionParam)
	if err != nil {
		return nil, fmt.Errorf("error escaping selection (%s): %w", selectionParam, err)
	}

	for _, dimension := range dimensions {
		if err == nil && selection == strings.ToLower(dimension.Name) {
			return &dimension, nil
		}
	}

	return nil, dimensionNotFoundErr{dimensions, selection}
}

// dimensionNotFoundErr is an error provided when a matching dimension cannot be found
// in the current filter using a dimension route parameter.
type dimensionNotFoundErr struct {
	list      []filter.ModelDimension
	dimension string
}

func (d dimensionNotFoundErr) Error() string {
	return fmt.Sprintf("could not find dimension with name %s in list (%+v)", d.dimension, d.list)
}

func (d dimensionNotFoundErr) Code() int {
	return http.StatusNotFound
}
