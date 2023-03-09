package mapper

import (
	"fmt"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
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
	req := httptest.NewRequest("", "/dimensions", nil)
	lang := "en"
	eb := getTestEmergencyBanner()
	sm := getTestServiceMessage()
	m := NewMapper(req, mdl, eb, lang, sm, "12345")
	filterJob := filter.GetFilterResponse{
		Dataset: filter.Dataset{
			DatasetID: "example",
			Edition:   "2021",
			Version:   1,
		},
		PopulationType: "UR",
		Custom:         helpers.ToBoolPtr(false),
	}
	filterDims := []model.FilterDimension{
		{
			Dimension: filter.Dimension{
				Name:       "Dim 1",
				IsAreaType: helpers.ToBoolPtr(false),
				Options:    []string{"Opt 1", "Opt 2"},
			},
			OptionsCount:        2,
			CategorisationCount: 1,
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
			OptionsCount:        20,
			CategorisationCount: 2,
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
			OptionsCount:        12,
			CategorisationCount: 2,
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
			OptionsCount:        10,
			CategorisationCount: 2,
		},
	}
	dimDescriptions := population.GetDimensionsResponse{
		Dimensions: []population.Dimension{
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
	sdc := cantabular.GetBlockedAreaCountResult{
		Passed:  100,
		Blocked: 0,
		Total:   100,
	}
	pops := population.GetPopulationTypesResponse{
		Items: []population.PopulationType{
			{
				Name:        "UR",
				Label:       "Usual residents",
				Description: "The description of usual residents",
			},
			{
				Name:        "AP",
				Label:       "Another population type",
				Description: "The description of another population type",
			},
		},
	}

	Convey("test filter flex overview maps correctly", t, func() {
		overview := m.CreateFilterFlexOverview(filterJob, filterDims, dimDescriptions, pops, sdc, false)
		So(overview.BetaBannerEnabled, ShouldBeTrue)
		So(overview.Type, ShouldEqual, "review_changes")
		So(overview.Metadata.Title, ShouldEqual, "Review changes")
		So(overview.Breadcrumb[0].Title, ShouldEqual, "Back")
		So(overview.Breadcrumb[0].URI, ShouldEqual, fmt.Sprintf("/datasets/%s/editions/%s/versions/%s",
			filterJob.Dataset.DatasetID,
			filterJob.Dataset.Edition,
			strconv.Itoa(filterJob.Dataset.Version)))
		So(overview.Language, ShouldEqual, lang)
		So(overview.SearchNoIndexEnabled, ShouldBeTrue)
		So(overview.IsMultivariate, ShouldBeFalse)

		So(overview.Dimensions[0].Name, ShouldEqual, "Population type")
		So(overview.Dimensions[0].Options[0], ShouldEqual, pops.Items[0].Label)
		So(overview.Dimensions[0].ID, ShouldEqual, pops.Items[0].Name)
		So(overview.Dimensions[0].IsGeography, ShouldBeTrue)

		So(overview.Dimensions[1].Name, ShouldEqual, "Area type")
		So(overview.Dimensions[1].IsGeography, ShouldBeTrue)
		So(overview.Dimensions[1].OptionsCount, ShouldEqual, filterDims[3].OptionsCount)
		So(overview.Dimensions[1].ID, ShouldEqual, filterDims[3].ID)
		So(overview.Dimensions[1].URI, ShouldEqual, fmt.Sprintf("/dimensions/%s", filterDims[3].Name))
		So(overview.Dimensions[1].IsTruncated, ShouldBeFalse)

		So(overview.Dimensions[2].Name, ShouldEqual, "Coverage")
		So(overview.Dimensions[2].IsGeography, ShouldBeTrue)
		So(overview.Dimensions[2].Options, ShouldResemble, filterDims[3].Options)
		So(overview.Dimensions[2].URI, ShouldEqual, "/dimensions/geography/coverage")
		So(overview.Dimensions[2].IsTruncated, ShouldBeFalse)

		So(overview.Dimensions[3].Name, ShouldEqual, filterDims[0].Label)
		So(overview.Dimensions[3].IsGeography, ShouldBeFalse)
		So(overview.Dimensions[3].IsGeography, ShouldBeFalse)
		So(overview.Dimensions[3].ID, ShouldEqual, filterDims[0].ID)
		So(overview.Dimensions[3].URI, ShouldEqual, fmt.Sprintf("/dimensions/%s", filterDims[0].Name))
		So(overview.Dimensions[3].IsTruncated, ShouldBeFalse)

		So(overview.Dimensions[4].Name, ShouldEqual, filterDims[1].Label)
		So(overview.Dimensions[4].IsGeography, ShouldBeFalse)
		So(overview.Dimensions[4].IsGeography, ShouldBeFalse)
		So(overview.Dimensions[4].ID, ShouldEqual, filterDims[1].ID)
		So(overview.Dimensions[4].URI, ShouldEqual, fmt.Sprintf("/dimensions/%s", filterDims[1].Name))
		So(overview.Dimensions[4].IsTruncated, ShouldBeTrue)

		So(overview.EmergencyBanner, ShouldResemble, mappedEmergencyBanner())
		So(overview.ServiceMessage, ShouldEqual, sm)

		So(overview.HasSDC, ShouldBeFalse)
	})

	Convey("test filter flex overview maps for custom datasets correctly", t, func() {
		customFilterJob := filter.GetFilterResponse{
			Dataset: filter.Dataset{
				DatasetID: "example",
				Edition:   "2021",
				Version:   1,
			},
			PopulationType: "UR",
			Custom:         helpers.ToBoolPtr(true),
		}
		overview := m.CreateFilterFlexOverview(customFilterJob, filterDims, dimDescriptions, pops, sdc, false)
		So(overview.Metadata.Title, ShouldEqual, "Custom dataset")
	})

	Convey("test truncation maps as expected", t, func() {
		overview := m.CreateFilterFlexOverview(filterJob, filterDims, dimDescriptions, pops, sdc, false)
		So(overview.Dimensions[4].OptionsCount, ShouldEqual, filterDims[1].OptionsCount)
		So(overview.Dimensions[4].Options, ShouldHaveLength, 9)
		So(overview.Dimensions[4].Options[:3], ShouldResemble, []string{"Opt 1", "Opt 2", "Opt 3"})
		So(overview.Dimensions[4].Options[3:6], ShouldResemble, []string{"Opt 9", "Opt 10", "Opt 11"})
		So(overview.Dimensions[4].Options[6:], ShouldResemble, []string{"Opt 18", "Opt 19", "Opt 20"})
		So(overview.Dimensions[4].IsTruncated, ShouldBeTrue)

		So(overview.Dimensions[5].OptionsCount, ShouldEqual, filterDims[2].OptionsCount)
		So(overview.Dimensions[5].Options, ShouldHaveLength, 9)
		So(overview.Dimensions[5].Options[:3], ShouldResemble, []string{"Opt 1", "Opt 2", "Opt 3"})
		So(overview.Dimensions[5].Options[3:6], ShouldResemble, []string{"Opt 5", "Opt 6", "Opt 7"})
		So(overview.Dimensions[5].Options[6:], ShouldResemble, []string{"Opt 10", "Opt 11", "Opt 12"})
		So(overview.Dimensions[5].IsTruncated, ShouldBeTrue)
	})

	Convey("test truncation shows all when parameter given", t, func() {
		m.req = httptest.NewRequest("", "/?showAll=Truncated+dim+2", nil)
		overview := m.CreateFilterFlexOverview(filterJob, filterDims, dimDescriptions, pops, sdc, false)
		So(overview.Dimensions[5].OptionsCount, ShouldEqual, filterDims[2].OptionsCount)
		So(overview.Dimensions[5].Options, ShouldHaveLength, 12)
		So(overview.Dimensions[5].IsTruncated, ShouldBeFalse)
	})

	Convey("test area type dimension options do not truncate and map to 'coverage' dimension", t, func() {
		overview := m.CreateFilterFlexOverview(filterJob, filterDims, dimDescriptions, pops, sdc, false)
		So(overview.Dimensions[2].Options, ShouldHaveLength, 10)
		So(overview.Dimensions[2].IsTruncated, ShouldBeFalse)
		So(overview.Dimensions[2].IsGeography, ShouldBeTrue)
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

		overview := m.CreateFilterFlexOverview(filterJob, newFilterDims, dimDescriptions, pops, sdc, false)
		So(overview.Dimensions[6].Name, ShouldEqual, "Example")
	})

	Convey("Given area type selection", t, func() {
		Convey("When area types are selected", func() {
			overview := m.CreateFilterFlexOverview(filterJob, filterDims, dimDescriptions, pops, sdc, false)
			Convey("Then area selection is displayed", func() {
				So(overview.Dimensions[2].Options, ShouldResemble, filterDims[3].Options)
			})
		})
		Convey("When no area types are selected", func() {
			filterDims[3].Options = []string{}
			overview := m.CreateFilterFlexOverview(filterJob, filterDims, dimDescriptions, pops, sdc, false)
			Convey("Then the default coverage is displayed", func() {
				So(overview.Dimensions[2].Options[0], ShouldResemble, "England and Wales")
			})
		})
	})

	Convey("Given a filter based on a multivariate dataset", t, func() {
		Convey("When there are blocked areas greater than zero", func() {
			sdc.Blocked = 10
			sdc.Passed = 15
			sdc.Total = 25
			overview := m.CreateFilterFlexOverview(filterJob, filterDims, dimDescriptions, pops, sdc, true)
			Convey("Then the bool isMultivariate is true", func() {
				So(overview.IsMultivariate, ShouldBeTrue)
			})
			Convey("Then the sdc bool is true", func() {
				So(overview.HasSDC, ShouldBeTrue)
			})
			Convey("Then the sdc panel is displayed", func() {
				mockPanel := model.Panel{
					Type:       model.Pending,
					CssClasses: []string{"ons-u-mb-s"},
					SafeHTML:   []string{"15 of 25 areas are available", "Protecting personal data will prevent 10 areas from being published"},
					Language:   lang,
				}
				So(overview.Panel, ShouldResemble, mockPanel)
			})
			Convey("Then the 'how to improve your results collapsible' is populated", func() {
				So(overview.ImproveResults.CollapsibleItems, ShouldHaveLength, 1)
			})
		})
		Convey("When all areas are available", func() {
			sdc.Blocked = 0
			sdc.Passed = 25
			sdc.Total = 25
			overview := m.CreateFilterFlexOverview(filterJob, filterDims, dimDescriptions, pops, sdc, true)
			Convey("Then the bool isMultivariate is true", func() {
				So(overview.IsMultivariate, ShouldBeTrue)
			})
			Convey("Then the sdc bool is true", func() {
				So(overview.HasSDC, ShouldBeTrue)
			})
			Convey("Then the sdc panel is displayed", func() {
				mockPanel := model.Panel{
					Type:       model.Success,
					CssClasses: []string{"ons-u-mb-l"},
					SafeHTML:   []string{"All areas available"},
					Language:   lang,
				}
				So(overview.Panel, ShouldResemble, mockPanel)
			})
			Convey("Then the 'how to improve your results' collapsible is empty", func() {
				So(overview.ImproveResults.CollapsibleItems, ShouldHaveLength, 0)
			})
		})
	})

	Convey("test IsChangeVisible parameter", t, func() {
		Convey("when isMultivariate is false", func() {
			overview := m.CreateFilterFlexOverview(filterJob, filterDims, dimDescriptions, pops, sdc, false)
			Convey("then IsChangeCategories is false for all", func() {
				So(overview.Dimensions[3].HasChange, ShouldBeFalse)
				So(overview.Dimensions[4].HasChange, ShouldBeFalse)
				So(overview.Dimensions[5].HasChange, ShouldBeFalse)
			})
		})

		Convey("when isMultivariate is true", func() {
			overview := m.CreateFilterFlexOverview(filterJob, filterDims, dimDescriptions, pops, sdc, true)
			Convey("then IsChangeCategories is false if categorisation is only one available", func() {
				So(overview.Dimensions[3].HasChange, ShouldBeFalse)
				So(overview.Dimensions[4].HasChange, ShouldBeTrue)
				So(overview.Dimensions[5].HasChange, ShouldBeTrue)
			})
		})
	})

	Convey("test EnableGetData boolean", t, func() {
		Convey("when isMultivariate is false", func() {
			overview := m.CreateFilterFlexOverview(filterJob, filterDims, dimDescriptions, pops, sdc, false)
			Convey("then EnableGetData should be true", func() {
				So(overview.EnableGetData, ShouldBeTrue)
			})
		})

		Convey("when isMultivariate is true and one or more dimensions are added", func() {
			overview := m.CreateFilterFlexOverview(filterJob, filterDims, dimDescriptions, pops, sdc, true)
			Convey("then EnableGetData should be true", func() {
				So(overview.EnableGetData, ShouldBeTrue)
			})
		})

		Convey("when isMultivariate is true and only the area type dimension is added", func() {
			filterDims = []model.FilterDimension{
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
					OptionsCount:        10,
					CategorisationCount: 2,
				},
			}
			overview := m.CreateFilterFlexOverview(filterJob, filterDims, dimDescriptions, pops, sdc, true)
			Convey("then EnableGetData should be false", func() {
				So(overview.EnableGetData, ShouldBeFalse)
			})
		})
	})
}
