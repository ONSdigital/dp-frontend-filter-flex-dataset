package mapper

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-cookies/cookies"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/helpers"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/model"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/pagination"
	"github.com/ONSdigital/dp-renderer/helper"
	coreModel "github.com/ONSdigital/dp-renderer/model"
	"github.com/ONSdigital/log.go/v2/log"
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
func CreateFilterFlexOverview(req *http.Request, basePage coreModel.Page, lang, path string, queryStrValues []string, filterJob filter.GetFilterResponse, filterDims []model.FilterDimension, datasetDims dataset.VersionDimensions, hasNoAreaOptions, isMultivariate bool, eb zebedee.EmergencyBanner, serviceMsg string) model.Overview {
	p := model.Overview{
		Page: basePage,
	}
	mapCommonProps(req, &p.Page, reviewPageType, "Review changes", lang, serviceMsg, eb)
	p.FilterID = filterJob.FilterID
	dataset := filterJob.Dataset
	p.IsMultivariate = isMultivariate

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
		pageDim.Name = cleanDimensionLabel(dim.Label)
		pageDim.IsAreaType = *dim.IsAreaType
		pageDim.OptionsCount = dim.OptionsCount
		pageDim.ID = dim.ID
		pageDim.URI = fmt.Sprintf("%s/%s", path, dim.Name)
		pageDim.IsChangeCategories = isMultivariate
		q := url.Values{}
		midFloor, midCeiling := getTruncationMidRange(dim.OptionsCount)

		var displayedOptions []string
		if len(dim.Options) > 9 && !helpers.HasStringInSlice(dim.Name, queryStrValues) && !*dim.IsAreaType {
			displayedOptions = append(displayedOptions, dim.Options[:3]...)
			displayedOptions = append(displayedOptions, dim.Options[midFloor:midCeiling]...)
			displayedOptions = append(displayedOptions, dim.Options[len(dim.Options)-3:]...)
			q.Add(queryStrKey, dim.Name)
			helpers.PersistExistingParams(queryStrValues, queryStrKey, dim.Name, q)
			pageDim.IsTruncated = true
		} else {
			helpers.PersistExistingParams(queryStrValues, queryStrKey, dim.Name, q)
			displayedOptions = dim.Options
			pageDim.IsTruncated = false
		}

		pageDim.Options = append(pageDim.Options, displayedOptions...)
		pageDim.TruncateLink = generateTruncatePath(path, dim.ID, q)
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

	// TODO: Temporarily removing mapping as new endpoints are required to return reliable dataset information

	return p
}

// CreateCategorisationsSelector maps data to the Selector model
func CreateCategorisationsSelector(req *http.Request, basePage coreModel.Page, dimLabel, lang, filterID, dimId, serviceMsg string, eb zebedee.EmergencyBanner, cats population.GetCategorisationsResponse, isValidationError bool) model.Selector {
	p := model.Selector{
		Page: basePage,
	}
	mapCommonProps(req, &p.Page, "filter-flex-selector", cleanDimensionLabel(dimLabel), lang, serviceMsg, eb)
	p.Breadcrumb = []coreModel.TaxonomyNode{
		{
			Title: helper.Localise("Back", lang, 1),
			URI:   fmt.Sprintf("/filters/%s/dimensions", filterID),
		},
	}
	p.LeadText = helper.Localise("SelectCategoriesLeadText", lang, 1)
	p.InitialSelection = dimId

	var selections []model.Selection
	for _, cat := range cats.Items {
		cats := []string{}
		for _, c := range sortCategoriesByID(cat.Categories) {
			cats = append(cats, c.Label)
		}
		selections = append(selections, mapCats(cats, req.URL.Query()["showAll"], lang, req.URL.Path, cat.ID))
	}
	p.Selections = selections

	if isValidationError {
		p.Page.Error = coreModel.Error{
			Title: p.Page.Metadata.Title,
			ErrorItems: []coreModel.ErrorItem{
				{
					Description: coreModel.Localisation{
						LocaleKey: "SelectCategoriesError",
						Plural:    1,
					},
					URL: "#categories-error",
				},
			},
			Language: lang,
		}
		p.ErrorId = "categories-error"
	}

	return p
}

