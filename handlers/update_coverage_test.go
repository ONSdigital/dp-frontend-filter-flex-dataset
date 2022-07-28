package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUpdateCoverageHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	Convey("Update coverage", t, func() {
		stubFormData := url.Values{}
		stubFormData.Add("dimension", "geography")
		stubFormData.Add("option", "0")

		Convey("Given a valid geography", func() {
			Convey("When the user is redirected to the get coverage screen", func() {
				const filterID = "1234"

				filterClient := NewMockFilterClient(mockCtrl)
				filterClient.
					EXPECT().
					AddDimensionValue(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", nil)

				w := runUpdateCoverage(filterID, "geography", stubFormData, UpdateCoverage(filterClient))

				Convey("Then the location header should match the get coverage screen", func() {
					So(w.Header().Get("Location"), ShouldEqual, fmt.Sprintf("/filters/%s/dimensions/geography/coverage", filterID))
				})

				Convey("And the status code should be 301", func() {
					So(w.Code, ShouldEqual, http.StatusMovedPermanently)
				})
			})

			Convey("When the filter API client responds with an error", func() {
				filterClient := NewMockFilterClient(mockCtrl)
				filterClient.
					EXPECT().
					AddDimensionValue(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", errors.New("internal error"))

				w := runUpdateCoverage("test", "test", stubFormData, UpdateCoverage(filterClient))

				Convey("Then the client should not be redirected", func() {
					So(w.Header().Get("Location"), ShouldBeEmpty)
				})

				Convey("And the status code should be 500", func() {
					So(w.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})
		})

		Convey("Given an invalid request", func() {
			Convey("When the request is missing the hidden required form values", func() {
				tests := map[string]url.Values{
					"Missing dimension": {"option": []string{"0"}},
					"Missing option":    {"dimension": []string{"geography"}},
				}

				for name, formData := range tests {
					Convey(name, func() {
						w := runUpdateCoverage("test", "test", formData, UpdateCoverage(NewMockFilterClient(mockCtrl)))

						Convey("Then the client should not be redirected", func() {
							So(w.Header().Get("Location"), ShouldBeEmpty)
						})

						Convey("And the status code should be 400", func() {
							So(w.Code, ShouldEqual, http.StatusBadRequest)
						})
					})
				}
			})
		})
	})
}

func runUpdateCoverage(filterID, dimension string, formData url.Values, handler http.HandlerFunc) *httptest.ResponseRecorder {
	encodedFormData := formData.Encode()
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/filters/%s/dimensions/geography/coverage", filterID), strings.NewReader(encodedFormData))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(encodedFormData)))

	w := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/filters/{filterID}/dimensions/geography/coverage", handler)
	router.ServeHTTP(w, req)

	return w
}
