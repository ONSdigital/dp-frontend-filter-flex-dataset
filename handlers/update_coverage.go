package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-net/v2/handlers"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// UpdateCoverage Handler
func UpdateCoverage(fc FilterClient) http.HandlerFunc {
	return handlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		updateCoverage(w, req, fc, accessToken, collectionID)
	})
}

func updateCoverage(w http.ResponseWriter, req *http.Request, fc FilterClient, accessToken, collectionID string) {
	ctx := req.Context()
	vars := mux.Vars(req)
	filterID := vars["filterID"]

	form, err := parseUpdateCoverageForm(req)
	if err != nil {
		log.Error(ctx, "failed to parse update coverage form", err, log.Data{
			"filter_id": filterID,
		})
		setStatusCode(req, w, err)
		return
	}

	if _, err = fc.AddDimensionValue(ctx, accessToken, "", collectionID, filterID, form.Dimension, form.Option, ""); err != nil {
		log.Error(ctx, "failed to add dimension value", err, log.Data{
			"dimension": form.Dimension,
			"option":    form.Option,
		})
		setStatusCode(req, w, err)
		return
	}

	http.Redirect(w, req, fmt.Sprint(req.URL), http.StatusMovedPermanently)
}

// updateCoverageForm represents form-data for the UpdateCoverage handler.
type updateCoverageForm struct {
	Dimension string
	Option    string
}

// parseUpdateCoverageForm parses form data from a http.Request into a updateCoverageForm.
func parseUpdateCoverageForm(req *http.Request) (updateCoverageForm, error) {
	err := req.ParseForm()
	if err != nil {
		return updateCoverageForm{}, fmt.Errorf("error parsing form: %w", err)
	}

	dimension := req.FormValue("dimension")
	if dimension == "" {
		return updateCoverageForm{}, &clientErr{errors.New("missing required value 'dimension'")}
	}

	option := req.FormValue("option")
	if option == "" {
		return updateCoverageForm{}, &clientErr{errors.New("missing required value 'option'")}
	}

	return updateCoverageForm{
		Dimension: dimension,
		Option:    option,
	}, nil
}
