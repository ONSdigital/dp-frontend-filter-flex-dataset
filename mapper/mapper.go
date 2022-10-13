package mapper

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-cookies/cookies"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/helpers"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/model"
	"github.com/ONSdigital/dp-renderer/helper"
	coreModel "github.com/ONSdigital/dp-renderer/model"
	"github.com/ONSdigital/log.go/v2/log"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Constants...
const (
	queryStrKey           = "showAll"
	areaTypePrefix        = "AreaType"
	areaTypeTitle         = "Area type"
	pluralInt             = 4
	nameSearch            = "name-search"
	parentSearch          = "parent-search"
	nameSearchFieldName   = "q"
	parentSearchFieldName = "pq"
	coveragePageType      = "coverage_options"
	coverageTitle         = "Coverage"
	areaPageType          = "area_type_options"
	reviewPageType        = "review_changes"
)

// CreateFilterFlexOverview maps data to the Overview model
func CreateFilterFlexOverview(req *http.Request, basePage coreModel.Page, lang, path string, queryStrValues []string, filterJob filter.GetFilterResponse, filterDims []model.FilterDimension, datasetDims dataset.VersionDimensions, hasNoAreaOptions bool) model.Overview {
	p := model.Overview{
		Page: basePage,
	}
	mapCommonProps(req, &p.Page, reviewPageType, "Review changes", lang)
	p.FilterID = filterJob.FilterID
	dataset := filterJob.Dataset

	p.Breadcrumb = []coreModel.TaxonomyNode{
		{
			Title: helper.Localise("Back", lang, 1),
			URI: fmt.Sprintf("/datasets/%s/editions/%s/versions/%s",
				dataset.DatasetID,
				dataset.Edition,
				strconv.Itoa(dataset.Version)),
		},
	}

	for _, dim := range filterDims {
		pageDim := model.Dimension{}
		pageDim.Name = dim.Label
		pageDim.IsAreaType = *dim.IsAreaType
		pageDim.OptionsCount = dim.OptionsCount
		pageDim.ID = dim.ID
		pageDim.URI = fmt.Sprintf("%s/%s", path, dim.Name)
		q := url.Values{}

		if len(dim.Options) > 9 && !helpers.HasStringInSlice(dim.Name, queryStrValues) && !*dim.IsAreaType {
			firstSlice := dim.Options[:3]

			mid := len(dim.Options) / 2
			midFloor := mid - 2
			midCeiling := midFloor + 3
			midSlice := dim.Options[midFloor:midCeiling]

			lastSlice := dim.Options[len(dim.Options)-3:]
			pageDim.Options = append(pageDim.Options, firstSlice...)
			pageDim.Options = append(pageDim.Options, midSlice...)
			pageDim.Options = append(pageDim.Options, lastSlice...)

			q.Add(queryStrKey, dim.Name)
			helpers.PersistExistingParams(queryStrValues, queryStrKey, dim.Name, q)
			pageDim.IsTruncated = true
		} else {
			helpers.PersistExistingParams(queryStrValues, queryStrKey, dim.Name, q)
			pageDim.Options = dim.Options
			pageDim.IsTruncated = false
		}

		truncatePath := path
		if q.Encode() != "" {
			truncatePath += fmt.Sprintf("?%s", q.Encode())
		}
		if dim.ID != "" {
			truncatePath += fmt.Sprintf("#%s", dim.ID)
		}
		pageDim.TruncateLink = truncatePath

		p.Dimensions = append(p.Dimensions, pageDim)
	}

	sort.Slice(p.Dimensions, func(i, j int) bool {
		return p.Dimensions[i].IsAreaType
	})

	coverage := []model.Dimension{
		{
			IsCoverage:        true,
			IsDefaultCoverage: hasNoAreaOptions,
			URI:               fmt.Sprintf("%s/geography/coverage", path),
			Options:           p.Dimensions[0].Options,
			ID:                "coverage",
		},
	}
	temp := append(coverage, p.Dimensions[1:]...)
	p.Dimensions = append(p.Dimensions[:1], temp...)

	collapsibleContentItems := mapCollapsible(datasetDims.Items)
	p.Collapsible = coreModel.Collapsible{
		Title: coreModel.Localisation{
			LocaleKey: "VariableExplanation",
			Plural:    4,
		},
		CollapsibleItems: collapsibleContentItems,
	}

	return p
}

