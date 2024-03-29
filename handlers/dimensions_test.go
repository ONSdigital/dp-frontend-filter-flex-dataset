package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/helpers"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/mocks"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/model"
	"github.com/ONSdigital/dp-renderer/v2/helper"
	coreModel "github.com/ONSdigital/dp-renderer/v2/model"
	gomock "github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDimensionsHandler(t *testing.T) {
	helper.InitialiseLocalisationsHelper(mocks.MockAssetFunction)
	mockCtrl := gomock.NewController(t)
	cfg := initialiseMockConfig()
	mockDataset := dataset.DatasetDetails{
		ID:    "Mock-Dataset-ID",
		Title: "Mock dataset title",
	}
	mockVersion1 := dataset.Version{
		ID:          "1",
		ReleaseDate: "2022/11/29",
	}

	Convey("Dimensions Selector", t, func() {
		Convey("Given a valid dimension param for a filter", func() {
			Convey("When the filter is not multivariate and the dimension is not an area type ", func() {
				const dimensionName = "Number Of Siblings"

				mockFilter := NewMockFilterClient(mockCtrl)
				mockFilter.
					EXPECT().
					GetJobState(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Model{}, "", nil)
				mockFilter.
					EXPECT().
					GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimension{Name: dimensionName}, "", nil).
					AnyTimes()

				mockRend := NewMockRenderClient(mockCtrl)
				mockRend.
					EXPECT().
					NewBasePageModel().
					Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain)).
					AnyTimes()

				mockDc := NewMockDatasetClient(mockCtrl)
				mockDc.EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(mockDataset, nil).AnyTimes()
				mockDc.EXPECT().
					GetVersion(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(mockVersion1, nil).AnyTimes()

				mockZc := NewMockZebedeeClient(mockCtrl)
				mockZc.
					EXPECT().
					GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(zebedee.HomepageContent{}, nil)

				ff := NewFilterFlex(mockRend, mockFilter, mockDc, NewMockPopulationClient(mockCtrl), mockZc, cfg)
				w := runDimensionsSelector(
					"number+of+siblings",
					ff.DimensionSelector(),
				)

				Convey("Then status code should be 400", func() {
					So(w.Code, ShouldEqual, http.StatusBadRequest)
				})
			})

			Convey("When the filter is multivariate and the dimension is not an area type ", func() {
				const dimensionName = "Number Of Siblings"
				mockDataset.Type = "multivariate"
				cfg.EnableMultivariate = true

				mockFilter := NewMockFilterClient(mockCtrl)
				mockFilter.
					EXPECT().
					GetJobState(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Model{}, "", nil)
				mockFilter.
					EXPECT().
					GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimension{Name: dimensionName}, "", nil).
					AnyTimes()

				mockRend := NewMockRenderClient(mockCtrl)
				mockRend.
					EXPECT().
					NewBasePageModel().
					Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain)).
					AnyTimes()
				mockRend.
					EXPECT().
					BuildPage(gomock.Any(), gomock.Any(), "selector")

				mockDc := NewMockDatasetClient(mockCtrl)
				mockDc.EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(mockDataset, nil).AnyTimes()
				mockDc.EXPECT().
					GetVersion(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(mockVersion1, nil).AnyTimes()

				mockZc := NewMockZebedeeClient(mockCtrl)
				mockZc.
					EXPECT().
					GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(zebedee.HomepageContent{}, nil)

				mockPc := NewMockPopulationClient(mockCtrl)
				mockPc.
					EXPECT().
					GetCategorisations(gomock.Any(), gomock.Any()).
					Return(population.GetCategorisationsResponse{}, nil)

				ff := NewFilterFlex(mockRend, mockFilter, mockDc, mockPc, mockZc, cfg)
				w := runDimensionsSelector(
					"number+of+siblings",
					ff.DimensionSelector(),
				)

				Convey("Then status code should be 200", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
				})
			})
		})

		Convey("Given a dimension param which is missing from a filter", func() {
			mockFilter := NewMockFilterClient(mockCtrl)
			mockFilter.
				EXPECT().
				GetJobState(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(filter.Model{}, "", nil) // No filter dimensions
			mockFilter.
				EXPECT().
				GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(filter.Dimension{}, "", &filter.ErrInvalidFilterAPIResponse{ExpectedCode: http.StatusOK, ActualCode: http.StatusNotFound}).
				AnyTimes()

			mockDc := NewMockDatasetClient(mockCtrl)
			mockDc.EXPECT().
				Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(mockDataset, nil).AnyTimes()
			mockDc.EXPECT().
				GetVersion(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(mockVersion1, nil).AnyTimes()

			mockZc := NewMockZebedeeClient(mockCtrl)
			mockZc.
				EXPECT().
				GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(zebedee.HomepageContent{}, nil)

			ff := NewFilterFlex(NewMockRenderClient(mockCtrl), mockFilter, mockDc, NewMockPopulationClient(mockCtrl), mockZc, cfg)
			w := runDimensionsSelector(
				"city",
				ff.DimensionSelector(),
			)

			Convey("Then the status code should be 404", func() {
				So(w.Code, ShouldEqual, http.StatusNotFound)
			})
		})

		Convey("Given an area type", func() {
			const dimensionName = "city"

			stubAreaTypeDimension := filter.Dimension{
				Name:       dimensionName,
				IsAreaType: helpers.ToBoolPtr(true),
			}

			Convey("When area types are returned", func() {
				Convey("Then the page should contain a sorted list of area type selections", func() {

					unsortedAreaTypes := []population.AreaType{
						{
							ID:         "ladcd",
							Label:      "Local authority code",
							TotalCount: 100,
						},
						{
							ID:         "country",
							Label:      "Country",
							TotalCount: 1,
						},
						{
							ID:         "region",
							Label:      "Region",
							TotalCount: 10,
						}}

					mockFilter := NewMockFilterClient(mockCtrl)
					mockFilter.EXPECT().
						GetJobState(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(filter.Model{}, "", nil).
						AnyTimes()
					mockFilter.
						EXPECT().
						GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(stubAreaTypeDimension, "", nil).
						AnyTimes()
					mockFilter.
						EXPECT().
						GetDimensionOptions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(filter.DimensionOptions{}, "", nil).
						AnyTimes()

					mockPc := NewMockPopulationClient(mockCtrl)
					mockPc.EXPECT().
						GetAreaTypes(gomock.Any(), gomock.Any()).
						Return(
							population.GetAreaTypesResponse{
								AreaTypes: unsortedAreaTypes,
							},
							nil,
						).
						AnyTimes()

					mockRend := NewMockRenderClient(mockCtrl)
					mockRend.EXPECT().
						NewBasePageModel().
						Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain)).
						AnyTimes()

					// Validate page data contains sorted selections
					sortedByCountAscending := []model.Selection{
						{
							Value:      "country",
							Label:      "Country",
							TotalCount: 1,
						},
						{
							Value:      "region",
							Label:      "Region",
							TotalCount: 10,
						},
						{
							Value:      "ladcd",
							Label:      "Local authority code",
							TotalCount: 100,
						}}

					mockRend.EXPECT().
						BuildPage(
							gomock.Any(),
							pageMatchesSelections{
								selections: sortedByCountAscending,
							},
							"selector",
						)

					mockDc := NewMockDatasetClient(mockCtrl)
					mockDc.EXPECT().
						Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(mockDataset, nil).AnyTimes()
					mockDc.
						EXPECT().
						GetVersion(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(dataset.Version{
							LowestGeography: "",
							Version:         1,
						}, nil).
						AnyTimes()

					mockZc := NewMockZebedeeClient(mockCtrl)
					mockZc.
						EXPECT().
						GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(zebedee.HomepageContent{}, nil)

					ff := NewFilterFlex(mockRend, mockFilter, mockDc, mockPc, mockZc, cfg)
					w := runDimensionsSelector(dimensionName, ff.DimensionSelector())

					Convey("And the status code should be 200", func() {
						So(w.Code, ShouldEqual, http.StatusOK)
					})
				})

				Convey("Then the page should limit selections by lowest geography", func() {

					areaTypes := []population.AreaType{
						{
							ID:         "country",
							Label:      "Country",
							TotalCount: 1,
						},
						{
							ID:         "ladcd",
							Label:      "Local authority code",
							TotalCount: 100,
						},
						{
							ID:         "region",
							Label:      "Region",
							TotalCount: 10,
						},
					}

					lowest_geography := "region"

					filteredSelections := []model.Selection{
						{
							Value:      "country",
							Label:      "Country",
							TotalCount: 1,
						},
						{
							Value:      "region",
							Label:      "Region",
							TotalCount: 10,
						}}

					mockFilter := NewMockFilterClient(mockCtrl)
					mockFilter.EXPECT().
						GetJobState(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(filter.Model{}, "", nil).
						AnyTimes()
					mockFilter.
						EXPECT().
						GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(stubAreaTypeDimension, "", nil).
						AnyTimes()
					mockFilter.
						EXPECT().
						GetDimensionOptions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(filter.DimensionOptions{}, "", nil).
						AnyTimes()

					mockPc := NewMockPopulationClient(mockCtrl)
					mockPc.EXPECT().
						GetAreaTypes(gomock.Any(), gomock.Any()).
						Return(
							population.GetAreaTypesResponse{
								AreaTypes: areaTypes,
							},
							nil,
						).
						AnyTimes()

					mockRend := NewMockRenderClient(mockCtrl)
					mockRend.EXPECT().
						NewBasePageModel().
						Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain)).
						AnyTimes()

					mockRend.EXPECT().
						BuildPage(
							gomock.Any(),
							pageMatchesSelections{
								selections: filteredSelections,
							},
							"selector",
						)

					mockDc := NewMockDatasetClient(mockCtrl)
					mockDc.EXPECT().
						Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(mockDataset, nil).AnyTimes()
					mockDc.
						EXPECT().
						GetVersion(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(dataset.Version{
							LowestGeography: lowest_geography,
							Version:         1,
						}, nil).
						AnyTimes()

					mockZc := NewMockZebedeeClient(mockCtrl)
					mockZc.
						EXPECT().
						GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(zebedee.HomepageContent{}, nil)

					ff := NewFilterFlex(mockRend, mockFilter, mockDc, mockPc, mockZc, cfg)
					w := runDimensionsSelector(dimensionName, ff.DimensionSelector())

					Convey("And the status code should be 200", func() {
						So(w.Code, ShouldEqual, http.StatusOK)
					})
				})

				Convey("Then the dimensions API client should request area types using the cantabular ID", func() {
					const cantabularID = "cantabular"

					mockFilter := NewMockFilterClient(mockCtrl)
					mockFilter.EXPECT().
						GetJobState(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(filter.Model{PopulationType: cantabularID}, "", nil).
						AnyTimes()
					mockFilter.
						EXPECT().
						GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(stubAreaTypeDimension, "", nil).
						AnyTimes()
					mockFilter.
						EXPECT().
						GetDimensionOptions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(filter.DimensionOptions{}, "", nil).
						AnyTimes()

					mockPc := NewMockPopulationClient(mockCtrl)
					mockPc.EXPECT().
						GetAreaTypes(gomock.Any(), gomock.Any()).
						Return(population.GetAreaTypesResponse{}, nil)

					mockRend := NewMockRenderClient(mockCtrl)
					mockRend.EXPECT().
						NewBasePageModel().
						Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain)).
						AnyTimes()

					mockRend.
						EXPECT().
						BuildPage(gomock.Any(), gomock.Any(), "selector").
						AnyTimes()

					mockDc := NewMockDatasetClient(mockCtrl)
					mockDc.EXPECT().
						Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(mockDataset, nil).AnyTimes()
					mockDc.
						EXPECT().
						GetVersion(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(dataset.Version{
							LowestGeography: "",
							Version:         1,
						}, nil).
						AnyTimes()

					mockZc := NewMockZebedeeClient(mockCtrl)
					mockZc.
						EXPECT().
						GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(zebedee.HomepageContent{}, nil)

					ff := NewFilterFlex(mockRend, mockFilter, mockDc, mockPc, mockZc, cfg)
					w := runDimensionsSelector(dimensionName, ff.DimensionSelector())

					Convey("And the status code should be 200", func() {
						So(w.Code, ShouldEqual, http.StatusOK)
					})
				})

				Convey("Then the page should have the area type bool set to true", func() {
					mockFilter := NewMockFilterClient(mockCtrl)
					mockFilter.EXPECT().
						GetJobState(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(filter.Model{}, "", nil).
						AnyTimes()
					mockFilter.
						EXPECT().
						GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(stubAreaTypeDimension, "", nil).
						AnyTimes()
					mockFilter.
						EXPECT().
						GetDimensionOptions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(filter.DimensionOptions{}, "", nil).
						AnyTimes()

					mockPc := NewMockPopulationClient(mockCtrl)
					mockPc.EXPECT().
						GetAreaTypes(gomock.Any(), gomock.Any()).
						Return(population.GetAreaTypesResponse{}, nil)

					mockRend := NewMockRenderClient(mockCtrl)
					mockRend.EXPECT().
						NewBasePageModel().
						Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain)).
						AnyTimes()

					mockRend.
						EXPECT().
						// Assert that the area type boolean is true
						BuildPage(gomock.Any(), pageIsAreaType{true}, gomock.Any())

					mockDc := NewMockDatasetClient(mockCtrl)
					mockDc.EXPECT().
						Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(mockDataset, nil).AnyTimes()
					mockDc.
						EXPECT().
						GetVersion(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(dataset.Version{
							LowestGeography: "",
							Version:         1,
						}, nil).
						AnyTimes()

					mockZc := NewMockZebedeeClient(mockCtrl)
					mockZc.
						EXPECT().
						GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(zebedee.HomepageContent{}, nil)

					ff := NewFilterFlex(mockRend, mockFilter, mockDc, mockPc, mockZc, cfg)
					w := runDimensionsSelector(dimensionName, ff.DimensionSelector())

					Convey("And the status code should be 200", func() {
						So(w.Code, ShouldEqual, http.StatusOK)
					})
				})
			})

			Convey("Given a truthy error query param", func() {
				req := httptest.NewRequest(http.MethodGet, "/filters/1234/dimensions/city?error=true", nil)

				Convey("Then the page should contain a populated error", func() {
					mockFilter := NewMockFilterClient(mockCtrl)
					mockFilter.
						EXPECT().
						GetJobState(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(filter.Model{}, "", nil)
					mockFilter.
						EXPECT().
						GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(filter.Dimension{IsAreaType: helpers.ToBoolPtr(true)}, "", nil).
						AnyTimes()
					mockFilter.
						EXPECT().
						GetDimensionOptions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(filter.DimensionOptions{}, "", nil).
						AnyTimes()

					mockPc := NewMockPopulationClient(mockCtrl)
					mockPc.EXPECT().
						GetAreaTypes(gomock.Any(), gomock.Any()).
						Return(population.GetAreaTypesResponse{}, nil).
						AnyTimes()

					mockRend := NewMockRenderClient(mockCtrl)
					mockRend.
						EXPECT().
						NewBasePageModel().
						Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain)).
						AnyTimes()

					mockRend.
						EXPECT().
						// Confirm the page contains an error
						BuildPage(gomock.Any(), pageHasError{true}, gomock.Any())

					mockDc := NewMockDatasetClient(mockCtrl)
					mockDc.EXPECT().
						Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(mockDataset, nil).AnyTimes()
					mockDc.
						EXPECT().
						GetVersion(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(dataset.Version{
							LowestGeography: "",
							Version:         1,
						}, nil).
						AnyTimes()

					mockZc := NewMockZebedeeClient(mockCtrl)
					mockZc.
						EXPECT().
						GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(zebedee.HomepageContent{}, nil)

					ff := NewFilterFlex(mockRend, mockFilter, mockDc, mockPc, mockZc, cfg)
					selector := ff.DimensionSelector()

					w := httptest.NewRecorder()
					router := mux.NewRouter()
					router.HandleFunc("/filters/{filterID}/dimensions/{name}", selector)
					router.ServeHTTP(w, req)

					Convey("And the status code should be 200", func() {
						So(w.Code, ShouldEqual, http.StatusOK)
					})
				})
			})

			Convey("When the dimension API responds with an error", func() {
				mockFilter := NewMockFilterClient(mockCtrl)
				mockFilter.EXPECT().
					GetJobState(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Model{}, "", nil)
				mockFilter.
					EXPECT().
					GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(stubAreaTypeDimension, "", nil).
					AnyTimes()
				mockFilter.
					EXPECT().
					GetDimensionOptions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.DimensionOptions{}, "", nil).
					AnyTimes()

				mockPc := NewMockPopulationClient(mockCtrl)
				mockPc.EXPECT().
					GetAreaTypes(gomock.Any(), gomock.Any()).
					Return(population.GetAreaTypesResponse{}, errors.New("oh no"))

				mockRend := NewMockRenderClient(mockCtrl)
				mockRend.EXPECT().
					NewBasePageModel().
					Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain)).
					AnyTimes()

				mockDc := NewMockDatasetClient(mockCtrl)
				mockDc.EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(mockDataset, nil).AnyTimes()
				mockDc.
					EXPECT().
					GetVersion(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dataset.Version{
						LowestGeography: "",
						Version:         1,
					}, nil).
					AnyTimes()

				mockZc := NewMockZebedeeClient(mockCtrl)
				mockZc.
					EXPECT().
					GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(zebedee.HomepageContent{}, nil)

				ff := NewFilterFlex(mockRend, mockFilter, mockDc, mockPc, mockZc, cfg)
				w := runDimensionsSelector(dimensionName, ff.DimensionSelector())

				Convey("Then the status code should be 500", func() {
					So(w.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})

			Convey("When the zebedee API responds with an error", func() {
				mockFilter := NewMockFilterClient(mockCtrl)
				mockFilter.EXPECT().
					GetJobState(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Model{}, "", nil)
				mockFilter.
					EXPECT().
					GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(stubAreaTypeDimension, "", nil).
					AnyTimes()
				mockFilter.
					EXPECT().
					GetDimensionOptions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.DimensionOptions{}, "", nil).
					AnyTimes()

				mockPc := NewMockPopulationClient(mockCtrl)
				mockPc.EXPECT().
					GetAreaTypes(gomock.Any(), gomock.Any()).
					Return(population.GetAreaTypesResponse{}, nil)

				mockRend := NewMockRenderClient(mockCtrl)
				mockRend.EXPECT().
					NewBasePageModel().
					Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain)).
					AnyTimes()
				mockRend.
					EXPECT().
					BuildPage(gomock.Any(), gomock.Any(), gomock.Any())

				mockDc := NewMockDatasetClient(mockCtrl)
				mockDc.EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(mockDataset, nil).AnyTimes()
				mockDc.
					EXPECT().
					GetVersion(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dataset.Version{
						LowestGeography: "",
						Version:         1,
					}, nil).
					AnyTimes()

				mockZc := NewMockZebedeeClient(mockCtrl)
				mockZc.
					EXPECT().
					GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(zebedee.HomepageContent{}, errors.New("Internal error"))

				ff := NewFilterFlex(mockRend, mockFilter, mockDc, mockPc, mockZc, cfg)
				w := runDimensionsSelector(dimensionName, ff.DimensionSelector())

				Convey("Then the status code should be 200", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
				})
			})

			Convey("When the filter API responds with an error on GetDimensionOptions", func() {
				mockFilter := NewMockFilterClient(mockCtrl)
				mockFilter.EXPECT().
					GetJobState(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Model{}, "", nil)
				mockFilter.
					EXPECT().
					GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(stubAreaTypeDimension, "", nil).
					AnyTimes()
				mockFilter.
					EXPECT().
					GetDimensionOptions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.DimensionOptions{}, "", errors.New("uh oh")).
					AnyTimes()

				mockRend := NewMockRenderClient(mockCtrl)
				mockRend.EXPECT().
					NewBasePageModel().
					Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain)).
					AnyTimes()

				mockDc := NewMockDatasetClient(mockCtrl)
				mockDc.EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(mockDataset, nil).AnyTimes()
				mockDc.
					EXPECT().
					GetVersion(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dataset.Version{
						LowestGeography: "",
						Version:         1,
					}, nil).
					AnyTimes()

				mockZc := NewMockZebedeeClient(mockCtrl)
				mockZc.
					EXPECT().
					GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(zebedee.HomepageContent{}, nil)

				ff := NewFilterFlex(mockRend, mockFilter, mockDc, NewMockPopulationClient(mockCtrl), mockZc, cfg)
				w := runDimensionsSelector(dimensionName, ff.DimensionSelector())

				Convey("Then the status code should be 500", func() {
					So(w.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})

			Convey("When the dataset client responds with an error on Get", func() {
				mockDataset.Type = "multivariate"

				mockFilter := NewMockFilterClient(mockCtrl)
				mockFilter.EXPECT().
					GetJobState(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Model{}, "", nil)
				mockFilter.
					EXPECT().
					GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimension{}, "", nil).
					AnyTimes()

				mockRend := NewMockRenderClient(mockCtrl)
				mockRend.EXPECT().
					NewBasePageModel().
					Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain)).
					AnyTimes()

				mockDc := NewMockDatasetClient(mockCtrl)
				mockDc.EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(mockDataset, errors.New("Internal error"))

				mockZc := NewMockZebedeeClient(mockCtrl)
				mockZc.
					EXPECT().
					GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(zebedee.HomepageContent{}, nil)

				ff := NewFilterFlex(mockRend, mockFilter, mockDc, NewMockPopulationClient(mockCtrl), mockZc, cfg)
				w := runDimensionsSelector(dimensionName, ff.DimensionSelector())

				Convey("Then the status code should be 500", func() {
					So(w.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})

			Convey("When the population client responds with an error on GetCategorisations", func() {
				mockDataset.Type = "multivariate"

				mockFilter := NewMockFilterClient(mockCtrl)
				mockFilter.EXPECT().
					GetJobState(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Model{}, "", nil)
				mockFilter.
					EXPECT().
					GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimension{}, "", nil).
					AnyTimes()

				mockRend := NewMockRenderClient(mockCtrl)
				mockRend.EXPECT().
					NewBasePageModel().
					Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain)).
					AnyTimes()

				mockDc := NewMockDatasetClient(mockCtrl)
				mockDc.EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(mockDataset, nil)

				mockZc := NewMockZebedeeClient(mockCtrl)
				mockZc.
					EXPECT().
					GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(zebedee.HomepageContent{}, nil)

				mockPc := NewMockPopulationClient(mockCtrl)
				mockPc.
					EXPECT().
					GetCategorisations(gomock.Any(), gomock.Any()).
					Return(population.GetCategorisationsResponse{}, errors.New("Internal error"))

				ff := NewFilterFlex(mockRend, mockFilter, mockDc, mockPc, mockZc, cfg)
				w := runDimensionsSelector(dimensionName, ff.DimensionSelector())

				Convey("Then the status code should be 500", func() {
					So(w.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})
		})

	})

	Convey("Lowest geography override", t, func() {
		Convey("Test the override function", func() {
			hardcodedPopulation := "UR_CE"
			hardcodedLG := "msoa"
			defaultLG := "default"

			So(overrideLowestGeography(defaultLG, hardcodedPopulation, true), ShouldEqual, hardcodedLG)
			So(overrideLowestGeography(defaultLG, hardcodedPopulation, false), ShouldEqual, defaultLG)
			So(overrideLowestGeography(defaultLG, "normal population", true), ShouldEqual, defaultLG)
			So(overrideLowestGeography(defaultLG, "normal population", false), ShouldEqual, defaultLG)
		})
	})
}

