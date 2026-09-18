package mapper

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"

	"github.com/ONSdigital/dis-design-system-go/v2/helper"
	core "github.com/ONSdigital/dis-design-system-go/v2/model"
	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/config"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/helpers"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/model"
)

// CreateFilterFlexOverview maps data to the Overview model
func (m *Mapper) CreateFilterFlexOverview(filterJob filter.GetFilterResponse, filterDims []model.FilterDimension, dimDescriptions population.GetDimensionsResponse, pops population.GetPopulationTypeResponse, sdc cantabular.GetBlockedAreaCountResult, isMultivariate bool) model.Overview {
	cfg, _ := config.Get()

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
	p.FeatureFlags.FeedbackAPIURL = cfg.FeedbackAPIURL

	p.Breadcrumb = buildBreadcrumb(dataset, helpers.IsBoolPtr(filterJob.Custom), m.lang)

	pop := model.Dimension{
		Name:        "Population type",
		ID:          pops.PopulationType.Name,
		Options:     []string{pops.PopulationType.Label},
		IsGeography: true,
	}

	coverage := model.Dimension{
		Name:        helper.Localise("AreaTypeCoverageTitle", m.lang, 1),
		IsGeography: true,
		HasChange:   true,
		URI:         fmt.Sprintf("%s/geography/coverage", path),
		ID:          "coverage",
	}

	var area model.Dimension
	for i := range filterDims {
		dim := &filterDims[i]
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

	p.DimensionDescriptions = core.Collapsible{
		Title: core.Localisation{
			LocaleKey: "VariableExplanation",
			Plural:    4,
		},
		CollapsibleItems: mapDescriptionsCollapsible(dimDescriptions, p.Dimensions),
	}

	if isMultivariate {
		maxCellsError := isMaxCellsError(&sdc)
		switch {
		case sdc.Blocked > 0 || maxCellsError: // areas blocked
			p.HasSDC = true
			p.Panel = *m.mapBlockedAreasPanel(&sdc, maxCellsError, model.Pending)

			areaTypeUri, dimNames := mapImproveResultsCollapsible(p.Dimensions)
			p.ImproveResults = core.Collapsible{
				Title: core.Localisation{
					LocaleKey: "ImproveResultsTitle",
					Plural:    4,
				},
				CollapsibleItems: []core.CollapsibleItem{
					{
						Subheading: helper.Localise("ImproveResultsSubHeading", m.lang, 1),
						SafeHTML: core.Localisation{
							Text: helper.Localise("ImproveResultsList", m.lang, 1, areaTypeUri, dimNames),
						},
					},
				},
			}
		case sdc.Passed == sdc.Total && sdc.Total > 0: // all areas passing
			p.HasSDC = true
			p.Panel = *m.mapBlockedAreasPanel(&sdc, maxCellsError, model.Success)
		}

		p.ShowGetDataButton = len(p.Dimensions) > 3 // all geography dimensions (population type, area type and coverage)
	} else {
		p.ShowGetDataButton = true
	}

	if isMaxVariablesError(&sdc) {
		p.Error = core.Error{
			Title: helper.Localise("MaximumVariablesErrorTitle", m.lang, 1),
			ErrorItems: []core.ErrorItem{
				{
					Description: core.Localisation{
						LocaleKey: "MaximumVariablesErrorDescription",
						Plural:    1,
					},
					URL: fmt.Sprintf("%s/change#dimensions--added", path),
				},
			},
			Language: m.lang,
		}
		p.MaxVariableError = true
	} else {
		p.MaxVariableError = false
	}

	p.DisableGetDataButton = isMultivariate && sdc.Passed == 0

	return p
}

func buildBreadcrumb(dataset filter.Dataset, isCustom bool, lang string) []core.TaxonomyNode {
	if isCustom {
		return []core.TaxonomyNode{
			{
				Title: helper.Localise("CustomBack", lang, 1),
				URI:   "/datasets/create",
			},
		}
	} else {
		return []core.TaxonomyNode{
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
