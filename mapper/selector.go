package mapper

import (
	"fmt"
	"strconv"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/model"
	"github.com/ONSdigital/dp-renderer/v2/helper"
	coreModel "github.com/ONSdigital/dp-renderer/v2/model"
)

// CreateCategorisationsSelector maps data to the Selector model
func (m *Mapper) CreateCategorisationsSelector(dimLabel, dimId string, cats population.GetCategorisationsResponse) model.Selector {
	p := model.Selector{
		Page: m.basePage,
	}
	mapCommonProps(m.req, &p.Page, "filter-flex-selector", cleanDimensionLabel(dimLabel), m.lang, m.serviceMsg, m.eb)
	p.Breadcrumb = []coreModel.TaxonomyNode{
		{
			Title: helper.Localise("Back", m.lang, 1),
			URI:   fmt.Sprintf("/filters/%s/dimensions", m.fid),
		},
	}
	p.LeadText = helper.Localise("SelectCategoriesLeadText", m.lang, 1)
	p.InitialSelection = dimId

	var selections []model.Selection
	for _, cat := range cats.Items {
		cats := []string{}
		for _, c := range sortCategoriesByID(cat.Categories) {
			cats = append(cats, c.Label)
		}
		selections = append(selections, mapCats(cats, m.req.URL.Query()["showAll"], m.lang, m.req.URL.Path, cat.ID, cat.DefaultCategorisation))
	}
	p.Selections = selections

	isValidationError, _ := strconv.ParseBool(m.req.URL.Query().Get("error"))
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
			Language: m.lang,
		}
		p.ErrorId = "categories-error"
	}

	return p
}

// CreateAreaTypeSelector maps data to the Selector model
func (m *Mapper) CreateAreaTypeSelector(areaType []population.AreaType, fDim filter.Dimension, lowest_geography, releaseDate string, dataset dataset.DatasetDetails, hasOpts bool) model.Selector {
	p := model.Selector{
		Page: m.basePage,
	}
	mapCommonProps(m.req, &p.Page, areaPageType, areaTypeTitle, m.lang, m.serviceMsg, m.eb)
	p.Breadcrumb = []coreModel.TaxonomyNode{
		{
			Title: helper.Localise("Back", m.lang, 1),
			URI:   fmt.Sprintf("/filters/%s/dimensions", m.fid),
		},
	}
	p.LeadText = helper.Localise("SelectAreaTypeLeadText", m.lang, 1)
	if hasOpts {
		p.Panel = mapPanel(coreModel.Localisation{
			LocaleKey: "ChangeAreaTypeWarning",
			Plural:    1,
		}, m.lang, []string{"ons-u-mb-l"})
	}

	isValidationError, _ := strconv.ParseBool(m.req.URL.Query().Get("error"))
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
			Language: m.lang,
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
