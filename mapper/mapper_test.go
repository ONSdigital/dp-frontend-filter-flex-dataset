package mapper

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	coreModel "github.com/ONSdigital/dp-renderer/model"
	. "github.com/smartystreets/goconvey/convey"
)

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

func TestMaxVariablesError(t *testing.T) {
	Convey("Returns true if the sdc error string contains Maximum variables", t, func() {
		So(isMaxVariablesError(&cantabular.GetBlockedAreaCountResult{
			TableError: "Maximum variables at the start",
		}), ShouldBeTrue)
		So(isMaxVariablesError(&cantabular.GetBlockedAreaCountResult{
			TableError: "Not at the start but still contains Maximum variables",
		}), ShouldBeTrue)

		So(isMaxVariablesError(&cantabular.GetBlockedAreaCountResult{
			TableError: "maximum variables in lower case",
		}), ShouldBeFalse)
		So(isMaxVariablesError(&cantabular.GetBlockedAreaCountResult{
			TableError: "doesn't contain string at all",
		}), ShouldBeFalse)
	})
}

func TestMaxCellsError(t *testing.T) {
	Convey("Returns true if the sdc error string contains withinMaxCells", t, func() {
		So(isMaxCellsError(&cantabular.GetBlockedAreaCountResult{
			TableError: "withinMaxCells",
		}), ShouldBeTrue)
		So(isMaxCellsError(&cantabular.GetBlockedAreaCountResult{
			TableError: "withinmaxcells is case sensitive",
		}), ShouldBeFalse)
		So(isMaxCellsError(&cantabular.GetBlockedAreaCountResult{
			TableError: "doesn't contain string at all",
		}), ShouldBeFalse)
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
