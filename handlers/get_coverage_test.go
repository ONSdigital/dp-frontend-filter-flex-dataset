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

				mockPc := NewMockPopulationClient(mockCtrl)

				router := mux.NewRouter()
				router.HandleFunc("/filters/12345/dimensions/geography/coverage", GetCoverage(mockRend, mockFc, mockPc))
				router.ServeHTTP(w, req)

				Convey("And the status code should be 200", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
				})
			})

			Convey("When the user performs a search", func() {
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

				mockPc := NewMockPopulationClient(mockCtrl)
				mockPc.
					EXPECT().
					GetAreas(gomock.Any(), gomock.Any()).
					Return(population.GetAreasResponse{}, nil)

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

		Convey("When the GetAreas API call responds with an error", func() {
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

			mockPc := NewMockPopulationClient(mockCtrl)
			mockPc.
				EXPECT().
				GetAreas(gomock.Any(), gomock.Any()).
				Return(population.GetAreasResponse{}, errors.New("sorry"))

			router := mux.NewRouter()
			router.HandleFunc("/filters/12345/dimensions/geography/coverage", GetCoverage(NewMockRenderClient(mockCtrl), mockFc, mockPc))
			router.ServeHTTP(w, req)

			Convey("Then the status code should be 500", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})
}
