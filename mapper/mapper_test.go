package mapper

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
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
	filterDims := []model.FilterDimension{
		{
			Dimension: filter.Dimension{
				Name:       "Dim 1",
				IsAreaType: helpers.ToBoolPtr(false),
				Options:    []string{"Opt 1", "Opt 2"},
			},
			OptionsCount: 2,
		},
		{
			Dimension: filter.Dimension{
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
			OptionsCount: 20,
		},
		{
			Dimension: filter.Dimension{
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
			OptionsCount: 12,
		},
		{
			Dimension: filter.Dimension{
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
			},
			OptionsCount: 10,
		},
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
	eb := getTestEmergencyBanner()
	sm := getTestServiceMessage()

	Convey("test filter flex overview maps correctly", t, func() {
		m := CreateFilterFlexOverview(req, mdl, lang, "", showAll, filterJob, filterDims, datasetDims, false, true, eb, sm)
		So(m.BetaBannerEnabled, ShouldBeTrue)
		So(m.Type, ShouldEqual, "review_changes")
		So(m.Metadata.Title, ShouldEqual, "Review changes")
		So(m.Breadcrumb[0].Title, ShouldEqual, "Back")
		So(m.Breadcrumb[0].URI, ShouldEqual, fmt.Sprintf("/datasets/%s/editions/%s/versions/%s",
			filterJob.Dataset.DatasetID,
			filterJob.Dataset.Edition,
			strconv.Itoa(filterJob.Dataset.Version)))
		So(m.Language, ShouldEqual, lang)
		So(m.SearchNoIndexEnabled, ShouldBeTrue)
		So(m.IsMultivariate, ShouldBeTrue)

		So(m.Dimensions[0].Name, ShouldEqual, filterDims[3].Label)
		So(m.Dimensions[0].IsAreaType, ShouldBeTrue)
		So(m.Dimensions[0].IsCoverage, ShouldBeFalse)
		So(m.Dimensions[0].Options, ShouldResemble, filterDims[3].Options)
		So(m.Dimensions[0].OptionsCount, ShouldEqual, filterDims[3].OptionsCount)
		So(m.Dimensions[0].ID, ShouldEqual, filterDims[3].ID)
		So(m.Dimensions[0].URI, ShouldEqual, fmt.Sprintf("%s/%s", "", filterDims[3].Name))
		So(m.Dimensions[0].IsTruncated, ShouldBeFalse)

		So(m.Dimensions[1].Name, ShouldBeBlank)
		So(m.Dimensions[1].IsAreaType, ShouldBeFalse)
		So(m.Dimensions[1].IsCoverage, ShouldBeTrue)
		So(m.Dimensions[1].IsDefaultCoverage, ShouldBeFalse)
		So(m.Dimensions[1].Options, ShouldResemble, filterDims[3].Options)
		So(m.Dimensions[1].URI, ShouldEqual, fmt.Sprintf("%s/%s", "", "geography/coverage"))
		So(m.Dimensions[1].IsTruncated, ShouldBeFalse)

		So(m.Dimensions[2].Name, ShouldEqual, filterDims[0].Label)
		So(m.Dimensions[2].IsAreaType, ShouldBeFalse)
		So(m.Dimensions[2].IsCoverage, ShouldBeFalse)
		So(m.Dimensions[2].ID, ShouldEqual, filterDims[0].ID)
		So(m.Dimensions[2].URI, ShouldEqual, fmt.Sprintf("%s/%s", "", filterDims[0].Name))
		So(m.Dimensions[2].IsTruncated, ShouldBeFalse)

		So(m.Dimensions[3].Name, ShouldEqual, filterDims[1].Label)
		So(m.Dimensions[3].IsAreaType, ShouldBeFalse)
		So(m.Dimensions[3].IsCoverage, ShouldBeFalse)
		So(m.Dimensions[3].ID, ShouldEqual, filterDims[1].ID)
		So(m.Dimensions[3].URI, ShouldEqual, fmt.Sprintf("%s/%s", "", filterDims[1].Name))
		So(m.Dimensions[3].IsTruncated, ShouldBeTrue)

		So(m.EmergencyBanner, ShouldResemble, mappedEmergencyBanner())
		So(m.ServiceMessage, ShouldEqual, sm)

		// TODO: Removing test coverage until endpoint is created
	})

	Convey("test truncation maps as expected", t, func() {
		m := CreateFilterFlexOverview(req, mdl, lang, "", showAll, filterJob, filterDims, datasetDims, false, false, zebedee.EmergencyBanner{}, "")
		So(m.Dimensions[3].OptionsCount, ShouldEqual, filterDims[1].OptionsCount)
		So(m.Dimensions[3].Options, ShouldHaveLength, 9)
		So(m.Dimensions[3].Options[:3], ShouldResemble, []string{"Opt 1", "Opt 2", "Opt 3"})
		So(m.Dimensions[3].Options[3:6], ShouldResemble, []string{"Opt 9", "Opt 10", "Opt 11"})
		So(m.Dimensions[3].Options[6:], ShouldResemble, []string{"Opt 18", "Opt 19", "Opt 20"})
		So(m.Dimensions[3].IsTruncated, ShouldBeTrue)

		So(m.Dimensions[4].OptionsCount, ShouldEqual, filterDims[2].OptionsCount)
		So(m.Dimensions[4].Options, ShouldHaveLength, 9)
		So(m.Dimensions[4].Options[:3], ShouldResemble, []string{"Opt 1", "Opt 2", "Opt 3"})
		So(m.Dimensions[4].Options[3:6], ShouldResemble, []string{"Opt 5", "Opt 6", "Opt 7"})
		So(m.Dimensions[4].Options[6:], ShouldResemble, []string{"Opt 10", "Opt 11", "Opt 12"})
		So(m.Dimensions[4].IsTruncated, ShouldBeTrue)
	})

	Convey("test truncation shows all when parameter given", t, func() {
		m := CreateFilterFlexOverview(req, mdl, lang, "", []string{"Truncated dim 2"}, filterJob, filterDims, datasetDims, false, false, zebedee.EmergencyBanner{}, "")
		So(m.Dimensions[4].OptionsCount, ShouldEqual, filterDims[2].OptionsCount)
		So(m.Dimensions[4].Options, ShouldHaveLength, 12)
		So(m.Dimensions[4].IsTruncated, ShouldBeFalse)
	})

	Convey("test area type dimension options do not truncate and map to 'coverage' dimension", t, func() {
		m := CreateFilterFlexOverview(req, mdl, lang, "", showAll, filterJob, filterDims, datasetDims, false, false, zebedee.EmergencyBanner{}, "")
		So(m.Dimensions[1].Options, ShouldHaveLength, 10)
		So(m.Dimensions[1].IsTruncated, ShouldBeFalse)
		So(m.Dimensions[1].IsCoverage, ShouldBeTrue)
	})

	Convey("test filter dims format labels using cleanDimensionLabel", t, func() {
		newFilterDims := append([]model.FilterDimension{}, filterDims...)
		newFilterDims = append(newFilterDims, []model.FilterDimension{
			{
				Dimension: filter.Dimension{
					Label:      "Example (21 categories)",
					IsAreaType: helpers.ToBoolPtr(false),
				},
			},
		}...)

		m := CreateFilterFlexOverview(req, mdl, lang, "", showAll, filterJob, newFilterDims, datasetDims, false, false, zebedee.EmergencyBanner{}, "")
		So(m.Dimensions[5].Name, ShouldEqual, "Example")
	})

	Convey("given hasNoAreaOptions parameter", t, func() {
		Convey("when parameter is true", func() {
			m := CreateFilterFlexOverview(req, mdl, lang, "", []string{""}, filterJob, filterDims, datasetDims, true, false, zebedee.EmergencyBanner{}, "")
			Convey("then isDefaultCoverage is set to true", func() {
				So(m.Dimensions[1].IsDefaultCoverage, ShouldBeTrue)
			})
		})
		Convey("when parameter is false", func() {
			m := CreateFilterFlexOverview(req, mdl, lang, "", []string{""}, filterJob, filterDims, datasetDims, false, false, zebedee.EmergencyBanner{}, "")
			Convey("then isDefaultCoverage is set to false", func() {
				So(m.Dimensions[1].IsDefaultCoverage, ShouldBeFalse)
			})
		})
	})
}

func TestCreateCategorisationsSelector(t *testing.T) {
	helper.InitialiseLocalisationsHelper(mocks.MockAssetFunction)
	mdl := coreModel.Page{}
	req := httptest.NewRequest("", "/", nil)
	lang := "en"
	eb := getTestEmergencyBanner()
	sm := getTestServiceMessage()
	Convey("Given a request to the CreateCategorisationsSelector", t, func() {
		Convey("When valid parameters are provided", func() {
			cats := population.GetCategorisationsResponse{
				Items: []population.Dimension{
					{
						ID: "cat_3a",
						Categories: []population.Category{
							{
								ID:    "1",
								Label: "Cat one",
							},
							{
								ID:    "2",
								Label: "Cat two",
							},
							{
								ID:    "3",
								Label: "Cat three",
							},
						},
					},
					{
						ID: "cat_2a",
						Categories: []population.Category{
							{
								ID:    "1",
								Label: "Cat one",
							},
							{
								ID:    "2",
								Label: "Cat two",
							},
						},
					},
					{
						ID: "cat_4a",
						Categories: []population.Category{
							{
								ID:    "1",
								Label: "Cat one",
							},
							{
								ID:    "2",
								Label: "Cat two",
							},
							{
								ID:    "3",
								Label: "Cat three",
							},
							{
								ID:    "4",
								Label: "Cat four",
							},
						},
					},
				},
			}
			m := CreateCategorisationsSelector(req, mdl, "Dimension", lang, "12345", "dim1234", sm, eb, cats, false)

			Convey("Then it maps the page metadata", func() {
				So(m.BetaBannerEnabled, ShouldBeTrue)
				So(m.Type, ShouldEqual, "filter-flex-selector")
				So(m.Metadata.Title, ShouldEqual, "Dimension")
				So(m.Language, ShouldEqual, lang)
				So(m.Breadcrumb[0].URI, ShouldEqual, "/filters/12345/dimensions")
				So(m.Breadcrumb[0].Title, ShouldEqual, "Back")
			})

			Convey("Then it sets the lead text", func() {
				So(m.LeadText, ShouldEqual, "Select categories")
			})

			Convey("Then it sets SearchNoIndexEnabled to false", func() {
				So(m.SearchNoIndexEnabled, ShouldBeTrue)
			})

			Convey("Then it sets InitialSelection to dim1234", func() {
				So(m.InitialSelection, ShouldEqual, "dim1234")
			})

			Convey("Then it maps the service message", func() {
				So(m.ServiceMessage, ShouldEqual, sm)
			})

			Convey("Then it maps the emergency banner", func() {
				So(m.EmergencyBanner, ShouldResemble, mappedEmergencyBanner())
			})

			Convey("Then it maps the categories", func() {
				mockedCats := []model.Selection{
					{
						Value: cats.Items[0].ID,
						Label: "3 categories",
						Categories: []string{
							cats.Items[0].Categories[0].Label,
							cats.Items[0].Categories[1].Label,
							cats.Items[0].Categories[2].Label,
						},
						IsTruncated:     false,
						TruncateLink:    "/#cat_3a",
						CategoriesCount: 3,
					},
					{
						Value: cats.Items[1].ID,
						Label: "2 categories",
						Categories: []string{
							cats.Items[1].Categories[0].Label,
							cats.Items[1].Categories[1].Label,
						},
						IsTruncated:     false,
						TruncateLink:    "/#cat_2a",
						CategoriesCount: 2,
					},
					{
						Value: cats.Items[2].ID,
						Label: "4 categories",
						Categories: []string{
							cats.Items[2].Categories[0].Label,
							cats.Items[2].Categories[1].Label,
							cats.Items[2].Categories[2].Label,
							cats.Items[2].Categories[3].Label,
						},
						IsTruncated:     false,
						TruncateLink:    "/#cat_4a",
						CategoriesCount: 4,
					},
				}
				So(m.Selections, ShouldResemble, mockedCats)
			})
		})
		Convey("When a form validation error occurs", func() {
			m := CreateCategorisationsSelector(req, mdl, "Dimension", lang, "12345", "dim1234", sm, eb, population.GetCategorisationsResponse{}, true)
			Convey("Then it sets the error title", func() {
				So(m.Error.Title, ShouldEqual, "Dimension")
			})

			Convey("Then it populates the error items struct", func() {
				So(m.Error.ErrorItems, ShouldHaveLength, 1)
				So(m.Error.ErrorItems[0].Description.LocaleKey, ShouldEqual, "SelectCategoriesError")
				So(m.Error.ErrorItems[0].URL, ShouldEqual, "#categories-error")
			})

			Convey("Then it sets the ErrorId", func() {
				So(m.ErrorId, ShouldEqual, "categories-error")
			})
		})

		Convey("When categories are greater than 9", func() {
			cats := population.GetCategorisationsResponse{
				Items: []population.Dimension{
					{
						ID: "cat_12a",
						Categories: []population.Category{
							{
								ID:    "1",
								Label: "Cat one",
							},
							{
								ID:    "2",
								Label: "Cat two",
							},
							{
								ID:    "3",
								Label: "Cat three",
							},
							{
								ID:    "4",
								Label: "Cat four",
							},
							{
								ID:    "5",
								Label: "Cat five",
							},
							{
								ID:    "6",
								Label: "Cat six",
							},
							{
								ID:    "7",
								Label: "Cat seven",
							},
							{
								ID:    "8",
								Label: "Cat eight",
							},
							{
								ID:    "9",
								Label: "Cat nine",
							},
							{
								ID:    "10",
								Label: "Cat ten",
							},
							{
								ID:    "11",
								Label: "Cat eleven",
							},
							{
								ID:    "12",
								Label: "Cat twelve",
							},
						},
					},
				},
			}
			Convey("Then categories are truncated as expected", func() {
				m := CreateCategorisationsSelector(req, mdl, "Dimension", lang, "12345", "dim1234", "", zebedee.EmergencyBanner{}, cats, false)
				truncCat := []model.Selection{
					{
						Value: cats.Items[0].ID,
						Label: "12 categories",
						Categories: []string{
							cats.Items[0].Categories[0].Label,
							cats.Items[0].Categories[1].Label,
							cats.Items[0].Categories[2].Label,
							cats.Items[0].Categories[4].Label,
							cats.Items[0].Categories[5].Label,
							cats.Items[0].Categories[6].Label,
							cats.Items[0].Categories[9].Label,
							cats.Items[0].Categories[10].Label,
							cats.Items[0].Categories[11].Label,
						},
						CategoriesCount: 12,
						IsTruncated:     true,
						TruncateLink:    "/?showAll=cat_12a#cat_12a",
					},
				}
				So(m.Selections, ShouldResemble, truncCat)
			})

			Convey("Then a showAll request shows all categories as expected", func() {
				req := httptest.NewRequest("", "/?showAll=cat_12a", nil)
				m := CreateCategorisationsSelector(req, mdl, "Dimension", lang, "12345", "dim1234", "", zebedee.EmergencyBanner{}, cats, false)
				allCats := []model.Selection{
					{
						Value: cats.Items[0].ID,
						Label: "12 categories",
						Categories: []string{
							cats.Items[0].Categories[0].Label,
							cats.Items[0].Categories[1].Label,
							cats.Items[0].Categories[2].Label,
							cats.Items[0].Categories[3].Label,
							cats.Items[0].Categories[4].Label,
							cats.Items[0].Categories[5].Label,
							cats.Items[0].Categories[6].Label,
							cats.Items[0].Categories[7].Label,
							cats.Items[0].Categories[8].Label,
							cats.Items[0].Categories[9].Label,
							cats.Items[0].Categories[10].Label,
							cats.Items[0].Categories[11].Label,
						},
						CategoriesCount: 12,
						IsTruncated:     false,
						TruncateLink:    "/#cat_12a",
					},
				}
				So(m.Selections, ShouldResemble, allCats)
			})
		})
	})
}

func TestCreateAreaTypeSelector(t *testing.T) {
	helper.InitialiseLocalisationsHelper(mocks.MockAssetFunction)
	eb := getTestEmergencyBanner()
	sm := getTestServiceMessage()
	Convey("Given a slice of geography areas", t, func() {
		areas := []population.AreaType{
			{ID: "one", Label: "One", Description: "One description", TotalCount: 1},
			{ID: "two", Label: "Two", Description: "Two description", TotalCount: 2},
		}

		req := httptest.NewRequest("", "/", nil)
		changeDimension := CreateAreaTypeSelector(req, false, coreModel.Page{}, "en", "12345", areas, filter.Dimension{}, "", "", dataset.DatasetDetails{}, false, false, "", zebedee.EmergencyBanner{})

		expectedSelections := []model.Selection{
			{Value: "one", Label: "One", Description: "One description", TotalCount: 1},
			{Value: "two", Label: "Two", Description: "Two description", TotalCount: 2},
		}

		Convey("Maps each geography dimension into a selection", func() {
			So(changeDimension.Selections, ShouldResemble, expectedSelections)
		})
	})

	Convey("Given a slice of geography areas", t, func() {
		areas := []population.AreaType{
			{ID: "nat", Label: "Nation", TotalCount: 1},
			{ID: "ctry", Label: "Country", TotalCount: 2},
			{ID: "rgn", Label: "Region", TotalCount: 3},
			{ID: "utla", Label: "UTLA", TotalCount: 4},
		}

		req := httptest.NewRequest("", "/", nil)
		changeDimension := CreateAreaTypeSelector(req, false, coreModel.Page{}, "en", "12345", areas, filter.Dimension{}, "", "", dataset.DatasetDetails{}, false, false, "", zebedee.EmergencyBanner{})

		expectedSelections := []model.Selection{
			{Value: "nat", Label: "Nation", TotalCount: 1},
			{Value: "ctry", Label: "Country", TotalCount: 2},
			{Value: "rgn", Label: "Region", TotalCount: 3},
			{Value: "utla", Label: "UTLA", TotalCount: 4},
		}

		Convey("Maps each geography dimension into a selection", func() {
			So(changeDimension.Selections, ShouldResemble, expectedSelections)
		})
	})

	Convey("Given a slice of standard geography areas out of order", t, func() {
		areas := []population.AreaType{
			{ID: "rgn", Label: "Region", TotalCount: 11},
			{ID: "ctry", Label: "Country", TotalCount: 33},
			{ID: "nat", Label: "Nation", TotalCount: 1},
			{ID: "utla", Label: "UTLA", TotalCount: 7},
		}

		req := httptest.NewRequest("", "/", nil)
		changeDimension := CreateAreaTypeSelector(req, true, coreModel.Page{}, "en", "12345", areas, filter.Dimension{}, "", "", dataset.DatasetDetails{}, false, false, "", zebedee.EmergencyBanner{})

		Convey("Sorts selections ascending by standard order", func() {
			expectedSelections := []model.Selection{
				{Value: "nat", Label: "Nation", TotalCount: 1},
				{Value: "ctry", Label: "Country", TotalCount: 33},
				{Value: "rgn", Label: "Region", TotalCount: 11},
				{Value: "utla", Label: "UTLA", TotalCount: 7},
			}

			So(changeDimension.Selections, ShouldResemble, expectedSelections)
		})
	})

	Convey("Given a slice of non-standard geography areas", t, func() {
		areas := []population.AreaType{
			{ID: "three", Label: "Three", TotalCount: 3},
			{ID: "two", Label: "Two", TotalCount: 2},
			{ID: "one", Label: "One", TotalCount: 1},
		}

		req := httptest.NewRequest("", "/", nil)
		changeDimension := CreateAreaTypeSelector(req, false, coreModel.Page{}, "en", "12345", areas, filter.Dimension{}, "", "", dataset.DatasetDetails{}, false, false, "", zebedee.EmergencyBanner{})

		Convey("Sorts known items by order then unknown items by TotalCount", func() {
			expectedSelections := []model.Selection{
				{Value: "one", Label: "One", TotalCount: 1},
				{Value: "two", Label: "Two", TotalCount: 2},
				{Value: "three", Label: "Three", TotalCount: 3},
			}

			So(changeDimension.Selections, ShouldResemble, expectedSelections)
		})
	})

	Convey("Given a mixed slice of known and unknown geography areas", t, func() {
		areas := []population.AreaType{
			{ID: "three", Label: "Three", TotalCount: 3},
			{ID: "two", Label: "Two", TotalCount: 2},
			{ID: "nat", Label: "Nation", TotalCount: 3},
			{ID: "ctry", Label: "Country", TotalCount: 2},
			{ID: "rgn", Label: "Region", TotalCount: 1},
			{ID: "one", Label: "One", TotalCount: 1},
		}

		req := httptest.NewRequest("", "/", nil)
		changeDimension := CreateAreaTypeSelector(req, true, coreModel.Page{}, "en", "12345", areas, filter.Dimension{}, "", "", dataset.DatasetDetails{}, false, false, "", zebedee.EmergencyBanner{})

		Convey("Sorts selections ascending by TotalCount", func() {
			expectedSelections := []model.Selection{
				{Value: "nat", Label: "Nation", TotalCount: 3},
				{Value: "ctry", Label: "Country", TotalCount: 2},
				{Value: "rgn", Label: "Region", TotalCount: 1},
				{Value: "one", Label: "One", TotalCount: 1},
				{Value: "two", Label: "Two", TotalCount: 2},
				{Value: "three", Label: "Three", TotalCount: 3},
			}

			So(changeDimension.Selections, ShouldResemble, expectedSelections)
		})
	})

	Convey("Given an unsorted slice of geography areas and lowest_level of geography", t, func() {
		areas := []population.AreaType{
			{ID: "rgn", Label: "Region", TotalCount: 11},
			{ID: "utla", Label: "UTLA", TotalCount: 7},
			{ID: "ctry", Label: "Country", TotalCount: 33},
			{ID: "nat", Label: "Nation", TotalCount: 1},
		}
		lowest_geography := "rgn"

		req := httptest.NewRequest("", "/", nil)
		changeDimension := CreateAreaTypeSelector(req, true, coreModel.Page{}, "en", "12345", areas, filter.Dimension{}, lowest_geography, "", dataset.DatasetDetails{}, false, false, "", zebedee.EmergencyBanner{})

		Convey("Returns the sorted selections stopping at the lowest_level", func() {
			expectedSelections := []model.Selection{
				{Value: "nat", Label: "Nation", TotalCount: 1},
				{Value: "ctry", Label: "Country", TotalCount: 33},
				{Value: "rgn", Label: "Region", TotalCount: 11},
			}

			So(changeDimension.Selections, ShouldResemble, expectedSelections)
		})
	})

	Convey("Given a valid page", t, func() {
		const lang = "en"
		req := httptest.NewRequest("", "/", nil)
		changeDimension := CreateAreaTypeSelector(req, false, coreModel.Page{}, lang, "12345", nil, filter.Dimension{}, "", "", dataset.DatasetDetails{}, false, false, sm, eb)

		Convey("it sets page metadata", func() {
			So(changeDimension.BetaBannerEnabled, ShouldBeTrue)
			So(changeDimension.Type, ShouldEqual, "area_type_options")
			So(changeDimension.Language, ShouldEqual, lang)
		})

		Convey("it sets the title to Area type", func() {
			So(changeDimension.Metadata.Title, ShouldEqual, "Area type")
		})

		Convey("it sets IsAreaType to true", func() {
			So(changeDimension.IsAreaType, ShouldBeTrue)
		})

		Convey("it sets SearchNoIndexEnabled to true", func() {
			So(changeDimension.SearchNoIndexEnabled, ShouldBeTrue)
		})

		Convey("it maps the emergency banner", func() {
			So(changeDimension.EmergencyBanner, ShouldResemble, mappedEmergencyBanner())
		})

		Convey("it maps the service message", func() {
			So(changeDimension.ServiceMessage, ShouldEqual, sm)
		})
	})

	Convey("Given the current filter dimension", t, func() {
		const selectionName = "test"
		req := httptest.NewRequest("", "/", nil)
		changeDimension := CreateAreaTypeSelector(req, false, coreModel.Page{}, "en", "12345", nil, filter.Dimension{ID: selectionName}, "", "", dataset.DatasetDetails{}, false, false, "", zebedee.EmergencyBanner{})

		Convey("it returns the value as an initial selection", func() {
			So(changeDimension.InitialSelection, ShouldEqual, selectionName)
		})
	})

	Convey("Given a validation error", t, func() {
		req := httptest.NewRequest("", "/", nil)
		changeDimension := CreateAreaTypeSelector(req, false, coreModel.Page{}, "en", "12345", nil, filter.Dimension{}, "", "", dataset.DatasetDetails{}, true, false, "", zebedee.EmergencyBanner{})

		Convey("it returns a populated error", func() {
			So(changeDimension.Error.Title, ShouldNotBeEmpty)
		})
	})

	Convey("Given saved options", t, func() {
		req := httptest.NewRequest("", "/", nil)
		changeDimension := CreateAreaTypeSelector(req, false, coreModel.Page{}, "en", "12345", nil, filter.Dimension{}, "", "", dataset.DatasetDetails{}, false, true, "", zebedee.EmergencyBanner{})

		Convey("it returns a warning that saved options will be removed", func() {
			So(changeDimension.HasOptions, ShouldBeTrue)
		})
	})

	Convey("Given analytics metadata", t, func() {
		req := httptest.NewRequest("", "/", nil)
		releaseDate := "2022/11/29"
		dataset := dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"}
		changeDimension := CreateAreaTypeSelector(req, false, coreModel.Page{}, "en", "12345", nil, filter.Dimension{}, "", releaseDate, dataset, false, true, "", zebedee.EmergencyBanner{})

		Convey("it sets DatasetID, DatasetTitle and ReleaseData", func() {
			So(changeDimension.DatasetId, ShouldEqual, dataset.ID)
			So(changeDimension.DatasetTitle, ShouldEqual, dataset.Title)
			So(changeDimension.ReleaseDate, ShouldEqual, releaseDate)
		})
	})
}

func TestGetCoverage(t *testing.T) {
	helper.InitialiseLocalisationsHelper(mocks.MockAssetFunction)
	Convey("Given a valid page", t, func() {
		const lang = "en"
		req := httptest.NewRequest("", "/", nil)
		eb := getTestEmergencyBanner()
		sm := getTestServiceMessage()

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
				"",
				"dim",
				"geogID",
				"2022/11/29",
				sm,
				eb,
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				population.GetAreasResponse{},
				[]model.SelectableElement{},
				population.GetAreaTypeParentsResponse{},
				false,
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
				"",
				"",
				"",
				zebedee.EmergencyBanner{},
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				population.GetAreasResponse{},
				[]model.SelectableElement{},
				parents,
				false,
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
				"",
				"",
				"",
				zebedee.EmergencyBanner{},
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				population.GetAreasResponse{},
				[]model.SelectableElement{},
				population.GetAreaTypeParentsResponse{},
				false,
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
				"",
				"",
				"",
				zebedee.EmergencyBanner{},
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				population.GetAreasResponse{},
				[]model.SelectableElement{},
				parents,
				false,
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
				"",
				"",
				"",
				zebedee.EmergencyBanner{},
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				population.GetAreasResponse{},
				[]model.SelectableElement{},
				parents,
				false,
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
				"",
				"",
				"",
				zebedee.EmergencyBanner{},
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				population.GetAreasResponse{},
				[]model.SelectableElement{},
				population.GetAreaTypeParentsResponse{},
				false,
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
				"name-search",
				"",
				"",
				"",
				"",
				zebedee.EmergencyBanner{},
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				mockedSearchResults,
				[]model.SelectableElement{},
				population.GetAreaTypeParentsResponse{},
				false,
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
				"name-search",
				"",
				"",
				"",
				"",
				zebedee.EmergencyBanner{},
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				mockedSearchResults,
				[]model.SelectableElement{},
				population.GetAreaTypeParentsResponse{},
				false,
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

			coverage := CreateGetCoverage(
				req,
				coreModel.Page{},
				lang,
				"12345",
				"Unknown geography",
				"",
				"search",
				"",
				"parent",
				"parent-search",
				"",
				"",
				"",
				"",
				zebedee.EmergencyBanner{},
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				mockedSearchResults,
				[]model.SelectableElement{},
				population.GetAreaTypeParentsResponse{},
				false,
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
				"name-search",
				"",
				"",
				"",
				"",
				zebedee.EmergencyBanner{},
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				population.GetAreasResponse{},
				[]model.SelectableElement{},
				population.GetAreaTypeParentsResponse{},
				false,
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
				"",
				"",
				"",
				zebedee.EmergencyBanner{},
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				population.GetAreasResponse{},
				[]model.SelectableElement{},
				population.GetAreaTypeParentsResponse{},
				false,
				true,
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
			coverage := CreateGetCoverage(
				req,
				coreModel.Page{},
				lang,
				"12345",
				"Unknown geography",
				"",
				"search",
				"",
				"",
				"parent-search",
				"",
				"",
				"",
				"",
				zebedee.EmergencyBanner{},
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				population.GetAreasResponse{},
				[]model.SelectableElement{},
				population.GetAreaTypeParentsResponse{},
				false,
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
				"",
				"",
				"",
				zebedee.EmergencyBanner{},
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				population.GetAreasResponse{},
				mockedOpt,
				population.GetAreaTypeParentsResponse{},
				false,
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
				"",
				"",
				"",
				zebedee.EmergencyBanner{},
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				population.GetAreasResponse{},
				mockedOpt,
				population.GetAreaTypeParentsResponse{},
				true,
				false,
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
				"",
				"",
				"",
				zebedee.EmergencyBanner{},
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				mockedSearchResults,
				mockedOpt,
				population.GetAreaTypeParentsResponse{},
				false,
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
				"",
				"",
				"",
				zebedee.EmergencyBanner{},
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				mockedSearchResults,
				mockedOpt,
				population.GetAreaTypeParentsResponse{},
				true,
				false,
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
				"name-search",
				"",
				"",
				"",
				"",
				zebedee.EmergencyBanner{},
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				mockedSearchResults,
				mockedOpt,
				population.GetAreaTypeParentsResponse{},
				false,
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
				"parent-search",
				"",
				"",
				"",
				"",
				zebedee.EmergencyBanner{},
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				mockedSearchResults,
				mockedOpt,
				population.GetAreaTypeParentsResponse{},
				true,
				false,
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
				"parent-search",
				"",
				"",
				"",
				"",
				zebedee.EmergencyBanner{},
				dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"},
				mockedSearchResults,
				[]model.SelectableElement{},
				population.GetAreaTypeParentsResponse{},
				true,
				false,
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

func TestGetChangeDimensions(t *testing.T) {
	helper.InitialiseLocalisationsHelper(mocks.MockAssetFunction)
	Convey("Given a valid page request", t, func() {
		const lang = "en"
		req := httptest.NewRequest("", "/", nil)
		eb := getTestEmergencyBanner()
		sm := getTestServiceMessage()

		Convey("When the parameters are valid", func() {
			mockFds := []model.FilterDimension{
				{
					Dimension: filter.Dimension{
						Name:       "dim-1",
						ID:         "dim-1",
						Label:      "dim one",
						IsAreaType: helpers.ToBoolPtr(false),
					},
				},
				{
					Dimension: filter.Dimension{
						Name:       "dim-2",
						ID:         "dim-2",
						Label:      "dim two",
						IsAreaType: helpers.ToBoolPtr(true),
					},
				},
			}
			mockPds := population.GetDimensionsResponse{
				Dimensions: []population.Dimension{
					{
						ID:          "dim-1",
						Label:       "dim one (100 categories)",
						Description: "description one",
					},
					{
						ID:          "dim-a",
						Label:       "dim a (1 category)",
						Description: "description a",
					},
					{
						ID:          "dim-b",
						Label:       "dim b",
						Description: "description b",
					},
					{
						ID:          "dim-c",
						Label:       "dim c",
						Description: "description c",
					},
				},
			}
			mockPdsR := population.GetDimensionsResponse{
				Dimensions: []population.Dimension{
					{
						ID:          "dim-a",
						Label:       "dim a",
						Description: "description a",
					},
				},
			}
			p := CreateGetChangeDimensions(
				req,
				coreModel.Page{},
				lang,
				"12345",
				"dim-a",
				"",
				sm,
				eb,
				mockFds,
				mockPds,
				mockPdsR,
			)
			Convey("Then it maps page metadata", func() {
				So(p.BetaBannerEnabled, ShouldBeTrue)
				So(p.Type, ShouldEqual, "change_variables")
				So(p.Language, ShouldEqual, lang)
				So(p.URI, ShouldEqual, "/")
				So(p.Metadata.Title, ShouldEqual, "Add or remove variables")
			})

			Convey("Then it maps non area-type filter dimensions", func() {
				mockDims := []model.SelectableElement{
					{
						Text:  "dim one",
						Value: "dim-1",
						Name:  "delete-option",
					},
				}
				So(p.Output.Selections, ShouldResemble, mockDims)
				So(p.Output.Selections, ShouldHaveLength, 1)
			})

			Convey("Then it maps available population types dimensions", func() {
				mockPds := []model.SelectableElement{
					{
						Text:       "dim a",
						Value:      "dim-a",
						Name:       "add-dimension",
						IsSelected: false,
						InnerText:  "description a",
					},
					{
						Text:       "dim b",
						Value:      "dim-b",
						Name:       "add-dimension",
						IsSelected: false,
						InnerText:  "description b",
					},
					{
						Text:       "dim c",
						Value:      "dim-c",
						Name:       "add-dimension",
						IsSelected: false,
						InnerText:  "description c",
					},
					{
						Text:       "dim one",
						Value:      "dim-1",
						Name:       "delete-option",
						IsSelected: true,
						InnerText:  "description one",
					},
				}
				So(p.Output.Results, ShouldResemble, mockPds)
				So(p.Output.Results, ShouldHaveLength, 4)
			})

			Convey("Then it maps available dimensions search results", func() {
				mockPds := []model.SelectableElement{
					{
						Text:       "dim a",
						Value:      "dim-a",
						Name:       "add-dimension",
						IsSelected: false,
						InnerText:  "description a",
					},
				}
				So(p.SearchOutput.Results, ShouldResemble, mockPds)
				So(p.SearchOutput.Results, ShouldHaveLength, 1)
				So(p.SearchOutput.HasNoResults, ShouldBeFalse)
			})

			Convey("Then it sets HasNoResults", func() {
				mockPds := []model.SelectableElement{
					{
						Text:       "dim a",
						Value:      "dim-a",
						Name:       "add-dimension",
						IsSelected: false,
						InnerText:  "description a",
					},
				}
				So(p.SearchOutput.Results, ShouldResemble, mockPds)
				So(p.SearchOutput.Results, ShouldHaveLength, 1)
				So(p.SearchOutput.HasNoResults, ShouldBeFalse)
			})

			Convey("Then it maps the emergency banner", func() {
				So(p.EmergencyBanner, ShouldResemble, mappedEmergencyBanner())
			})

			Convey("Then it maps the service message", func() {
				So(p.ServiceMessage, ShouldEqual, sm)
			})
		})

		Convey("when a valid search with no results is performed", func() {
			p := CreateGetChangeDimensions(
				req,
				coreModel.Page{},
				lang,
				"12345",
				"dim-a",
				"search",
				"",
				zebedee.EmergencyBanner{},
				[]model.FilterDimension{},
				population.GetDimensionsResponse{},
				population.GetDimensionsResponse{},
			)
			Convey("then it sets HasNoResults to true", func() {
				So(p.SearchOutput.HasNoResults, ShouldBeTrue)
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

func TestCleanDimensionsLabel(t *testing.T) {
	Convey("Removes categories count from label - case insensitive", t, func() {
		So(cleanDimensionLabel("Example (100 categories)"), ShouldEqual, "Example")
		So(cleanDimensionLabel("Example (7 Categories)"), ShouldEqual, "Example")
		So(cleanDimensionLabel("Example (1 category)"), ShouldEqual, "Example")
		So(cleanDimensionLabel("Example (1 Category)"), ShouldEqual, "Example")
		So(cleanDimensionLabel(""), ShouldEqual, "")
		So(cleanDimensionLabel("Example 1 category"), ShouldEqual, "Example 1 category")
		So(cleanDimensionLabel("Example (something in brackets) (1 Category)"), ShouldEqual, "Example (something in brackets)")
	})
}

func TestSortCategoriesByID(t *testing.T) {

	Convey("Population categories are sorted", t, func() {
		getIDList := func(items []population.Category) []string {
			results := []string{}
			for _, item := range items {
				results = append(results, item.ID)
			}
			return results
		}

		Convey("given non-numeric options", func() {
			nonNumeric := []population.Category{
				{
					ID:    "dim_2",
					Label: "option 2",
				},
				{
					ID:    "dim_1",
					Label: "option 1",
				},
			}
			Convey("when they are sorted", func() {
				sorted := sortCategoriesByID(nonNumeric)

				Convey("then options are sorted alphabetically", func() {
					actual := getIDList(sorted)
					expected := []string{"dim_1", "dim_2"}
					So(actual, ShouldResemble, expected)
				})
			})
		})

		Convey("given simple numeric options", func() {
			numeric := []population.Category{
				{
					ID:    "10",
					Label: "option 10",
				}, {
					ID:    "2",
					Label: "option 2",
				}, {
					ID:    "1",
					Label: "option 1",
				},
			}
			Convey("when they are sorted", func() {
				sorted := sortCategoriesByID(numeric)

				Convey("then options are sorted numerically", func() {
					actual := getIDList(sorted)
					expected := []string{"1", "2", "10"}
					So(actual, ShouldResemble, expected)
				})
			})
		})

		Convey("given numeric options with negatives", func() {
			numericWithNegatives := []population.Category{
				{
					ID:    "10",
					Label: "option 10",
				}, {
					ID:    "2",
					Label: "option 2",
				}, {
					ID:    "-1",
					Label: "option -1",
				},
				{
					ID:    "1",
					Label: "option 1",
				}, {
					ID:    "-10",
					Label: "option -10",
				},
			}

			Convey("when they are sorted", func() {
				sorted := sortCategoriesByID(numericWithNegatives)

				Convey("then options are sorted numerically with negatives at the end", func() {
					actual := getIDList(sorted)
					expected := []string{"1", "2", "10", "-1", "-10"}
					So(actual, ShouldResemble, expected)
				})
			})
		})

		Convey("given mixed numeric and non-numeric options", func() {
			alphanumeric := []population.Category{
				{
					ID:    "10",
					Label: "option 10",
				}, {
					ID:    "2nd Option",
					Label: "option 2",
				}, {
					ID:    "1",
					Label: "option 1",
				},
			}
			Convey("when they are sorted", func() {
				sorted := sortCategoriesByID(alphanumeric)

				Convey("then options are sorted alphanumerically", func() {
					actual := getIDList(sorted)
					expected := []string{"1", "10", "2nd Option"}
					So(actual, ShouldResemble, expected)
				})
			})
		})
	})
}

func getTestEmergencyBanner() zebedee.EmergencyBanner {
	return zebedee.EmergencyBanner{
		Type:        "notable_death",
		Title:       "This is not not an emergency",
		Description: "Something has gone wrong",
		URI:         "google.com",
		LinkText:    "More info",
	}
}

func getTestServiceMessage() string {
	return "Test service message"
}

func mappedEmergencyBanner() coreModel.EmergencyBanner {
	return coreModel.EmergencyBanner{
		Type:        "notable-death",
		Title:       "This is not not an emergency",
		Description: "Something has gone wrong",
		URI:         "google.com",
		LinkText:    "More info",
	}
}
