package mapper

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/helpers"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/model"
	"github.com/ONSdigital/dp-renderer/helper"
	coreModel "github.com/ONSdigital/dp-renderer/model"
)

// CreateFilterFlexOverview maps data to the Overview model
func (m *Mapper) CreateFilterFlexOverview(filterJob filter.GetFilterResponse, filterDims []model.FilterDimension, dimDescriptions population.GetDimensionsResponse, pops population.GetPopulationTypesResponse, sdc cantabular.GetBlockedAreaCountResult, isMultivariate bool) model.Overview {
	queryStrValues := m.req.URL.Query()["showAll"]
	path := m.req.URL.Path

	p := model.Overview{
		Page: m.basePage,
	}

	title := helper.Localise("OverviewTitle", m.lang, 1)
	if helpers.IsBoolPtr(filterJob.Custom) {
		title = helper.Localise("OverviewCustomTitle", m.lang, 1)
	}

	mapCommonProps(m.req, &p.Page, reviewPageType, title, m.lang, m.serviceMsg, m.eb)
	p.FilterID = filterJob.FilterID
	dataset := filterJob.Dataset
	p.IsMultivariate = isMultivariate

	p.Breadcrumb = buildBreadcrumb(dataset, helpers.IsBoolPtr(filterJob.Custom), m.lang)

	pop := model.Dimension{
		Name:        "Population type",
		ID:          filterJob.PopulationType,
		IsGeography: true,
	}

	for _, population := range pops.Items {
		if population.Name == filterJob.PopulationType {
			pop.Options = []string{population.Label}
			break
		}
	}

	coverage := model.Dimension{
		Name:        helper.Localise("AreaTypeCoverageTitle", m.lang, 1),
		IsGeography: true,
		HasChange:   true,
		URI:         fmt.Sprintf("%s/geography/coverage", path),
		ID:          "coverage",
	}

	var area model.Dimension
	for _, dim := range filterDims {
		if *dim.IsAreaType {
			area.Name = helper.Localise("AreaTypeDescription", m.lang, 1)
			area.Options = []string{cleanDimensionLabel(dim.Label)}
			area.IsGeography = true
			area.OptionsCount = dim.OptionsCount
			coverage.Options = dim.Options
			area.ID = dim.ID
			area.URI = fmt.Sprintf("%s/%s", path, dim.Name)
			area.HasChange = true
		} else {
			pageDim := model.Dimension{}
			pageDim.Name = cleanDimensionLabel(dim.Label)
			pageDim.OptionsCount = dim.OptionsCount
			pageDim.IsGeography = *dim.IsAreaType
			pageDim.ID = dim.ID
			pageDim.URI = fmt.Sprintf("%s/%s", path, dim.Name)
			pageDim.HasChange = isMultivariate && dim.CategorisationCount > 1
			pageDim.HasCategories = true
			q := url.Values{}
			midFloor, midCeiling := getTruncationMidRange(dim.OptionsCount)

			var displayedOptions []string
			if len(dim.Options) > 9 && !helpers.HasStringInSlice(dim.Name, queryStrValues) {
				displayedOptions = append(displayedOptions, dim.Options[:3]...)
				displayedOptions = append(displayedOptions, dim.Options[midFloor:midCeiling]...)
				displayedOptions = append(displayedOptions, dim.Options[len(dim.Options)-3:]...)
				q.Add(queryStrKey, dim.Name)
				helpers.PersistExistingParams(queryStrValues, queryStrKey, dim.Name, q)
				pageDim.IsTruncated = true
			} else {
				helpers.PersistExistingParams(queryStrValues, queryStrKey, dim.Name, q)
				displayedOptions = dim.Options
				pageDim.IsTruncated = false
			}

			pageDim.Options = append(pageDim.Options, displayedOptions...)
			pageDim.TruncateLink = generateTruncatePath(path, dim.ID, q)
			p.Dimensions = append(p.Dimensions, pageDim)
		}
	}

	if len(coverage.Options) == 0 {
		coverage.Options = []string{helper.Localise("AreaTypeDefaultCoverage", m.lang, 1)}
	}

	sort.Slice(p.Dimensions, func(i, j int) bool {
		return p.Dimensions[i].Name < p.Dimensions[j].Name
	})

	p.Dimensions = append([]model.Dimension{
		pop,
		area,
		coverage,
	}, p.Dimensions...)

	p.DimensionDescriptions = coreModel.Collapsible{
		Title: coreModel.Localisation{
			LocaleKey: "VariableExplanation",
			Plural:    4,
		},
		CollapsibleItems: mapDescriptionsCollapsible(dimDescriptions, p.Dimensions),
	}

	if isMultivariate {
		switch {
		case sdc.Blocked > 0: // areas blocked
			p.HasSDC = true
			p.Panel = *m.mapBlockedAreasPanel(&sdc, model.Pending)

			areaTypeUri, dimNames := mapImproveResultsCollapsible(p.Dimensions)
			p.ImproveResults = coreModel.Collapsible{
				Title: coreModel.Localisation{
					LocaleKey: "ImproveResultsTitle",
					Plural:    4,
				},
				CollapsibleItems: []coreModel.CollapsibleItem{
					{
						Subheading: helper.Localise("ImproveResultsSubHeading", m.lang, 1),
						SafeHTML: coreModel.Localisation{
							Text: helper.Localise("ImproveResultsList", m.lang, 1, areaTypeUri, dimNames),
						},
					},
				},
			}
		case sdc.Passed == sdc.Total && sdc.Total > 0: // all areas passing
			p.HasSDC = true
			p.Panel = *m.mapBlockedAreasPanel(&sdc, model.Success)
		}

		p.EnableGetData = len(p.Dimensions) > 3 // all geography dimensions (population type, area type and coverage)
	} else {
		p.EnableGetData = true
	}

	if isMaxVariablesError(&sdc) {
		p.Page.Error = coreModel.Error{
			Title: helper.Localise("MaximumVariablesErrorTitle", m.lang, 1),
			ErrorItems: []coreModel.ErrorItem{
				{
					Description: coreModel.Localisation{
						LocaleKey: "MaximumVariablesErrorDescription",
						Plural:    1,
					},
					URL: fmt.Sprintf("%s/change", path),
				},
			},
			Language: m.lang,
		}
		p.MaxVariableError = true
	} else {
		p.MaxVariableError = false
	}

	return p
}

func buildBreadcrumb(dataset filter.Dataset, isCustom bool, lang string) []coreModel.TaxonomyNode {
	if isCustom {
		return []coreModel.TaxonomyNode{
			{
				Title: helper.Localise("CustomBack", lang, 1),
				URI:   "/datasets/create",
			},
		}
	} else {
		return []coreModel.TaxonomyNode{
			{
				Title: helper.Localise("Back", lang, 1),
				URI: fmt.Sprintf("/datasets/%s/editions/%s/versions/%s",
					dataset.DatasetID,
					dataset.Edition,
					strconv.Itoa(dataset.Version)),
			},
		}
	}
}
