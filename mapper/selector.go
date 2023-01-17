package mapper

import (
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/model"
	"github.com/ONSdigital/dp-renderer/helper"
	coreModel "github.com/ONSdigital/dp-renderer/model"
)

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
	p.LeadText = helper.Localise("SelectAreaTypeLeadText", lang, 1)
	if hasOpts {
		p.Panel = mapPanel(coreModel.Localisation{
			LocaleKey: "ChangeAreaTypeWarning",
			Plural:    1,
		}, lang, []string{"ons-u-mb-l"})
	}

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