// CreateAreaTypeSelector maps data to the Selector model
func CreateAreaTypeSelector(req *http.Request, basePage coreModel.Page, lang, filterID string, areaType []population.AreaType, fDim filter.Dimension, lowest_geography, releaseDate string, dataset dataset.DatasetDetails, isValidationError, hasOpts bool, serviceMsg string, eb zebedee.EmergencyBanner) model.Selector {
	p := model.Selector{
		Page: basePage,
	}
	mapCommonProps(req, &p.Page, areaPageType, areaTypeTitle, lang, serviceMsg, eb)
	p.Breadcrumb = []coreModel.TaxonomyNode{
		{
			Title: helper.Localise("Back", lang, 1),
			URI:   fmt.Sprintf("/filters/%s/dimensions", filterID),
		},
	}
	p.HasOptions = hasOpts
	p.LeadText = helper.Localise("SelectAreaTypeLeadText", lang, 1)

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
		p.ErrorId = "area-type-error"
	}

	selections := mapAreaTypesToSelection(sortAreaTypes(areaType))

	if lowest_geography != "" {
		var filtered_selections []model.Selection
		var lowest_found = false
		for _, selection := range selections {
			if !lowest_found {
				filtered_selections = append(filtered_selections, selection)
				lowest_found = selection.Value == lowest_geography
			}
		}
		p.Selections = filtered_selections
	} else {
		p.Selections = selections
	}

	p.InitialSelection = fDim.ID
	p.IsAreaType = true

	p.DatasetId = dataset.ID
	p.DatasetTitle = dataset.Title
	p.ReleaseDate = releaseDate

	return p
}

