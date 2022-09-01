package mapper

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/helpers"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/mocks"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/model"
	"github.com/ONSdigital/dp-renderer/helper"
	coreModel "github.com/ONSdigital/dp-renderer/model"
	. "github.com/smartystreets/goconvey/convey"
)

func TestOverview(t *testing.T) {
	helper.InitialiseLocalisationsHelper(mocks.MockAssetFunction)
	mdl := coreModel.Page{}
	req := httptest.NewRequest("", "/", nil)
	lang := "en"
	showAll := []string{}
	filterJob := filter.GetFilterResponse{
		Dataset: filter.Dataset{
			DatasetID: "example",
			Edition:   "2021",
			Version:   1,
		},
	}
	filterDims := filter.Dimensions{
		Items: []filter.Dimension{
			{
				Name:       "Dim 1",
				IsAreaType: helpers.ToBoolPtr(false),
				Options:    []string{"Opt 1", "Opt 2"},
			},
			{
				Name:       "Truncated dim 1",
				IsAreaType: helpers.ToBoolPtr(false),
				Options: []string{"Opt 1",
					"Opt 2",
					"Opt 3",
					"Opt 4",
					"Opt 5",
					"Opt 6",
					"Opt 7",
					"Opt 8",
					"Opt 9",
					"Opt 10",
					"Opt 11",
					"Opt 12",
					"Opt 13",
					"Opt 14",
					"Opt 15",
					"Opt 16",
					"Opt 17",
					"Opt 18",
					"Opt 19",
					"Opt 20",
				},
			},
			{
				Name:       "Truncated dim 2",
				IsAreaType: helpers.ToBoolPtr(false),
				Options: []string{"Opt 1",
					"Opt 2",
					"Opt 3",
					"Opt 4",
					"Opt 5",
					"Opt 6",
					"Opt 7",
					"Opt 8",
					"Opt 9",
					"Opt 10",
					"Opt 11",
					"Opt 12",
				},
			},
			{
				Name:       "An area dim",
				IsAreaType: helpers.ToBoolPtr(true),
				Options: []string{
					"area 1",
					"area 2",
					"area 3",
					"area 4",
					"area 5",
					"area 6",
					"area 7",
					"area 8",
					"area 9",
					"area 10",
				},
			}},
	}
	datasetDims := dataset.VersionDimensions{
		Items: []dataset.VersionDimension{
			{
				Description: "A description on one line",
				Label:       "Dimension 1",
			},
			{
				Description: "A description on one line \n Then a line break",
				Label:       "Dimension 2",
			},
			{
				Description: "",
				Label:       "Only a name - I shouldn't map",
			},
		},
	}

	Convey("test filter flex overview maps correctly", t, func() {
		m := CreateFilterFlexOverview(req, mdl, lang, "", showAll, filterJob, filterDims, datasetDims, false)
		So(m.BetaBannerEnabled, ShouldBeTrue)
		So(m.Type, ShouldEqual, "filter-flex-overview")
		So(m.Metadata.Title, ShouldEqual, "Review changes")
		So(m.Breadcrumb[0].Title, ShouldEqual, "Back")
		So(m.Breadcrumb[0].URI, ShouldEqual, fmt.Sprintf("/datasets/%s/editions/%s/versions/%s",
			filterJob.Dataset.DatasetID,
			filterJob.Dataset.Edition,
			strconv.Itoa(filterJob.Dataset.Version)))
		So(m.Language, ShouldEqual, lang)

		So(m.Dimensions[0].Name, ShouldEqual, filterDims.Items[3].Label)
		So(m.Dimensions[0].IsAreaType, ShouldBeTrue)
		So(m.Dimensions[0].IsCoverage, ShouldBeFalse)
		So(m.Dimensions[0].Options, ShouldResemble, filterDims.Items[3].Options)
		So(m.Dimensions[0].OptionsCount, ShouldEqual, 10)
		So(m.Dimensions[0].ID, ShouldEqual, filterDims.Items[3].ID)
		So(m.Dimensions[0].URI, ShouldEqual, fmt.Sprintf("%s/%s", "", filterDims.Items[3].Name))
		So(m.Dimensions[0].IsTruncated, ShouldBeFalse)

		So(m.Dimensions[1].Name, ShouldBeBlank)
		So(m.Dimensions[1].IsAreaType, ShouldBeFalse)
		So(m.Dimensions[1].IsCoverage, ShouldBeTrue)
		So(m.Dimensions[1].IsDefaultCoverage, ShouldBeFalse)
		So(m.Dimensions[1].Options, ShouldResemble, filterDims.Items[3].Options)
		So(m.Dimensions[1].URI, ShouldEqual, fmt.Sprintf("%s/%s", "", "geography/coverage"))
		So(m.Dimensions[1].IsTruncated, ShouldBeFalse)

		So(m.Dimensions[2].Name, ShouldEqual, filterDims.Items[0].Label)
		So(m.Dimensions[2].IsAreaType, ShouldBeFalse)
		So(m.Dimensions[2].IsCoverage, ShouldBeFalse)
		So(m.Dimensions[2].ID, ShouldEqual, filterDims.Items[0].ID)
		So(m.Dimensions[2].URI, ShouldEqual, fmt.Sprintf("%s/%s", "", filterDims.Items[0].Name))
		So(m.Dimensions[2].IsTruncated, ShouldBeFalse)

		So(m.Dimensions[3].Name, ShouldEqual, filterDims.Items[1].Label)
		So(m.Dimensions[3].IsAreaType, ShouldBeFalse)
		So(m.Dimensions[3].IsCoverage, ShouldBeFalse)
		So(m.Dimensions[3].ID, ShouldEqual, filterDims.Items[1].ID)
		So(m.Dimensions[3].URI, ShouldEqual, fmt.Sprintf("%s/%s", "", filterDims.Items[1].Name))
		So(m.Dimensions[3].IsTruncated, ShouldBeTrue)

		So(m.Collapsible.CollapsibleItems[0].Subheading, ShouldEqual, datasetDims.Items[0].Label)
		So(m.Collapsible.CollapsibleItems[0].Content[0], ShouldEqual, datasetDims.Items[0].Description)
		So(m.Collapsible.CollapsibleItems[1].Subheading, ShouldEqual, datasetDims.Items[1].Label)
		So(m.Collapsible.CollapsibleItems[1].Content, ShouldResemble, strings.Split(datasetDims.Items[1].Description, "\n"))
		So(m.Collapsible.CollapsibleItems, ShouldHaveLength, 2)
	})

	Convey("test truncation maps as expected", t, func() {
		m := CreateFilterFlexOverview(req, mdl, lang, "", showAll, filterJob, filterDims, datasetDims, false)
		So(m.Dimensions[3].OptionsCount, ShouldEqual, len(filterDims.Items[1].Options))
		So(m.Dimensions[3].Options, ShouldHaveLength, 9)
		So(m.Dimensions[3].Options[:3], ShouldResemble, []string{"Opt 1", "Opt 2", "Opt 3"})
		So(m.Dimensions[3].Options[3:6], ShouldResemble, []string{"Opt 9", "Opt 10", "Opt 11"})
		So(m.Dimensions[3].Options[6:], ShouldResemble, []string{"Opt 18", "Opt 19", "Opt 20"})
		So(m.Dimensions[3].IsTruncated, ShouldBeTrue)

		So(m.Dimensions[4].OptionsCount, ShouldEqual, len(filterDims.Items[2].Options))
		So(m.Dimensions[4].Options, ShouldHaveLength, 9)
		So(m.Dimensions[4].Options[:3], ShouldResemble, []string{"Opt 1", "Opt 2", "Opt 3"})
		So(m.Dimensions[4].Options[3:6], ShouldResemble, []string{"Opt 5", "Opt 6", "Opt 7"})
		So(m.Dimensions[4].Options[6:], ShouldResemble, []string{"Opt 10", "Opt 11", "Opt 12"})
		So(m.Dimensions[4].IsTruncated, ShouldBeTrue)
	})

	Convey("test truncation shows all when parameter given", t, func() {
		m := CreateFilterFlexOverview(req, mdl, lang, "", []string{"Truncated dim 2"}, filterJob, filterDims, datasetDims, false)
		So(m.Dimensions[4].OptionsCount, ShouldEqual, len(filterDims.Items[2].Options))
		So(m.Dimensions[4].Options, ShouldHaveLength, 12)
		So(m.Dimensions[4].IsTruncated, ShouldBeFalse)
	})

	Convey("test area type dimension options do not truncate and map to 'coverage' dimension", t, func() {
		m := CreateFilterFlexOverview(req, mdl, lang, "", showAll, filterJob, filterDims, datasetDims, false)
		So(m.Dimensions[1].Options, ShouldHaveLength, 10)
		So(m.Dimensions[1].IsTruncated, ShouldBeFalse)
		So(m.Dimensions[1].IsCoverage, ShouldBeTrue)
	})

	Convey("given hasNoAreaOptions parameter", t, func() {
		Convey("when parameter is true", func() {
			m := CreateFilterFlexOverview(req, mdl, lang, "", []string{""}, filterJob, filterDims, datasetDims, true)
			Convey("then isDefaultCoverage is set to true", func() {
				So(m.Dimensions[1].IsDefaultCoverage, ShouldBeTrue)
			})
		})
		Convey("when parameter is false", func() {
			m := CreateFilterFlexOverview(req, mdl, lang, "", []string{""}, filterJob, filterDims, datasetDims, false)
			Convey("then isDefaultCoverage is set to false", func() {
				So(m.Dimensions[1].IsDefaultCoverage, ShouldBeFalse)
			})
		})
	})
}

