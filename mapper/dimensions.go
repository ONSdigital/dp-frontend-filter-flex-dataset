package mapper

import (
	"fmt"

	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/model"
	"github.com/ONSdigital/dp-renderer/helper"
	coreModel "github.com/ONSdigital/dp-renderer/model"
)

// CreateGetChangeDimensions maps data to the ChangeDimensions model
func (m *Mapper) CreateGetChangeDimensions(q, formAction string, dims []model.FilterDimension, pDims, results population.GetDimensionsResponse) model.ChangeDimensions {
	p := model.ChangeDimensions{
		Page: m.basePage,
	}
	p.Breadcrumb = []coreModel.TaxonomyNode{
		{
			Title: helper.Localise("Back", m.lang, 1),
			URI:   fmt.Sprintf("/filters/%s/dimensions", m.fid),
		},
	}
	mapCommonProps(m.req, &p.Page, "change_variables", "Add or remove variables", m.lang, m.serviceMsg, m.eb)
	p.Panel = mapPanel(coreModel.Localisation{
		LocaleKey: "DimensionsChangeWarning",
		Plural:    1,
	}, m.lang, []string{"ons-u-mb-s"})
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
		Language: m.lang,
		Value:    q,
		Label:    helper.Localise("DimensionsSearchLabel", m.lang, 1),
	}

	browseResults := mapDimensionsResponse(pDims, &selections)
	searchResults := mapDimensionsResponse(results, &selections)

	p.Output.Results = browseResults
	p.SearchOutput.Results = searchResults
	p.SearchOutput.HasNoResults = len(p.SearchOutput.Results) == 0 && formAction == "search"

	return p
}