// CreateGetCoverage maps data to the coverage model
func CreateGetCoverage(req *http.Request, basePage coreModel.Page, lang, filterID, geogName, nameQ, parentQ, parentArea, setParent, coverage, dim, geogID, releaseDate, serviceMsg string, eb zebedee.EmergencyBanner, dataset dataset.DatasetDetails, areas population.GetAreasResponse, opts []model.SelectableElement, parents population.GetAreaTypeParentsResponse, hasFilterByParent, hasValidationErr bool, currentPage int) model.Coverage {
	p := model.Coverage{
		Page: basePage,
	}
	mapCommonProps(req, &p.Page, coveragePageType, coverageTitle, lang, serviceMsg, eb)
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
	p.SetParent = setParent
	p.NameSearch = model.SearchField{
		Name:     nameSearchFieldName,
		ID:       nameSearch,
		Value:    nameQ,
		Language: lang,
		Label:    helper.Localise("CoverageSearchLabel", lang, 1),
	}
	p.ParentSearch = model.SearchField{
		Name:     parentSearchFieldName,
		ID:       parentSearch,
		Value:    parentQ,
		Language: lang,
		Label:    helper.Localise("CoverageSearchLabel", lang, 1),
	}

	p.DatasetId = dataset.ID
	p.DatasetTitle = dataset.Title
	p.ReleaseDate = releaseDate

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

	totalPages := pagination.GetTotalPages(areas.TotalCount, areas.Limit)
	var paginatedResults coreModel.Pagination
	if totalPages > 1 {
		paginatedResults = coreModel.Pagination{
			CurrentPage:       currentPage,
			TotalPages:        totalPages,
			Limit:             areas.Limit,
			FirstAndLastPages: pagination.GetFirstAndLastPages(req, totalPages),
			PagesToDisplay:    pagination.GetPagesToDisplay(currentPage, totalPages, req),
		}
	}

	if len(opts) > 0 && hasFilterByParent {
		p.CoverageType = parentSearch
		p.ParentSearchOutput.Selections = opts
		p.ParentSearchOutput.SelectionsTitle = helper.Localise("AreasAddedTitle", lang, len(opts))
		p.OptionType = parentSearch
	} else if len(opts) > 0 {
		p.CoverageType = nameSearch
		p.NameSearchOutput.Selections = opts
		p.ParentSearchOutput.SelectionsTitle = helper.Localise("AreasAddedTitle", lang, len(opts))
		p.OptionType = nameSearch
	}

	switch coverage {
	case nameSearch:
		p.CoverageType = nameSearch
		p.NameSearchOutput.Results = results
		p.NameSearchOutput.HasNoResults = len(p.NameSearchOutput.Results) == 0
		p.NameSearchOutput.Pagination = paginatedResults
	case parentSearch:
		p.CoverageType = parentSearch
		p.ParentSearchOutput.Results = results
		p.ParentSearchOutput.HasNoResults = len(p.ParentSearchOutput.Results) == 0 && !hasValidationErr
		p.ParentSearchOutput.Pagination = paginatedResults
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

// CreateGetChangeDimensions maps data to the ChangeDimensions model
func CreateGetChangeDimensions(req *http.Request, basePage coreModel.Page, lang, fid, q, formAction, serviceMsg string, eb zebedee.EmergencyBanner, dims []model.FilterDimension, pDims, results population.GetDimensionsResponse) model.ChangeDimensions {
	p := model.ChangeDimensions{
		Page: basePage,
	}
	p.Breadcrumb = []coreModel.TaxonomyNode{
		{
			Title: helper.Localise("Back", lang, 1),
			URI:   fmt.Sprintf("/filters/%s/dimensions", fid),
		},
	}
	mapCommonProps(req, &p.Page, "change_variables", "Add or remove variables", lang, serviceMsg, eb)
	p.FormAction = formAction

	selections := []model.SelectableElement{}
	for _, dim := range dims {
		if !*dim.IsAreaType {
			selections = append(selections, model.SelectableElement{
				Text:  cleanDimensionLabel(dim.Label),
				Value: dim.ID,
				Name:  "delete-option",
			})
		}
	}
	p.Output.Selections = selections
	p.Output.SelectionsTitle = "Variables added"
	p.Search = model.SearchField{
		Name:     "q",
		ID:       "dimensions-search",
		Language: lang,
		Value:    q,
		Label:    helper.Localise("DimensionsSearchLabel", lang, 1),
	}

	browseResults := mapDimensionsResponse(pDims, &selections)
	searchResults := mapDimensionsResponse(results, &selections)

	p.Output.Results = browseResults
	p.SearchOutput.Results = searchResults
	p.SearchOutput.HasNoResults = len(p.SearchOutput.Results) == 0 && formAction == "search"

	return p
}

// mapDimensionsResponse returns a sorted array of selectable elements
func mapDimensionsResponse(pDims population.GetDimensionsResponse, selections *[]model.SelectableElement) []model.SelectableElement {
	results := []model.SelectableElement{}
	for _, pDim := range pDims.Dimensions {
		var sel model.SelectableElement
		sel.Name = "add-dimension"
		sel.Text = cleanDimensionLabel(pDim.Label)
		sel.InnerText = pDim.Description
		sel.Value = pDim.ID
		pDimId := helpers.TrimCategoryValue(pDim.ID)
		for _, dim := range *selections {
			dimV := helpers.TrimCategoryValue(dim.Value)
			if strings.EqualFold(dimV, pDimId) {
				sel.IsSelected = true
				sel.Name = "delete-option"
				sel.Value = dim.Value
				break
			}
		}
		results = append(results, sel)
	}
	sort.SliceStable(results, func(i, j int) bool {
		return results[i].Text < results[j].Text
	})
	return results
}

// cleanDimensionlabel is a helper function that parses dimension labels from cantabular into display text
func cleanDimensionLabel(label string) string {
	matcher := regexp.MustCompile(`(\(\d+ ((C|c)ategories|(C|c)ategory)\))`)
	result := matcher.ReplaceAllString(label, "")
	return strings.TrimSpace(result)
}

// getAddOptionStr is a helper function to determine which add option string should be returned
func getAddOptionStr(isParentSearch bool) string {
	if isParentSearch {
		return "add-parent-option"
	}
	return "add-option"
}

// mapCommonProps maps common properties on all filter/flex pages
func mapCommonProps(req *http.Request, p *coreModel.Page, pageType, title, lang, serviceMsg string, eb zebedee.EmergencyBanner) {
	mapCookiePreferences(req, &p.CookiesPreferencesSet, &p.CookiesPolicy)
	p.BetaBannerEnabled = true
	p.Type = pageType
	p.Metadata.Title = title
	p.Language = lang
	p.URI = req.URL.Path
	p.SearchNoIndexEnabled = true
	p.ServiceMessage = serviceMsg
	p.EmergencyBanner = mapEmergencyBanner(eb)
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

// mapEmergencyBanner maps relevant emergency banner data
func mapEmergencyBanner(bannerData zebedee.EmergencyBanner) coreModel.EmergencyBanner {
	var mappedEmergencyBanner coreModel.EmergencyBanner
	emptyBannerObj := zebedee.EmergencyBanner{}
	if bannerData != emptyBannerObj {
		mappedEmergencyBanner.Title = bannerData.Title
		mappedEmergencyBanner.Type = strings.Replace(bannerData.Type, "_", "-", -1)
		mappedEmergencyBanner.Description = bannerData.Description
		mappedEmergencyBanner.URI = bannerData.URI
		mappedEmergencyBanner.LinkText = bannerData.LinkText
	}
	return mappedEmergencyBanner
}

// getTruncationMidRange returns ints that can be used as the truncation mid range
func getTruncationMidRange(total int) (int, int) {
	mid := total / 2
	midFloor := mid - 2
	midCeiling := midFloor + 3
	if midFloor < 0 {
		midFloor = 0
	}
	return midFloor, midCeiling
}

// generateTruncatePath returns the path to truncate or show all
func generateTruncatePath(path, dimID string, q url.Values) string {
	truncatePath := path
	if q.Encode() != "" {
		truncatePath += fmt.Sprintf("?%s", q.Encode())
	}
	if dimID != "" {
		truncatePath += fmt.Sprintf("#%s", dimID)
	}
	return truncatePath
}

// mapCats is a helper function that returns either truncated or untruncated mapped categories
func mapCats(cats, queryStrValues []string, lang, path, catID string) model.Selection {
	q := url.Values{}
	catsLength := len(cats)
	midFloor, midCeiling := getTruncationMidRange(catsLength)

	var displayedCats []string
	var isTruncated bool
	if catsLength > 9 && !helpers.HasStringInSlice(catID, queryStrValues) {
		displayedCats = append(displayedCats, cats[:3]...)
		displayedCats = append(displayedCats, cats[midFloor:midCeiling]...)
		displayedCats = append(displayedCats, cats[catsLength-3:]...)
		q.Add(queryStrKey, catID)
		helpers.PersistExistingParams(queryStrValues, queryStrKey, catID, q)
		isTruncated = true
	} else {
		helpers.PersistExistingParams(queryStrValues, queryStrKey, catID, q)
		displayedCats = cats
		isTruncated = false
	}
	return model.Selection{
		Value:           catID,
		Label:           fmt.Sprintf("%d %s", catsLength, helper.Localise("Category", lang, catsLength)),
		Categories:      displayedCats,
		CategoriesCount: catsLength,
		IsTruncated:     isTruncated,
		TruncateLink:    generateTruncatePath((path), catID, q),
	}
}

func mapAreaTypesToSelection(areaTypes []population.AreaType) []model.Selection {
	var selections []model.Selection
	for _, area := range areaTypes {
		selections = append(selections, model.Selection{
			Value:       area.ID,
			Label:       area.Label,
			Description: area.Description,
			TotalCount:  area.TotalCount,
		})
	}
	return selections
}
