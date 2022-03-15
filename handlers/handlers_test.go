package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/config"
	coreModel "github.com/ONSdigital/dp-renderer/model"
	gomock "github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

type testCliError struct{}

func (e *testCliError) Error() string { return "client error" }
func (e *testCliError) Code() int     { return http.StatusNotFound }

func TestUnitHandlers(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	cfg := initialiseMockConfig()
	ctx := gomock.Any()
	var mockServiceAuthToken, mockDownloadToken, mockUserAuthToken, mockCollectionID, mockFilterID string
	mockFm := []filter.Model{
		{
			Dimensions: []filter.ModelDimension{
				{
					Name:       "Dim 1",
					IsAreaType: new(bool),
					Options:    []string{"option 1", "option 2"},
				},
			},
		},
		{
			Dimensions: []filter.ModelDimension{
				{
					Name:       "Dim 2",
					IsAreaType: new(bool),
				},
			},
		},
	}
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

	Convey("test setStatusCode", t, func() {

		Convey("test status code handles 404 response from client", func() {
			req := httptest.NewRequest("GET", "http://localhost:20100", nil)
			w := httptest.NewRecorder()
			err := &testCliError{}

			setStatusCode(req, w, err)

			So(w.Code, ShouldEqual, http.StatusNotFound)
		})

		Convey("test status code handles internal server error", func() {
			req := httptest.NewRequest("GET", "http://localhost:20100", nil)
			w := httptest.NewRecorder()
			err := errors.New("internal server error")

			setStatusCode(req, w, err)

			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})
	})

	Convey("test filter flex overview", t, func() {
		Convey("test filter flex overview page is successful", func() {
			mockRend := NewMockRenderClient(mockCtrl)
			mockDc := NewMockDatasetClient(mockCtrl)
			mockFc := NewMockFilterClient(mockCtrl)

			Convey("options on filter job no additional call to get options", func() {
				mockRend.EXPECT().NewBasePageModel().Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain))
				mockRend.EXPECT().BuildPage(gomock.Any(), gomock.Any(), "overview")
				mockFc.EXPECT().GetJobState(ctx, mockUserAuthToken, mockServiceAuthToken, mockDownloadToken, mockCollectionID, mockFilterID).Return(mockFm[0], "", nil)

				w := httptest.NewRecorder()
				req := httptest.NewRequest("GET", "/filters/12345/dimensions", nil)

				router := mux.NewRouter()
				router.HandleFunc("/filters/12345/dimensions", FilterFlexOverview(mockRend, mockFc, mockDc))

				router.ServeHTTP(w, req)

				So(w.Code, ShouldEqual, http.StatusOK)
			})

			Convey("no options on filter job additional call to get options", func() {
				mockRend.EXPECT().NewBasePageModel().Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain))
				mockRend.EXPECT().BuildPage(gomock.Any(), gomock.Any(), "overview")
				mockFc.EXPECT().GetJobState(ctx, mockUserAuthToken, mockServiceAuthToken, mockDownloadToken, mockCollectionID, mockFilterID).Return(mockFm[1], "", nil)
				mockDc.EXPECT().GetOptions(ctx, mockUserAuthToken, mockServiceAuthToken, mockCollectionID, "", "", "0", "Dim 2",
					&dataset.QueryParams{Offset: 0, Limit: 1000}).Return(mockOpts[0], nil)

				w := httptest.NewRecorder()
				req := httptest.NewRequest("GET", "/filters/12345/dimensions", nil)

				router := mux.NewRouter()
				router.HandleFunc("/filters/12345/dimensions", FilterFlexOverview(mockRend, mockFc, mockDc))

				router.ServeHTTP(w, req)

				So(w.Code, ShouldEqual, http.StatusOK)
			})

		})

		Convey("test filter flex overview errors", func() {
			mockRend := NewMockRenderClient(mockCtrl)

			Convey("test FilterFlexOverview returns 500 if client GetJobState returns an error", func() {
				mockFc := NewMockFilterClient(mockCtrl)
				mockDc := NewMockDatasetClient(mockCtrl)
				mockFc.EXPECT().GetJobState(ctx, mockUserAuthToken, mockServiceAuthToken, mockDownloadToken, mockCollectionID, mockFilterID).Return(filter.Model{}, "", errors.New("sorry"))

				w := httptest.NewRecorder()
				req := httptest.NewRequest("GET", "/filters/12345/dimensions", nil)

				router := mux.NewRouter()
				router.HandleFunc("/filters/12345/dimensions", FilterFlexOverview(mockRend, mockFc, mockDc))

				router.ServeHTTP(w, req)

				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})

	Convey("test dimension selector", t, func() {
		Convey("test dimension selector page is successful", func() {
			mockRend := NewMockRenderClient(mockCtrl)
			mockRend.EXPECT().NewBasePageModel().Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain))
			mockRend.EXPECT().BuildPage(gomock.Any(), gomock.Any(), "selector")

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/filters/1234/dimensions/test-name", nil)

			router := mux.NewRouter()
			router.HandleFunc("/filters/1234/dimensions/test-name", DimensionsSelector(mockRend))

			router.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)
		})
	})
}

func initialiseMockConfig() config.Config {
	return config.Config{
		PatternLibraryAssetsPath: "http://localhost:9000/dist",
		SiteDomain:               "ons",
		SupportedLanguages:       []string{"en", "cy"},
	}
}
