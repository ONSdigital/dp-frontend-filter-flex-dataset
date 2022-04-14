package handlers

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	gomock "github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

// TestSubmitHandler unit tests
func TestSubmitHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	ctx := gomock.Any()

	Convey("test submit handler", t, func() {
		Convey("test Submit handler, starts a filter-outputs job and redirects on success", func() {
			mockFc := NewMockFilterClient(mockCtrl)
			mockJobStateModel := &filter.GetFilterResponse{
				Dataset: filter.Dataset{
					DatasetID: "5678",
					Edition:   "2021",
					Version:   1,
				},
			}
			mockFilterOutputModel := &filter.SubmitFilterResponse{}
			mockFilterOutputModel.Links.FilterOutputs.ID = "abcde12345"
			mockFc.EXPECT().GetFilter(ctx, gomock.Any()).Return(mockJobStateModel, nil)
			mockFc.EXPECT().SubmitFilter(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockFilterOutputModel, "", nil)

			w := testResponse(http.StatusFound, "/filters/12345/submit", mockFc)

			location := w.Header().Get("Location")
			So(location, ShouldNotBeEmpty)
			So(location, ShouldEqual, "/datasets/5678/editions/2021/versions/1/filter-outputs/abcde12345")
		})

		Convey("test Submit handler returns 500 if unable to get job state", func() {
			mockClient := NewMockFilterClient(mockCtrl)
			mockClient.EXPECT().GetFilter(ctx, gomock.Any()).Return(nil, errors.New("failed to get job state"))
			testResponse(http.StatusInternalServerError, "/filters/12345/submit", mockClient)
		})

		Convey("test Submit handler returns 500 if unable to update flex blueprint", func() {
			mockFc := NewMockFilterClient(mockCtrl)
			mockJobStateModel := &filter.GetFilterResponse{
				Dataset: filter.Dataset{
					DatasetID: "5678",
					Edition:   "2021",
					Version:   1,
				},
			}
			mockFc.EXPECT().GetFilter(ctx, gomock.Any()).Return(mockJobStateModel, nil)
			mockFc.EXPECT().SubmitFilter(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, "", errors.New("failed to submit filter blueprint"))
			testResponse(http.StatusInternalServerError, "/filters/12345/submit", mockFc)
		})
	})
}

func testResponse(code int, url string, fc FilterClient) *httptest.ResponseRecorder {
	req := httptest.NewRequest("POST", url, nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/filters/{filterID}/submit", Submit(fc))
	router.ServeHTTP(w, req)

	So(w.Code, ShouldEqual, code)

	b, err := ioutil.ReadAll(w.Body)
	So(err, ShouldBeNil)
	// Writer body should be empty, we don't write a response
	So(b, ShouldBeEmpty)

	return w
}
