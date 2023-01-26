package mapper

import (
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/helpers"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/mocks"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/model"
	"github.com/ONSdigital/dp-renderer/helper"
	coreModel "github.com/ONSdigital/dp-renderer/model"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetChangeDimensions(t *testing.T) {
	helper.InitialiseLocalisationsHelper(mocks.MockAssetFunction)
	Convey("Given a valid page request", t, func() {
		const lang = "en"
		req := httptest.NewRequest("", "/", nil)
		eb := getTestEmergencyBanner()
		sm := getTestServiceMessage()
		m := NewMapper(req, coreModel.Page{}, eb, lang, sm, "12345")

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
			p := m.CreateGetChangeDimensions(
				"dim-a",
				"",
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

			Convey("Then it maps the information panel", func() {
				So(p.Panel.Body, ShouldEqual, "Dimensions change warning")
				So(p.Panel.CssClasses, ShouldResemble, []string{"ons-u-mb-s"})
				So(p.Panel.Language, ShouldEqual, lang)
			})
		})

		Convey("when a valid search with no results is performed", func() {
			p := m.CreateGetChangeDimensions(
				"dim-a",
				"search",
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
