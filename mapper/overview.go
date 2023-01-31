package mapper

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/helpers"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/model"
	"github.com/ONSdigital/dp-renderer/helper"
	coreModel "github.com/ONSdigital/dp-renderer/model"
)

// CreateFilterFlexOverview maps data to the Overview model
func (m *Mapper) CreateFilterFlexOverview(filterJob filter.GetFilterResponse, filterDims []model.FilterDimension, dimDescriptions population.GetDimensionsResponse, hasNoAreaOptions, isMultivariate bool) model.Overview {
	queryStrValues := m.req.URL.Query()["showAll"]
	path := m.req.URL.Path

	p := model.Overview{
		Page: m.basePage,
	}
	mapCommonProps(m.req, &p.Page, reviewPageType, "Review changes", m.lang, m.serviceMsg, m.eb)
	p.FilterID = filterJob.FilterID
	dataset := filterJob.Dataset
	p.IsMultivariate = isMultivariate

	p.Breadcrumb = []coreModel.TaxonomyNode{
		{
			Title: helper.Localise("Back", m.lang, 1),
			URI: fmt.Sprintf("/datasets/%s/editions/%s/versions/%s",
				dataset.DatasetID,
				dataset.Edition,
				strconv.Itoa(dataset.Version)),
		},
	}

	for _, dim := range filterDims {
		pageDim := model.Dimension{}
		pageDim.Name = cleanDimensionLabel(dim.Label)
		pageDim.IsAreaType = *dim.IsAreaType
		pageDim.OptionsCount = dim.OptionsCount
		pageDim.ID = dim.ID
		pageDim.URI = fmt.Sprintf("%s/%s", path, dim.Name)
		pageDim.IsChangeCategories = isMultivariate
		q := url.Values{}
		midFloor, midCeiling := getTruncationMidRange(dim.OptionsCount)

		var displayedOptions []string
		if len(dim.Options) > 9 && !helpers.HasStringInSlice(dim.Name, queryStrValues) && !*dim.IsAreaType {
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

	sort.Slice(p.Dimensions, func(i, j int) bool {
		return p.Dimensions[i].IsAreaType
	})

	coverage := []model.Dimension{
		{
			IsCoverage:        true,
			IsDefaultCoverage: hasNoAreaOptions,
			URI:               fmt.Sprintf("%s/geography/coverage", path),
			Options:           p.Dimensions[0].Options,
			ID:                "coverage",
		},
	}
	temp := append(coverage, p.Dimensions[1:]...)
	p.Dimensions = append(p.Dimensions[:1], temp...)

	p.Collapsible = coreModel.Collapsible{
		Title: coreModel.Localisation{
			LocaleKey: "VariableExplanation",
			Plural:    4,
		},
		CollapsibleItems: mapCollapsible(dimDescriptions, p.Dimensions),
	}

	return p
}

func mapCollapsible(dimDescriptions population.GetDimensionsResponse, dims []model.Dimension) []coreModel.CollapsibleItem {
	var collapsibleContentItems []coreModel.CollapsibleItem
	var areaItem coreModel.CollapsibleItem

	for _, dim := range dims {
		for _, dimDescription := range dimDescriptions.Dimensions {
			if dim.ID == dimDescription.ID && !dim.IsAreaType {
				collapsibleContentItems = append(collapsibleContentItems, coreModel.CollapsibleItem{
					Subheading: cleanDimensionLabel(dimDescription.Label),
					Content:    strings.Split(dimDescription.Description, "\n"),
				})
			} else if dim.ID == dimDescription.ID && dim.IsAreaType {
				areaItem.Subheading = cleanDimensionLabel(dimDescription.Label)
				areaItem.Content = strings.Split(dimDescription.Description, "\n")
			}
		}
	}

	collapsibleContentItems = append([]coreModel.CollapsibleItem{
		{
			Subheading: areaTypeTitle,
			SafeHTML: coreModel.Localisation{
				LocaleKey: "VariableInfoAreaType",
				Plural:    1,
			},
		},
		areaItem,
		{
			Subheading: coverageTitle,
			SafeHTML: coreModel.Localisation{
				LocaleKey: "VariableInfoCoverage",
				Plural:    1,
			},
		},
	}, collapsibleContentItems...)

	return collapsibleContentItems
}
