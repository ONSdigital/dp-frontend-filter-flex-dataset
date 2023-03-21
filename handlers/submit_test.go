package handlers

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/helpers"
	gomock "github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

// TestSubmitHandler unit tests
func TestSubmitHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	ctx := gomock.Any()
	cfg := initialiseMockConfig()

	Convey("test submit handler", t, func() {
		Convey("test Submit handler, starts a filter-outputs job and redirects to specific filter-outputs page on success", func() {
			mockFc := NewMockFilterClient(mockCtrl)
			mockFilter := &filter.GetFilterResponse{
				Dataset: filter.Dataset{
					DatasetID: "5678",
					Edition:   "2021",
					Version:   1,
				},
			}
			mockFilterResp := &filter.SubmitFilterResponse{}
			mockFilterResp.FilterOutputID = "abcde12345"
			mockFc.EXPECT().GetFilter(ctx, gomock.Any()).Return(mockFilter, nil)
			mockFc.EXPECT().SubmitFilter(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockFilterResp, "", nil)

			ff := NewFilterFlex(
				NewMockRenderClient(mockCtrl),
				mockFc,
				NewMockDatasetClient(mockCtrl),
				NewMockPopulationClient(mockCtrl),
				NewMockZebedeeClient(mockCtrl),
				cfg)
			w := testResponse(http.StatusFound, "/filters/12345/submit", ff)

			location := w.Header().Get("Location")
			So(location, ShouldNotBeEmpty)
			So(location, ShouldEqual, "/datasets/5678/editions/2021/versions/1/filter-outputs/abcde12345#get-data")
		})

		Convey("test Submit handler, starts a filter-outputs job and redirects to /datasets/create filter-outputs page on success if custom is true", func() {
			mockFc := NewMockFilterClient(mockCtrl)
			mockFilter := &filter.GetFilterResponse{
				Dataset: filter.Dataset{
					DatasetID: "5678",
					Edition:   "2021",
					Version:   1,
				},
				Custom: helpers.ToBoolPtr(true),
			}
			mockFilterResp := &filter.SubmitFilterResponse{}
			mockFilterResp.FilterOutputID = "abcde12345"
			mockFc.EXPECT().GetFilter(ctx, gomock.Any()).Return(mockFilter, nil)
			mockFc.EXPECT().SubmitFilter(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockFilterResp, "", nil)

			ff := NewFilterFlex(
				NewMockRenderClient(mockCtrl),
				mockFc,
				NewMockDatasetClient(mockCtrl),
				NewMockPopulationClient(mockCtrl),
				NewMockZebedeeClient(mockCtrl),
				cfg)
			w := testResponse(http.StatusFound, "/filters/12345/submit", ff)

			location := w.Header().Get("Location")
			So(location, ShouldNotBeEmpty)
			So(location, ShouldEqual, "/datasets/create/filter-outputs/abcde12345#get-data")
		})

		Convey("test Submit handler returns 500 if unable to get job state", func() {
			mockFc := NewMockFilterClient(mockCtrl)
			mockFc.EXPECT().GetFilter(ctx, gomock.Any()).Return(nil, errors.New("failed to get job state"))

			ff := NewFilterFlex(
				NewMockRenderClient(mockCtrl),
				mockFc,
				NewMockDatasetClient(mockCtrl),
				NewMockPopulationClient(mockCtrl),
				NewMockZebedeeClient(mockCtrl),
				cfg)

			testResponse(http.StatusInternalServerError, "/filters/12345/submit", ff)
		})

		Convey("test Submit handler returns 500 if unable to update flex blueprint", func() {
			mockFc := NewMockFilterClient(mockCtrl)
			mockFilter := &filter.GetFilterResponse{
				Dataset: filter.Dataset{
					DatasetID: "5678",
					Edition:   "2021",
					Version:   1,
				},
			}
			mockFc.EXPECT().GetFilter(ctx, gomock.Any()).Return(mockFilter, nil)
			mockFc.EXPECT().SubmitFilter(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, "", errors.New("failed to submit filter blueprint"))

			ff := NewFilterFlex(
				NewMockRenderClient(mockCtrl),
				mockFc,
				NewMockDatasetClient(mockCtrl),
				NewMockPopulationClient(mockCtrl),
				NewMockZebedeeClient(mockCtrl),
				cfg)
			testResponse(http.StatusInternalServerError, "/filters/12345/submit", ff)
		})
	})
}

func testResponse(code int, url string, ff *FilterFlex) *httptest.ResponseRecorder {
	req := httptest.NewRequest("POST", url, nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/filters/{filterID}/submit", ff.Submit())
	router.ServeHTTP(w, req)

	So(w.Code, ShouldEqual, code)

	b, err := ioutil.ReadAll(w.Body)
	So(err, ShouldBeNil)
	// Writer body should be empty, we don't write a response
	So(b, ShouldBeEmpty)

	return w
}
