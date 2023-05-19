package mapper

import (
	"fmt"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/mocks"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/model"
	"github.com/ONSdigital/dp-renderer/v2/helper"
	coreModel "github.com/ONSdigital/dp-renderer/v2/model"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBuildLinksString(t *testing.T) {
	tc := []struct {
		given    []Link
		expected string
	}{
		{
			given:    []Link{},
			expected: "",
		},
		{
			given: []Link{
				{
					Uri:  "/page/here",
					Text: "a human name",
				},
			},
			expected: "<a href=\"/page/here\">a human name</a>",
		},
		{
			given: []Link{
				{
					Uri:  "/page/here",
					Text: "a human name",
				},
				{
					Uri:  "/page/there",
					Text: "another human name",
				},
			},
			expected: "<a href=\"/page/here\">a human name</a> or <a href=\"/page/there\">another human name</a>",
		},
		{
			given: []Link{
				{
					Uri:  "/page/here",
					Text: "a human name",
				},
				{
					Uri:  "/page/there",
					Text: "another human name",
				},
				{
					Uri:  "/page/everywhere",
					Text: "this human name",
				},
			},
			expected: "<a href=\"/page/here\">a human name</a>, <a href=\"/page/there\">another human name</a> or <a href=\"/page/everywhere\">this human name</a>",
		},
		{
			given: []Link{
				{
					Uri:  "/page/here",
					Text: "a human name",
				},
				{
					Uri:  "/page/there",
					Text: "another human name",
				},
				{
					Uri:  "/page/everywhere",
					Text: "this human name",
				},
				{
					Uri:  "/page/somewhere",
					Text: "human name",
				},
			},
			expected: "<a href=\"/page/here\">a human name</a>, <a href=\"/page/there\">another human name</a>, <a href=\"/page/everywhere\">this human name</a> or <a href=\"/page/somewhere\">human name</a>",
		},
	}

	Convey("Given a link", t, func() {
		Convey("When the buildLinksString function is called", func() {
			for i, test := range tc {
				Convey(fmt.Sprintf("Then the given link (test index %d) returns %s", i, test.expected), func() {
					So(buildLinksString(test.given), ShouldEqual, test.expected)
				})
			}
		})
	})
}

func TestMapImproveResultsCollapsible(t *testing.T) {
	Convey("Given page dimensions", t, func() {
		Convey("When the mapImproveResultsCollapsible function is called", func() {
			mockDims := []model.Dimension{
				{
					Name:        "Test area dim",
					ID:          "test_dim_1",
					URI:         "/test_dim_1",
					IsGeography: true,
				},
				{
					Name:        "Test dim 2",
					ID:          "test_dim_2",
					URI:         "/test_dim_2",
					IsGeography: false,
					HasChange:   true,
				},
			}
			areaUri, dimLink := mapImproveResultsCollapsible(mockDims)

			Convey("Then the area type URI is returned", func() {
				So(areaUri, ShouldEqual, "/test_dim_1")
			})
			Convey("Then the dimension link and text is returned", func() {
				So(dimLink, ShouldEqual, "<a href=\"/test_dim_2\">Test dim 2</a>")
			})
		})
	})

	Convey("Given page dimensions where some do not have change pages", t, func() {
		Convey("When the mapImproveResultsCollapsible function is called", func() {
			mockDims := []model.Dimension{
				{
					Name:        "Test area dim",
					ID:          "test_dim_1",
					URI:         "/test_dim_1",
					IsGeography: true,
				},
				{
					Name:        "Test dim 2",
					ID:          "test_dim_2",
					URI:         "/test_dim_2",
					IsGeography: false,
					HasChange:   false,
				},
				{
					Name:        "Test dim 3",
					ID:          "test_dim_3",
					URI:         "/test_dim_3",
					IsGeography: false,
					HasChange:   true,
				},
			}
			_, dimLink := mapImproveResultsCollapsible(mockDims)

			Convey("Then the dimension link should only include the changeable categories", func() {
				So(dimLink, ShouldEqual, "<a href=\"/test_dim_3\">Test dim 3</a>")
			})
		})
	})
}

func TestMapDescriptionsCollapsible(t *testing.T) {
	helper.InitialiseLocalisationsHelper(mocks.MockAssetFunction)
	Convey("Given a dimension descriptions and page dimensions", t, func() {
		Convey("When the mapDescriptionsCollapsible function is called", func() {
			mockDescriptions := population.GetDimensionsResponse{
				Dimensions: []population.Dimension{
					{
						Description: "A description on one line",
						Label:       "Test area dim",
						ID:          "test_dim_1",
					},
					{
						Description: "A description on one line\nThen a line break",
						Label:       "Test dim 2",
						ID:          "test_dim_2",
					},
					{
						Description: "",
						Label:       "Only a name - I shouldn't map",
						ID:          "test_dim_3",
					},
				},
			}

			mockPageDims := []model.Dimension{
				{
					Name:        "Test area dim",
					ID:          "test_dim_1",
					URI:         "/test_dim_1",
					IsGeography: true,
				},
				{
					Name:        "Test dim 2",
					ID:          "test_dim_2",
					URI:         "/test_dim_2",
					IsGeography: false,
				},
			}
			sut := mapDescriptionsCollapsible(mockDescriptions, mockPageDims)

			Convey("Then the collapsible items are mapped as expected", func() {
				mockedCollapsible := []coreModel.CollapsibleItem{
					{
						Subheading: "Area type",
						Content:    []string(nil),
						SafeHTML: coreModel.Localisation{
							LocaleKey: "VariableInfoAreaType",
							Plural:    1,
						},
					},
					{
						Subheading: "Test area dim",
						Content:    []string{"A description on one line"},
						SafeHTML:   coreModel.Localisation{},
					},
					{
						Subheading: "Coverage",
						Content:    []string(nil),
						SafeHTML: coreModel.Localisation{
							LocaleKey: "VariableInfoCoverage",
							Plural:    1,
						},
					},
					{
						Subheading: "Test dim 2",
						Content:    []string{"A description on one line", "Then a line break"},
						SafeHTML:   coreModel.Localisation{},
					},
				}
				So(sut, ShouldResemble, mockedCollapsible)
			})
		})
	})
}
