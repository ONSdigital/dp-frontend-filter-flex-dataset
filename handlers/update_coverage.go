package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ONSdigital/dp-net/v2/handlers"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

type FormAction int

const (
	Unknown FormAction = iota
	CoverageAll
	Add
	Delete
	Search
	CoverageSearch
	CoverageDefault = "default"
	NameSearch      = "name-search"
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

	switch form.Action {
	case CoverageSearch:
		http.Redirect(w, req, fmt.Sprintf("/filters/%s/dimensions", filterID), http.StatusMovedPermanently)
		return
	case Search:
		v := url.Values{}
		v.Set("q", form.Value)
		req.URL.RawQuery = v.Encode()
	case Delete:
		_, err := fc.RemoveDimensionValue(ctx, accessToken, "", collectionID, filterID, form.Dimension, form.Value, "")
		if err != nil {
			log.Error(ctx, "failed to remove dimension value", err, log.Data{
				"dimension": form.Dimension,
				"option":    form.Value,
			})
			setStatusCode(req, w, err)
			return
		}
	case Add:
		_, err := fc.AddDimensionValue(ctx, accessToken, "", collectionID, filterID, form.Dimension, form.Value, "")
		if err != nil {
			log.Error(ctx, "failed to add dimension value", err, log.Data{
				"dimension": form.Dimension,
				"option":    form.Value,
			})
			setStatusCode(req, w, err)
			return
		}
	case CoverageAll:
		_, err := fc.DeleteDimensionOptions(ctx, accessToken, "", collectionID, filterID, form.Dimension)
		if err != nil {
			log.Error(ctx, "failed to delete dimension options", err, log.Data{
				"dimension": form.Dimension,
			})
			setStatusCode(req, w, err)
			return
		}
		http.Redirect(w, req, fmt.Sprintf("/filters/%s/dimensions", filterID), http.StatusMovedPermanently)
		return
	}

	http.Redirect(w, req, fmt.Sprint(req.URL), http.StatusMovedPermanently)
}

// updateCoverageForm represents form-data for the UpdateCoverage handler.
type updateCoverageForm struct {
	Action    FormAction
	Value     string
	Dimension string
}

// parseUpdateCoverageForm parses form data from a http.Request into a updateCoverageForm.
func parseUpdateCoverageForm(req *http.Request) (updateCoverageForm, error) {
	var action FormAction
	var value string

	err := req.ParseForm()
	if err != nil {
		return updateCoverageForm{}, fmt.Errorf("error parsing form: %w", err)
	}

	dimension := req.FormValue("dimension")
	if dimension == "" {
		return updateCoverageForm{}, &clientErr{errors.New("missing required value 'dimension'")}
	}

	coverage := req.FormValue("coverage")
	if coverage == "" {
		return updateCoverageForm{}, &clientErr{errors.New("missing required value 'coverage'")}
	}

	switch coverage {
	case CoverageDefault:
		action = CoverageAll
		value = coverage
	case NameSearch:
		action = CoverageSearch
		value = coverage
	default:
		return updateCoverageForm{}, &clientErr{errors.New("unknown coverage type")}
	}

	isSearch, err := strconv.ParseBool(req.FormValue("is-search"))
	if isSearch {
		q := req.FormValue("q")
		action = Search
		value = q
	}

	addOption := req.FormValue("add-option")
	if addOption != "" {
		action = Add
		value = addOption
	}

	deleteOption := req.FormValue("delete-option")
	if deleteOption != "" {
		action = Delete
		value = deleteOption
	}

	return updateCoverageForm{
		Action:    action,
		Value:     value,
		Dimension: dimension,
	}, nil
}
