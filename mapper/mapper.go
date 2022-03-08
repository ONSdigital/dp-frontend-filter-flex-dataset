package mapper

import (
	"net/http"
	"strings"

	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-cookies/cookies"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/model"
	coreModel "github.com/ONSdigital/dp-renderer/model"
)

// CreateFilterFlexOverview maps data to the Overview model
func CreateFilterFlexOverview(req *http.Request, basePage coreModel.Page, lang string, dims filter.Dimensions, showAll []string) model.Overview {
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

	for _, dim := range dims.Dimensions {
		var dimension model.Dimension
		dimension.Name = dim.Name
		dimension.IsAreaType = dim.IsAreaType
		dimension.Options = append(dimension.Options, dim.Options...)
		p.Dimensions = append(p.Dimensions, dimension)
	}

	// p.Dimensions = []model.Dimension{
	// 	{
	// 		Options: []model.Option{
	// 			{
	// 				Code:  "EW",
	// 				Label: "Electoral Wards and Divisions",
	// 			},
	// 		},
	// 		Name:         "Area type",
	// 		EncodedName:  url.QueryEscape("Area type"),
	// 		OptionsCount: 7878,
	// 		IsAreaType:   true,
	// 		URI:          fmt.Sprintf("/flex/1234/dimensions/%s", "area_type"),
	// 	},
	// 	{
	// 		Options: []model.Option{
	// 			{
	// 				Code:  "E",
	// 				Label: "England",
	// 			},
	// 			{
	// 				Code:  "W",
	// 				Label: "Wales",
	// 			},
	// 		},
	// 		Name:         "Coverage",
	// 		EncodedName:  url.QueryEscape("Coverage"),
	// 		OptionsCount: 2,
	// 		IsAreaType:   true,
	// 		URI:          fmt.Sprintf("/flex/1234/dimensions/%s", strings.ToLower("Coverage")),
	// 	},
	// 	{
	// 		Options: []model.Option{
	// 			{
	// 				Code:  "B",
	// 				Label: "Bob",
	// 			},
	// 			{
	// 				Code:  "W",
	// 				Label: "Bill",
	// 			},
	// 			{
	// 				Code:  "W",
	// 				Label: "fred",
	// 			},
	// 			{
	// 				Code:  "B",
	// 				Label: "Bob",
	// 			},
	// 			{
	// 				Code:  "W",
	// 				Label: "Bill",
	// 			},
	// 			{
	// 				Code:  "W",
	// 				Label: "fred",
	// 			},
	// 			{
	// 				Code:  "B",
	// 				Label: "Bob",
	// 			},
	// 			{
	// 				Code:  "W",
	// 				Label: "Bill",
	// 			},
	// 			{
	// 				Code:  "W",
	// 				Label: "fred",
	// 			},
	// 			// {
	// 			// 	Code:  "B",
	// 			// 	Label: "Bob",
	// 			// },
	// 			// {
	// 			// 	Code:  "W",
	// 			// 	Label: "Bill",
	// 			// },
	// 			// {
	// 			// 	Code:  "W",
	// 			// 	Label: "fred",
	// 			// },
	// 		},
	// 		Name:         "Names",
	// 		EncodedName:  url.QueryEscape("Names"),
	// 		OptionsCount: 12,
	// 		IsAreaType:   false,
	// 		URI:          fmt.Sprintf("/flex/1234/dimensions/%s", strings.ToLower("Coverage")),
	// 	},
	// 	{
	// 		Options: []model.Option{
	// 			{
	// 				Code:  "1",
	// 				Label: "first",
	// 			},
	// 			{
	// 				Code:  "2",
	// 				Label: "second",
	// 			},
	// 			{
	// 				Code:  "3",
	// 				Label: "third",
	// 			},
	// 			{
	// 				Code:  "1",
	// 				Label: "fourth",
	// 			},
	// 			{
	// 				Code:  "2",
	// 				Label: "fifth",
	// 			},
	// 			{
	// 				Code:  "3",
	// 				Label: "sixth",
	// 			},
	// 			{
	// 				Code:  "3",
	// 				Label: "seventh",
	// 			},
	// 			{
	// 				Code:  "3",
	// 				Label: "eight",
	// 			},
	// 			{
	// 				Code:  "3",
	// 				Label: "nine",
	// 			},
	// 			{
	// 				Code:  "3",
	// 				Label: "ten",
	// 			},
	// 			{
	// 				Code:  "3",
	// 				Label: "eleven",
	// 			},
	// 			{
	// 				Code:  "1",
	// 				Label: "twelve",
	// 			},
	// 			{
	// 				Code:  "2",
	// 				Label: "thirteen",
	// 			},
	// 			{
	// 				Code:  "3",
	// 				Label: "fourteen",
	// 			},
	// 			{
	// 				Code:  "3",
	// 				Label: "fifteen",
	// 			},
	// 			{
	// 				Code:  "3",
	// 				Label: "sixteen",
	// 			},
	// 			{
	// 				Code:  "3",
	// 				Label: "seventeen",
	// 			},
	// 			{
	// 				Code:  "3",
	// 				Label: "eighteen",
	// 			},
	// 			{
	// 				Code:  "3",
	// 				Label: "19",
	// 			},
	// 			{
	// 				Code:  "3",
	// 				Label: "20",
	// 			},
	// 			{
	// 				Code:  "3",
	// 				Label: "21",
	// 			},
	// 			{
	// 				Code:  "1",
	// 				Label: "22",
	// 			},
	// 			{
	// 				Code:  "2",
	// 				Label: "23",
	// 			},
	// 			{
	// 				Code:  "3",
	// 				Label: "24",
	// 			},
	// 			{
	// 				Code:  "3",
	// 				Label: "25",
	// 			},
	// 			{
	// 				Code:  "3",
	// 				Label: "26",
	// 			},
	// 			{
	// 				Code:  "3",
	// 				Label: "27",
	// 			},
	// 			{
	// 				Code:  "last-3",
	// 				Label: "last-3",
	// 			},
	// 			{
	// 				Code:  "last-2",
	// 				Label: "last-2",
	// 			},
	// 			{
	// 				Code:  "last-1",
	// 				Label: "last-1",
	// 			},
	// 		},
	// 		Name:         "Ethnic group",
	// 		EncodedName:  url.QueryEscape("Ethnic group"),
	// 		OptionsCount: 18,
	// 		URI:          fmt.Sprintf("/flex/1234/dimensions/%s", strings.ToLower("Ethnic group")),
	// 	},
	// }

	// for i, dim := range p.Dimensions {
	// 	if len(dim.Options) > 9 && !stringInSlice(dim.Name, showAll) {
	// 		firstSlice := dim.Options[:3]

	// 		mid := len(dim.Options) / 2
	// 		midFloor := mid - 2
	// 		midCeiling := midFloor + 3
	// 		midSlice := dim.Options[midFloor:midCeiling]

	// 		lastSlice := dim.Options[len(dim.Options)-3:]
	// 		p.Dimensions[i].TruncatedOptions = append(p.Dimensions[i].TruncatedOptions, firstSlice...)
	// 		p.Dimensions[i].TruncatedOptions = append(p.Dimensions[i].TruncatedOptions, midSlice...)
	// 		p.Dimensions[i].TruncatedOptions = append(p.Dimensions[i].TruncatedOptions, lastSlice...)

	// 		encodedName := url.QueryEscape(p.Dimensions[i].Name)
	// 		q := &url.Values{}
	// 		q.Add("showAll", p.Dimensions[i].Name)
	// 		if len(showAll) > 0 {
	// 			for _, name := range showAll {
	// 				if name != dim.Name {
	// 					q.Add("showAll", name)
	// 				}
	// 			}
	// 		}
	// 		p.Dimensions[i].TruncateLink = fmt.Sprintf("?%s#%s", q.Encode(), encodedName)
	// 	}

	// 	if len(dim.Options) > 9 && stringInSlice(dim.Name, showAll) {
	// 		q := &url.Values{}
	// 		if len(showAll) > 0 {
	// 			for _, name := range showAll {
	// 				if name != dim.Name {
	// 					q.Add("showAll", name)
	// 				}
	// 			}
	// 		}
	// 		encodedName := url.QueryEscape(p.Dimensions[i].Name)
	// 		p.Dimensions[i].TruncateLink = fmt.Sprintf("?%s#%s", q.Encode(), encodedName)
	// 	}
	// }

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

// mapCookiePreferences reads cookie policy and preferences cookies and then maps the values to the page model
func mapCookiePreferences(req *http.Request, preferencesIsSet *bool, policy *coreModel.CookiesPolicy) {
	preferencesCookie := cookies.GetCookiePreferences(req)
	*preferencesIsSet = preferencesCookie.IsPreferenceSet
	*policy = coreModel.CookiesPolicy{
		Essential: preferencesCookie.Policy.Essential,
		Usage:     preferencesCookie.Policy.Usage,
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
