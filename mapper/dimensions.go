package mapper

import (
	"fmt"

	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/model"
	"github.com/ONSdigital/dp-renderer/helper"
	coreModel "github.com/ONSdigital/dp-renderer/model"
)

// CreateGetChangeDimensions maps data to the ChangeDimensions model
func (m *Mapper) CreateGetChangeDimensions(q, formAction string, dims []model.FilterDimension, pDims, results population.GetDimensionsResponse, sdc *population.GetBlockedAreaCountResult) model.ChangeDimensions {
	p := model.ChangeDimensions{
		Page: m.basePage,
	}
	p.Breadcrumb = []coreModel.TaxonomyNode{
		{
			Title: helper.Localise("Back", m.lang, 1),
			URI:   fmt.Sprintf("/filters/%s/dimensions", m.fid),
		},
	}
	p.FormAction = formAction

	selections := []model.SelectableElement{}
	pageDims := []model.Dimension{}
	for _, dim := range dims {
		if !*dim.IsAreaType {
			selections = append(selections, model.SelectableElement{
				Text:  cleanDimensionLabel(dim.Label),
				Value: dim.ID,
				Name:  "delete-option",
			})
		}
		pageDims = append(pageDims, model.Dimension{
			IsGeography: *dim.IsAreaType,
			ID:          dim.ID,
			URI:         fmt.Sprintf("/filters/%s/dimensions/%s", m.fid, dim.Name),
			Name:        cleanDimensionLabel(dim.Label),
		})
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
	var title string
	if len(selections) > 0 {
		title = "Add or remove variables"
	} else {
		title = "Add variables"
	}
	mapCommonProps(m.req, &p.Page, "change_variables", title, m.lang, m.serviceMsg, m.eb)

	browseResults := mapDimensionsResponse(pDims, &selections, m.lang)
	searchResults := mapDimensionsResponse(results, &selections, m.lang)

	p.Output.Results = browseResults
	p.SearchOutput.Results = searchResults
	p.SearchOutput.HasNoResults = len(p.SearchOutput.Results) == 0 && formAction == "search"

	if sdc.Blocked > 0 {
		p.HasSDC = true
		p.Panel = *m.mapBlockedAreasPanel(sdc, model.Pending)

		areaTypeUri, dimNames := mapImproveResultsCollapsible(pageDims)
		p.ImproveResults = coreModel.Collapsible{
			Title: coreModel.Localisation{
				LocaleKey: "ImproveResultsTitle",
				Plural:    4,
			},
			Language: m.lang,
			CollapsibleItems: []coreModel.CollapsibleItem{
				{
					Subheading: helper.Localise("ImproveResultsSubHeading", m.lang, 1),
					SafeHTML: coreModel.Localisation{
						Text: helper.Localise("ImproveResultsListVariant", m.lang, 1, areaTypeUri, dimNames),
					},
				},
			},
		}
	} else {
		p.Panel = mapPanel(coreModel.Localisation{
			LocaleKey: "DimensionsChangeWarning",
			Plural:    1,
		}, m.lang, []string{"ons-u-mb-s"})
	}

	return p
}
