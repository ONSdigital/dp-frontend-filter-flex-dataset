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

// ChangeDimension Handler
func PostChangeDimensions(fc FilterClient) http.HandlerFunc {
	return handlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		postChangeDimensions(w, req, fc, accessToken, collectionID)
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

	switch form.Action {
	case Search:
		v.Set("q", form.Value)
	}
	req.URL.RawQuery = v.Encode()
	http.Redirect(w, req, fmt.Sprint(req.URL), http.StatusSeeOther)
}

// changeDimensionsForm represents form-data for the UpdateCoverage handler.
type changeDimensionsForm struct {
	Action               FormAction
	Value, PrimaryAction string
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

	delOption := req.FormValue("delete-dimension")
	if delOption != "" {
		action = Delete
		value = delOption
	}

	isSearch, _ := strconv.ParseBool(req.FormValue("is-search"))
	if isSearch {
		action = Search
		q := req.FormValue("q")
		value = q
	}

	return changeDimensionsForm{
		Action:        action,
		Value:         value,
		PrimaryAction: dimensions,
	}, nil
}
