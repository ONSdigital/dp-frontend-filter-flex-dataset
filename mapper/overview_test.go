package mapper

import (
	"fmt"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
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
