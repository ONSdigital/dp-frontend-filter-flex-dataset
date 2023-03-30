package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/helpers"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/mocks"
	"github.com/ONSdigital/dp-renderer/helper"
	coreModel "github.com/ONSdigital/dp-renderer/model"
	gomock "github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

func TestOverviewHandler(t *testing.T) {
	helper.InitialiseLocalisationsHelper(mocks.MockAssetFunction)
	mockCtrl := gomock.NewController(t)
	cfg := initialiseMockConfig()
	ctx := gomock.Any()
	mockDimensionCategories := []population.DimensionCategory{
		{
			Categories: []population.DimensionCategoryItem{
				{
					Label: "an option",
				},
			},
		},
		{
			Categories: []population.DimensionCategoryItem{},
		},
	}

	Convey("test filter flex overview", t, func() {
		Convey("test filter flex overview page is successful", func() {
			dims := filter.Dimensions{
				Items: []filter.Dimension{
					{
						Name:       "Test",
						IsAreaType: new(bool),
						Options:    []string{},
					},
				},
			}

			Convey("options on filter job no additional call to get options", func() {
				mockRend := NewMockRenderClient(mockCtrl)
				mockDc := NewMockDatasetClient(mockCtrl)
				mockPc := NewMockPopulationClient(mockCtrl)

				mockFc := NewMockFilterClient(mockCtrl)
				mockFilterDims := filter.Dimensions{
					Items: []filter.Dimension{
						{
							Name:       "Test",
							IsAreaType: new(bool),
							Options:    []string{"an option", "and another"},
						},
					},
				}
				mockRend.EXPECT().NewBasePageModel().Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain))
				mockRend.EXPECT().BuildPage(gomock.Any(), gomock.Any(), "overview")
				mockFc.EXPECT().GetFilter(ctx, gomock.Any()).Return(&filter.GetFilterResponse{}, nil)
				mockFc.EXPECT().GetDimensions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockFilterDims, "", nil)
				mockFc.EXPECT().GetDimension(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockFilterDims.Items[0], "", nil)
				mockDc.EXPECT().Get(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(dataset.DatasetDetails{}, nil)
				mockZc := NewMockZebedeeClient(mockCtrl)
				mockZc.
					EXPECT().
					GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(zebedee.HomepageContent{}, nil)
				mockPc.EXPECT().GetCategorisations(ctx, gomock.Any()).Return(population.GetCategorisationsResponse{
					PaginationResponse: population.PaginationResponse{
						TotalCount: 2,
					},
				}, nil).AnyTimes()

				mockPc.
					EXPECT().
					GetDimensionsDescription(ctx, gomock.Any()).
					Return(population.GetDimensionsResponse{}, nil)
				mockPc.
					EXPECT().
					GetDimensionCategories(ctx, gomock.Any()).
					Return(population.GetDimensionCategoriesResponse{
						PaginationResponse: population.PaginationResponse{TotalCount: 1},
						Categories:         mockDimensionCategories,
					}, nil).AnyTimes()
				mockPc.
					EXPECT().
					GetPopulationType(ctx, gomock.Any()).
					Return(population.GetPopulationTypeResponse{}, nil)

				w := httptest.NewRecorder()
				req := httptest.NewRequest("GET", "/filters/12345/dimensions", nil)

				ff := NewFilterFlex(mockRend, mockFc, mockDc, mockPc, mockZc, cfg)
				router := mux.NewRouter()
				router.HandleFunc("/filters/12345/dimensions", ff.FilterFlexOverview())
				router.ServeHTTP(w, req)

				So(w.Code, ShouldEqual, http.StatusOK)
			})

			Convey("when the zebedee.GetHomepageContent api method responds with an error", func() {
				mockRend := NewMockRenderClient(mockCtrl)
				mockDc := NewMockDatasetClient(mockCtrl)
				mockPc := NewMockPopulationClient(mockCtrl)

				mockFc := NewMockFilterClient(mockCtrl)
				mockFilterDims := filter.Dimensions{
					Items: []filter.Dimension{
						{
							Name:       "Test",
							IsAreaType: new(bool),
							Options:    []string{"an option", "and another"},
						},
					},
				}
				mockRend.EXPECT().NewBasePageModel().Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain))
				mockRend.EXPECT().BuildPage(gomock.Any(), gomock.Any(), "overview")
				mockFc.EXPECT().GetFilter(ctx, gomock.Any()).Return(&filter.GetFilterResponse{}, nil)
				mockFc.EXPECT().GetDimensions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockFilterDims, "", nil)
				mockFc.EXPECT().GetDimension(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockFilterDims.Items[0], "", nil)
				mockDc.EXPECT().Get(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(dataset.DatasetDetails{}, nil)
				mockZc := NewMockZebedeeClient(mockCtrl)
				mockZc.
					EXPECT().
					GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(zebedee.HomepageContent{}, errors.New("Internal error"))
				mockPc.EXPECT().GetCategorisations(ctx, gomock.Any()).Return(population.GetCategorisationsResponse{
					PaginationResponse: population.PaginationResponse{
						TotalCount: 2,
					},
				}, nil).AnyTimes()

				mockPc.
					EXPECT().
					GetDimensionsDescription(ctx, gomock.Any()).
					Return(population.GetDimensionsResponse{}, nil)
				mockPc.EXPECT().GetDimensionCategories(ctx, gomock.Any()).
					Return(population.GetDimensionCategoriesResponse{
						PaginationResponse: population.PaginationResponse{TotalCount: 1},
						Categories:         mockDimensionCategories,
					}, nil).AnyTimes()
				mockPc.
					EXPECT().
					GetPopulationType(ctx, gomock.Any()).
					Return(population.GetPopulationTypeResponse{}, nil)

				w := httptest.NewRecorder()
				req := httptest.NewRequest("GET", "/filters/12345/dimensions", nil)

				ff := NewFilterFlex(mockRend, mockFc, mockDc, mockPc, mockZc, cfg)
				router := mux.NewRouter()
				router.HandleFunc("/filters/12345/dimensions", ff.FilterFlexOverview())

				router.ServeHTTP(w, req)

				Convey("then the status is 200", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
				})
			})

			Convey("no options on filter job additional call to get options", func() {
				mockRend := NewMockRenderClient(mockCtrl)
				mockDc := NewMockDatasetClient(mockCtrl)
				mockFc := NewMockFilterClient(mockCtrl)
				mockPc := NewMockPopulationClient(mockCtrl)

				mockRend.EXPECT().NewBasePageModel().Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain))
				mockRend.EXPECT().BuildPage(gomock.Any(), gomock.Any(), "overview")
				mockFc.EXPECT().GetFilter(ctx, gomock.Any()).Return(&filter.GetFilterResponse{}, nil)
				mockFc.EXPECT().GetDimensions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(dims, "", nil)
				mockFc.EXPECT().GetDimension(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(dims.Items[0], "", nil)
				mockDc.EXPECT().Get(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(dataset.DatasetDetails{}, nil)
				mockZc := NewMockZebedeeClient(mockCtrl)
				mockZc.
					EXPECT().
					GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(zebedee.HomepageContent{}, nil)
				mockPc.EXPECT().GetCategorisations(ctx, gomock.Any()).Return(population.GetCategorisationsResponse{
					PaginationResponse: population.PaginationResponse{
						TotalCount: 2,
					},
				}, nil).AnyTimes()

				mockPc.
					EXPECT().
					GetDimensionsDescription(ctx, gomock.Any()).
					Return(population.GetDimensionsResponse{}, nil)
				mockPc.EXPECT().GetDimensionCategories(ctx, gomock.Any()).
					Return(population.GetDimensionCategoriesResponse{
						PaginationResponse: population.PaginationResponse{TotalCount: 1},
						Categories:         mockDimensionCategories,
					}, nil).AnyTimes()
				mockPc.
					EXPECT().
					GetPopulationType(ctx, gomock.Any()).
					Return(population.GetPopulationTypeResponse{}, nil)

				w := httptest.NewRecorder()
				req := httptest.NewRequest("GET", "/filters/12345/dimensions", nil)

				ff := NewFilterFlex(mockRend, mockFc, mockDc, mockPc, mockZc, cfg)
				router := mux.NewRouter()
				router.HandleFunc("/filters/12345/dimensions", ff.FilterFlexOverview())

				router.ServeHTTP(w, req)

				So(w.Code, ShouldEqual, http.StatusOK)
			})

			Convey("Given an area type dimension", func() {
				Convey("When the dimensions API responds with an error", func() {
					filterDim := filter.Dimension{
						Name:       "geography",
						ID:         "city",
						Label:      "City",
						IsAreaType: helpers.ToBoolPtr(true),
					}

					mockFc := NewMockFilterClient(mockCtrl)
					mockFc.
						EXPECT().
						GetFilter(ctx, gomock.Any()).
						Return(&filter.GetFilterResponse{}, nil).
						AnyTimes()
					mockFc.
						EXPECT().
						GetDimensions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(filter.Dimensions{Items: []filter.Dimension{filterDim}}, "", nil).
						AnyTimes()
					mockFc.
						EXPECT().
						GetDimension(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(filterDim, "", nil).
						AnyTimes()
					mockFc.
						EXPECT().
						GetDimensionOptions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(filter.DimensionOptions{}, "", nil)

					mockPc := NewMockPopulationClient(mockCtrl)
					mockPc.
						EXPECT().
						GetAreas(ctx, gomock.Any()).
						Return(population.GetAreasResponse{}, errors.New("internal error"))
					mockPc.
						EXPECT().
						GetDimensionsDescription(ctx, gomock.Any()).
						Return(population.GetDimensionsResponse{}, nil)
					mockPc.
						EXPECT().
						GetDimensionCategories(ctx, gomock.Any()).
						Return(population.GetDimensionCategoriesResponse{
							PaginationResponse: population.PaginationResponse{TotalCount: 1},
							Categories:         mockDimensionCategories,
						}, nil).AnyTimes()
					mockPc.EXPECT().
						GetCategorisations(ctx, gomock.Any()).
						Return(population.GetCategorisationsResponse{
							PaginationResponse: population.PaginationResponse{
								TotalCount: 2,
							},
						}, nil).AnyTimes()
					mockPc.
						EXPECT().
						GetPopulationType(ctx, gomock.Any()).
						Return(population.GetPopulationTypeResponse{}, nil)

					mockDc := NewMockDatasetClient(mockCtrl)
					mockDc.
						EXPECT().
						Get(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(dataset.DatasetDetails{}, nil)

					mockZc := NewMockZebedeeClient(mockCtrl)
					mockZc.
						EXPECT().
						GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(zebedee.HomepageContent{}, nil)

					w := httptest.NewRecorder()
					req := httptest.NewRequest(http.MethodGet, "/", nil)

					ff := NewFilterFlex(NewMockRenderClient(mockCtrl), mockFc, mockDc, mockPc, mockZc, cfg)
					ff.FilterFlexOverview().
						ServeHTTP(w, req)

					Convey("Then the status code should be 500", func() {
						So(w.Code, ShouldEqual, http.StatusInternalServerError)
					})
				})

				Convey("When the dimensions API responds successfully", func() {
					Convey("Then the dimensions API should be called with the population type and area type ID", func() {
						const (
							cantabularPopType = "cantabular-flexible-example"
							dimensionID       = "city"
						)

						filterDim := filter.Dimension{
							Name:       "geography",
							ID:         "city",
							Label:      "City",
							IsAreaType: helpers.ToBoolPtr(true),
						}

						mockFc := NewMockFilterClient(mockCtrl)
						mockFc.
							EXPECT().
							GetFilter(ctx, gomock.Any()).
							Return(&filter.GetFilterResponse{PopulationType: cantabularPopType}, nil).
							AnyTimes()
						mockFc.
							EXPECT().
							GetDimensions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
							Return(filter.Dimensions{Items: []filter.Dimension{filterDim}}, "", nil).
							AnyTimes()
						mockFc.
							EXPECT().
							GetDimension(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
							Return(filterDim, "", nil).
							AnyTimes()
						mockFc.
							EXPECT().
							GetDimensionOptions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
							Return(filter.DimensionOptions{}, "", nil)

						expGetAreasInput := population.GetAreasInput{
							PaginationParams: population.PaginationParams{
								Limit: 1000,
							},
							PopulationType: cantabularPopType,
							AreaTypeID:     dimensionID,
						}

						mockPc := NewMockPopulationClient(mockCtrl)
						mockPc.
							EXPECT().
							GetAreas(ctx, gomock.Eq(expGetAreasInput)).
							Return(population.GetAreasResponse{}, nil)
						mockPc.
							EXPECT().
							GetDimensionsDescription(ctx, gomock.Any()).
							Return(population.GetDimensionsResponse{}, nil)
						mockPc.EXPECT().
							GetCategorisations(ctx, gomock.Any()).
							Return(population.GetCategorisationsResponse{
								PaginationResponse: population.PaginationResponse{
									TotalCount: 2,
								},
							}, nil).AnyTimes()
						mockPc.
							EXPECT().
							GetPopulationType(ctx, gomock.Any()).
							Return(population.GetPopulationTypeResponse{}, nil)

						mockRend := NewMockRenderClient(mockCtrl)
						mockRend.
							EXPECT().
							NewBasePageModel().
							Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain)).
							AnyTimes()
						mockRend.
							EXPECT().
							BuildPage(gomock.Any(), gomock.Any(), "overview").
							AnyTimes()

						mockDc := NewMockDatasetClient(mockCtrl)
						mockDc.
							EXPECT().
							Get(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
							Return(dataset.DatasetDetails{}, nil)

						mockZc := NewMockZebedeeClient(mockCtrl)
						mockZc.
							EXPECT().
							GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
							Return(zebedee.HomepageContent{}, nil)

						w := httptest.NewRecorder()
						req := httptest.NewRequest(http.MethodGet, "/", nil)

						ff := NewFilterFlex(mockRend, mockFc, mockDc, mockPc, mockZc, cfg)
						ff.FilterFlexOverview().
							ServeHTTP(w, req)

						Convey("Then the status code should be 200", func() {
							So(w.Code, ShouldEqual, http.StatusOK)
						})
					})

					Convey("Then area type dimensions are used as options", func() {
						filterDim := filter.Dimension{
							Name:       "geography",
							ID:         "city",
							Label:      "City",
							IsAreaType: helpers.ToBoolPtr(true),
						}

						mockFc := NewMockFilterClient(mockCtrl)
						mockFc.
							EXPECT().
							GetFilter(ctx, gomock.Any()).
							Return(&filter.GetFilterResponse{}, nil).
							AnyTimes()
						mockFc.
							EXPECT().
							GetDimensions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
							Return(filter.Dimensions{Items: []filter.Dimension{filterDim}}, "", nil).
							AnyTimes()
						mockFc.
							EXPECT().
							GetDimension(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
							Return(filterDim, "", nil).
							AnyTimes()
						mockFc.
							EXPECT().
							GetDimensionOptions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
							Return(filter.DimensionOptions{}, "", nil)

						areas := population.GetAreasResponse{
							Areas: []population.Area{
								{
									ID:       "0",
									Label:    "London",
									AreaType: "city",
								},
							},
						}

						mockPc := NewMockPopulationClient(mockCtrl)
						mockPc.
							EXPECT().
							GetAreas(ctx, gomock.Any()).
							Return(areas, nil)

						mockPc.
							EXPECT().
							GetDimensionsDescription(ctx, gomock.Any()).
							Return(population.GetDimensionsResponse{}, nil)
						mockPc.EXPECT().
							GetCategorisations(ctx, gomock.Any()).
							Return(population.GetCategorisationsResponse{
								PaginationResponse: population.PaginationResponse{
									TotalCount: 2,
								},
							}, nil).AnyTimes()
						mockPc.
							EXPECT().
							GetPopulationType(ctx, gomock.Any()).
							Return(population.GetPopulationTypeResponse{}, nil)

						mockRend := NewMockRenderClient(mockCtrl)
						mockRend.
							EXPECT().
							NewBasePageModel().
							Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain)).
							AnyTimes()
						mockRend.
							EXPECT().
							BuildPage(gomock.Any(), gomock.Any(), "overview")

						mockDc := NewMockDatasetClient(mockCtrl)
						mockDc.
							EXPECT().
							Get(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
							Return(dataset.DatasetDetails{}, nil)

						mockZc := NewMockZebedeeClient(mockCtrl)
						mockZc.
							EXPECT().
							GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
							Return(zebedee.HomepageContent{}, nil)

						w := httptest.NewRecorder()
						req := httptest.NewRequest(http.MethodGet, "/test", nil)

						ff := NewFilterFlex(mockRend, mockFc, mockDc, mockPc, mockZc, cfg)
						ff.FilterFlexOverview().
							ServeHTTP(w, req)

						Convey("Then the status code should be 200", func() {
							So(w.Code, ShouldEqual, http.StatusOK)
						})
					})

					Convey("Then additional call to GetArea when GetDimensionOptions contains data", func() {
						filterDim := filter.Dimension{
							Name:       "geography",
							ID:         "city",
							Label:      "City",
							IsAreaType: helpers.ToBoolPtr(true),
						}

						mockFc := NewMockFilterClient(mockCtrl)
						mockFc.
							EXPECT().
							GetFilter(ctx, gomock.Any()).
							Return(&filter.GetFilterResponse{}, nil).
							AnyTimes()
						mockFc.
							EXPECT().
							GetDimensions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
							Return(filter.Dimensions{Items: []filter.Dimension{filterDim}}, "", nil).
							AnyTimes()
						mockFc.
							EXPECT().
							GetDimension(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
							Return(filterDim, "", nil).
							AnyTimes()
						mockFc.
							EXPECT().
							GetDimensionOptions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
							Return(filter.DimensionOptions{
								Items: []filter.DimensionOption{
									{
										Option: "0",
									},
								},
								Count:      1,
								TotalCount: 1,
							}, "", nil)

						area := population.GetAreaResponse{
							Area: population.Area{
								ID:       "0",
								Label:    "London",
								AreaType: "city",
							},
						}

						mockPc := NewMockPopulationClient(mockCtrl)
						mockPc.
							EXPECT().
							GetArea(ctx, gomock.Any()).
							Return(area, nil)

						mockPc.
							EXPECT().
							GetDimensionsDescription(ctx, gomock.Any()).
							Return(population.GetDimensionsResponse{}, nil)
						mockPc.
							EXPECT().
							GetDimensionCategories(ctx, gomock.Any()).
							Return(population.GetDimensionCategoriesResponse{
								PaginationResponse: population.PaginationResponse{TotalCount: 1},
								Categories:         mockDimensionCategories,
							}, nil).AnyTimes()
						mockPc.EXPECT().
							GetCategorisations(ctx, gomock.Any()).
							Return(population.GetCategorisationsResponse{
								PaginationResponse: population.PaginationResponse{
									TotalCount: 2,
								},
							}, nil).AnyTimes()
						mockPc.
							EXPECT().
							GetPopulationType(ctx, gomock.Any()).
							Return(population.GetPopulationTypeResponse{}, nil)

						mockRend := NewMockRenderClient(mockCtrl)
						mockRend.
							EXPECT().
							NewBasePageModel().
							Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain)).
							AnyTimes()
						mockRend.
							EXPECT().
							BuildPage(gomock.Any(), gomock.Any(), "overview")

						mockDc := NewMockDatasetClient(mockCtrl)
						mockDc.
							EXPECT().
							Get(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
							Return(dataset.DatasetDetails{}, nil)

						mockZc := NewMockZebedeeClient(mockCtrl)
						mockZc.
							EXPECT().
							GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
							Return(zebedee.HomepageContent{}, nil)

						w := httptest.NewRecorder()
						req := httptest.NewRequest(http.MethodGet, "/test", nil)

						ff := NewFilterFlex(mockRend, mockFc, mockDc, mockPc, mockZc, cfg)
						ff.FilterFlexOverview().
							ServeHTTP(w, req)

						Convey("Then the status code should be 200", func() {
							So(w.Code, ShouldEqual, http.StatusOK)
						})
					})

					Convey("Then additional call to GetBlockedAreaCount for multivariate dataset types", func() {
						filterDim := filter.Dimension{
							Name:           "geography",
							ID:             "city",
							Label:          "City",
							IsAreaType:     helpers.ToBoolPtr(true),
							FilterByParent: "england",
						}

						mockFc := NewMockFilterClient(mockCtrl)
						mockFc.
							EXPECT().
							GetFilter(ctx, gomock.Any()).
							Return(&filter.GetFilterResponse{}, nil).
							AnyTimes()
						mockFc.
							EXPECT().
							GetDimensions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
							Return(filter.Dimensions{Items: []filter.Dimension{filterDim}}, "", nil).
							AnyTimes()
						mockFc.
							EXPECT().
							GetDimension(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
							Return(filterDim, "", nil).
							AnyTimes()
						mockFc.
							EXPECT().
							GetDimensionOptions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
							Return(filter.DimensionOptions{
								Items: []filter.DimensionOption{
									{
										Option: "0",
									},
								},
								Count:      1,
								TotalCount: 1,
							}, "", nil)

						area := population.GetAreaResponse{
							Area: population.Area{
								ID:       "0",
								Label:    "London",
								AreaType: "city",
							},
						}

						mockPc := NewMockPopulationClient(mockCtrl)
						mockPc.
							EXPECT().
							GetArea(ctx, gomock.Any()).
							Return(area, nil)

						mockPc.
							EXPECT().
							GetDimensionsDescription(ctx, gomock.Any()).
							Return(population.GetDimensionsResponse{}, nil)
						mockPc.
							EXPECT().
							GetBlockedAreaCount(gomock.Any(), gomock.Any()).
							Return(&cantabular.GetBlockedAreaCountResult{}, nil)
						mockPc.
							EXPECT().
							GetDimensionCategories(ctx, gomock.Any()).
							Return(population.GetDimensionCategoriesResponse{
								PaginationResponse: population.PaginationResponse{TotalCount: 1},
								Categories:         mockDimensionCategories,
							}, nil).AnyTimes()

						mockPc.EXPECT().
							GetCategorisations(ctx, gomock.Any()).
							Return(population.GetCategorisationsResponse{
								PaginationResponse: population.PaginationResponse{
									TotalCount: 2,
								},
							}, nil).AnyTimes()
						mockPc.
							EXPECT().
							GetPopulationType(ctx, gomock.Any()).
							Return(population.GetPopulationTypeResponse{}, nil)

						// TODO: pc.GetParentAreaCount is causing production issues
						// mockPc.
						// 	EXPECT().
						// 	GetParentAreaCount(ctx, gomock.Any()).
						// 	Return(0, nil)

						mockRend := NewMockRenderClient(mockCtrl)
						mockRend.
							EXPECT().
							NewBasePageModel().
							Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain)).
							AnyTimes()
						mockRend.
							EXPECT().
							BuildPage(gomock.Any(), gomock.Any(), "overview")

						mockDc := NewMockDatasetClient(mockCtrl)
						mockDc.
							EXPECT().
							Get(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
							Return(dataset.DatasetDetails{
								Type: "multivariate",
							}, nil)

						mockZc := NewMockZebedeeClient(mockCtrl)
						mockZc.
							EXPECT().
							GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
							Return(zebedee.HomepageContent{}, nil)

						w := httptest.NewRecorder()
						req := httptest.NewRequest(http.MethodGet, "/test", nil)

						ff := NewFilterFlex(mockRend, mockFc, mockDc, mockPc, mockZc, cfg)
						ff.FilterFlexOverview().
							ServeHTTP(w, req)

						Convey("Then the status code should be 200", func() {
							So(w.Code, ShouldEqual, http.StatusOK)
						})
					})
				})
			})
		})

		Convey("test filter flex overview errors", func() {
			mockRend := NewMockRenderClient(mockCtrl)

			Convey("test FilterFlexOverview returns 500 if client GetJobState returns an error", func() {
				mockFc := NewMockFilterClient(mockCtrl)
				mockFc.
					EXPECT().
					GetFilter(ctx, gomock.Any()).
					Return(nil, errors.New("sorry"))
				mockFc.
					EXPECT().
					GetDimensions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimensions{}, "", nil)

				mockZc := NewMockZebedeeClient(mockCtrl)
				mockZc.
					EXPECT().
					GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(zebedee.HomepageContent{}, nil)

				mockPc := NewMockPopulationClient(mockCtrl)
				mockPc.
					EXPECT().
					GetDimensionCategories(ctx, gomock.Any()).
					Return(population.GetDimensionCategoriesResponse{
						PaginationResponse: population.PaginationResponse{TotalCount: 1},
						Categories:         mockDimensionCategories,
					}, nil).AnyTimes()
				mockPc.EXPECT().
					GetCategorisations(ctx, gomock.Any()).
					Return(population.GetCategorisationsResponse{
						PaginationResponse: population.PaginationResponse{
							TotalCount: 2,
						},
					}, nil).AnyTimes()

				w := httptest.NewRecorder()
				req := httptest.NewRequest("GET", "/filters/12345/dimensions", nil)

				ff := NewFilterFlex(mockRend, mockFc, NewMockDatasetClient(mockCtrl), mockPc, mockZc, cfg)
				router := mux.NewRouter()
				router.HandleFunc("/filters/12345/dimensions", ff.FilterFlexOverview())
				router.ServeHTTP(w, req)

				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})

			Convey("test FilterFlexOverview returns 500 if client dc.Get returns an error", func() {
				mockFc := NewMockFilterClient(mockCtrl)
				mockFc.
					EXPECT().
					GetFilter(ctx, gomock.Any()).
					Return(&filter.GetFilterResponse{}, nil)
				mockFc.
					EXPECT().
					GetDimensions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimensions{}, "", nil)

				mockDc := NewMockDatasetClient(mockCtrl)
				mockDc.
					EXPECT().
					Get(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dataset.DatasetDetails{}, errors.New("Internal error"))

				mockZc := NewMockZebedeeClient(mockCtrl)
				mockZc.
					EXPECT().
					GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(zebedee.HomepageContent{}, nil)

				mockPc := NewMockPopulationClient(mockCtrl)
				mockPc.
					EXPECT().
					GetDimensionCategories(ctx, gomock.Any()).
					Return(population.GetDimensionCategoriesResponse{
						PaginationResponse: population.PaginationResponse{TotalCount: 1},
						Categories:         mockDimensionCategories,
					}, nil).AnyTimes()
				mockPc.EXPECT().
					GetCategorisations(ctx, gomock.Any()).
					Return(population.GetCategorisationsResponse{
						PaginationResponse: population.PaginationResponse{
							TotalCount: 2,
						},
					}, nil).AnyTimes()

				w := httptest.NewRecorder()
				req := httptest.NewRequest("GET", "/filters/12345/dimensions", nil)

				ff := NewFilterFlex(mockRend, mockFc, mockDc, mockPc, mockZc, cfg)
				router := mux.NewRouter()
				router.HandleFunc("/filters/12345/dimensions", ff.FilterFlexOverview())

				router.ServeHTTP(w, req)

				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})

			Convey("test FilterFlexOverview returns 500 if client GetDimensionsDescription returns an error", func() {
				mockFc := NewMockFilterClient(mockCtrl)
				mockDc := NewMockDatasetClient(mockCtrl)
				mockPc := NewMockPopulationClient(mockCtrl)

				mockFilterDims := filter.Dimensions{
					Items: []filter.Dimension{
						{
							Name:       "Test",
							IsAreaType: new(bool),
						},
					},
				}

				mockFc.
					EXPECT().
					GetFilter(ctx, gomock.Any()).
					Return(&filter.GetFilterResponse{}, nil)
				mockFc.
					EXPECT().
					GetDimensions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(mockFilterDims, "", nil)
				mockDc.
					EXPECT().
					Get(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dataset.DatasetDetails{}, nil)

				mockZc := NewMockZebedeeClient(mockCtrl)
				mockZc.
					EXPECT().
					GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(zebedee.HomepageContent{}, nil)

				mockPc.
					EXPECT().
					GetDimensionsDescription(ctx, gomock.Any()).
					Return(population.GetDimensionsResponse{}, errors.New("Internal error"))
				mockPc.
					EXPECT().
					GetDimensionCategories(ctx, gomock.Any()).
					Return(population.GetDimensionCategoriesResponse{
						PaginationResponse: population.PaginationResponse{TotalCount: 1},
						Categories:         mockDimensionCategories,
					}, nil).AnyTimes()
				mockPc.EXPECT().
					GetCategorisations(ctx, gomock.Any()).
					Return(population.GetCategorisationsResponse{
						PaginationResponse: population.PaginationResponse{
							TotalCount: 2,
						},
					}, nil).AnyTimes()
				mockPc.
					EXPECT().
					GetPopulationType(ctx, gomock.Any()).
					Return(population.GetPopulationTypeResponse{}, nil)

				w := httptest.NewRecorder()
				req := httptest.NewRequest("GET", "/filters/12345/dimensions", nil)

				ff := NewFilterFlex(mockRend, mockFc, mockDc, mockPc, mockZc, cfg)
				router := mux.NewRouter()
				router.HandleFunc("/filters/12345/dimensions", ff.FilterFlexOverview())
				router.ServeHTTP(w, req)

				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})

			Convey("test FilterFlexOverview returns 500 if client GetBlockedAreaCount returns an error", func() {
				filterDim := filter.Dimension{
					Name:           "geography",
					ID:             "city",
					Label:          "City",
					IsAreaType:     helpers.ToBoolPtr(true),
					FilterByParent: "england",
				}

				mockFc := NewMockFilterClient(mockCtrl)
				mockFc.
					EXPECT().
					GetFilter(ctx, gomock.Any()).
					Return(&filter.GetFilterResponse{}, nil).
					AnyTimes()
				mockFc.
					EXPECT().
					GetDimensions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimensions{Items: []filter.Dimension{filterDim}}, "", nil).
					AnyTimes()
				mockFc.
					EXPECT().
					GetDimension(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filterDim, "", nil).
					AnyTimes()
				mockFc.
					EXPECT().
					GetDimensionOptions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.DimensionOptions{
						Items: []filter.DimensionOption{
							{
								Option: "0",
							},
						},
						Count:      1,
						TotalCount: 1,
					}, "", nil)

				area := population.GetAreaResponse{
					Area: population.Area{
						ID:       "0",
						Label:    "London",
						AreaType: "city",
					},
				}

				mockPc := NewMockPopulationClient(mockCtrl)
				mockPc.
					EXPECT().
					GetArea(ctx, gomock.Any()).
					Return(area, nil)
				mockPc.
					EXPECT().
					GetDimensionsDescription(ctx, gomock.Any()).
					Return(population.GetDimensionsResponse{}, nil)
				mockPc.
					EXPECT().
					GetBlockedAreaCount(gomock.Any(), gomock.Any()).
					Return(&cantabular.GetBlockedAreaCountResult{}, errors.New("Sorry"))
				mockPc.
					EXPECT().
					GetDimensionCategories(ctx, gomock.Any()).
					Return(population.GetDimensionCategoriesResponse{
						PaginationResponse: population.PaginationResponse{TotalCount: 1},
						Categories:         mockDimensionCategories,
					}, nil).AnyTimes()
				mockPc.EXPECT().
					GetCategorisations(ctx, gomock.Any()).
					Return(population.GetCategorisationsResponse{
						PaginationResponse: population.PaginationResponse{
							TotalCount: 2,
						},
					}, nil).AnyTimes()
				mockPc.
					EXPECT().
					GetPopulationType(ctx, gomock.Any()).
					Return(population.GetPopulationTypeResponse{}, nil)

				mockRend := NewMockRenderClient(mockCtrl)

				mockDc := NewMockDatasetClient(mockCtrl)
				mockDc.
					EXPECT().
					Get(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dataset.DatasetDetails{
						Type: "multivariate",
					}, nil)

				mockZc := NewMockZebedeeClient(mockCtrl)
				mockZc.
					EXPECT().
					GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(zebedee.HomepageContent{}, nil)

				w := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodGet, "/test", nil)

				ff := NewFilterFlex(mockRend, mockFc, mockDc, mockPc, mockZc, cfg)
				ff.FilterFlexOverview().
					ServeHTTP(w, req)

				Convey("Then the status code should be 500", func() {
					So(w.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})

			Convey("test FilterFlexOverview returns 500 if client GetDimensions returns an error", func() {
				mockFc := NewMockFilterClient(mockCtrl)
				mockDc := NewMockDatasetClient(mockCtrl)

				mockFc.
					EXPECT().
					GetFilter(ctx, gomock.Any()).
					Return(&filter.GetFilterResponse{}, nil)
				mockFc.
					EXPECT().
					GetDimensions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimensions{}, "", errors.New("sorry"))

				mockDc.
					EXPECT().
					Get(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dataset.DatasetDetails{}, nil)

				mockZc := NewMockZebedeeClient(mockCtrl)
				mockZc.
					EXPECT().
					GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(zebedee.HomepageContent{}, nil)

				mockPc := NewMockPopulationClient(mockCtrl)
				mockPc.
					EXPECT().
					GetDimensionCategories(ctx, gomock.Any()).
					Return(population.GetDimensionCategoriesResponse{
						PaginationResponse: population.PaginationResponse{TotalCount: 1},
						Categories:         mockDimensionCategories,
					}, nil).AnyTimes()
				mockPc.EXPECT().
					GetCategorisations(ctx, gomock.Any()).
					Return(population.GetCategorisationsResponse{
						PaginationResponse: population.PaginationResponse{
							TotalCount: 2,
						},
					}, nil).AnyTimes()

				w := httptest.NewRecorder()
				req := httptest.NewRequest("GET", "/filters/12345/dimensions", nil)

				ff := NewFilterFlex(mockRend, mockFc, mockDc, mockPc, mockZc, cfg)
				router := mux.NewRouter()
				router.HandleFunc("/filters/12345/dimensions", ff.FilterFlexOverview())
				router.ServeHTTP(w, req)

				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})

			Convey("test FilterFlexOverview returns 500 if client GetDimension returns an error", func() {
				mockFc := NewMockFilterClient(mockCtrl)
				mockDc := NewMockDatasetClient(mockCtrl)

				mockFilterDims := filter.Dimensions{
					Items: []filter.Dimension{
						{
							Name:       "Test",
							IsAreaType: new(bool),
						},
					},
				}

				mockFc.
					EXPECT().
					GetFilter(ctx, gomock.Any()).
					Return(&filter.GetFilterResponse{}, nil)
				mockFc.
					EXPECT().
					GetDimensions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(mockFilterDims, "", nil)
				mockFc.
					EXPECT().
					GetDimension(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), mockFilterDims.Items[0].Name).
					Return(filter.Dimension{}, "", errors.New("sorry"))

				mockDc.
					EXPECT().
					Get(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dataset.DatasetDetails{}, nil)

				mockPc := NewMockPopulationClient(mockCtrl)
				mockPc.
					EXPECT().
					GetDimensionCategories(ctx, gomock.Any()).
					Return(population.GetDimensionCategoriesResponse{
						PaginationResponse: population.PaginationResponse{TotalCount: 1},
						Categories:         mockDimensionCategories,
					}, nil).AnyTimes()
				mockPc.
					EXPECT().
					GetDimensionsDescription(ctx, gomock.Any()).
					Return(population.GetDimensionsResponse{}, nil)
				mockPc.EXPECT().
					GetCategorisations(ctx, gomock.Any()).
					Return(population.GetCategorisationsResponse{
						PaginationResponse: population.PaginationResponse{
							TotalCount: 2,
						},
					}, nil).AnyTimes()

				mockZc := NewMockZebedeeClient(mockCtrl)
				mockZc.
					EXPECT().
					GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(zebedee.HomepageContent{}, nil)

				mockPc.EXPECT().
					GetCategorisations(ctx, gomock.Any()).
					Return(population.GetCategorisationsResponse{
						PaginationResponse: population.PaginationResponse{
							TotalCount: 2,
						},
					}, nil).AnyTimes()
				mockPc.
					EXPECT().
					GetPopulationType(ctx, gomock.Any()).
					Return(population.GetPopulationTypeResponse{}, nil)

				w := httptest.NewRecorder()
				req := httptest.NewRequest("GET", "/filters/12345/dimensions", nil)

				ff := NewFilterFlex(mockRend, mockFc, mockDc, mockPc, mockZc, cfg)
				router := mux.NewRouter()
				router.HandleFunc("/filters/12345/dimensions", ff.FilterFlexOverview())
				router.ServeHTTP(w, req)

				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})

			Convey("test FilterFlexOverview returns 500 if pc.GetPopulationType returns an error", func() {
				mockFilterDims := filter.Dimensions{
					Items: []filter.Dimension{
						{
							Name:       "Test",
							IsAreaType: new(bool),
						},
					},
				}

				mockFc := NewMockFilterClient(mockCtrl)
				mockFc.
					EXPECT().
					GetFilter(ctx, gomock.Any()).
					Return(&filter.GetFilterResponse{}, nil)
				mockFc.
					EXPECT().
					GetDimensions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(mockFilterDims, "", nil)

				mockDc := NewMockDatasetClient(mockCtrl)
				mockDc.
					EXPECT().
					Get(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dataset.DatasetDetails{}, nil)

				mockPc := NewMockPopulationClient(mockCtrl)
				mockPc.
					EXPECT().
					GetDimensionCategories(ctx, gomock.Any()).
					Return(population.GetDimensionCategoriesResponse{
						PaginationResponse: population.PaginationResponse{TotalCount: 1},
						Categories:         mockDimensionCategories,
					}, nil).AnyTimes()
				mockPc.
					EXPECT().
					GetDimensionsDescription(gomock.Any(), gomock.Any()).
					Return(population.GetDimensionsResponse{}, nil)
				mockPc.
					EXPECT().
					GetPopulationType(ctx, gomock.Any()).
					Return(population.GetPopulationTypeResponse{}, errors.New("Sorry"))

				mockZc := NewMockZebedeeClient(mockCtrl)
				mockZc.
					EXPECT().
					GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(zebedee.HomepageContent{}, nil)

				w := httptest.NewRecorder()
				req := httptest.NewRequest("GET", "/filters/12345/dimensions", nil)

				ff := NewFilterFlex(mockRend, mockFc, mockDc, mockPc, mockZc, cfg)
				router := mux.NewRouter()
				router.HandleFunc("/filters/12345/dimensions", ff.FilterFlexOverview())
				router.ServeHTTP(w, req)

				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})

			// TODO: pc.GetParentAreaCount is causing production issues
			// Convey("test FilterFlexOverview returns 500 if client GetParentAreaCount returns an error", func() {
			// 	mockFc := NewMockFilterClient(mockCtrl)
			// 	mockDc := NewMockDatasetClient(mockCtrl)
			// 	mockPc := NewMockPopulationClient(mockCtrl)

			// 	filterDim := filter.Dimension{
			// 		Name:           "geography",
			// 		ID:             "city",
			// 		Label:          "City",
			// 		IsAreaType:     helpers.ToBoolPtr(true),
			// 		FilterByParent: "england",
			// 	}
			// 	mockDatasetDims := dataset.VersionDimensions{
			// 		Items: []dataset.VersionDimension{},
			// 	}

			// 	mockFc.
			// 		EXPECT().
			// 		GetFilter(ctx, gomock.Any()).
			// 		Return(&filter.GetFilterResponse{}, nil).
			// 		AnyTimes()
			// 	mockFc.
			// 		EXPECT().
			// 		GetDimensions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			// 		Return(filter.Dimensions{Items: []filter.Dimension{filterDim}}, "", nil).
			// 		AnyTimes()
			// 	mockFc.
			// 		EXPECT().
			// 		GetDimension(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			// 		Return(filterDim, "", nil).
			// 		AnyTimes()
			// 	mockFc.
			// 		EXPECT().
			// 		GetDimensionOptions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			// 		Return(filter.DimensionOptions{
			// 			Items: []filter.DimensionOption{
			// 				{
			// 					Option: "0",
			// 				},
			// 			},
			// 			Count:      1,
			// 			TotalCount: 1,
			// 		}, "", nil)

			// 	area := population.GetAreaResponse{
			// 		Area: population.Area{
			// 			ID:       "0",
			// 			Label:    "London",
			// 			AreaType: "city",
			// 		},
			// 	}

			// 	mockPc.
			// 		EXPECT().
			// 		GetArea(ctx, gomock.Any()).
			// 		Return(area, nil)
			// 	mockPc.
			// 		EXPECT().
			// 		GetParentAreaCount(ctx, gomock.Any()).
			// 		Return(0, errors.New("sorry"))

			// 	mockRend := NewMockRenderClient(mockCtrl)

			// 	mockDc.
			// 		EXPECT().
			// 		GetVersionDimensions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			// 		Return(mockDatasetDims, nil)
			// 	mockDc.
			// 		EXPECT().
			// 		Get(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			// 		Return(dataset.DatasetDetails{}, nil)

			// 	mockZc := NewMockZebedeeClient(mockCtrl)
			// 	mockZc.
			// 		EXPECT().
			// 		GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			// 		Return(zebedee.HomepageContent{}, nil)

			// 	w := httptest.NewRecorder()
			// 	req := httptest.NewRequest("GET", "/filters/12345/dimensions", nil)

			// 	router := mux.NewRouter()
			// 	router.HandleFunc("/filters/12345/dimensions", FilterFlexOverview(mockRend, mockFc, mockDc, mockPc, mockZc, cfg))

			// 	router.ServeHTTP(w, req)

			// 	So(w.Code, ShouldEqual, http.StatusInternalServerError)
			// })
		})
	})

	Convey("Population categories are sorted", t, func() {
		getCategoryList := func(items []population.DimensionCategoryItem) []string {
			results := []string{}
			for _, item := range items {
				results = append(results, item.ID)
			}
			return results
		}

		Convey("given non-numeric options", func() {
			nonNumeric := []population.DimensionCategoryItem{
				{
					ID: "option 2",
				},
				{
					ID: "option 1",
				},
			}
			Convey("when they are sorted", func() {
				sorted := sortCategoriesByID(nonNumeric)

				Convey("then categories are sorted alphabetically", func() {
					actual := getCategoryList(sorted)
					expected := []string{"option 1", "option 2"}
					So(actual, ShouldResemble, expected)
				})
			})
		})

		Convey("given simple numeric options", func() {
			simpleNumeric := []population.DimensionCategoryItem{
				{
					ID: "2",
				},
				{
					ID: "10",
				},
				{
					ID: "1",
				},
			}
			Convey("when they are sorted", func() {
				sorted := sortCategoriesByID(simpleNumeric)

				Convey("then options are sorted numerically", func() {
					actual := getCategoryList(sorted)
					expected := []string{"1", "2", "10"}
					So(actual, ShouldResemble, expected)
				})
			})
		})

		Convey("given numeric options with negatives", func() {
			numeric := []population.DimensionCategoryItem{
				{
					ID: "2",
				},
				{
					ID: "-1",
				},
				{
					ID: "10",
				},
				{
					ID: "-10",
				},
				{
					ID: "1",
				},
			}

			Convey("when they are sorted", func() {
				sorted := sortCategoriesByID(numeric)

				Convey("then options are sorted numerically with negatives at the end", func() {
					actual := getCategoryList(sorted)
					expected := []string{"1", "2", "10", "-1", "-10"}
					So(actual, ShouldResemble, expected)
				})
			})
		})

		Convey("given mixed numeric and non-numeric options", func() {
			mixed := []population.DimensionCategoryItem{
				{
					ID: "2nd Option",
				},
				{
					ID: "1",
				},
				{
					ID: "10",
				},
			}
			Convey("when they are sorted", func() {
				sorted := sortCategoriesByID(mixed)

				Convey("then options are sorted alphanumerically", func() {
					actual := getCategoryList(sorted)
					expected := []string{"1", "10", "2nd Option"}
					So(actual, ShouldResemble, expected)
				})
			})
		})
	})
}
