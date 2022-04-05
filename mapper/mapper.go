package mapper

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ONSdigital/dp-api-clients-go/v2/dimension"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-cookies/cookies"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/helpers"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/model"
	coreModel "github.com/ONSdigital/dp-renderer/model"
)

// Constants...
const queryStrKey = "showAll"

// CreateFilterFlexOverview maps data to the Overview model
func CreateFilterFlexOverview(req *http.Request, basePage coreModel.Page, lang, path string, queryStrValues []string, filterJob filter.Model) model.Overview {
	p := model.Overview{
		Page: basePage,
	}
	mapCookiePreferences(req, &p.Page.CookiesPreferencesSet, &p.Page.CookiesPolicy)

	p.BetaBannerEnabled = true
	p.Type = "filter-flex-overview"
	p.Metadata.Title = "Review changes"
	p.Language = lang
	p.FilterID = filterJob.FilterID

	p.Breadcrumb = []coreModel.TaxonomyNode{
		{
			Title: "Back",
			URI:   "#",
		},
	}

	for _, dim := range filterJob.Dimensions {
		pageDim := model.Dimension{}
		encodedName := url.QueryEscape(dim.Name)
		pageDim.Name = dim.Name
		pageDim.IsAreaType = *dim.IsAreaType
		pageDim.OptionsCount = len(dim.Options)
		pageDim.EncodedName = strings.ToLower(encodedName)
		pageDim.URI = fmt.Sprintf("%s/%s", path, strings.ToLower(encodedName))
		q := url.Values{}

		if len(dim.Options) > 9 && !helpers.HasStringInSlice(dim.Name, queryStrValues) {
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
		if encodedName != "" {
			truncatePath += fmt.Sprintf("#%s", strings.ToLower(encodedName))
		}
		pageDim.TruncateLink = truncatePath

		p.Dimensions = append(p.Dimensions, pageDim)
	}

	// TODO: Get this from dataset client
	p.Collapsible = coreModel.Collapsible{
		LocaliseKey:       "VariableExplanation",
		LocalisePluralInt: 4,
		CollapsibleItems: []coreModel.CollapsibleItem{
			{
				Subheading: "This is a subheading",
				Content:    []string{"a string"},
			},
			{
				Subheading: "This is another subheading",
				Content:    []string{"another string", "and another"},
			},
		},
	}

	return p
}

// CreateSelector maps data to the Selector model
func CreateSelector(req *http.Request, basePage coreModel.Page, dimName, lang string) model.Selector {
	p := model.Selector{
		Page: basePage,
	}
	mapCookiePreferences(req, &p.Page.CookiesPreferencesSet, &p.Page.CookiesPolicy)

	p.BetaBannerEnabled = true
	p.Type = "filter-flex-selector"
	p.Metadata.Title = strings.Title(dimName)
	p.Language = lang

	p.Breadcrumb = []coreModel.TaxonomyNode{
		{
			Title: "Back",
			URI:   "../dimensions",
		},
	}

	return p
}

// CreateAreaTypeSelector maps data to the Overview model
func CreateAreaTypeSelector(req *http.Request, basePage coreModel.Page, lang string, areaType []dimension.AreaType, selectionName string) model.Selector {
	p := CreateSelector(req, basePage, selectionName, lang)
	p.Page.Metadata.Title = "Area Type"

	var selections []model.Selection
	for _, area := range areaType {
		selections = append(selections, model.Selection{
			// Currently, labels are used instead of ID's, since dimensions are stored/queried using their
			// display name. Once that changes we can use the area-type ID, knowing it will match the imported dimension.
			Value:      area.Label,
			Label:      area.Label,
			TotalCount: area.TotalCount,
		})
	}

	p.Selections = selections
	p.InitialSelection = selectionName
	p.IsAreaType = true

	return p
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
