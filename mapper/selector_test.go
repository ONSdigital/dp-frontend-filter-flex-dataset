package mapper

import (
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/mocks"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/model"
	"github.com/ONSdigital/dp-renderer/helper"
	coreModel "github.com/ONSdigital/dp-renderer/model"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateCategorisationsSelector(t *testing.T) {
	helper.InitialiseLocalisationsHelper(mocks.MockAssetFunction)
	mdl := coreModel.Page{}
	req := httptest.NewRequest("", "/", nil)
	lang := "en"
	eb := getTestEmergencyBanner()
	sm := getTestServiceMessage()
	m := NewMapper(req, mdl, eb, lang, sm, "12345")
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
			selector := m.CreateCategorisationsSelector("Dimension", "dim1234", cats)

			Convey("Then it maps the page metadata", func() {
				So(selector.BetaBannerEnabled, ShouldBeTrue)
				So(selector.Type, ShouldEqual, "filter-flex-selector")
				So(selector.Metadata.Title, ShouldEqual, "Dimension")
				So(selector.Language, ShouldEqual, lang)
				So(selector.Breadcrumb[0].URI, ShouldEqual, "/filters/12345/dimensions")
				So(selector.Breadcrumb[0].Title, ShouldEqual, "Back")
			})

			Convey("Then it sets the lead text", func() {
				So(selector.LeadText, ShouldEqual, "Select categories")
			})

			Convey("Then it sets SearchNoIndexEnabled to false", func() {
				So(selector.SearchNoIndexEnabled, ShouldBeTrue)
			})

			Convey("Then it sets InitialSelection to dim1234", func() {
				So(selector.InitialSelection, ShouldEqual, "dim1234")
			})

			Convey("Then it maps the service message", func() {
				So(selector.ServiceMessage, ShouldEqual, sm)
			})

			Convey("Then it maps the emergency banner", func() {
				So(selector.EmergencyBanner, ShouldResemble, mappedEmergencyBanner())
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
				So(selector.Selections, ShouldResemble, mockedCats)
			})
		})
		Convey("When a form validation error occurs", func() {
			m.req = httptest.NewRequest("", "/?error=true", nil)
			selector := m.CreateCategorisationsSelector("Dimension", "dim1234", population.GetCategorisationsResponse{})
			Convey("Then it sets the error title", func() {
				So(selector.Error.Title, ShouldEqual, "Dimension")
			})

			Convey("Then it populates the error items struct", func() {
				So(selector.Error.ErrorItems, ShouldHaveLength, 1)
				So(selector.Error.ErrorItems[0].Description.LocaleKey, ShouldEqual, "SelectCategoriesError")
				So(selector.Error.ErrorItems[0].URL, ShouldEqual, "#categories-error")
			})

			Convey("Then it sets the ErrorId", func() {
				So(selector.ErrorId, ShouldEqual, "categories-error")
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
				selector := m.CreateCategorisationsSelector("Dimension", "dim1234", cats)
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
				So(selector.Selections, ShouldResemble, truncCat)
			})

			Convey("Then a showAll request shows all categories as expected", func() {
				m.req = httptest.NewRequest("", "/?showAll=cat_12a", nil)
				selector := m.CreateCategorisationsSelector("Dimension", "dim1234", cats)
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
				So(selector.Selections, ShouldResemble, allCats)
			})
		})
	})
}

