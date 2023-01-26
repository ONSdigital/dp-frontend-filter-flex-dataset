package mapper

import (
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/mocks"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/model"
	"github.com/ONSdigital/dp-renderer/helper"
	coreModel "github.com/ONSdigital/dp-renderer/model"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetCoverage(t *testing.T) {
	helper.InitialiseLocalisationsHelper(mocks.MockAssetFunction)
	Convey("Given a valid page", t, func() {
		const lang = "en"
		req := httptest.NewRequest("", "/", nil)
		eb := getTestEmergencyBanner()
		sm := getTestServiceMessage()
		m := NewMapper(req, coreModel.Page{}, eb, lang, sm, "12345")

		Convey("When the parameters are valid", func() {
			coverage := m.CreateGetCoverage(
				"Country",
				"",
				"",
				"",
				"",
				"",
				"dim",
				"geogID",
				"2022/11/29",
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				population.GetAreasResponse{},
				[]model.SelectableElement{},
				population.GetAreaTypeParentsResponse{},
				false,
				1)
			Convey("it sets page metadata", func() {
				So(coverage.BetaBannerEnabled, ShouldBeTrue)
				So(coverage.Type, ShouldEqual, "coverage_options")
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

			Convey("it sets the SearchNoIndexEnabled to true", func() {
				So(coverage.SearchNoIndexEnabled, ShouldBeTrue)
			})

			Convey("it sets analytics values", func() {
				So(coverage.DatasetId, ShouldEqual, "dataset-id")
				So(coverage.DatasetTitle, ShouldEqual, "Dataset title")
				So(coverage.ReleaseDate, ShouldEqual, "2022/11/29")
			})

			Convey("it sets the emergency banner values", func() {
				So(coverage.EmergencyBanner, ShouldResemble, mappedEmergencyBanner())
			})

			Convey("it maps the service message value", func() {
				So(coverage.ServiceMessage, ShouldEqual, sm)
			})
		})

		Convey("When parent types is populated", func() {
			parents := population.GetAreaTypeParentsResponse{
				AreaTypes: []population.AreaType{
					{
						Label: "Area 1",
						ID:    "id",
					},
				},
			}
			coverage := m.CreateGetCoverage(
				"geography",
				"",
				"",
				"",
				"",
				"",
				"",
				"",
				"",
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				population.GetAreasResponse{},
				[]model.SelectableElement{},
				parents,
				false,
				1)
			Convey("Then it maps to the ParentSelect property", func() {
				So(coverage.ParentSelect[0].Text, ShouldEqual, parents.AreaTypes[0].Label)
				So(coverage.ParentSelect[0].Value, ShouldEqual, parents.AreaTypes[0].ID)
				So(coverage.ParentSelect[0].IsDisabled, ShouldBeFalse)
				So(coverage.ParentSelect[0].IsSelected, ShouldBeFalse)
			})
			Convey("Then it sets the IsSelectParent property", func() {
				So(coverage.IsSelectParents, ShouldBeTrue)
			})
		})

		Convey("When parent types returns an empty list", func() {
			coverage := m.CreateGetCoverage(
				"geography",
				"",
				"",
				"",
				"",
				"",
				"",
				"",
				"",
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				population.GetAreasResponse{},
				[]model.SelectableElement{},
				population.GetAreaTypeParentsResponse{},
				false,
				1)
			Convey("Then it sets the IsSelectParent property", func() {
				So(coverage.IsSelectParents, ShouldBeFalse)
			})
		})

		Convey("When parent type is selected", func() {
			parents := population.GetAreaTypeParentsResponse{
				AreaTypes: []population.AreaType{
					{
						Label: "Area 1",
						ID:    "id",
					},
				},
			}
			coverage := m.CreateGetCoverage(
				"geography",
				"",
				"",
				"id",
				"",
				"",
				"",
				"",
				"",
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				population.GetAreasResponse{},
				[]model.SelectableElement{},
				parents,
				false,
				1)
			Convey("Then it sets the IsSelected property", func() {
				So(coverage.ParentSelect[0].IsSelected, ShouldBeTrue)
			})
			Convey("Then it sets the IsSelectParent property", func() {
				So(coverage.IsSelectParents, ShouldBeTrue)
			})
		})

		Convey("When more than one parent type is returned", func() {
			parents := population.GetAreaTypeParentsResponse{
				AreaTypes: []population.AreaType{
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
			coverage := m.CreateGetCoverage(
				"geography",
				"",
				"",
				"",
				"",
				"",
				"",
				"",
				"",
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				population.GetAreasResponse{},
				[]model.SelectableElement{},
				parents,
				false,
				1)
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
			Convey("Then it sets the IsSelectParent property", func() {
				So(coverage.IsSelectParents, ShouldBeTrue)
			})
		})

		Convey("When an unknown geography parameter is given", func() {
			coverage := m.CreateGetCoverage(
				"Unknown geography",
				"",
				"",
				"",
				"",
				"",
				"",
				"",
				"",
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				population.GetAreasResponse{},
				[]model.SelectableElement{},
				population.GetAreaTypeParentsResponse{},
				false,
				1)
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

			coverage := m.CreateGetCoverage(
				"Unknown geography",
				"search",
				"",
				"",
				"",
				"name-search",
				"",
				"",
				"",
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				mockedSearchResults,
				[]model.SelectableElement{},
				population.GetAreaTypeParentsResponse{},
				false,
				1)
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
				So(coverage.NameSearchOutput.Results, ShouldResemble, expectedResult)
			})

			Convey("Then it sets the search input field value", func() {
				So(coverage.NameSearch.Value, ShouldEqual, "search")
			})
		})

		Convey("When a valid name search is performed with paginated results", func() {
			mockedSearchResults := population.GetAreasResponse{
				Areas: []population.Area{
					{
						Label: "area one",
						ID:    "area ID",
					},
				},
				PaginationResponse: population.PaginationResponse{
					PaginationParams: population.PaginationParams{
						Offset: 0,
						Limit:  50,
					},
					Count:      0,
					TotalCount: 101,
				},
			}

			coverage := m.CreateGetCoverage(
				"Unknown geography",
				"search",
				"",
				"",
				"",
				"name-search",
				"",
				"",
				"",
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				mockedSearchResults,
				[]model.SelectableElement{},
				population.GetAreaTypeParentsResponse{},
				false,
				2)
			Convey("Then it sets HasNoResults property", func() {
				So(coverage.NameSearchOutput.HasNoResults, ShouldBeFalse)
			})

			Convey("Then it paginates the search results", func() {
				expectedPagination := coreModel.Pagination{
					CurrentPage: 2,
					PagesToDisplay: []coreModel.PageToDisplay{
						{
							PageNumber: 1,
							URL:        "/?page=1",
						},
						{
							PageNumber: 2,
							URL:        "/?page=2",
						},
						{
							PageNumber: 3,
							URL:        "/?page=3",
						},
					},
					FirstAndLastPages: []coreModel.PageToDisplay{
						{
							PageNumber: 1,
							URL:        "/?page=1",
						},
						{
							PageNumber: 3,
							URL:        "/?page=3",
						},
					},
					TotalPages: 3,
					Limit:      50,
				}
				So(coverage.NameSearchOutput.Pagination, ShouldResemble, expectedPagination)
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

			coverage := m.CreateGetCoverage(
				"Unknown geography",
				"",
				"search",
				"",
				"parent",
				"parent-search",
				"",
				"",
				"",
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				mockedSearchResults,
				[]model.SelectableElement{},
				population.GetAreaTypeParentsResponse{},
				false,
				1)
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
				So(coverage.ParentSearchOutput.Results, ShouldResemble, expectedResult)
			})

			Convey("Then it sets the search input field value", func() {
				So(coverage.ParentSearch.Value, ShouldEqual, "search")
			})

			Convey("Then it sets the set parent field value", func() {
				So(coverage.SetParent, ShouldEqual, "parent")
			})
		})

		Convey("When an invalid name search is performed", func() {
			coverage := m.CreateGetCoverage(
				"Unknown geography",
				"search",
				"",
				"",
				"",
				"name-search",
				"",
				"",
				"",
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				population.GetAreasResponse{},
				[]model.SelectableElement{},
				population.GetAreaTypeParentsResponse{},
				false,
				1)
			Convey("Then it sets HasNoResults property correctly", func() {
				So(coverage.NameSearchOutput.HasNoResults, ShouldBeTrue)
			})

			Convey("Then search results struct is empty", func() {
				So(coverage.NameSearchOutput.Results, ShouldResemble, []model.SelectableElement(nil))
			})

			Convey("Then it sets the search input field value", func() {
				So(coverage.NameSearch.Value, ShouldEqual, "search")
			})
		})

		Convey("When an invalid parent search is performed", func() {
			mockErrStruct := coreModel.Error{
				Title: "Coverage",
				ErrorItems: []coreModel.ErrorItem{
					{
						Description: coreModel.Localisation{
							LocaleKey: "CoverageSelectDefault",
							Plural:    1,
						},
						URL: "#coverage-error",
					},
				},
				Language: lang,
			}
			m.req = httptest.NewRequest("", "/?error=true", nil)
			coverage := m.CreateGetCoverage(
				"Unknown geography",
				"",
				"search",
				"",
				"parent-search",
				"",
				"",
				"",
				"",
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				population.GetAreasResponse{},
				[]model.SelectableElement{},
				population.GetAreaTypeParentsResponse{},
				false,
				1)
			Convey("Then it sets HasNoResults property correctly", func() {
				So(coverage.ParentSearchOutput.HasNoResults, ShouldBeFalse)
			})

			Convey("Then search results struct is empty", func() {
				So(coverage.ParentSearchOutput.Results, ShouldResemble, []model.SelectableElement(nil))
			})

			Convey("Then it sets the search input field value", func() {
				So(coverage.ParentSearch.Value, ShouldEqual, "search")
			})

			Convey("Then it sets the page Error struct", func() {
				So(coverage.Error, ShouldResemble, mockErrStruct)
			})
		})

		Convey("When a 'no results' parent search is performed", func() {
			coverage := m.CreateGetCoverage(
				"Unknown geography",
				"",
				"search",
				"",
				"",
				"parent-search",
				"",
				"",
				"",
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				population.GetAreasResponse{},
				[]model.SelectableElement{},
				population.GetAreaTypeParentsResponse{},
				false,
				1)
			Convey("Then it sets HasNoResults property correctly", func() {
				So(coverage.ParentSearchOutput.HasNoResults, ShouldBeTrue)
			})

			Convey("Then search results struct is empty", func() {
				So(coverage.ParentSearchOutput.Results, ShouldResemble, []model.SelectableElement(nil))
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
			coverage := m.CreateGetCoverage(
				"Unknown geography",
				"search",
				"",
				"",
				"",
				"",
				"",
				"",
				"",
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				population.GetAreasResponse{},
				mockedOpt,
				population.GetAreaTypeParentsResponse{},
				false,
				1)
			Convey("Then it sets Options property correctly", func() {
				So(coverage.NameSearchOutput.Selections, ShouldResemble, mockedOpt)
			})
		})

		Convey("When a parent option is added", func() {
			mockedOpt := []model.SelectableElement{
				{
					Text:  "label",
					Value: "0",
				},
			}
			coverage := m.CreateGetCoverage(
				"",
				"",
				"",
				"",
				"",
				"",
				"",
				"",
				"",
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				population.GetAreasResponse{},
				mockedOpt,
				population.GetAreaTypeParentsResponse{},
				true,
				1)
			Convey("Then it sets Options property correctly", func() {
				So(coverage.ParentSearchOutput.Selections, ShouldResemble, mockedOpt)
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
			coverage := m.CreateGetCoverage(
				"Unknown geography",
				"search",
				"",
				"",
				"",
				"",
				"",
				"",
				"",
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				mockedSearchResults,
				mockedOpt,
				population.GetAreaTypeParentsResponse{},
				false,
				1)
			Convey("Then it sets Options property correctly", func() {
				So(coverage.NameSearchOutput.Selections, ShouldResemble, mockedOpt)
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
			coverage := m.CreateGetCoverage(
				"",
				"",
				"",
				"",
				"",
				"",
				"",
				"",
				"",
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				mockedSearchResults,
				mockedOpt,
				population.GetAreaTypeParentsResponse{},
				true,
				1)
			Convey("Then it sets Options property correctly", func() {
				So(coverage.ParentSearchOutput.Selections, ShouldResemble, mockedOpt)
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
			coverage := m.CreateGetCoverage(
				"Unknown geography",
				"search",
				"",
				"",
				"",
				"name-search",
				"",
				"",
				"",
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				mockedSearchResults,
				mockedOpt,
				population.GetAreaTypeParentsResponse{},
				false,
				1)
			Convey("Then it sets Options property correctly", func() {
				So(coverage.NameSearchOutput.Selections, ShouldResemble, mockedOpt)
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
				So(coverage.NameSearchOutput.Results, ShouldResemble, expectedResults)
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
			coverage := m.CreateGetCoverage(
				"Unknown geography",
				"",
				"",
				"",
				"",
				"parent-search",
				"",
				"",
				"",
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				mockedSearchResults,
				mockedOpt,
				population.GetAreaTypeParentsResponse{},
				true,
				1)
			Convey("Then it sets Options property correctly", func() {
				So(coverage.ParentSearchOutput.Selections, ShouldResemble, mockedOpt)
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
				So(coverage.ParentSearchOutput.Results, ShouldResemble, expectedResults)
			})
		})

		Convey("When a parent search is performed with paginated results", func() {
			mockedSearchResults := population.GetAreasResponse{
				Areas: []population.Area{
					{
						Label: "parent area one",
						ID:    "0",
					},
				},
				PaginationResponse: population.PaginationResponse{
					PaginationParams: population.PaginationParams{
						Offset: 0,
						Limit:  50,
					},
					Count:      0,
					TotalCount: 101,
				},
			}
			coverage := m.CreateGetCoverage(
				"Unknown geography",
				"",
				"",
				"",
				"",
				"parent-search",
				"",
				"",
				"",
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				mockedSearchResults,
				[]model.SelectableElement{},
				population.GetAreaTypeParentsResponse{},
				true,
				2)

			Convey("Then it sets HasNoResults property", func() {
				So(coverage.ParentSearchOutput.HasNoResults, ShouldBeFalse)
			})

			Convey("Then it paginates the search results", func() {
				expectedPagination := coreModel.Pagination{
					CurrentPage: 2,
					PagesToDisplay: []coreModel.PageToDisplay{
						{
							PageNumber: 1,
							URL:        "/?page=1",
						},
						{
							PageNumber: 2,
							URL:        "/?page=2",
						},
						{
							PageNumber: 3,
							URL:        "/?page=3",
						},
					},
					FirstAndLastPages: []coreModel.PageToDisplay{
						{
							PageNumber: 1,
							URL:        "/?page=1",
						},
						{
							PageNumber: 3,
							URL:        "/?page=3",
						},
					},
					TotalPages: 3,
					Limit:      50,
				}
				So(coverage.ParentSearchOutput.Pagination, ShouldResemble, expectedPagination)
			})
		})
	})
}