func mapCollapsible(datasetDims []dataset.VersionDimension) []coreModel.CollapsibleItem {
	var collapsibleContentItems []coreModel.CollapsibleItem
	collapsibleContentItems = append(collapsibleContentItems, coreModel.CollapsibleItem{
		Subheading: areaTypeTitle,
		SafeHTML: coreModel.Localisation{
			LocaleKey: "VariableInfoAreaType",
			Plural:    1,
		},
	})
	collapsibleContentItems = append(collapsibleContentItems, coreModel.CollapsibleItem{
		Subheading: coverageTitle,
		SafeHTML: coreModel.Localisation{
			LocaleKey: "VariableInfoCoverage",
			Plural:    1,
		},
	})
	for _, dims := range datasetDims {
		if dims.Description != "" {
			var collapsibleContent coreModel.CollapsibleItem
			collapsibleContent.Subheading = dims.Label
			collapsibleContent.Content = strings.Split(dims.Description, "\n")
			collapsibleContentItems = append(collapsibleContentItems, collapsibleContent)
		}
	}
	return collapsibleContentItems
}

// CreateSelector maps data to the Selector model
func CreateSelector(req *http.Request, basePage coreModel.Page, dimName, lang, filterID string) model.Selector {
	p := model.Selector{
		Page: basePage,
	}
	mapCommonProps(req, &p.Page, "filter-flex-selector", cases.Title(language.English).String(dimName), lang)
	p.Breadcrumb = []coreModel.TaxonomyNode{
		{
			Title: helper.Localise("Back", lang, 1),
			URI:   fmt.Sprintf("/filters/%s/dimensions", filterID),
		},
	}

	return p
}

// CreateAreaTypeSelector maps data to the Selector model
func CreateAreaTypeSelector(req *http.Request, basePage coreModel.Page, lang, filterID string, areaType []population.AreaType, fDim filter.Dimension, isValidationError bool) model.Selector {
	p := CreateSelector(req, basePage, fDim.Label, lang, filterID)
	p.Page.Metadata.Title = areaTypeTitle
	p.Page.Type = areaPageType

	if isValidationError {
		p.Page.Error = coreModel.Error{
			Title: p.Page.Metadata.Title,
			ErrorItems: []coreModel.ErrorItem{
				{
					Description: coreModel.Localisation{
						LocaleKey: "SelectAreaTypeError",
						Plural:    1,
					},
					URL: "#area-type-error",
				},
			},
			Language: lang,
		}
	}

	var selections []model.Selection
	for _, area := range areaType {
		selections = append(selections, model.Selection{
			Value:      area.ID,
			Label:      area.Label,
			TotalCount: area.TotalCount,
		})
	}

	p.Selections = selections
	p.InitialSelection = fDim.ID
	p.IsAreaType = true

	return p
}