func TestCreateAreaTypeSelector(t *testing.T) {
	helper.InitialiseLocalisationsHelper(mocks.MockAssetFunction)
	eb := getTestEmergencyBanner()
	sm := getTestServiceMessage()
	req := httptest.NewRequest("", "/", nil)
	m := NewMapper(req, coreModel.Page{}, eb, "en", sm, "12345")
	Convey("Given a slice of geography areas", t, func() {
		areas := []population.AreaType{
			{ID: "one", Label: "One", Description: "One description", TotalCount: 1},
			{ID: "two", Label: "Two", Description: "Two description", TotalCount: 2},
		}

		changeDimension := m.CreateAreaTypeSelector(areas, filter.Dimension{}, "", "", dataset.DatasetDetails{}, false)

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

		changeDimension := m.CreateAreaTypeSelector(areas, filter.Dimension{}, "", "", dataset.DatasetDetails{}, false)

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
			{ID: "rgn", Label: "Region", TotalCount: 11, Hierarchy_Order: 800},
			{ID: "ctry", Label: "Country", TotalCount: 33, Hierarchy_Order: 900},
			{ID: "nat", Label: "Nation", TotalCount: 1, Hierarchy_Order: 1000},
			{ID: "utla", Label: "UTLA", TotalCount: 7, Hierarchy_Order: 700},
		}

		changeDimension := m.CreateAreaTypeSelector(areas, filter.Dimension{}, "", "", dataset.DatasetDetails{}, false)

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

	Convey("Given an unsorted slice of geography areas and lowest_level of geography", t, func() {
		areas := []population.AreaType{
			{ID: "rgn", Label: "Region", TotalCount: 11, Hierarchy_Order: 800},
			{ID: "ctry", Label: "Country", TotalCount: 33, Hierarchy_Order: 900},
			{ID: "nat", Label: "Nation", TotalCount: 1, Hierarchy_Order: 1000},
			{ID: "utla", Label: "UTLA", TotalCount: 7, Hierarchy_Order: 700},
		}
		lowest_geography := "rgn"

		changeDimension := m.CreateAreaTypeSelector(areas, filter.Dimension{}, lowest_geography, "", dataset.DatasetDetails{}, false)

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
		changeDimension := m.CreateAreaTypeSelector(nil, filter.Dimension{}, "", "", dataset.DatasetDetails{}, false)

		Convey("it sets page metadata", func() {
			So(changeDimension.BetaBannerEnabled, ShouldBeTrue)
			So(changeDimension.Type, ShouldEqual, "area_type_options")
			So(changeDimension.Language, ShouldEqual, "en")
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
		changeDimension := m.CreateAreaTypeSelector(nil, filter.Dimension{ID: selectionName}, "", "", dataset.DatasetDetails{}, false)

		Convey("it returns the value as an initial selection", func() {
			So(changeDimension.InitialSelection, ShouldEqual, selectionName)
		})
	})

	Convey("Given a validation error", t, func() {
		m.req = httptest.NewRequest("", "/?error=true", nil)
		changeDimension := m.CreateAreaTypeSelector(nil, filter.Dimension{}, "", "", dataset.DatasetDetails{}, false)

		Convey("it returns a populated error", func() {
			So(changeDimension.Error.Title, ShouldNotBeEmpty)
		})
	})

	Convey("Given saved options", t, func() {
		changeDimension := m.CreateAreaTypeSelector(nil, filter.Dimension{}, "", "", dataset.DatasetDetails{}, true)

		Convey("it maps a warning that saved options will be removed", func() {
			So(changeDimension.Panel.Body, ShouldEqual, "Saved options warning")
			So(changeDimension.Panel.CssClasses, ShouldResemble, []string{"ons-u-mb-l"})
			So(changeDimension.Panel.Language, ShouldEqual, "en")
		})
	})

	Convey("Given analytics metadata", t, func() {
		releaseDate := "2022/11/29"
		dataset := dataset.DatasetDetails{ID: "dataset-id", Title: "Dataset title"}
		changeDimension := m.CreateAreaTypeSelector(nil, filter.Dimension{}, "", releaseDate, dataset, true)

		Convey("it sets DatasetID, DatasetTitle and ReleaseData", func() {
			So(changeDimension.DatasetId, ShouldEqual, dataset.ID)
			So(changeDimension.DatasetTitle, ShouldEqual, dataset.Title)
			So(changeDimension.ReleaseDate, ShouldEqual, releaseDate)
		})
	})
}
