package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
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
	mockOpts := []dataset.Options{
		{
			Items: []dataset.Option{
				{
					Label: "an option",
				},
			},
		},
		{
			Items: []dataset.Option{},
		},
	}

	Convey("test filter flex overview", t, func() {
		Convey("test filter flex overview page is successful", func() {
			mockDatasetDims := dataset.VersionDimensions{
				Items: []dataset.VersionDimension{},
			}
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
				mockDc.EXPECT().GetVersionDimensions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockDatasetDims, nil)
				mockDc.EXPECT().GetOptions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), "", "", "0", dims.Items[0].Name,
					&dataset.QueryParams{Offset: 0, Limit: 1000}).Return(mockOpts[0], nil)

				w := httptest.NewRecorder()
				req := httptest.NewRequest("GET", "/filters/12345/dimensions", nil)

				router := mux.NewRouter()
				router.HandleFunc("/filters/12345/dimensions", FilterFlexOverview(mockRend, mockFc, mockDc, NewMockPopulationClient(mockCtrl)))

				router.ServeHTTP(w, req)

				So(w.Code, ShouldEqual, http.StatusOK)
			})

			Convey("no options on filter job additional call to get options", func() {
				mockRend := NewMockRenderClient(mockCtrl)
				mockDc := NewMockDatasetClient(mockCtrl)
				mockFc := NewMockFilterClient(mockCtrl)

				mockRend.EXPECT().NewBasePageModel().Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain))
				mockRend.EXPECT().BuildPage(gomock.Any(), gomock.Any(), "overview")
				mockFc.EXPECT().GetFilter(ctx, gomock.Any()).Return(&filter.GetFilterResponse{}, nil)
				mockFc.EXPECT().GetDimensions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(dims, "", nil)
				mockFc.EXPECT().GetDimension(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(dims.Items[0], "", nil)
				mockDc.EXPECT().GetOptions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), "", "", "0", dims.Items[0].Name,
					&dataset.QueryParams{Offset: 0, Limit: 1000}).Return(mockOpts[0], nil)
				mockDc.EXPECT().GetVersionDimensions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockDatasetDims, nil)

				w := httptest.NewRecorder()
				req := httptest.NewRequest("GET", "/filters/12345/dimensions", nil)

				router := mux.NewRouter()
				router.HandleFunc("/filters/12345/dimensions", FilterFlexOverview(mockRend, mockFc, mockDc, NewMockPopulationClient(mockCtrl)))

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

					mockDc := NewMockDatasetClient(mockCtrl)
					mockDc.
						EXPECT().
						GetVersionDimensions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(mockDatasetDims, nil)

					w := httptest.NewRecorder()
					req := httptest.NewRequest(http.MethodGet, "/", nil)

					FilterFlexOverview(NewMockRenderClient(mockCtrl), mockFc, mockDc, mockPc).
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
							GetVersionDimensions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
							Return(mockDatasetDims, nil)

						w := httptest.NewRecorder()
						req := httptest.NewRequest(http.MethodGet, "/", nil)

						FilterFlexOverview(mockRend, mockFc, mockDc, mockPc).
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
							GetVersionDimensions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
							Return(mockDatasetDims, nil)

						w := httptest.NewRecorder()
						req := httptest.NewRequest(http.MethodGet, "/test", nil)

						FilterFlexOverview(mockRend, mockFc, mockDc, mockPc).
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
							GetVersionDimensions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
							Return(mockDatasetDims, nil)

						w := httptest.NewRecorder()
						req := httptest.NewRequest(http.MethodGet, "/test", nil)

						FilterFlexOverview(mockRend, mockFc, mockDc, mockPc).
							ServeHTTP(w, req)

						Convey("Then the status code should be 200", func() {
							So(w.Code, ShouldEqual, http.StatusOK)
						})
					})

					// TODO: Hotfix due to graphQL error
					// Convey("Then additional call to GetParentAreaCount when dimension.FilterByParent contains data", func() {
					// 	filterDim := filter.Dimension{
					// 		Name:           "geography",
					// 		ID:             "city",
					// 		Label:          "City",
					// 		IsAreaType:     helpers.ToBoolPtr(true),
					// 		FilterByParent: "england",
					// 	}

					// 	mockFc := NewMockFilterClient(mockCtrl)
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

					// 	mockPc := NewMockPopulationClient(mockCtrl)
					// 	mockPc.
					// 		EXPECT().
					// 		GetArea(ctx, gomock.Any()).
					// 		Return(area, nil)
					// 	mockPc.
					// 		EXPECT().
					// 		GetParentAreaCount(ctx, gomock.Any()).
					// 		Return(0, nil)

					// 	mockRend := NewMockRenderClient(mockCtrl)
					// 	mockRend.
					// 		EXPECT().
					// 		NewBasePageModel().
					// 		Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain)).
					// 		AnyTimes()
					// 	mockRend.
					// 		EXPECT().
					// 		BuildPage(gomock.Any(), gomock.Any(), "overview")

					// 	mockDc := NewMockDatasetClient(mockCtrl)
					// 	mockDc.
					// 		EXPECT().
					// 		GetVersionDimensions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					// 		Return(mockDatasetDims, nil)

					// 	w := httptest.NewRecorder()
					// 	req := httptest.NewRequest(http.MethodGet, "/test", nil)

					// 	FilterFlexOverview(mockRend, mockFc, mockDc, mockPc).
					// 		ServeHTTP(w, req)

					// 	Convey("Then the status code should be 200", func() {
					// 		So(w.Code, ShouldEqual, http.StatusOK)
					// 	})
					// })
				})
			})
		})

		Convey("test filter flex overview errors", func() {
			mockRend := NewMockRenderClient(mockCtrl)

			Convey("test FilterFlexOverview returns 500 if client GetJobState returns an error", func() {
				mockFc := NewMockFilterClient(mockCtrl)
				mockDc := NewMockDatasetClient(mockCtrl)
				mockFc.
					EXPECT().
					GetFilter(ctx, gomock.Any()).
					Return(nil, errors.New("sorry"))
				mockFc.
					EXPECT().
					GetDimensions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimensions{}, "", nil)

				w := httptest.NewRecorder()
				req := httptest.NewRequest("GET", "/filters/12345/dimensions", nil)

				router := mux.NewRouter()
				router.HandleFunc("/filters/12345/dimensions", FilterFlexOverview(mockRend, mockFc, mockDc, NewMockPopulationClient(mockCtrl)))

				router.ServeHTTP(w, req)

				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})

			Convey("test FilterFlexOverview returns 500 if client GetVersionDimensions returns an error", func() {
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
					GetVersionDimensions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dataset.VersionDimensions{}, errors.New("sorry"))

				w := httptest.NewRecorder()
				req := httptest.NewRequest("GET", "/filters/12345/dimensions", nil)

				router := mux.NewRouter()
				router.HandleFunc("/filters/12345/dimensions", FilterFlexOverview(mockRend, mockFc, mockDc, mockPc))

				router.ServeHTTP(w, req)

				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})

			Convey("test FilterFlexOverview returns 500 if client GetDimensions returns an error", func() {
				mockFc := NewMockFilterClient(mockCtrl)
				mockDc := NewMockDatasetClient(mockCtrl)
				mockPc := NewMockPopulationClient(mockCtrl)

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
					GetVersionDimensions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dataset.VersionDimensions{}, nil)

				w := httptest.NewRecorder()
				req := httptest.NewRequest("GET", "/filters/12345/dimensions", nil)

				router := mux.NewRouter()
				router.HandleFunc("/filters/12345/dimensions", FilterFlexOverview(mockRend, mockFc, mockDc, mockPc))

				router.ServeHTTP(w, req)

				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})

			Convey("test FilterFlexOverview returns 500 if client GetDimension returns an error", func() {
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
				mockFc.
					EXPECT().
					GetDimension(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), mockFilterDims.Items[0].Name).
					Return(filter.Dimension{}, "", errors.New("sorry"))

				mockDc.
					EXPECT().
					GetVersionDimensions(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dataset.VersionDimensions{}, nil)

				w := httptest.NewRecorder()
				req := httptest.NewRequest("GET", "/filters/12345/dimensions", nil)

				router := mux.NewRouter()
				router.HandleFunc("/filters/12345/dimensions", FilterFlexOverview(mockRend, mockFc, mockDc, mockPc))

				router.ServeHTTP(w, req)

				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})

			// TODO: Hotfix due to graphQL error
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

			// 	w := httptest.NewRecorder()
			// 	req := httptest.NewRequest("GET", "/filters/12345/dimensions", nil)

			// 	router := mux.NewRouter()
			// 	router.HandleFunc("/filters/12345/dimensions", FilterFlexOverview(mockRend, mockFc, mockDc, mockPc))

			// 	router.ServeHTTP(w, req)

			// 	So(w.Code, ShouldEqual, http.StatusInternalServerError)
			// })
		})
	})
}
