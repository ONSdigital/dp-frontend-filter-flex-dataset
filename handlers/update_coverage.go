package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/helpers"
	"github.com/ONSdigital/dp-net/v2/handlers"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// UpdateCoverage Handler
func (f *FilterFlex) UpdateCoverage() http.HandlerFunc {
	return handlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		updateCoverage(w, req, f.FilterClient, accessToken, collectionID)
	})
}

func updateCoverage(w http.ResponseWriter, req *http.Request, fc FilterClient, accessToken, collectionID string) {
	ctx := req.Context()
	vars := mux.Vars(req)
	filterID := vars["filterID"]

	form, err := parseUpdateCoverageForm(req)
	if isValidationErr(err) {
		v := url.Values{}
		v.Set("c", ParentSearch)
		v.Set("error", "true")
		req.URL.RawQuery = v.Encode()
		http.Redirect(w, req, fmt.Sprint(req.URL), http.StatusMovedPermanently)
		return
	}
	if err != nil {
		log.Error(ctx, "failed to parse update coverage form", err, log.Data{
			"filter_id": filterID,
		})
		setStatusCode(req, w, err)
		return
	}

	switch form.Action {
	case Continue:
		http.Redirect(w, req, fmt.Sprintf("/filters/%s/dimensions", filterID), http.StatusMovedPermanently)
		return
	case Search:
		v := url.Values{}
		if form.Coverage == ParentSearch {
			v.Set("p", form.LargerArea)
			v.Set("pq", form.Value)
		} else {
			v.Set("q", form.Value)
		}
		v.Set("c", form.Coverage)
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
		opts, _, err := fc.GetDimensionOptions(ctx, accessToken, "", collectionID, filterID, form.Dimension, &filter.QueryParams{Limit: 500})
		if err != nil {
			log.Error(ctx, "failed to get dimension options", err, log.Data{"dimension_name": form.Dimension})
			setStatusCode(req, w, err)
			return
		}

		if opts.TotalCount > 0 && form.Coverage != form.OptionType || opts.TotalCount > 0 && form.SetParent != form.LargerArea {
			log.Info(ctx, "invalid options combination, removing existing options", log.Data{"filter_id": filterID})
			_, err := fc.DeleteDimensionOptions(ctx, accessToken, "", collectionID, filterID, form.Dimension)
			if err != nil {
				log.Error(ctx, "failed to delete dimension options", err, log.Data{
					"dimension": form.Dimension,
				})
				setStatusCode(req, w, err)
				return
			}
			opts = filter.DimensionOptions{}
		}

		var options []string
		for _, opt := range opts.Items {
			options = append(options, opt.Option)
		}
		options = append(options, form.Value)

		dim := filter.Dimension{
			Name:           form.Dimension,
			ID:             form.GeographyID,
			IsAreaType:     helpers.ToBoolPtr(true),
			Options:        options,
			FilterByParent: form.LargerArea,
		}
		_, _, err = fc.UpdateDimensions(ctx, accessToken, "", collectionID, filterID, form.Dimension, "", dim)
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

	switch form.Coverage {
	case ParentSearch:
		req.URL.Fragment = "search--parent"
	case NameSearch:
		req.URL.Fragment = "search--name"
	}

	http.Redirect(w, req, fmt.Sprint(req.URL), http.StatusMovedPermanently)
}

// updateCoverageForm represents form-data for the UpdateCoverage handler.
type updateCoverageForm struct {
	Action      FormAction
	Value       string
	Dimension   string
	LargerArea  string
	SetParent   string
	Coverage    string
	GeographyID string
	OptionType  string
}

// parseUpdateCoverageForm parses form data from a http.Request into a updateCoverageForm.
func parseUpdateCoverageForm(req *http.Request) (updateCoverageForm, error) {
	var action FormAction
	var value, largerArea string

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

	geogID := req.FormValue("geog-id")
	if geogID == "" {
		return updateCoverageForm{}, &clientErr{errors.New("missing required value 'geog-id'")}
	}

	switch coverage {
	case CoverageDefault:
		action = CoverageAll
		value = coverage
	case NameSearch:
		action = Continue
		value = coverage
	case ParentSearch:
		action = Continue
		value = coverage
	default:
		return updateCoverageForm{}, &clientErr{errors.New("unknown coverage type")}
	}

	parent := req.FormValue("larger-area")
	setParent := req.FormValue("set-parent")

	isSearch, _ := strconv.ParseBool(req.FormValue("is-search"))
	if isSearch {
		action = Search
		q := req.FormValue("q")
		value = q
	}
	if isSearch && coverage == ParentSearch {
		largerArea = parent
		pq := req.FormValue("pq")
		value = pq
		action = Search
	}
	if isSearch && coverage == ParentSearch && parent == "" {
		return updateCoverageForm{}, &validationErr{errors.New("missing required value 'larger-area'")}
	}

	addOption := req.FormValue("add-option")
	if addOption != "" {
		action = Add
		value = addOption
	}

	addParentOption := req.FormValue("add-parent-option")
	if addParentOption != "" {
		action = Add
		value = addParentOption
		largerArea = parent
	}

	deleteOption := req.FormValue("delete-option")
	if deleteOption != "" {
		action = Delete
		value = deleteOption
	}

	optType := req.FormValue("option-type")

	return updateCoverageForm{
		Action:      action,
		Value:       value,
		Dimension:   dimension,
		LargerArea:  largerArea,
		SetParent:   setParent,
		Coverage:    coverage,
		GeographyID: geogID,
		OptionType:  optType,
	}, nil
}
