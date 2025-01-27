package mapper

import (
	"fmt"
	"sort"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/config"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/model"
	"github.com/ONSdigital/dp-renderer/v2/helper"
	coreModel "github.com/ONSdigital/dp-renderer/v2/model"
)

// CreateGetChangeDimensions maps data to the ChangeDimensions model
func (m *Mapper) CreateGetChangeDimensions(q, formAction string, dims []model.FilterDimension, pDims, results population.GetDimensionsResponse, sdc *cantabular.GetBlockedAreaCountResult) model.ChangeDimensions {
	cfg, _ := config.Get()

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
	p.FeatureFlags.EnableFeedbackAPI = cfg.EnableFeedbackAPI
	p.FeatureFlags.FeedbackAPIURL = cfg.FeedbackAPIURL

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
			HasChange:   dim.CategorisationCount > 1,
		})
	}
	sort.Slice(selections, func(i, j int) bool {
		return selections[i].Text < selections[j].Text
	})
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

	maxCellsError := isMaxCellsError(sdc)
	if sdc.Blocked > 0 || maxCellsError {
		p.HasSDC = true
		p.Panel = *m.mapBlockedAreasPanel(sdc, maxCellsError, model.Pending)
		sort.Slice(pageDims, func(i, j int) bool {
			return pageDims[i].Name < pageDims[j].Name
		})
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

	if isMaxVariablesError(sdc) {
		p.Page.Error = coreModel.Error{
			Title: helper.Localise("MaximumVariablesErrorTitle", m.lang, 1),
			ErrorItems: []coreModel.ErrorItem{
				{
					Description: coreModel.Localisation{
						LocaleKey: "MaximumVariablesErrorDescription",
						Plural:    1,
					},
					URL: "#dimensions--added",
				},
			},
			Language: m.lang,
		}
		p.MaxVariableError = true
		p.Output.HasValidationError = true
	} else {
		p.MaxVariableError = false
	}

	return p
}
