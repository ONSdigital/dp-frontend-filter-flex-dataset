package mapper

import (
	"net/http"
	"strings"

	"github.com/ONSdigital/dp-cookies/cookies"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/model"
	coreModel "github.com/ONSdigital/dp-renderer/model"
)

// CreateFilterFlexOverview maps data to the Overview model
func CreateFilterFlexOverview(req *http.Request, basePage coreModel.Page, lang string) model.Overview {
	p := model.Overview{
		Page: basePage,
	}
	mapCookiePreferences(req, &p.Page.CookiesPreferencesSet, &p.Page.CookiesPolicy)

	p.BetaBannerEnabled = true
	p.Type = "filter-flex-overview"
	p.Metadata.Title = "Review changes"
	p.Language = lang

	p.Breadcrumb = []coreModel.TaxonomyNode{
		{
			Title: "Back",
			URI:   "#",
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

// mapCookiePreferences reads cookie policy and preferences cookies and then maps the values to the page model
func mapCookiePreferences(req *http.Request, preferencesIsSet *bool, policy *coreModel.CookiesPolicy) {
	preferencesCookie := cookies.GetCookiePreferences(req)
	*preferencesIsSet = preferencesCookie.IsPreferenceSet
	*policy = coreModel.CookiesPolicy{
		Essential: preferencesCookie.Policy.Essential,
		Usage:     preferencesCookie.Policy.Usage,
	}
}
