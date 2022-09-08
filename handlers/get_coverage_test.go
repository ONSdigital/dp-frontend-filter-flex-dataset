package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

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

func TestGetCoverageHandler(t *testing.T) {
	helper.InitialiseLocalisationsHelper(mocks.MockAssetFunction)
	mockCtrl := gomock.NewController(t)
	cfg := initialiseMockConfig()
	mockFilterDims := filter.Dimensions{
		Items: []filter.Dimension{
			{
				Name:       "Test",
				IsAreaType: new(bool),
			},
			{
				Name:       "Test 2",
				IsAreaType: helpers.ToBoolPtr(true),
				ID:         "city",
			},
		},
	}
	mockParentFilterDims := filter.Dimensions{
		Items: []filter.Dimension{
			{
				Name:       "Test",
				IsAreaType: new(bool),
			},
			{
				Name:           "Test 2",
				IsAreaType:     helpers.ToBoolPtr(true),
				FilterByParent: "country",
				ID:             "city",
			},
		},
	}

	Convey("Get coverage", t, func() {
		Convey("Given a valid request", func() {
			Convey("When the user is redirected to the change coverage screen", func() {
				w := httptest.NewRecorder()
				req := httptest.NewRequest("GET", "/filters/12345/dimensions/geography/coverage", nil)

				mockRend := NewMockRenderClient(mockCtrl)
				mockRend.
					EXPECT().
					NewBasePageModel().
					Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain))
				mockRend.
					EXPECT().
					BuildPage(gomock.Any(), gomock.Any(), "coverage")

				mockFc := NewMockFilterClient(mockCtrl)
				mockFc.
					EXPECT().
					GetDimensions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(mockFilterDims, "", nil)
				mockFc.
					EXPECT().
					GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), mockFilterDims.Items[0].Name).
					Return(mockFilterDims.Items[0], "", nil)
				mockFc.
					EXPECT().
					GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), mockFilterDims.Items[1].Name).
					Return(mockFilterDims.Items[1], "", nil)
				mockFc.EXPECT().
					GetFilter(gomock.Any(), gomock.Any()).
					Return(&filter.GetFilterResponse{}, nil)
				mockFc.EXPECT().
					GetDimensionOptions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.DimensionOptions{}, "", nil)

				mockPc := NewMockPopulationClient(mockCtrl)
				mockPc.EXPECT().
					GetAreaTypeParents(gomock.Any(), gomock.Any()).
					Return(population.GetAreaTypeParentsResponse{}, nil)

				router := mux.NewRouter()
				router.HandleFunc("/filters/12345/dimensions/geography/coverage", GetCoverage(mockRend, mockFc, mockPc))
				router.ServeHTTP(w, req)

				Convey("And the status code should be 200", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
				})
			})

			Convey("When the user performs a name search", func() {
				w := httptest.NewRecorder()
				req := httptest.NewRequest("GET", "/filters/12345/dimensions/geography/coverage?q=name", nil)

				mockRend := NewMockRenderClient(mockCtrl)
				mockRend.
					EXPECT().
					NewBasePageModel().
					Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain))
				mockRend.
					EXPECT().
					BuildPage(gomock.Any(), gomock.Any(), "coverage")

				mockFc := NewMockFilterClient(mockCtrl)
				mockFc.
					EXPECT().
					GetDimensions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(mockFilterDims, "", nil)
				mockFc.
					EXPECT().
					GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), mockFilterDims.Items[0].Name).
					Return(mockFilterDims.Items[0], "", nil)
				mockFc.
					EXPECT().
					GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), mockFilterDims.Items[1].Name).
					Return(mockFilterDims.Items[1], "", nil)
				mockFc.EXPECT().
					GetFilter(gomock.Any(), gomock.Any()).
					Return(&filter.GetFilterResponse{}, nil)
				mockFc.EXPECT().
					GetDimensionOptions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.DimensionOptions{}, "", nil)

				mockPc := NewMockPopulationClient(mockCtrl)
				mockPc.
					EXPECT().
					GetAreas(gomock.Any(), gomock.Any()).
					Return(population.GetAreasResponse{}, nil)
				mockPc.EXPECT().
					GetAreaTypeParents(gomock.Any(), gomock.Any()).
					Return(population.GetAreaTypeParentsResponse{}, nil)

				router := mux.NewRouter()
				router.HandleFunc("/filters/12345/dimensions/geography/coverage", GetCoverage(mockRend, mockFc, mockPc))
				router.ServeHTTP(w, req)

				Convey("And the status code should be 200", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
				})
			})

			Convey("When the user performs a parent search", func() {
				w := httptest.NewRecorder()
				req := httptest.NewRequest("GET", "/filters/12345/dimensions/geography/coverage?p=parent&pq=name", nil)

				mockRend := NewMockRenderClient(mockCtrl)
				mockRend.
					EXPECT().
					NewBasePageModel().
					Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain))
				mockRend.
					EXPECT().
					BuildPage(gomock.Any(), gomock.Any(), "coverage")

				mockFc := NewMockFilterClient(mockCtrl)
				mockFc.
					EXPECT().
					GetDimensions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(mockParentFilterDims, "", nil)
				mockFc.
					EXPECT().
					GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), mockParentFilterDims.Items[0].Name).
					Return(mockParentFilterDims.Items[0], "", nil)
				mockFc.
					EXPECT().
					GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), mockParentFilterDims.Items[1].Name).
					Return(mockParentFilterDims.Items[1], "", nil)
				mockFc.EXPECT().
					GetFilter(gomock.Any(), gomock.Any()).
					Return(&filter.GetFilterResponse{}, nil)
				mockFc.EXPECT().
					GetDimensionOptions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.DimensionOptions{}, "", nil)

				mockPc := NewMockPopulationClient(mockCtrl)
				mockPc.
					EXPECT().
					GetAreas(gomock.Any(), gomock.Any()).
					Return(population.GetAreasResponse{}, nil)
				mockPc.EXPECT().
					GetAreaTypeParents(gomock.Any(), gomock.Any()).
					Return(population.GetAreaTypeParentsResponse{}, nil)

				router := mux.NewRouter()
				router.HandleFunc("/filters/12345/dimensions/geography/coverage", GetCoverage(mockRend, mockFc, mockPc))
				router.ServeHTTP(w, req)

				Convey("And the status code should be 200", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
				})
			})

			Convey("When the user has saved options", func() {
				w := httptest.NewRecorder()
				req := httptest.NewRequest("GET", "/filters/12345/dimensions/geography/coverage", nil)

				mockRend := NewMockRenderClient(mockCtrl)
				mockRend.
					EXPECT().
					NewBasePageModel().
					Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain))
				mockRend.
					EXPECT().
					BuildPage(gomock.Any(), gomock.Any(), "coverage")

				mockFc := NewMockFilterClient(mockCtrl)
				mockFc.
					EXPECT().
					GetDimensions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(mockFilterDims, "", nil)
				mockFc.
					EXPECT().
					GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), mockFilterDims.Items[0].Name).
					Return(mockFilterDims.Items[0], "", nil)
				mockFc.
					EXPECT().
					GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), mockFilterDims.Items[1].Name).
					Return(mockFilterDims.Items[1], "", nil)
				mockFc.EXPECT().
					GetFilter(gomock.Any(), gomock.Any()).
					Return(&filter.GetFilterResponse{}, nil)
				mockFc.EXPECT().
					GetDimensionOptions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.DimensionOptions{
						Items: []filter.DimensionOption{
							{
								Option: "Option 1",
							},
						}}, "", nil)

				mockPc := NewMockPopulationClient(mockCtrl)
				mockPc.
					EXPECT().
					GetAreas(gomock.Any(), gomock.Any()).
					Return(population.GetAreasResponse{
						Areas: []population.Area{
							{
								ID:       "1",
								Label:    "Label",
								AreaType: "Geography",
							},
						},
					}, nil)
				mockPc.EXPECT().
					GetAreaTypeParents(gomock.Any(), gomock.Any()).
					Return(population.GetAreaTypeParentsResponse{}, nil)

				router := mux.NewRouter()
				router.HandleFunc("/filters/12345/dimensions/geography/coverage", GetCoverage(mockRend, mockFc, mockPc))
				router.ServeHTTP(w, req)

				Convey("And the status code should be 200", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
				})
			})

			Convey("When the user has saved parent options", func() {
				w := httptest.NewRecorder()
				req := httptest.NewRequest("GET", "/filters/12345/dimensions/geography/coverage", nil)

				mockRend := NewMockRenderClient(mockCtrl)
				mockRend.
					EXPECT().
					NewBasePageModel().
					Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain))
				mockRend.
					EXPECT().
					BuildPage(gomock.Any(), gomock.Any(), "coverage")

				mockFc := NewMockFilterClient(mockCtrl)
				mockFc.
					EXPECT().
					GetDimensions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(mockParentFilterDims, "", nil)
				mockFc.
					EXPECT().
					GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), mockParentFilterDims.Items[0].Name).
					Return(mockParentFilterDims.Items[0], "", nil)
				mockFc.
					EXPECT().
					GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), mockParentFilterDims.Items[1].Name).
					Return(mockParentFilterDims.Items[1], "", nil)
				mockFc.EXPECT().
					GetFilter(gomock.Any(), gomock.Any()).
					Return(&filter.GetFilterResponse{}, nil)
				mockFc.EXPECT().
					GetDimensionOptions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.DimensionOptions{
						Items: []filter.DimensionOption{
							{
								Option: "parent 1",
							},
						}}, "", nil)

				mockPc := NewMockPopulationClient(mockCtrl)
				mockPc.
					EXPECT().
					GetAreas(gomock.Any(), gomock.Any()).
					Return(population.GetAreasResponse{
						Areas: []population.Area{
							{
								ID:       "1",
								Label:    "Label",
								AreaType: "Geography",
							},
						},
					}, nil)
				mockPc.EXPECT().
					GetAreaTypeParents(gomock.Any(), gomock.Any()).
					Return(population.GetAreaTypeParentsResponse{}, nil)

				router := mux.NewRouter()
				router.HandleFunc("/filters/12345/dimensions/geography/coverage", GetCoverage(mockRend, mockFc, mockPc))
				router.ServeHTTP(w, req)

				Convey("And the status code should be 200", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
				})
			})
		})

		Convey("When the GetFilter API call responds with an error", func() {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/filters/12345/dimensions/geography/coverage", nil)

			mockFc := NewMockFilterClient(mockCtrl)
			mockFc.EXPECT().
				GetFilter(gomock.Any(), gomock.Any()).
				Return(&filter.GetFilterResponse{}, errors.New("sorry"))
			mockFc.
				EXPECT().
				GetDimensions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(filter.Dimensions{}, "", nil)

			router := mux.NewRouter()
			router.HandleFunc("/filters/12345/dimensions/geography/coverage", GetCoverage(NewMockRenderClient(mockCtrl), mockFc, NewMockPopulationClient(mockCtrl)))
			router.ServeHTTP(w, req)

			Convey("Then the status code should be 500", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})

		Convey("When the GetDimensions API call responds with an error", func() {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/filters/12345/dimensions/geography/coverage", nil)

			mockFc := NewMockFilterClient(mockCtrl)
			mockFc.EXPECT().
				GetFilter(gomock.Any(), gomock.Any()).
				Return(&filter.GetFilterResponse{}, nil)
			mockFc.
				EXPECT().
				GetDimensions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(filter.Dimensions{}, "", errors.New("sorry"))

			router := mux.NewRouter()
			router.HandleFunc("/filters/12345/dimensions/geography/coverage", GetCoverage(NewMockRenderClient(mockCtrl), mockFc, NewMockPopulationClient(mockCtrl)))
			router.ServeHTTP(w, req)

			Convey("Then the status code should be 500", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})

		Convey("When the subsequent GetDimension API call responds with an error", func() {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/filters/12345/dimensions/geography/coverage", nil)

			mockFc := NewMockFilterClient(mockCtrl)
			mockFc.EXPECT().
				GetFilter(gomock.Any(), gomock.Any()).
				Return(&filter.GetFilterResponse{}, nil)
			mockFc.
				EXPECT().
				GetDimensions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(mockFilterDims, "", nil)
			mockFc.
				EXPECT().
				GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), mockFilterDims.Items[0].Name).
				Return(filter.Dimension{}, "", errors.New("sorry"))

			router := mux.NewRouter()
			router.HandleFunc("/filters/12345/dimensions/geography/coverage", GetCoverage(NewMockRenderClient(mockCtrl), mockFc, NewMockPopulationClient(mockCtrl)))
			router.ServeHTTP(w, req)

			Convey("Then the status code should be 500", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})

		Convey("When the GetDimensionOptions API call responds with an error", func() {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/filters/12345/dimensions/geography/coverage", nil)

			mockFc := NewMockFilterClient(mockCtrl)
			mockFc.
				EXPECT().
				GetDimensions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(mockFilterDims, "", nil)
			mockFc.
				EXPECT().
				GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), mockFilterDims.Items[0].Name).
				Return(mockFilterDims.Items[0], "", nil)
			mockFc.
				EXPECT().
				GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), mockFilterDims.Items[1].Name).
				Return(mockFilterDims.Items[1], "", nil)
			mockFc.EXPECT().
				GetFilter(gomock.Any(), gomock.Any()).
				Return(&filter.GetFilterResponse{}, nil)
			mockFc.EXPECT().
				GetDimensionOptions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(filter.DimensionOptions{}, "", errors.New("sorry"))

			mockPc := NewMockPopulationClient(mockCtrl)
			mockPc.EXPECT().
				GetAreaTypeParents(gomock.Any(), gomock.Any()).
				Return(population.GetAreaTypeParentsResponse{}, nil)

			router := mux.NewRouter()
			router.HandleFunc("/filters/12345/dimensions/geography/coverage", GetCoverage(NewMockRenderClient(mockCtrl), mockFc, mockPc))
			router.ServeHTTP(w, req)

			Convey("Then the status code should be 500", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})

		Convey("When the GetAreaTypeParents API call responds with an error", func() {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/filters/12345/dimensions/geography/coverage", nil)

			mockFc := NewMockFilterClient(mockCtrl)
			mockFc.
				EXPECT().
				GetDimensions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(mockFilterDims, "", nil)
			mockFc.
				EXPECT().
				GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), mockFilterDims.Items[0].Name).
				Return(mockFilterDims.Items[0], "", nil)
			mockFc.
				EXPECT().
				GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), mockFilterDims.Items[1].Name).
				Return(mockFilterDims.Items[1], "", nil)
			mockFc.EXPECT().
				GetFilter(gomock.Any(), gomock.Any()).
				Return(&filter.GetFilterResponse{}, nil)
			mockFc.EXPECT().
				GetDimensionOptions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(filter.DimensionOptions{}, "", nil)

			mockPc := NewMockPopulationClient(mockCtrl)
			mockPc.EXPECT().
				GetAreaTypeParents(gomock.Any(), gomock.Any()).
				Return(population.GetAreaTypeParentsResponse{}, errors.New("sorry"))

			router := mux.NewRouter()
			router.HandleFunc("/filters/12345/dimensions/geography/coverage", GetCoverage(NewMockRenderClient(mockCtrl), mockFc, mockPc))
			router.ServeHTTP(w, req)

			Convey("Then the status code should be 500", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})

		Convey("When the GetAreas API call via the options loop responds with an error", func() {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/filters/12345/dimensions/geography/coverage?q=test", nil)

			mockFc := NewMockFilterClient(mockCtrl)
			mockFc.
				EXPECT().
				GetDimensions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(mockFilterDims, "", nil)
			mockFc.
				EXPECT().
				GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), mockFilterDims.Items[0].Name).
				Return(mockFilterDims.Items[0], "", nil)
			mockFc.
				EXPECT().
				GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), mockFilterDims.Items[1].Name).
				Return(mockFilterDims.Items[1], "", nil)
			mockFc.EXPECT().
				GetFilter(gomock.Any(), gomock.Any()).
				Return(&filter.GetFilterResponse{}, nil)
			mockFc.EXPECT().
				GetDimensionOptions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(filter.DimensionOptions{
					Items: []filter.DimensionOption{
						{
							Option: "option 1",
						},
					},
				}, "", nil)

			mockPc := NewMockPopulationClient(mockCtrl)
			mockPc.
				EXPECT().
				GetAreas(gomock.Any(), gomock.Any()).
				Return(population.GetAreasResponse{}, errors.New("sorry"))
			mockPc.EXPECT().
				GetAreaTypeParents(gomock.Any(), gomock.Any()).
				Return(population.GetAreaTypeParentsResponse{}, nil)

			router := mux.NewRouter()
			router.HandleFunc("/filters/12345/dimensions/geography/coverage", GetCoverage(NewMockRenderClient(mockCtrl), mockFc, mockPc))
			router.ServeHTTP(w, req)

			Convey("Then the status code should be 500", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})

		Convey("When the GetAreas API call via the name search responds with an error", func() {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/filters/12345/dimensions/geography/coverage?q=test", nil)

			mockFc := NewMockFilterClient(mockCtrl)
			mockFc.
				EXPECT().
				GetDimensions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(mockFilterDims, "", nil)
			mockFc.
				EXPECT().
				GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), mockFilterDims.Items[0].Name).
				Return(mockFilterDims.Items[0], "", nil)
			mockFc.
				EXPECT().
				GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), mockFilterDims.Items[1].Name).
				Return(mockFilterDims.Items[1], "", nil)
			mockFc.EXPECT().
				GetFilter(gomock.Any(), gomock.Any()).
				Return(&filter.GetFilterResponse{}, nil)
			mockFc.EXPECT().
				GetDimensionOptions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(filter.DimensionOptions{}, "", nil)

			mockPc := NewMockPopulationClient(mockCtrl)
			mockPc.
				EXPECT().
				GetAreas(gomock.Any(), gomock.Any()).
				Return(population.GetAreasResponse{}, errors.New("sorry"))
			mockPc.EXPECT().
				GetAreaTypeParents(gomock.Any(), gomock.Any()).
				Return(population.GetAreaTypeParentsResponse{}, nil)

			router := mux.NewRouter()
			router.HandleFunc("/filters/12345/dimensions/geography/coverage", GetCoverage(NewMockRenderClient(mockCtrl), mockFc, mockPc))
			router.ServeHTTP(w, req)

			Convey("Then the status code should be 500", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})

		Convey("When the GetAreas API call via the parent search responds with an error", func() {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/filters/12345/dimensions/geography/coverage?p=country&pq=test", nil)

			mockFc := NewMockFilterClient(mockCtrl)
			mockFc.
				EXPECT().
				GetDimensions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(mockFilterDims, "", nil)
			mockFc.
				EXPECT().
				GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), mockFilterDims.Items[0].Name).
				Return(mockFilterDims.Items[0], "", nil)
			mockFc.
				EXPECT().
				GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), mockFilterDims.Items[1].Name).
				Return(mockFilterDims.Items[1], "", nil)
			mockFc.EXPECT().
				GetFilter(gomock.Any(), gomock.Any()).
				Return(&filter.GetFilterResponse{}, nil)
			mockFc.EXPECT().
				GetDimensionOptions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(filter.DimensionOptions{}, "", nil)

			mockPc := NewMockPopulationClient(mockCtrl)
			mockPc.
				EXPECT().
				GetAreas(gomock.Any(), gomock.Any()).
				Return(population.GetAreasResponse{}, errors.New("sorry"))
			mockPc.EXPECT().
				GetAreaTypeParents(gomock.Any(), gomock.Any()).
				Return(population.GetAreaTypeParentsResponse{}, nil)

			router := mux.NewRouter()
			router.HandleFunc("/filters/12345/dimensions/geography/coverage", GetCoverage(NewMockRenderClient(mockCtrl), mockFc, mockPc))
			router.ServeHTTP(w, req)

			Convey("Then the status code should be 500", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})
}
