package mapper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/helpers"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/model"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/pagination"
	"github.com/ONSdigital/dp-renderer/v2/helper"
	coreModel "github.com/ONSdigital/dp-renderer/v2/model"
	"github.com/ONSdigital/log.go/v2/log"
)

// CreateGetCoverage maps data to the coverage model
func (m *Mapper) CreateGetCoverage(geogName, nameQ, parentQ, parentArea, setParent, coverage, dim, geogID, releaseDate string, dataset dataset.DatasetDetails, areas population.GetAreasResponse, opts []model.SelectableElement, parents population.GetAreaTypeParentsResponse, hasFilterByParent bool, currentPage int) model.Coverage {
	hasValidationErr, _ := strconv.ParseBool(m.req.URL.Query().Get("error"))
	p := model.Coverage{
		Page: m.basePage,
	}
	mapCommonProps(m.req, &p.Page, coveragePageType, coverageTitle, m.lang, m.serviceMsg, m.eb)
	p.Breadcrumb = []coreModel.TaxonomyNode{
		{
			Title: helper.Localise("Back", m.lang, 1),
			URI:   fmt.Sprintf("/filters/%s/dimensions", m.fid),
		},
	}
	geography := helpers.Pluralise(m.req, geogName, m.lang, areaTypePrefix, pluralInt)
	if geography == "" {
		log.Info(m.req.Context(), "pluralisation lookup failed, reverting to initial input", log.Data{
			"initial_input": geogName,
		})
		geography = geogName
	}

	p.Geography = strings.ToLower(geography)
	p.CoverageType = coverage
	p.Dimension = dim
	p.GeographyID = geogID
	p.SetParent = setParent
	p.NameSearch = model.SearchField{
		Name:     nameSearchFieldName,
		ID:       nameSearch,
		Value:    nameQ,
		Language: m.lang,
		Label:    helper.Localise("CoverageSearchLabel", m.lang, 1),
	}
	p.ParentSearch = model.SearchField{
		Name:     parentSearchFieldName,
		ID:       parentSearch,
		Value:    parentQ,
		Language: m.lang,
		Label:    helper.Localise("CoverageSearchLabel", m.lang, 1),
	}

	p.DatasetId = dataset.ID
	p.DatasetTitle = dataset.Title
	p.ReleaseDate = releaseDate

	if len(parents.AreaTypes) > 1 && parentArea == "" {
		p.ParentSelect = []model.SelectableElement{
			{
				Text:       helper.Localise("CoverageSelectDefault", m.lang, 1),
				IsSelected: true,
				IsDisabled: true,
			},
		}
	}
	for _, parent := range parents.AreaTypes {
		var sel model.SelectableElement
		sel.Text = parent.Label
		sel.Value = parent.ID
		if parentArea == parent.ID {
			sel.IsSelected = true
		}
		p.ParentSelect = append(p.ParentSelect, sel)
	}

	var isParentSearch bool
	if coverage == parentSearch {
		isParentSearch = true
	}
	var results []model.SelectableElement
	for _, area := range areas.Areas {
		var result model.SelectableElement
		result.Text = area.Label
		result.Value = area.ID
		result.Name = getAddOptionStr(isParentSearch)
		for _, opt := range opts {
			if opt.Value == area.ID {
				result.IsSelected = true
				result.Name = "delete-option"
				break
			}
		}
		results = append(results, result)
	}

	totalPages := pagination.GetTotalPages(areas.TotalCount, areas.Limit)
	var paginatedResults coreModel.Pagination
	if totalPages > 1 {
		paginatedResults = coreModel.Pagination{
			CurrentPage:       currentPage,
			TotalPages:        totalPages,
			Limit:             areas.Limit,
			FirstAndLastPages: pagination.GetFirstAndLastPages(m.req, totalPages),
			PagesToDisplay:    pagination.GetPagesToDisplay(currentPage, totalPages, m.req),
		}
	}

	if len(opts) > 0 && hasFilterByParent {
		p.CoverageType = parentSearch
		p.ParentSearchOutput.Selections = opts
		p.ParentSearchOutput.SelectionsTitle = helper.Localise("AreasAddedTitle", m.lang, len(opts))
		p.OptionType = parentSearch
	} else if len(opts) > 0 {
		p.CoverageType = nameSearch
		p.NameSearchOutput.Selections = opts
		p.NameSearchOutput.SelectionsTitle = helper.Localise("AreasAddedTitle", m.lang, len(opts))
		p.OptionType = nameSearch
	}

	switch coverage {
	case nameSearch:
		p.CoverageType = nameSearch
		p.NameSearchOutput.Results = results
		p.NameSearchOutput.HasNoResults = len(p.NameSearchOutput.Results) == 0
		p.NameSearchOutput.Pagination = paginatedResults
	case parentSearch:
		p.CoverageType = parentSearch
		p.ParentSearchOutput.Results = results
		p.ParentSearchOutput.HasNoResults = len(p.ParentSearchOutput.Results) == 0 && !hasValidationErr
		p.ParentSearchOutput.Pagination = paginatedResults
	}

	if hasValidationErr {
		p.Page.Error = coreModel.Error{
			Title: p.Metadata.Title,
			ErrorItems: []coreModel.ErrorItem{
				{
					Description: coreModel.Localisation{
						LocaleKey: "CoverageSelectDefault",
						Plural:    1,
					},
					URL: "#coverage-error",
				},
			},
			Language: m.lang,
		}
	}

	p.IsSelectParents = len(parents.AreaTypes) > 0

	return p
}
