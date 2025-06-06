package mapper

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-cookies/cookies"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/helpers"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/model"
	"github.com/ONSdigital/dp-renderer/v2/helper"
	coreModel "github.com/ONSdigital/dp-renderer/v2/model"
)

// Mapper represents the core mappings required for all pages
type Mapper struct {
	req        *http.Request
	basePage   coreModel.Page
	eb         zebedee.EmergencyBanner
	lang       string
	serviceMsg string
	fid        string
}

// NewMapper creates a new instance of Mapper
func NewMapper(request *http.Request, basePage coreModel.Page, emergencyBanner zebedee.EmergencyBanner, language, serviceMsg, filterId string) *Mapper {
	return &Mapper{
		req:        request,
		basePage:   basePage,
		eb:         emergencyBanner,
		lang:       language,
		serviceMsg: serviceMsg,
		fid:        filterId,
	}
}

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
	maxVariableErrorStr   = "Maximum variables"
	maxCellsErrorStr      = "withinMaxCells"
)

// mapDimensionsResponse returns a sorted array of selectable elements
func mapDimensionsResponse(pDims population.GetDimensionsResponse, selections *[]model.SelectableElement, lang string) []model.SelectableElement {
	results := []model.SelectableElement{}
	for _, pDim := range pDims.Dimensions {
		var sel model.SelectableElement
		sel.Name = "add-dimension"
		sel.Text = cleanDimensionLabel(pDim.Label)
		sel.InnerText = pDim.Description
		sel.Value = pDim.ID
		sel.QualityStatement = model.Panel{
			CssClasses: []string{"ons-u-mt-s", "ons-u-mb-xs"},
			Body:       pDim.QualityStatementText,
			Language:   lang,
		}
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

// cleanDimensionLabel is a helper function that parses dimension labels from cantabular into display text
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
	preferencesCookie := cookies.GetONSCookiePreferences(req)
	*preferencesIsSet = preferencesCookie.IsPreferenceSet
	*policy = coreModel.CookiesPolicy{
		Communications: preferencesCookie.Policy.Campaigns,
		Essential:      preferencesCookie.Policy.Essential,
		Settings:       preferencesCookie.Policy.Settings,
		Usage:          preferencesCookie.Policy.Usage,
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
func mapCats(cats, queryStrValues []string, lang, path, catID string, isSuggested bool) model.Selection {
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
		IsSuggested:     isSuggested,
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

// mapPanel is a helper function that returns a mapped panel
func mapPanel(locale coreModel.Localisation, language string, utilityCssClasses []string) model.Panel {
	return model.Panel{
		Body:       helper.Localise(locale.LocaleKey, language, locale.Plural),
		Language:   language,
		CssClasses: utilityCssClasses,
	}
}

// mapBlockedAreasPanel is a helper function that returns the blocked areas panel
func (m *Mapper) mapBlockedAreasPanel(sdc *cantabular.GetBlockedAreaCountResult, isMaxCellsError bool, panelType model.PanelType) (p *model.Panel) {
	if isMaxCellsError {
		return &model.Panel{
			Type:       model.Error,
			CssClasses: []string{"ons-u-mb-s"},
			Language:   m.lang,
			SafeHTML: []string{
				helper.Localise("MaxCellsErrorPanelDescription", m.lang, 1),
			},
		}
	}

	switch panelType {
	case model.Pending:
		p = &model.Panel{
			Type:       model.Pending,
			CssClasses: []string{"ons-u-mb-s"},
			Language:   m.lang,
			SafeHTML: []string{
				helper.Localise("SDCAreasAvailable", m.lang, 1, helper.ThousandsSeparator(sdc.Passed), helper.ThousandsSeparator(sdc.Total)),
				helper.Localise("SDCRestrictedAreas", m.lang, 4, helper.ThousandsSeparator(sdc.Blocked)),
			},
		}
	case model.Success:
		p = &model.Panel{
			Type:       model.Success,
			CssClasses: []string{"ons-u-mb-l"},
			Language:   m.lang,
			SafeHTML: []string{
				helper.Localise("SDCAllAreasAvailable", m.lang, sdc.Total, helper.ThousandsSeparator(sdc.Total)),
			},
		}
	}
	return p
}

// isMaxVariablesError returns true if the sdc result is returning a maximum variables exceeded TableError
func isMaxVariablesError(sdc *cantabular.GetBlockedAreaCountResult) bool {
	return strings.Contains(sdc.TableError, maxVariableErrorStr)
}

// isMaxCellsError returns true if the sdc result is returning a maximum cells exceeded TableError
func isMaxCellsError(sdc *cantabular.GetBlockedAreaCountResult) bool {
	return strings.Contains(sdc.TableError, maxCellsErrorStr)
}
