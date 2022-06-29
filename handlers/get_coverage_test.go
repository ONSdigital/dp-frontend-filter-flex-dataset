package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

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

				rend := NewMockRenderClient(mockCtrl)
				rend.EXPECT().NewBasePageModel().Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain))
				rend.EXPECT().BuildPage(gomock.Any(), gomock.Any(), "coverage")

				router := mux.NewRouter()
				router.HandleFunc("/filters/12345/dimensions/geography/coverage", GetCoverage(rend))
				router.ServeHTTP(w, req)

				Convey("And the status code should be 200", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
				})
			})
		})
	})
}