func TestCreateSelector(t *testing.T) {
	helper.InitialiseLocalisationsHelper(mocks.MockAssetFunction)
	mdl := coreModel.Page{}
	req := httptest.NewRequest("", "/", nil)
	lang := "en"
	Convey("test create selector maps correctly", t, func() {
		m := CreateSelector(req, mdl, "dimension Name", lang, "12345")
		So(m.BetaBannerEnabled, ShouldBeTrue)
		So(m.Type, ShouldEqual, "filter-flex-selector")
		So(m.Metadata.Title, ShouldEqual, "Dimension Name")
		So(m.Language, ShouldEqual, lang)
		So(m.Breadcrumb[0].URI, ShouldEqual, "/filters/12345/dimensions")
		So(m.Breadcrumb[0].Title, ShouldEqual, "Back")
	})
}

func TestCreateAreaTypeSelector(t *testing.T) {
	Convey("Given a slice of geography areas", t, func() {
		areas := []population.AreaType{
			{ID: "one", Label: "One", TotalCount: 1},
			{ID: "two", Label: "Two", TotalCount: 2},
		}

		req := httptest.NewRequest("", "/", nil)
		changeDimension := CreateAreaTypeSelector(req, coreModel.Page{}, "en", "12345", areas, filter.Dimension{}, false)

		expectedSelections := []model.Selection{
			{Value: "one", Label: "One", TotalCount: 1},
			{Value: "two", Label: "Two", TotalCount: 2},
		}

		Convey("Maps each geography dimension into a selection", func() {
			So(changeDimension.Selections, ShouldResemble, expectedSelections)
		})
	})

	Convey("Given a valid page", t, func() {
		const lang = "en"
		req := httptest.NewRequest("", "/", nil)
		changeDimension := CreateAreaTypeSelector(req, coreModel.Page{}, lang, "12345", nil, filter.Dimension{}, false)

		Convey("it sets page metadata", func() {
			So(changeDimension.BetaBannerEnabled, ShouldBeTrue)
			So(changeDimension.Type, ShouldEqual, "filter-flex-selector")
			So(changeDimension.Language, ShouldEqual, lang)
		})

		Convey("it sets the title to Area type", func() {
			So(changeDimension.Metadata.Title, ShouldEqual, "Area type")
		})

		Convey("it sets IsAreaType to true", func() {
			So(changeDimension.IsAreaType, ShouldBeTrue)
		})
	})

	Convey("Given the current filter dimension", t, func() {
		const selectionName = "test"
		req := httptest.NewRequest("", "/", nil)
		changeDimension := CreateAreaTypeSelector(req, coreModel.Page{}, "en", "12345", nil, filter.Dimension{ID: selectionName}, false)

		Convey("it returns the value as an initial selection", func() {
			So(changeDimension.InitialSelection, ShouldEqual, selectionName)
		})
	})

	Convey("Given a validation error", t, func() {
		req := httptest.NewRequest("", "/", nil)
		changeDimension := CreateAreaTypeSelector(req, coreModel.Page{}, "en", "12345", nil, filter.Dimension{}, true)

		Convey("it returns a populated error", func() {
			So(changeDimension.Error.Title, ShouldNotBeEmpty)
		})
	})
}

