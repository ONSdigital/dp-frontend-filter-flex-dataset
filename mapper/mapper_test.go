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

func TestUnitMapper(t *testing.T) {
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
				Options:    []string{},
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
		m := CreateFilterFlexOverview(req, mdl, lang, "", showAll, filterJob, filterDims, datasetDims)
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
		So(m.Dimensions[0].Options, ShouldResemble, filterDims.Items[3].Options)
		So(m.Dimensions[0].OptionsCount, ShouldEqual, 0)
		So(m.Dimensions[0].ID, ShouldEqual, filterDims.Items[3].ID)
		So(m.Dimensions[0].URI, ShouldEqual, fmt.Sprintf("%s/%s", "", filterDims.Items[3].Name))
		So(m.Dimensions[0].IsTruncated, ShouldBeFalse)
		So(m.Dimensions[1].Name, ShouldEqual, filterDims.Items[0].Label)
		So(m.Dimensions[1].IsAreaType, ShouldBeFalse)
		So(m.Dimensions[1].Options, ShouldResemble, filterDims.Items[0].Options)
		So(m.Dimensions[1].OptionsCount, ShouldEqual, 2)
		So(m.Dimensions[1].ID, ShouldEqual, filterDims.Items[0].ID)
		So(m.Dimensions[1].URI, ShouldEqual, fmt.Sprintf("%s/%s", "", filterDims.Items[0].Name))
		So(m.Dimensions[1].IsTruncated, ShouldBeFalse)
		So(m.Dimensions[2].Name, ShouldEqual, filterDims.Items[1].Label)
		So(m.Dimensions[2].IsAreaType, ShouldBeFalse)
		So(m.Dimensions[2].ID, ShouldEqual, filterDims.Items[1].ID)
		So(m.Dimensions[2].URI, ShouldEqual, fmt.Sprintf("%s/%s", "", filterDims.Items[1].Name))
		So(m.Dimensions[2].IsTruncated, ShouldBeTrue)
		So(m.Collapsible.CollapsibleItems[0].Subheading, ShouldEqual, datasetDims.Items[0].Label)
		So(m.Collapsible.CollapsibleItems[0].Content[0], ShouldEqual, datasetDims.Items[0].Description)
		So(m.Collapsible.CollapsibleItems[1].Subheading, ShouldEqual, datasetDims.Items[1].Label)
		So(m.Collapsible.CollapsibleItems[1].Content, ShouldResemble, strings.Split(datasetDims.Items[1].Description, "\n"))
		So(m.Collapsible.CollapsibleItems, ShouldHaveLength, 2)
	})

	Convey("test truncation maps as expected", t, func() {
		m := CreateFilterFlexOverview(req, mdl, lang, "", showAll, filterJob, filterDims, datasetDims)
		So(m.Dimensions[2].OptionsCount, ShouldEqual, len(filterDims.Items[1].Options))
		So(m.Dimensions[2].Options, ShouldHaveLength, 9)
		So(m.Dimensions[2].Options[:3], ShouldResemble, []string{"Opt 1", "Opt 2", "Opt 3"})
		So(m.Dimensions[2].Options[3:6], ShouldResemble, []string{"Opt 9", "Opt 10", "Opt 11"})
		So(m.Dimensions[2].Options[6:], ShouldResemble, []string{"Opt 18", "Opt 19", "Opt 20"})
		So(m.Dimensions[2].IsTruncated, ShouldBeTrue)

		So(m.Dimensions[3].OptionsCount, ShouldEqual, len(filterDims.Items[2].Options))
		So(m.Dimensions[3].Options, ShouldHaveLength, 9)
		So(m.Dimensions[3].Options[:3], ShouldResemble, []string{"Opt 1", "Opt 2", "Opt 3"})
		So(m.Dimensions[3].Options[3:6], ShouldResemble, []string{"Opt 5", "Opt 6", "Opt 7"})
		So(m.Dimensions[3].Options[6:], ShouldResemble, []string{"Opt 10", "Opt 11", "Opt 12"})
		So(m.Dimensions[3].IsTruncated, ShouldBeTrue)
	})

	Convey("test truncation shows all when parameter given", t, func() {
		m := CreateFilterFlexOverview(req, mdl, lang, "", []string{"Truncated dim 2"}, filterJob, filterDims, datasetDims)
		So(m.Dimensions[3].OptionsCount, ShouldEqual, len(filterDims.Items[2].Options))
		So(m.Dimensions[3].Options, ShouldHaveLength, 12)
		So(m.Dimensions[3].IsTruncated, ShouldBeFalse)
	})

	Convey("test create selector maps correctly", t, func() {
		m := CreateSelector(req, mdl, "dimensionName", lang, "12345")
		So(m.BetaBannerEnabled, ShouldBeTrue)
		So(m.Type, ShouldEqual, "filter-flex-selector")
		So(m.Metadata.Title, ShouldEqual, "DimensionName")
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

		Convey("it sets the title to Area Type", func() {
			So(changeDimension.Metadata.Title, ShouldEqual, "Area Type")
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
