package mapper

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-renderer/model"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitMapper(t *testing.T) {
	mdl := model.Page{}
	req := httptest.NewRequest("", "/", nil)
	lang := "en"

	Convey("test filter flex overview maps correctly", t, func() {
		m := CreateFilterFlexOverview(req, mdl, lang)
		So(m.BetaBannerEnabled, ShouldBeTrue)
		So(m.Type, ShouldEqual, "filter-flex-overview")
		So(m.Metadata.Title, ShouldEqual, "Review changes")
		So(m.Language, ShouldEqual, lang)
	})

	Convey("test create selector maps correctly", t, func() {
		m := CreateSelector(req, mdl, "dimensionName", lang)
		So(m.BetaBannerEnabled, ShouldBeTrue)
		So(m.Type, ShouldEqual, "filter-flex-selector")
		So(m.Metadata.Title, ShouldEqual, "DimensionName")
		So(m.Language, ShouldEqual, lang)
	})
}

func TestUnitMapCookiesPreferences(t *testing.T) {
	req := httptest.NewRequest("", "/", nil)
	pageModel := model.Page{
		CookiesPreferencesSet: false,
		CookiesPolicy: model.CookiesPolicy{
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
