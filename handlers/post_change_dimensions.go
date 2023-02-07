package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ONSdigital/dp-net/v2/handlers"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// PostChangeDimensions Handler
func (f *FilterFlex) PostChangeDimensions() http.HandlerFunc {
	return handlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		postChangeDimensions(w, req, f.FilterClient, accessToken, collectionID)
	})
}

func postChangeDimensions(w http.ResponseWriter, req *http.Request, fc FilterClient, accessToken, collectionID string) {
	ctx := req.Context()
	vars := mux.Vars(req)
	filterID := vars["filterID"]

	form, err := parseChangeDimensionsForm(req)
	if err != nil {
		log.Error(ctx, "failed to parse change dimensions form", err, log.Data{
			"filter_id": filterID,
		})
		setStatusCode(req, w, err)
		return
	}

	v := url.Values{}
	v.Set("f", form.PrimaryAction)
	if form.PrimaryAction == "search" {
		v.Set("q", form.SearchQ)
	}

	switch form.Action {
	case Add:
		_, err := fc.AddFlexDimension(ctx, accessToken, "", collectionID, filterID, form.Value, []string{}, false, "")
		if err != nil {
			log.Error(ctx, "failed to add flex dimension", err, log.Data{
				"filter_id": filterID,
				"name":      form.Value,
			})
			setStatusCode(req, w, err)
			return
		}
	case Delete:
		_, err := fc.RemoveDimension(ctx, accessToken, "", collectionID, filterID, form.Value, "")
		if err != nil {
			log.Error(ctx, "failed to remove dimension", err, log.Data{
				"filter_id": filterID,
				"name":      form.Value,
			})
			setStatusCode(req, w, err)
			return
		}
	}
	req.URL.RawQuery = v.Encode()
	http.Redirect(w, req, fmt.Sprint(req.URL), http.StatusSeeOther)
}

// changeDimensionsForm represents form-data for the UpdateCoverage handler.
type changeDimensionsForm struct {
	Action                        FormAction
	Value, PrimaryAction, SearchQ string
}

// parseChangeDimensionsForm parses form data from a http.Request into a updateCoverageForm.
func parseChangeDimensionsForm(req *http.Request) (changeDimensionsForm, error) {
	var action FormAction
	var value string

	err := req.ParseForm()
	if err != nil {
		return changeDimensionsForm{}, fmt.Errorf("error parsing form: %w", err)
	}

	dimensions := req.FormValue("dimensions")
	if dimensions == "" {
		return changeDimensionsForm{}, &clientErr{errors.New("missing required value 'dimensions'")}
	}

	addOption := req.FormValue("add-dimension")
	if addOption != "" {
		action = Add
		value = addOption
	}

	delOption := req.FormValue("delete-option")
	if delOption != "" {
		action = Delete
		value = delOption
	}

	return changeDimensionsForm{
		Action:        action,
		Value:         value,
		PrimaryAction: dimensions,
		SearchQ:       req.FormValue("q"),
	}, nil
}
