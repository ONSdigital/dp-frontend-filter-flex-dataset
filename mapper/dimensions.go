package mapper

import (
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/model"
	"github.com/ONSdigital/dp-renderer/helper"
	coreModel "github.com/ONSdigital/dp-renderer/model"
)

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
	p.Panel = mapPanel(coreModel.Localisation{
		LocaleKey: "DimensionsChangeWarning",
		Plural:    1,
	}, lang, []string{"ons-u-mb-s"})
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
