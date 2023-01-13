package mapper

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/helpers"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/model"
	"github.com/ONSdigital/dp-renderer/helper"
	coreModel "github.com/ONSdigital/dp-renderer/model"
)

// CreateFilterFlexOverview maps data to the Overview model
func CreateFilterFlexOverview(req *http.Request, basePage coreModel.Page, lang, path string, queryStrValues []string, filterJob filter.GetFilterResponse, filterDims []model.FilterDimension, datasetDims dataset.VersionDimensions, hasNoAreaOptions, isMultivariate bool, eb zebedee.EmergencyBanner, serviceMsg string) model.Overview {
	p := model.Overview{
		Page: basePage,
	}
	mapCommonProps(req, &p.Page, reviewPageType, "Review changes", lang, serviceMsg, eb)
	p.FilterID = filterJob.FilterID
	dataset := filterJob.Dataset
	p.IsMultivariate = isMultivariate

	p.Breadcrumb = []coreModel.TaxonomyNode{
		{
			Title: helper.Localise("Back", lang, 1),
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

	// TODO: Temporarily removing mapping as new endpoints are required to return reliable dataset information

	return p
}
