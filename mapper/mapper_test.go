package mapper

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/dimension"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/model"
	coreModel "github.com/ONSdigital/dp-renderer/model"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitMapper(t *testing.T) {
	mdl := coreModel.Page{}
	req := httptest.NewRequest("", "/", nil)
	lang := "en"
	showAll := []string{}
	filterJob := filter.GetFilterResponse{}
	dims := filter.Dimensions{
		Items: []filter.Dimension{
			{
				Name:       "Dim 1",
				IsAreaType: new(bool),
				Options:    []string{"Opt 1", "Opt 2"},
			},
			{
				Name:       "Truncated dim 1",
				IsAreaType: new(bool),
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
				IsAreaType: new(bool),
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
		},
	}

	Convey("test filter flex overview maps correctly", t, func() {
		m := CreateFilterFlexOverview(req, mdl, lang, "", showAll, filterJob, dims)
		mockEncodedName := url.QueryEscape(strings.ToLower(dims.Items[0].Name))
		So(m.BetaBannerEnabled, ShouldBeTrue)
		So(m.Type, ShouldEqual, "filter-flex-overview")
		So(m.Metadata.Title, ShouldEqual, "Review changes")
		So(m.Language, ShouldEqual, lang)
		So(m.Dimensions[0].Name, ShouldEqual, dims.Items[0].Name)
		So(m.Dimensions[0].IsAreaType, ShouldBeFalse)
		So(m.Dimensions[0].Options, ShouldResemble, dims.Items[0].Options)
		So(m.Dimensions[0].OptionsCount, ShouldEqual, 2)
		So(m.Dimensions[0].EncodedName, ShouldEqual, mockEncodedName)
		So(m.Dimensions[0].URI, ShouldEqual, fmt.Sprintf("%s/%s", "", mockEncodedName))
		So(m.Dimensions[0].IsTruncated, ShouldBeFalse)
	})

	Convey("test truncation maps as expected", t, func() {
		m := CreateFilterFlexOverview(req, mdl, lang, "", showAll, filterJob, dims)
		So(m.Dimensions[1].OptionsCount, ShouldEqual, len(dims.Items[1].Options))
		So(m.Dimensions[1].Options, ShouldHaveLength, 9)
		So(m.Dimensions[1].Options[:3], ShouldResemble, []string{"Opt 1", "Opt 2", "Opt 3"})
		So(m.Dimensions[1].Options[3:6], ShouldResemble, []string{"Opt 9", "Opt 10", "Opt 11"})
		So(m.Dimensions[1].Options[6:], ShouldResemble, []string{"Opt 18", "Opt 19", "Opt 20"})
		So(m.Dimensions[1].IsTruncated, ShouldBeTrue)

		So(m.Dimensions[2].OptionsCount, ShouldEqual, len(dims.Items[2].Options))
		So(m.Dimensions[2].Options, ShouldHaveLength, 9)
		So(m.Dimensions[2].Options[:3], ShouldResemble, []string{"Opt 1", "Opt 2", "Opt 3"})
		So(m.Dimensions[2].Options[3:6], ShouldResemble, []string{"Opt 5", "Opt 6", "Opt 7"})
		So(m.Dimensions[2].Options[6:], ShouldResemble, []string{"Opt 10", "Opt 11", "Opt 12"})
		So(m.Dimensions[2].IsTruncated, ShouldBeTrue)
	})

	Convey("test truncation shows all when parameter given", t, func() {
		m := CreateFilterFlexOverview(req, mdl, lang, "", []string{"Truncated dim 2"}, filterJob, dims)
		So(m.Dimensions[2].OptionsCount, ShouldEqual, len(dims.Items[2].Options))
		So(m.Dimensions[2].Options, ShouldHaveLength, 12)
		So(m.Dimensions[2].IsTruncated, ShouldBeFalse)
	})

	Convey("test create selector maps correctly", t, func() {
		m := CreateSelector(req, mdl, "dimensionName", lang)
		So(m.BetaBannerEnabled, ShouldBeTrue)
		So(m.Type, ShouldEqual, "filter-flex-selector")
		So(m.Metadata.Title, ShouldEqual, "DimensionName")
		So(m.Language, ShouldEqual, lang)
	})
}

func TestCreateAreaTypeSelector(t *testing.T) {
	Convey("Given a slice of geography areas", t, func() {
		areas := []dimension.AreaType{
			{ID: "one", Label: "One", TotalCount: 1},
			{ID: "two", Label: "Two", TotalCount: 2},
		}

		req := httptest.NewRequest("", "/", nil)
		changeDimension := CreateAreaTypeSelector(req, coreModel.Page{}, "en", areas, "", false)

		expectedSelections := []model.Selection{
			{Value: "One", Label: "One", TotalCount: 1},
			{Value: "Two", Label: "Two", TotalCount: 2},
		}

		Convey("Maps each geography dimension into a selection", func() {
			So(changeDimension.Selections, ShouldResemble, expectedSelections)
		})
	})

	Convey("Given a valid page", t, func() {
		const lang = "en"
		req := httptest.NewRequest("", "/", nil)
		changeDimension := CreateAreaTypeSelector(req, coreModel.Page{}, lang, nil, "", false)

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

	Convey("Given a selection name", t, func() {
		const selectionName = "test"
		req := httptest.NewRequest("", "/", nil)
		changeDimension := CreateAreaTypeSelector(req, coreModel.Page{}, "en", nil, selectionName, false)

		Convey("it returns the value as an initial selection", func() {
			So(changeDimension.InitialSelection, ShouldEqual, selectionName)
		})
	})

	Convey("Given a validation error", t, func() {
		req := httptest.NewRequest("", "/", nil)
		changeDimension := CreateAreaTypeSelector(req, coreModel.Page{}, "en", nil, "", true)

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
