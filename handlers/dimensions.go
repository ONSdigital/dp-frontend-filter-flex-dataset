package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"unicode"

	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/mapper"
	"github.com/ONSdigital/dp-net/v2/handlers"
	coreModel "github.com/ONSdigital/dp-renderer/model"
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

	basePage := rc.NewBasePageModel()

	if req.Method == http.MethodPost {
		if err := changeDimensionFn(req, fc, accessToken, collectionID); err == nil {
			http.Redirect(w, req, fmt.Sprintf("/filters/%s/dimensions/", filterID), http.StatusMovedPermanently)
		}

		basePage.Error = coreModel.Error{Title: "oh no"}
	}

	currentFilter, _, err := fc.GetJobState(ctx, accessToken, "", "", collectionID, filterID)
	if err != nil {
		log.Error(ctx, "failed to get job state", err, logData)
		setStatusCode(req, w, err)
		return
	}

	dimensionName, err := convertDimensionToName(nameParam)
	if err != nil {
		log.Error(ctx, "failed to parse dimension name", err, logData)
		setStatusCode(req, w, err)
		return
	}

	filterDimension, _, err := fc.GetDimension(ctx, accessToken, "", collectionID, filterID, dimensionName)
	if err != nil {
		log.Error(ctx, "failed to find dimension in filter", err, logData)
		setStatusCode(req, w, err)
		return
	}

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

	selector := mapper.CreateAreaTypeSelector(req, basePage, lang, areaTypes.AreaTypes, dimensionName)
	rc.BuildPage(w, selector, "selector")
}

func changeDimensionFn(req *http.Request, fc FilterClient, accessToken, collectionID string) error {
	ctx := req.Context()
	vars := mux.Vars(req)
	filterID := vars["filterID"]
	dimensionParam := vars["name"]

	dimensionName, err := convertDimensionToName(dimensionParam)
	if err != nil {
		return err
	}

	form, err := parseChangeDimensionForm(req)
	if err != nil {
		return err
	}

	dimension := filter.Dimension{
		Name:       form.Dimension,
		IsAreaType: toBoolPtr(form.IsAreaType),
	}

	if _, _, err = fc.UpdateDimensions(ctx, accessToken, "", collectionID, filterID, dimensionName, "", dimension); err != nil {
		return err
	}

	return nil
}

// isAreaType determines if the current dimension is an area type
func isAreaType(dimension filter.Dimension) bool {
	if dimension.IsAreaType == nil {
		return false
	}

	return *dimension.IsAreaType
}

// convertDimensionToName takes a URL-coded param for a dimension name and attempts to convert it to the
// pretty label. This is temporary/best-effort, and only a stopgap until we start using ID's/names everywhere.
// Changing everything to use these will involve updating several services, including the importer, so to
// restore the journey right now we're taking the hacky approach.
// Example: `number+of+siblings+%283+mappings%29` -> `Number of siblings (3 mappings)`
func convertDimensionToName(param string) (string, error) {
	selection, err := url.QueryUnescape(param)
	if err != nil {
		return "", err
	}

	// Sentence case the param
	runes := []rune(strings.ToLower(selection))
	return string(append([]rune{unicode.ToUpper(runes[0])}, runes[1:]...)), nil
}
