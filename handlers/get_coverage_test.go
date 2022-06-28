package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
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

	Convey("Get coverage", t, func() {
		Convey("Given a valid request", func() {
			Convey("When the user is redirected to the change coverage screen", func() {
				w := httptest.NewRecorder()
				req := httptest.NewRequest("GET", "/filters/12345/dimensions/geography/coverage", nil)

				mockRend := NewMockRenderClient(mockCtrl)
				mockRend.EXPECT().NewBasePageModel().Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain))
				mockRend.EXPECT().BuildPage(gomock.Any(), gomock.Any(), "coverage")

				mockFc := NewMockFilterClient(mockCtrl)
				mockFc.EXPECT().GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(filter.Dimension{}, "", nil)

				router := mux.NewRouter()
				router.HandleFunc("/filters/12345/dimensions/geography/coverage", GetCoverage(mockRend, mockFc))
				router.ServeHTTP(w, req)

				Convey("And the status code should be 200", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
				})
			})
		})

		Convey("When the filter API responds with an error", func() {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/filters/12345/dimensions/geography/coverage", nil)

			mockFc := NewMockFilterClient(mockCtrl)
			mockFc.EXPECT().GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(filter.Dimension{}, "", errors.New("sorry"))

			router := mux.NewRouter()
			router.HandleFunc("/filters/12345/dimensions/geography/coverage", GetCoverage(NewMockRenderClient(mockCtrl), mockFc))
			router.ServeHTTP(w, req)

			Convey("Then the status code should be 500", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})
}