// CreateGetCoverage maps data to the coverage model
func CreateGetCoverage(req *http.Request, basePage coreModel.Page, lang, filterID, geogName, nameQ, parentQ, parentArea, coverage, dim, geogID string, areas population.GetAreasResponse, opts []model.SelectableElement, parents population.GetAreaTypeParentsResponse, hasFilterByParent, hasValidationErr bool) model.Coverage {
	p := model.Coverage{
		Page: basePage,
	}
	mapCommonProps(req, &p.Page, coveragePageType, coverageTitle, lang)
	p.Breadcrumb = []coreModel.TaxonomyNode{
		{
			Title: helper.Localise("Back", lang, 1),
			URI:   fmt.Sprintf("/filters/%s/dimensions", filterID),
		},
	}
	geography := helpers.Pluralise(req, geogName, lang, areaTypePrefix, pluralInt)
	if geography == "" {
		log.Info(req.Context(), "pluralisation lookup failed, reverting to initial input", log.Data{
			"initial_input": geogName,
		})
		geography = geogName
	}

	p.Geography = strings.ToLower(geography)
	p.CoverageType = coverage
	p.Dimension = dim
	p.GeographyID = geogID
	p.NameSearch = model.SearchField{
		Name:  nameSearchFieldName,
		ID:    nameSearch,
		Value: nameQ,
	}
	p.ParentSearch = model.SearchField{
		Name:  parentSearchFieldName,
		ID:    parentSearch,
		Value: parentQ,
	}
	if len(parents.AreaTypes) > 1 && parentArea == "" {
		p.ParentSelect = []model.SelectableElement{
			{
				Text:       helper.Localise("CoverageSelectDefault", lang, 1),
				IsSelected: true,
				IsDisabled: true,
			},
		}
	}
	for _, parent := range parents.AreaTypes {
		var sel model.SelectableElement
		sel.Text = parent.Label
		sel.Value = parent.ID
		if parentArea == parent.ID {
			sel.IsSelected = true
		}
		p.ParentSelect = append(p.ParentSelect, sel)
	}

	var isParentSearch bool
	if coverage == parentSearch {
		isParentSearch = true
	}
	var results []model.SelectableElement
	for _, area := range areas.Areas {
		var result model.SelectableElement
		result.Text = area.Label
		result.Value = area.ID
		result.Name = getAddOptionStr(isParentSearch)
		for _, opt := range opts {
			if opt.Value == area.ID {
				result.IsSelected = true
				result.Name = "delete-option"
				break
			}
		}
		results = append(results, result)
	}

	if len(opts) > 0 && hasFilterByParent {
		p.CoverageType = parentSearch
		p.ParentSearchOutput.Options = opts
	} else if len(opts) > 0 {
		p.CoverageType = nameSearch
		p.NameSearchOutput.Options = opts
	}

	switch coverage {
	case nameSearch:
		p.NameSearchOutput.SearchResults = results
		p.NameSearchOutput.HasNoResults = len(p.NameSearchOutput.SearchResults) == 0
	case parentSearch:
		p.ParentSearchOutput.SearchResults = results
		p.ParentSearchOutput.HasNoResults = len(p.ParentSearchOutput.SearchResults) == 0 && !hasValidationErr
	}

	if hasValidationErr {
		p.Page.Error = coreModel.Error{
			Title: p.Metadata.Title,
			ErrorItems: []coreModel.ErrorItem{
				{
					Description: coreModel.Localisation{
						LocaleKey: "CoverageSelectDefault",
						Plural:    1,
					},
					URL: "#coverage-error",
				},
			},
			Language: lang,
		}
	}

	p.IsSelectParents = len(parents.AreaTypes) > 0

	return p
}

// getAddOptionStr is a helper function to determine which add option string should be returned
func getAddOptionStr(isParentSearch bool) string {
	if isParentSearch {
		return "add-parent-option"
	}
	return "add-option"
}

// mapCommonProps maps common properties on all filter/flex pages
func mapCommonProps(req *http.Request, p *coreModel.Page, pageType, title, lang string) {
	mapCookiePreferences(req, &p.CookiesPreferencesSet, &p.CookiesPolicy)
	p.BetaBannerEnabled = true
	p.Type = pageType
	p.Metadata.Title = title
	p.Language = lang
	p.URI = req.URL.Path
}

// mapCookiePreferences reads cookie policy and preferences cookies and then maps the values to the page model
func mapCookiePreferences(req *http.Request, preferencesIsSet *bool, policy *coreModel.CookiesPolicy) {
	preferencesCookie := cookies.GetCookiePreferences(req)
	*preferencesIsSet = preferencesCookie.IsPreferenceSet
	*policy = coreModel.CookiesPolicy{
		Essential: preferencesCookie.Policy.Essential,
		Usage:     preferencesCookie.Policy.Usage,
	}
}