func TestGetCoverage(t *testing.T) {
	helper.InitialiseLocalisationsHelper(mocks.MockAssetFunction)
	Convey("Given a valid page", t, func() {
		const lang = "en"
		req := httptest.NewRequest("", "/", nil)

		Convey("When the parameters are valid", func() {
			coverage := CreateGetCoverage(
				req,
				coreModel.Page{},
				lang,
				"12345",
				"Country",
				"",
				"",
				"",
				"",
				"dim",
				"geogID",
				population.GetAreasResponse{},
				[]model.SelectableElement{},
				population.GetAreaTypeParentsResponse{},
				false)
			Convey("it sets page metadata", func() {
				So(coverage.BetaBannerEnabled, ShouldBeTrue)
				So(coverage.Type, ShouldEqual, "filter-flex-coverage")
				So(coverage.Language, ShouldEqual, lang)
				So(coverage.URI, ShouldEqual, "/")
			})

			Convey("it sets the title to Coverage", func() {
				So(coverage.Metadata.Title, ShouldEqual, "Coverage")
			})

			Convey("it sets the geography to countries", func() {
				So(coverage.Geography, ShouldEqual, "countries")
			})

			Convey("it sets HasNoResults property", func() {
				So(coverage.NameSearchOutput.HasNoResults, ShouldBeFalse)
				So(coverage.ParentSearchOutput.HasNoResults, ShouldBeFalse)
			})

			Convey("it sets the Dimension property", func() {
				So(coverage.Dimension, ShouldEqual, "dim")
			})

			Convey("it sets the GeographyID property", func() {
				So(coverage.GeographyID, ShouldEqual, "geogID")
			})
		})

		Convey("When parent types is populated", func() {
			parents := population.GetAreaTypeParentsResponse{
				AreaTypes: []population.AreaTypes{
					{
						Label: "Area 1",
						ID:    "id",
					},
				},
			}
			coverage := CreateGetCoverage(
				req,
				coreModel.Page{},
				lang,
				"12345",
				"geography",
				"",
				"",
				"",
				"",
				"",
				"",
				population.GetAreasResponse{},
				[]model.SelectableElement{},
				parents,
				false)
			Convey("Then it maps to the ParentSelect property", func() {
				So(coverage.ParentSelect[0].Text, ShouldEqual, parents.AreaTypes[0].Label)
				So(coverage.ParentSelect[0].Value, ShouldEqual, parents.AreaTypes[0].ID)
				So(coverage.ParentSelect[0].IsDisabled, ShouldBeFalse)
				So(coverage.ParentSelect[0].IsSelected, ShouldBeFalse)
			})
		})

		Convey("When parent type is selected", func() {
			parents := population.GetAreaTypeParentsResponse{
				AreaTypes: []population.AreaTypes{
					{
						Label: "Area 1",
						ID:    "id",
					},
				},
			}
			coverage := CreateGetCoverage(
				req,
				coreModel.Page{},
				lang,
				"12345",
				"geography",
				"",
				"",
				"id",
				"",
				"",
				"",
				population.GetAreasResponse{},
				[]model.SelectableElement{},
				parents,
				false)
			Convey("Then it sets the IsSelected property", func() {
				So(coverage.ParentSelect[0].IsSelected, ShouldBeTrue)
			})
		})

		Convey("When more than one parent type is returned", func() {
			parents := population.GetAreaTypeParentsResponse{
				AreaTypes: []population.AreaTypes{
					{
						Label: "Area 1",
						ID:    "id_1",
					},
					{
						Label: "Area 2",
						ID:    "id_2",
					},
				},
			}
			coverage := CreateGetCoverage(
				req,
				coreModel.Page{},
				lang,
				"12345",
				"geography",
				"",
				"",
				"",
				"",
				"",
				"",
				population.GetAreasResponse{},
				[]model.SelectableElement{},
				parents,
				false)
			Convey("Then it maps the ParentSelect default option", func() {
				So(coverage.ParentSelect[0].Text, ShouldEqual, "Select")
				So(coverage.ParentSelect[0].IsDisabled, ShouldBeTrue)
				So(coverage.ParentSelect[0].IsSelected, ShouldBeTrue)
			})
			Convey("Then it maps the ParentSelect properties", func() {
				So(coverage.ParentSelect[1].Text, ShouldEqual, parents.AreaTypes[0].Label)
				So(coverage.ParentSelect[1].Value, ShouldEqual, parents.AreaTypes[0].ID)
				So(coverage.ParentSelect[1].IsDisabled, ShouldBeFalse)
				So(coverage.ParentSelect[1].IsSelected, ShouldBeFalse)
				So(coverage.ParentSelect[2].Text, ShouldEqual, parents.AreaTypes[1].Label)
				So(coverage.ParentSelect[2].Value, ShouldEqual, parents.AreaTypes[1].ID)
				So(coverage.ParentSelect[2].IsDisabled, ShouldBeFalse)
				So(coverage.ParentSelect[2].IsSelected, ShouldBeFalse)
			})
		})

		Convey("When an unknown geography parameter is given", func() {
			coverage := CreateGetCoverage(
				req,
				coreModel.Page{},
				lang,
				"12345",
				"Unknown geography",
				"",
				"",
				"",
				"",
				"",
				"",
				population.GetAreasResponse{},
				[]model.SelectableElement{},
				population.GetAreaTypeParentsResponse{},
				false)
			Convey("Then it sets the geography to unknown geography", func() {
				So(coverage.Geography, ShouldEqual, "unknown geography")
			})
		})

		Convey("When a valid name search is performed", func() {
			mockedSearchResults := population.GetAreasResponse{
				Areas: []population.Area{
					{
						Label: "area one",
						ID:    "area ID",
					},
				},
			}

			coverage := CreateGetCoverage(
				req,
				coreModel.Page{},
				lang,
				"12345",
				"Unknown geography",
				"search",
				"",
				"",
				"name-search",
				"",
				"",
				mockedSearchResults,
				[]model.SelectableElement{},
				population.GetAreaTypeParentsResponse{},
				false)
			Convey("Then it sets HasNoResults property", func() {
				So(coverage.NameSearchOutput.HasNoResults, ShouldBeFalse)
			})

			Convey("Then it maps the search results", func() {
				expectedResult := []model.SelectableElement{
					{
						Text:  mockedSearchResults.Areas[0].Label,
						Value: mockedSearchResults.Areas[0].ID,
						Name:  "add-option",
					},
				}
				So(coverage.NameSearchOutput.SearchResults, ShouldResemble, expectedResult)
			})

			Convey("Then it sets the search input field value", func() {
				So(coverage.NameSearch.Value, ShouldEqual, "search")
			})
		})

		Convey("When a valid parent search is performed", func() {
			mockedSearchResults := population.GetAreasResponse{
				Areas: []population.Area{
					{
						Label: "area one",
						ID:    "area ID",
					},
				},
			}

			coverage := CreateGetCoverage(
				req,
				coreModel.Page{},
				lang,
				"12345",
				"Unknown geography",
				"",
				"search",
				"",
				"parent-search",
				"",
				"",
				mockedSearchResults,
				[]model.SelectableElement{},
				population.GetAreaTypeParentsResponse{},
				false)
			Convey("Then it sets HasNoResults property", func() {
				So(coverage.ParentSearchOutput.HasNoResults, ShouldBeFalse)
			})

			Convey("Then it maps the search results", func() {
				expectedResult := []model.SelectableElement{
					{
						Text:  mockedSearchResults.Areas[0].Label,
						Value: mockedSearchResults.Areas[0].ID,
						Name:  "add-parent-option",
					},
				}
				So(coverage.ParentSearchOutput.SearchResults, ShouldResemble, expectedResult)
			})

			Convey("Then it sets the search input field value", func() {
				So(coverage.ParentSearch.Value, ShouldEqual, "search")
			})
		})

		Convey("When an invalid name search is performed", func() {
			coverage := CreateGetCoverage(
				req,
				coreModel.Page{},
				lang,
				"12345",
				"Unknown geography",
				"search",
				"",
				"",
				"name-search",
				"",
				"",
				population.GetAreasResponse{},
				[]model.SelectableElement{},
				population.GetAreaTypeParentsResponse{},
				false)
			Convey("Then it sets HasNoResults property correctly", func() {
				So(coverage.NameSearchOutput.HasNoResults, ShouldBeTrue)
			})

			Convey("Then search results struct is empty", func() {
				So(coverage.NameSearchOutput.SearchResults, ShouldResemble, []model.SelectableElement(nil))
			})

			Convey("Then it sets the search input field value", func() {
				So(coverage.NameSearch.Value, ShouldEqual, "search")
			})
		})

		Convey("When an invalid parent search is performed", func() {
			coverage := CreateGetCoverage(
				req,
				coreModel.Page{},
				lang,
				"12345",
				"Unknown geography",
				"",
				"search",
				"",
				"parent-search",
				"",
				"",
				population.GetAreasResponse{},
				[]model.SelectableElement{},
				population.GetAreaTypeParentsResponse{},
				false)
			Convey("Then it sets HasNoResults property correctly", func() {
				So(coverage.ParentSearchOutput.HasNoResults, ShouldBeTrue)
			})

			Convey("Then search results struct is empty", func() {
				So(coverage.ParentSearchOutput.SearchResults, ShouldResemble, []model.SelectableElement(nil))
			})

			Convey("Then it sets the search input field value", func() {
				So(coverage.ParentSearch.Value, ShouldEqual, "search")
			})
		})

		Convey("When an option is added", func() {
			mockedOpt := []model.SelectableElement{
				{
					Text:  "label",
					Value: "0",
					Name:  "delete-option",
				},
			}
			coverage := CreateGetCoverage(
				req,
				coreModel.Page{},
				lang,
				"12345",
				"Unknown geography",
				"search",
				"",
				"",
				"",
				"",
				"",
				population.GetAreasResponse{},
				mockedOpt,
				population.GetAreaTypeParentsResponse{},
				false)
			Convey("Then it sets Options property correctly", func() {
				So(coverage.NameSearchOutput.Options, ShouldResemble, mockedOpt)
			})
		})

		Convey("When a parent option is added", func() {
			mockedOpt := []model.SelectableElement{
				{
					Text:  "label",
					Value: "0",
				},
			}
			coverage := CreateGetCoverage(
				req,
				coreModel.Page{},
				lang,
				"12345",
				"",
				"",
				"",
				"",
				"",
				"",
				"",
				population.GetAreasResponse{},
				mockedOpt,
				population.GetAreaTypeParentsResponse{},
				true)
			Convey("Then it sets Options property correctly", func() {
				So(coverage.ParentSearchOutput.Options, ShouldResemble, mockedOpt)
			})
		})

		Convey("When an option is added during a search", func() {
			mockedSearchResults := population.GetAreasResponse{
				Areas: []population.Area{
					{
						Label: "area one",
						ID:    "area ID",
					},
				},
			}
			mockedOpt := []model.SelectableElement{
				{
					Text:  "label",
					Value: "0",
				},
			}
			coverage := CreateGetCoverage(
				req,
				coreModel.Page{},
				lang,
				"12345",
				"Unknown geography",
				"search",
				"",
				"",
				"",
				"",
				"",
				mockedSearchResults,
				mockedOpt,
				population.GetAreaTypeParentsResponse{},
				false)
			Convey("Then it sets Options property correctly", func() {
				So(coverage.NameSearchOutput.Options, ShouldResemble, mockedOpt)
			})

			Convey("Then it sets HasNoResults property", func() {
				So(coverage.NameSearchOutput.HasNoResults, ShouldBeFalse)
			})
		})

		Convey("When an option is added during a parent search", func() {
			mockedSearchResults := population.GetAreasResponse{
				Areas: []population.Area{
					{
						Label: "parent area one",
						ID:    "parent area ID",
					},
				},
			}
			mockedOpt := []model.SelectableElement{
				{
					Text:  "label",
					Value: "0",
					Name:  "delete-option",
				},
			}
			coverage := CreateGetCoverage(
				req,
				coreModel.Page{},
				lang,
				"12345",
				"",
				"",
				"",
				"",
				"",
				"",
				"",
				mockedSearchResults,
				mockedOpt,
				population.GetAreaTypeParentsResponse{},
				true)
			Convey("Then it sets Options property correctly", func() {
				So(coverage.ParentSearchOutput.Options, ShouldResemble, mockedOpt)
			})

			Convey("Then it sets HasNoResults property", func() {
				So(coverage.ParentSearchOutput.HasNoResults, ShouldBeFalse)
			})
		})

		Convey("When a search is performed with one of the results already added as an option", func() {
			mockedSearchResults := population.GetAreasResponse{
				Areas: []population.Area{
					{
						Label: "area one",
						ID:    "0",
					},
					{
						Label: "area two",
						ID:    "1",
					},
				},
			}
			mockedOpt := []model.SelectableElement{
				{
					Text:  "area one",
					Value: "0",
				},
			}
			coverage := CreateGetCoverage(
				req,
				coreModel.Page{},
				lang,
				"12345",
				"Unknown geography",
				"search",
				"",
				"",
				"name-search",
				"",
				"",
				mockedSearchResults,
				mockedOpt,
				population.GetAreaTypeParentsResponse{},
				false)
			Convey("Then it sets Options property correctly", func() {
				So(coverage.NameSearchOutput.Options, ShouldResemble, mockedOpt)
			})

			Convey("Then it sets HasNoResults property", func() {
				So(coverage.NameSearchOutput.HasNoResults, ShouldBeFalse)
			})

			Convey("Then it maps the search results", func() {
				expectedResults := []model.SelectableElement{
					{
						Text:       mockedSearchResults.Areas[0].Label,
						Value:      mockedSearchResults.Areas[0].ID,
						IsSelected: true,
						Name:       "delete-option",
					},
					{
						Text:       mockedSearchResults.Areas[1].Label,
						Value:      mockedSearchResults.Areas[1].ID,
						IsSelected: false,
						Name:       "add-option",
					},
				}
				So(coverage.NameSearchOutput.SearchResults, ShouldResemble, expectedResults)
			})
		})

		Convey("When a parent search is performed with one of the parent results already added as an option", func() {
			mockedSearchResults := population.GetAreasResponse{
				Areas: []population.Area{
					{
						Label: "parent area one",
						ID:    "0",
					},
					{
						Label: "parent area two",
						ID:    "1",
					},
				},
			}
			mockedOpt := []model.SelectableElement{
				{
					Text:  "parent area one",
					Value: "0",
				},
			}
			coverage := CreateGetCoverage(
				req,
				coreModel.Page{},
				lang,
				"12345",
				"Unknown geography",
				"",
				"",
				"",
				"parent-search",
				"",
				"",
				mockedSearchResults,
				mockedOpt,
				population.GetAreaTypeParentsResponse{},
				true)
			Convey("Then it sets Options property correctly", func() {
				So(coverage.ParentSearchOutput.Options, ShouldResemble, mockedOpt)
			})

			Convey("Then it sets HasNoResults property", func() {
				So(coverage.ParentSearchOutput.HasNoResults, ShouldBeFalse)
			})

			Convey("Then it maps the search results", func() {
				expectedResults := []model.SelectableElement{
					{
						Text:       mockedSearchResults.Areas[0].Label,
						Value:      mockedSearchResults.Areas[0].ID,
						IsSelected: true,
						Name:       "delete-option",
					},
					{
						Text:       mockedSearchResults.Areas[1].Label,
						Value:      mockedSearchResults.Areas[1].ID,
						IsSelected: false,
						Name:       "add-parent-option",
					},
				}
				So(coverage.ParentSearchOutput.SearchResults, ShouldResemble, expectedResults)
			})
		})
	})
}

