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

func TestPostChangeDimensionsHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	Convey("Post change dimensions", t, func() {
		Convey("Given a valid search request", func() {
			stubFormData := url.Values{}
			stubFormData.Add("dimensions", "search")
			stubFormData.Add("q", "dimension")
			stubFormData.Add("is-search", "true")

			Convey("When the user is redirected to the get change dimensions screen", func() {
				const filterID = "1234"

				filterClient := NewMockFilterClient(mockCtrl)
				w := runPostChangeDimensions(filterID, stubFormData, PostChangeDimensions(filterClient))

				Convey("Then the location header should match the get change dimensions screen with query persisted", func() {
					So(w.Header().Get("Location"), ShouldEqual, fmt.Sprintf("/filters/%s/dimensions/change?f=search&q=dimension", filterID))
				})

				Convey("And the status code should be 303", func() {
					So(w.Code, ShouldEqual, http.StatusSeeOther)
				})
			})
		})

		Convey("Given a valid add dimension request", func() {
			stubFormData := url.Values{}
			stubFormData.Add("dimensions", "browse")
			stubFormData.Add("add-dimension", "age")

			Convey("When the user is redirected to the get change dimensions screen", func() {
				const fid = "1234"

				mockFc := NewMockFilterClient(mockCtrl)
				mockFc.
					EXPECT().
					AddFlexDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", nil)

				w := runPostChangeDimensions(fid, stubFormData, PostChangeDimensions(mockFc))

				Convey("Then the location header should match the get change dimensions screen with form action persisted", func() {
					So(w.Header().Get("Location"), ShouldEqual, fmt.Sprintf("/filters/%s/dimensions/change?f=browse", fid))
				})

				Convey("And the status code should be 303", func() {
					So(w.Code, ShouldEqual, http.StatusSeeOther)
				})
			})
		})

		Convey("Given an invalid add dimension request", func() {
			stubFormData := url.Values{}
			stubFormData.Add("dimensions", "browse")
			stubFormData.Add("add-dimension", "age")

			Convey("When the fc.AddFlexDimension api responds with an error", func() {
				const fid = "1234"

				mockFc := NewMockFilterClient(mockCtrl)
				mockFc.
					EXPECT().
					AddFlexDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", errors.New("Internal error"))

				w := runPostChangeDimensions(fid, stubFormData, PostChangeDimensions(mockFc))

				Convey("Then the client should not be redirected", func() {
					So(w.Header().Get("Location"), ShouldBeEmpty)
				})

				Convey("And the status code should be 500", func() {
					So(w.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})
		})

		Convey("Given a valid remove dimension request", func() {
			stubFormData := url.Values{}
			stubFormData.Add("dimensions", "browse")
			stubFormData.Add("delete-option", "age")

			Convey("When the user is redirected to the get change dimensions screen", func() {
				const fid = "1234"

				mockFc := NewMockFilterClient(mockCtrl)
				mockFc.
					EXPECT().
					RemoveDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", nil)

				w := runPostChangeDimensions(fid, stubFormData, PostChangeDimensions(mockFc))

				Convey("Then the location header should match the get change dimensions screen with form action persisted", func() {
					So(w.Header().Get("Location"), ShouldEqual, fmt.Sprintf("/filters/%s/dimensions/change?f=browse", fid))
				})

				Convey("And the status code should be 303", func() {
					So(w.Code, ShouldEqual, http.StatusSeeOther)
				})
			})
		})

		Convey("Given an invalid remove dimension request", func() {
			stubFormData := url.Values{}
			stubFormData.Add("dimensions", "browse")
			stubFormData.Add("delete-option", "age")

			Convey("When the fc.RemoveDimension api responds with an error", func() {
				const fid = "1234"

				mockFc := NewMockFilterClient(mockCtrl)
				mockFc.
					EXPECT().
					RemoveDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", errors.New("Internal error"))

				w := runPostChangeDimensions(fid, stubFormData, PostChangeDimensions(mockFc))

				Convey("Then the client should not be redirected", func() {
					So(w.Header().Get("Location"), ShouldBeEmpty)
				})

				Convey("And the status code should be 500", func() {
					So(w.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})
		})

		Convey("Given an invalid request", func() {
			Convey("When the request is missing the hidden required form value", func() {
				stubFormData := url.Values{}
				Convey("Missing dimensions", func() {
					w := runUpdateCoverage("test", "test", stubFormData, PostChangeDimensions(NewMockFilterClient(mockCtrl)))

					Convey("Then the client should not be redirected", func() {
						So(w.Header().Get("Location"), ShouldBeEmpty)
					})

					Convey("And the status code should be 400", func() {
						So(w.Code, ShouldEqual, http.StatusBadRequest)
					})
				})
			})
		})
	})
}

func runPostChangeDimensions(filterID string, formData url.Values, handler http.HandlerFunc) *httptest.ResponseRecorder {
	encodedFormData := formData.Encode()
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/filters/%s/dimensions/change", filterID), strings.NewReader(encodedFormData))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(encodedFormData)))

	w := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/filters/{filterID}/dimensions/change", handler)
	router.ServeHTTP(w, req)

	return w
}
