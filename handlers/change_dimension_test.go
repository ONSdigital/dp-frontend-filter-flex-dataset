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

	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/helpers"
	gomock "github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

func TestChangeDimensionHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	cfg := initialiseMockConfig()

	Convey("Change dimension", t, func() {
		stubFormData := url.Values{}
		stubFormData.Add("dimension", "country")
		stubFormData.Add("is_area_type", "true")

		filterClient := NewMockFilterClient(mockCtrl)
		ff := NewFilterFlex(
			NewMockRenderClient(mockCtrl),
			filterClient,
			NewMockDatasetClient(mockCtrl),
			NewMockPopulationClient(mockCtrl),
			NewMockZebedeeClient(mockCtrl),
			cfg)

		Convey("Given a valid dimension", func() {
			Convey("When the user is redirected to the dimensions review screen", func() {
				const filterID = "1234"

				filterClient.
					EXPECT().
					UpdateDimensions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimension{}, "", nil).
					AnyTimes()

				filterClient.
					EXPECT().
					GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimension{}, "", nil)

				w := runChangeDimension(filterID, "city", stubFormData, ff.ChangeDimension())

				Convey("Then the location header should match the review screen", func() {
					So(w.Header().Get("Location"), ShouldEqual, fmt.Sprintf("/filters/%s/dimensions", filterID))
				})

				Convey("And the status code should be 301", func() {
					So(w.Code, ShouldEqual, http.StatusMovedPermanently)
				})
			})

			Convey("When the filter client's `UpdateDimensions` method is called, it is passed the new dimension", func() {
				const filterID = "1234"
				const dimensionName = "geography"
				const newDimension = "country"

				expDimension := filter.Dimension{
					Name:       newDimension,
					ID:         newDimension,
					IsAreaType: helpers.ToBoolPtr(true),
				}

				filterClient.
					EXPECT().
					UpdateDimensions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), filterID, dimensionName, gomock.Any(), gomock.Eq(expDimension)).
					Return(filter.Dimension{}, "", nil)
				filterClient.
					EXPECT().
					GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimension{}, "", nil)

				formData := url.Values{}
				formData.Add("dimension", newDimension)
				formData.Add("is_area_type", "true")

				runChangeDimension(filterID, dimensionName, formData, ff.ChangeDimension())
			})

			Convey("When the filter.UpdateDimensions responds with an error", func() {
				filterClient.
					EXPECT().
					UpdateDimensions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimension{}, "", errors.New("internal error"))
				filterClient.
					EXPECT().
					GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimension{}, "", nil)

				w := runChangeDimension("test", "test", stubFormData, ff.ChangeDimension())

				Convey("Then the client should not be redirected", func() {
					So(w.Header().Get("Location"), ShouldBeEmpty)
				})

				Convey("And the status code should be 500", func() {
					So(w.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})

			Convey("When the filter.GetDimension responds with an error", func() {
				filterClient.
					EXPECT().
					GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimension{}, "", errors.New("internal error"))

				w := runChangeDimension("test", "test", stubFormData, ff.ChangeDimension())

				Convey("Then the client should not be redirected", func() {
					So(w.Header().Get("Location"), ShouldBeEmpty)
				})

				Convey("And the status code should be 500", func() {
					So(w.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})
		})

		Convey("Given an invalid request", func() {
			ff := NewFilterFlex(
				NewMockRenderClient(mockCtrl),
				filterClient,
				NewMockDatasetClient(mockCtrl),
				NewMockPopulationClient(mockCtrl),
				NewMockZebedeeClient(mockCtrl),
				cfg)

			filterClient.
				EXPECT().
				GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(filter.Dimension{}, "", nil)

			Convey("When the area type has not been provided", func() {
				formData := url.Values{}
				formData.Add("is_area_type", "true")

				w := runChangeDimension("test", "test", formData, ff.ChangeDimension())

				Convey("Then the client should be redirected with the error query param", func() {
					location := w.Header().Get("Location")
					So(location, ShouldNotBeEmpty)

					parse, err := url.Parse(location)
					So(err, ShouldBeNil)

					query := parse.Query()
					So(query["error"], ShouldNotBeEmpty)
				})

				Convey("And the status code should be 301", func() {
					So(w.Code, ShouldEqual, http.StatusMovedPermanently)
				})
			})

			Convey("When the request is missing the hidden required form values", func() {
				tests := map[string]url.Values{
					"Missing is_area_type":       {"dimension": []string{"country"}},
					"Invalid is_area_type value": {"dimension": []string{"country"}, "is_area_type": []string{"no"}},
				}

				for name, formData := range tests {
					Convey(name, func() {
						w := runChangeDimension("test", "test", formData, ff.ChangeDimension())

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

func runChangeDimension(filterID, dimension string, formData url.Values, handler http.HandlerFunc) *httptest.ResponseRecorder {
	encodedFormData := formData.Encode()
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/filters/%s/dimensions/%s", filterID, dimension), strings.NewReader(encodedFormData))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(encodedFormData)))

	w := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/filters/{filterID}/dimensions/{name}", handler)
	router.ServeHTTP(w, req)

	return w
}