func TestUnitMapCookiesPreferences(t *testing.T) {
	req := httptest.NewRequest("", "/", nil)
	pageModel := coreModel.Page{
		CookiesPreferencesSet: false,
		CookiesPolicy: coreModel.CookiesPolicy{
			Essential: false,
			Usage:     false,
		},
	}

	Convey("cookies preferences initialise as false", t, func() {
		So(pageModel.CookiesPreferencesSet, ShouldBeFalse)
		So(pageModel.CookiesPolicy.Essential, ShouldBeFalse)
		So(pageModel.CookiesPolicy.Usage, ShouldBeFalse)
	})

	Convey("cookie preferences map to page model", t, func() {
		req.AddCookie(&http.Cookie{Name: "cookies_preferences_set", Value: "true"})
		req.AddCookie(&http.Cookie{Name: "cookies_policy", Value: "%7B%22essential%22%3Atrue%2C%22usage%22%3Atrue%7D"})
		mapCookiePreferences(req, &pageModel.CookiesPreferencesSet, &pageModel.CookiesPolicy)
		So(pageModel.CookiesPreferencesSet, ShouldBeTrue)
		So(pageModel.CookiesPolicy.Essential, ShouldBeTrue)
		So(pageModel.CookiesPolicy.Usage, ShouldBeTrue)
	})
}