func runDimensionsSelector(dimension string, selector func(http.ResponseWriter, *http.Request)) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", fmt.Sprintf("/filters/1234/dimensions/%s", dimension), nil)

	router := mux.NewRouter()
	router.HandleFunc("/filters/{filterID}/dimensions/{name}", selector)

	router.ServeHTTP(w, req)

	return w
}

// pageMatchesSelections is a gomock matcher that confirms a selection page
// contains the correct selections (i.e. radio buttons).
type pageMatchesSelections struct {
	selections []model.Selection
}

func (c pageMatchesSelections) Matches(x interface{}) bool {
	page, ok := x.(model.Selector)
	if !ok {
		return false
	}

	return reflect.DeepEqual(c.selections, page.Selections)
}

func (c pageMatchesSelections) String() string {
	return fmt.Sprintf("is equal to %+v", c.selections)
}

// pageMatchesSelections is a gomock matcher that confirms a selection page
// has the correct page title.
type pageHasTitle struct {
	title string
}

func (p pageHasTitle) Matches(x interface{}) bool {
	page, ok := x.(model.Selector)
	if !ok {
		return false
	}

	return p.title == page.Page.Metadata.Title
}

func (p pageHasTitle) String() string {
	return fmt.Sprintf("title is equal to \"%s\"", p.title)
}

// pageIsAreaType is a gomock matcher that confirms a selection page
// `IsAreaType` boolean is set to the expected value.
type pageIsAreaType struct {
	expected bool
}

func (c pageIsAreaType) Matches(x interface{}) bool {
	page, ok := x.(model.Selector)
	if !ok {
		return false
	}

	return page.IsAreaType == c.expected
}

func (c pageIsAreaType) String() string {
	return fmt.Sprintf("is equal to %+v", c.expected)
}

// pageHasError is a gomock matcher that confirms a selection page
// has a populated error.
type pageHasError struct {
	expected bool
}

func (p pageHasError) Matches(x interface{}) bool {
	page, ok := x.(model.Selector)
	if !ok {
		return false
	}

	if p.expected {
		return page.Error.Title != ""
	}

	return page.Error.Title == ""
}

func (p pageHasError) String() string {
	return fmt.Sprintf("is equal to %+v", p.expected)
}
