package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-net/v2/handlers"
	coreModel "github.com/ONSdigital/dp-renderer/model"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// ChangeDimension Handler
func ChangeDimension(rc RenderClient, fc FilterClient, dimsc DimensionClient) http.HandlerFunc {
	return handlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		changeDimension(w, req, rc, fc, dimsc, collectionID, accessToken, lang)
	})
}

func changeDimension(w http.ResponseWriter, req *http.Request, rc RenderClient, fc FilterClient, dimsc DimensionClient, collectionID, accessToken, lang string) {
	ctx := req.Context()
	vars := mux.Vars(req)
	filterID := vars["filterID"]
	dimensionParam := vars["name"]

	logData := log.Data{
		"filter_id": filterID,
	}

	dimensionName, err := convertDimensionToName(dimensionParam)
	if err != nil {
		log.Error(ctx, "failed to parse dimension name", err, logData)
		setStatusCode(req, w, err)
		return
	}

	form, err := parseChangeDimensionForm(req)
	if err != nil {
		selector, err := createDimensionsSelectorPage(req, rc, fc, dimsc, collectionID, accessToken, lang)
		if err != nil {
			setStatusCode(req, w, err)
			return
		}

		selector.Error = coreModel.Error{
			Title: "Oh no",
		}

		rc.BuildPage(w, selector, "selector")
		return
	}

	dimension := filter.Dimension{
		Name:       form.Dimension,
		IsAreaType: toBoolPtr(form.IsAreaType),
	}

	if _, _, err = fc.UpdateDimensions(ctx, accessToken, "", collectionID, filterID, dimensionName, "", dimension); err != nil {
		log.Error(ctx, "error updating filter dimension", err, logData)
		setStatusCode(req, w, err)
		return
	}

	http.Redirect(w, req, fmt.Sprintf("/filters/%s/dimensions/", filterID), http.StatusMovedPermanently)
}

// changeDimensionForm represents form-data for the ChangeDimension handler.
type changeDimensionForm struct {
	Dimension  string
	IsAreaType bool
}

// parseChangeDimensionForm parses form data from a http.Request into a changeDimensionForm.
func parseChangeDimensionForm(req *http.Request) (changeDimensionForm, error) {
	err := req.ParseForm()
	if err != nil {
		return changeDimensionForm{}, fmt.Errorf("erorr parsing form: %w", err)
	}

	dimension := req.FormValue("dimension")
	if dimension == "" {
		return changeDimensionForm{}, &validationErr{errors.New("missing required value 'dimension'")}
	}

	areaType, err := strconv.ParseBool(req.FormValue("is_area_type"))
	if err != nil {
		return changeDimensionForm{}, &validationErr{errors.New("missing or invalid value 'is_area_type', expected bool")}
	}

	return changeDimensionForm{
		Dimension:  dimension,
		IsAreaType: areaType,
	}, nil
}
