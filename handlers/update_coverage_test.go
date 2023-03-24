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
	gomock "github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUpdateCoverageHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	cfg := initialiseMockConfig()

	Convey("Update coverage", t, func() {
		Convey("Given a valid add option request", func() {
			stubFormData := url.Values{}
			stubFormData.Add("dimension", "geography")
			stubFormData.Add("add-option", "0")
			stubFormData.Add("coverage", "name-search")
			stubFormData.Add("geog-id", "city")

			Convey("When the user is redirected to the get coverage screen", func() {
				const filterID = "1234"

				filterClient := NewMockFilterClient(mockCtrl)
				filterClient.
					EXPECT().
					GetDimensionOptions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.DimensionOptions{}, "", nil)
				filterClient.
					EXPECT().
					UpdateDimensions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimension{}, "", nil)

				ff := NewFilterFlex(
					NewMockRenderClient(mockCtrl),
					filterClient,
					NewMockDatasetClient(mockCtrl),
					NewMockPopulationClient(mockCtrl),
					NewMockZebedeeClient(mockCtrl),
					cfg)
				w := runUpdateCoverage(filterID, "geography", stubFormData, ff.UpdateCoverage())

				Convey("Then the location header should match the get coverage screen", func() {
					So(w.Header().Get("Location"), ShouldEqual, fmt.Sprintf("/filters/%s/dimensions/geography/coverage#search--name", filterID))
				})

				Convey("And the status code should be 301", func() {
					So(w.Code, ShouldEqual, http.StatusMovedPermanently)
				})
			})

			Convey("When the GetDimensionOptions filter API client responds with an error", func() {
				filterClient := NewMockFilterClient(mockCtrl)
				filterClient.
					EXPECT().
					GetDimensionOptions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.DimensionOptions{}, "", errors.New("internal error"))

				ff := NewFilterFlex(
					NewMockRenderClient(mockCtrl),
					filterClient,
					NewMockDatasetClient(mockCtrl),
					NewMockPopulationClient(mockCtrl),
					NewMockZebedeeClient(mockCtrl),
					cfg)
				w := runUpdateCoverage("test", "test", stubFormData, ff.UpdateCoverage())

				Convey("Then the client should not be redirected", func() {
					So(w.Header().Get("Location"), ShouldBeEmpty)
				})

				Convey("And the status code should be 500", func() {
					So(w.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})

			Convey("When the UpdateDimensions filter API client responds with an error", func() {
				filterClient := NewMockFilterClient(mockCtrl)
				filterClient.
					EXPECT().
					GetDimensionOptions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.DimensionOptions{}, "", nil)
				filterClient.
					EXPECT().
					UpdateDimensions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimension{}, "", errors.New("internal error"))

				ff := NewFilterFlex(
					NewMockRenderClient(mockCtrl),
					filterClient,
					NewMockDatasetClient(mockCtrl),
					NewMockPopulationClient(mockCtrl),
					NewMockZebedeeClient(mockCtrl),
					cfg)
				w := runUpdateCoverage("test", "test", stubFormData, ff.UpdateCoverage())

				Convey("Then the client should not be redirected", func() {
					So(w.Header().Get("Location"), ShouldBeEmpty)
				})

				Convey("And the status code should be 500", func() {
					So(w.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})
		})

		Convey("Given a valid add parent option request", func() {
			stubFormData := url.Values{}
			stubFormData.Add("dimension", "geography")
			stubFormData.Add("add-parent-option", "0")
			stubFormData.Add("coverage", "parent-search")
			stubFormData.Add("larger-area", "country")
			stubFormData.Add("geog-id", "city")

			Convey("When the user is redirected to the get coverage screen", func() {
				const filterID = "1234"

				filterClient := NewMockFilterClient(mockCtrl)
				filterClient.
					EXPECT().
					GetDimensionOptions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.DimensionOptions{
						Items: []filter.DimensionOption{
							{
								Option: "Option 1",
							},
							{
								Option: "Option 2",
							},
						},
					}, "", nil)
				filterClient.
					EXPECT().
					UpdateDimensions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimension{}, "", nil)

				ff := NewFilterFlex(
					NewMockRenderClient(mockCtrl),
					filterClient,
					NewMockDatasetClient(mockCtrl),
					NewMockPopulationClient(mockCtrl),
					NewMockZebedeeClient(mockCtrl),
					cfg)
				w := runUpdateCoverage(filterID, "geography", stubFormData, ff.UpdateCoverage())

				Convey("Then the location header should match the get coverage screen", func() {
					So(w.Header().Get("Location"), ShouldEqual, fmt.Sprintf("/filters/%s/dimensions/geography/coverage#search--parent", filterID))
				})

				Convey("And the status code should be 301", func() {
					So(w.Code, ShouldEqual, http.StatusMovedPermanently)
				})
			})

			Convey("When the GetDimensionOptions filter API client responds with an error", func() {
				filterClient := NewMockFilterClient(mockCtrl)
				filterClient.
					EXPECT().
					GetDimensionOptions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.DimensionOptions{}, "", errors.New("internal error"))

				ff := NewFilterFlex(
					NewMockRenderClient(mockCtrl),
					filterClient,
					NewMockDatasetClient(mockCtrl),
					NewMockPopulationClient(mockCtrl),
					NewMockZebedeeClient(mockCtrl),
					cfg)
				w := runUpdateCoverage("test", "test", stubFormData, ff.UpdateCoverage())

				Convey("Then the client should not be redirected", func() {
					So(w.Header().Get("Location"), ShouldBeEmpty)
				})

				Convey("And the status code should be 500", func() {
					So(w.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})

			Convey("When the UpdateDimensions filter API client responds with an error", func() {
				filterClient := NewMockFilterClient(mockCtrl)
				filterClient.
					EXPECT().
					GetDimensionOptions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.DimensionOptions{}, "", nil)
				filterClient.
					EXPECT().
					UpdateDimensions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimension{}, "", errors.New("internal error"))

				ff := NewFilterFlex(
					NewMockRenderClient(mockCtrl),
					filterClient,
					NewMockDatasetClient(mockCtrl),
					NewMockPopulationClient(mockCtrl),
					NewMockZebedeeClient(mockCtrl),
					cfg)
				w := runUpdateCoverage("test", "test", stubFormData, ff.UpdateCoverage())

				Convey("Then the client should not be redirected", func() {
					So(w.Header().Get("Location"), ShouldBeEmpty)
				})

				Convey("And the status code should be 500", func() {
					So(w.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})
		})

		Convey("Given a valid add option request with a saved parent option", func() {
			stubFormData := url.Values{}
			stubFormData.Add("dimension", "geography")
			stubFormData.Add("add-option", "0")
			stubFormData.Add("coverage", "name-search")
			stubFormData.Add("geog-id", "city")
			stubFormData.Add("option-type", "parent-search")

			Convey("When additional call to DeleteDimensionOptions is made", func() {
				const filterID = "1234"

				filterClient := NewMockFilterClient(mockCtrl)
				filterClient.
					EXPECT().
					GetDimensionOptions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.DimensionOptions{
						Items: []filter.DimensionOption{
							{
								Option: "Option 1",
							},
							{
								Option: "Option 2",
							},
						},
						TotalCount: 2,
					}, "", nil)
				filterClient.
					EXPECT().
					UpdateDimensions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimension{}, "", nil)
				filterClient.
					EXPECT().
					DeleteDimensionOptions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", nil)

				ff := NewFilterFlex(
					NewMockRenderClient(mockCtrl),
					filterClient,
					NewMockDatasetClient(mockCtrl),
					NewMockPopulationClient(mockCtrl),
					NewMockZebedeeClient(mockCtrl),
					cfg)
				w := runUpdateCoverage(filterID, "geography", stubFormData, ff.UpdateCoverage())

				Convey("Then the location header should match the get coverage screen", func() {
					So(w.Header().Get("Location"), ShouldEqual, fmt.Sprintf("/filters/%s/dimensions/geography/coverage#search--name", filterID))
				})

				Convey("And the status code should be 301", func() {
					So(w.Code, ShouldEqual, http.StatusMovedPermanently)
				})
			})

			Convey("When the DeleteDimensionOptions filter API client responds with an error", func() {
				const filterID = "1234"

				filterClient := NewMockFilterClient(mockCtrl)
				filterClient.
					EXPECT().
					GetDimensionOptions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.DimensionOptions{
						Items: []filter.DimensionOption{
							{
								Option: "Option 1",
							},
							{
								Option: "Option 2",
							},
						},
						TotalCount: 2,
					}, "", nil)
				filterClient.
					EXPECT().
					DeleteDimensionOptions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", errors.New("internal error"))

				ff := NewFilterFlex(
					NewMockRenderClient(mockCtrl),
					filterClient,
					NewMockDatasetClient(mockCtrl),
					NewMockPopulationClient(mockCtrl),
					NewMockZebedeeClient(mockCtrl),
					cfg)
				w := runUpdateCoverage(filterID, "geography", stubFormData, ff.UpdateCoverage())

				Convey("Then the client should not be redirected", func() {
					So(w.Header().Get("Location"), ShouldBeEmpty)
				})

				Convey("And the status code should be 500", func() {
					So(w.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})
		})

		Convey("Given a valid add parent option request with a saved parent option of a different type", func() {
			stubFormData := url.Values{}
			stubFormData.Add("dimension", "geography")
			stubFormData.Add("add-parent-option", "0")
			stubFormData.Add("coverage", "parent-search")
			stubFormData.Add("larger-area", "country")
			stubFormData.Add("set-parent", "parent")
			stubFormData.Add("geog-id", "city")

			Convey("When additional call to DeleteDimensionOptions is made", func() {
				const filterID = "1234"

				filterClient := NewMockFilterClient(mockCtrl)
				filterClient.
					EXPECT().
					GetDimensionOptions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.DimensionOptions{
						Items: []filter.DimensionOption{
							{
								Option: "Option 1",
							},
							{
								Option: "Option 2",
							},
						},
						TotalCount: 2,
					}, "", nil)
				filterClient.
					EXPECT().
					UpdateDimensions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimension{}, "", nil)
				filterClient.
					EXPECT().
					DeleteDimensionOptions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", nil)

				ff := NewFilterFlex(
					NewMockRenderClient(mockCtrl),
					filterClient,
					NewMockDatasetClient(mockCtrl),
					NewMockPopulationClient(mockCtrl),
					NewMockZebedeeClient(mockCtrl),
					cfg)
				w := runUpdateCoverage(filterID, "geography", stubFormData, ff.UpdateCoverage())

				Convey("Then the location header should match the get coverage screen", func() {
					So(w.Header().Get("Location"), ShouldEqual, fmt.Sprintf("/filters/%s/dimensions/geography/coverage#search--parent", filterID))
				})

				Convey("And the status code should be 301", func() {
					So(w.Code, ShouldEqual, http.StatusMovedPermanently)
				})
			})

			Convey("When the DeleteDimensionOptions filter API client responds with an error", func() {
				const filterID = "1234"

				filterClient := NewMockFilterClient(mockCtrl)
				filterClient.
					EXPECT().
					GetDimensionOptions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.DimensionOptions{
						Items: []filter.DimensionOption{
							{
								Option: "Option 1",
							},
							{
								Option: "Option 2",
							},
						},
						TotalCount: 2,
					}, "", nil)
				filterClient.
					EXPECT().
					DeleteDimensionOptions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", errors.New("internal error"))

				ff := NewFilterFlex(
					NewMockRenderClient(mockCtrl),
					filterClient,
					NewMockDatasetClient(mockCtrl),
					NewMockPopulationClient(mockCtrl),
					NewMockZebedeeClient(mockCtrl),
					cfg)
				w := runUpdateCoverage(filterID, "geography", stubFormData, ff.UpdateCoverage())

				Convey("Then the client should not be redirected", func() {
					So(w.Header().Get("Location"), ShouldBeEmpty)
				})

				Convey("And the status code should be 500", func() {
					So(w.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})
		})

		Convey("Given a valid delete option request", func() {
			stubFormData := url.Values{}
			stubFormData.Add("dimension", "geography")
			stubFormData.Add("delete-option", "0")
			stubFormData.Add("coverage", "name-search")
			stubFormData.Add("geog-id", "city")

			Convey("When the user is redirected to the get coverage screen", func() {
				const filterID = "1234"

				filterClient := NewMockFilterClient(mockCtrl)
				filterClient.
					EXPECT().
					RemoveDimensionValue(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", nil)

				ff := NewFilterFlex(
					NewMockRenderClient(mockCtrl),
					filterClient,
					NewMockDatasetClient(mockCtrl),
					NewMockPopulationClient(mockCtrl),
					NewMockZebedeeClient(mockCtrl),
					cfg)
				w := runUpdateCoverage(filterID, "geography", stubFormData, ff.UpdateCoverage())

				Convey("Then the location header should match the get coverage screen", func() {
					So(w.Header().Get("Location"), ShouldEqual, fmt.Sprintf("/filters/%s/dimensions/geography/coverage#search--name", filterID))
				})

				Convey("And the status code should be 301", func() {
					So(w.Code, ShouldEqual, http.StatusMovedPermanently)
				})
			})

			Convey("When the filter API client responds with an error", func() {
				filterClient := NewMockFilterClient(mockCtrl)
				filterClient.
					EXPECT().
					RemoveDimensionValue(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", errors.New("internal error"))

				ff := NewFilterFlex(
					NewMockRenderClient(mockCtrl),
					filterClient,
					NewMockDatasetClient(mockCtrl),
					NewMockPopulationClient(mockCtrl),
					NewMockZebedeeClient(mockCtrl),
					cfg)
				w := runUpdateCoverage("test", "test", stubFormData, ff.UpdateCoverage())

				Convey("Then the client should not be redirected", func() {
					So(w.Header().Get("Location"), ShouldBeEmpty)
				})

				Convey("And the status code should be 500", func() {
					So(w.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})
		})

		Convey("Given a valid all geography request", func() {
			stubFormData := url.Values{}
			stubFormData.Add("dimension", "geography")
			stubFormData.Add("coverage", "default")
			stubFormData.Add("geog-id", "city")

			Convey("When the user selects the all geography option", func() {
				const filterID = "1234"

				filterClient := NewMockFilterClient(mockCtrl)
				filterClient.
					EXPECT().
					DeleteDimensionOptions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", nil)

				ff := NewFilterFlex(
					NewMockRenderClient(mockCtrl),
					filterClient,
					NewMockDatasetClient(mockCtrl),
					NewMockPopulationClient(mockCtrl),
					NewMockZebedeeClient(mockCtrl),
					cfg)
				w := runUpdateCoverage(filterID, "geography", stubFormData, ff.UpdateCoverage())

				Convey("Then the location header should match the review screen", func() {
					So(w.Header().Get("Location"), ShouldEqual, fmt.Sprintf("/filters/%s/dimensions", filterID))
				})

				Convey("And the status code should be 301", func() {
					So(w.Code, ShouldEqual, http.StatusMovedPermanently)
				})
			})

			Convey("When the filter API client responds with an error", func() {
				filterClient := NewMockFilterClient(mockCtrl)
				filterClient.
					EXPECT().
					DeleteDimensionOptions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", errors.New("internal error"))

				ff := NewFilterFlex(
					NewMockRenderClient(mockCtrl),
					filterClient,
					NewMockDatasetClient(mockCtrl),
					NewMockPopulationClient(mockCtrl),
					NewMockZebedeeClient(mockCtrl),
					cfg)
				w := runUpdateCoverage("test", "test", stubFormData, ff.UpdateCoverage())

				Convey("Then the client should not be redirected", func() {
					So(w.Header().Get("Location"), ShouldBeEmpty)
				})

				Convey("And the status code should be 500", func() {
					So(w.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})
		})

		Convey("Given a valid name search request", func() {
			stubFormData := url.Values{}
			stubFormData.Add("dimension", "geography")
			stubFormData.Add("coverage", "name-search")
			stubFormData.Add("q", "area")
			stubFormData.Add("is-search", "true")
			stubFormData.Add("geog-id", "city")

			Convey("When the user is redirected to the get coverage screen", func() {
				const filterID = "1234"

				ff := NewFilterFlex(
					NewMockRenderClient(mockCtrl),
					NewMockFilterClient(mockCtrl),
					NewMockDatasetClient(mockCtrl),
					NewMockPopulationClient(mockCtrl),
					NewMockZebedeeClient(mockCtrl),
					cfg)
				w := runUpdateCoverage(filterID, "geography", stubFormData, ff.UpdateCoverage())

				Convey("Then the location header should match the get coverage screen with query persisted", func() {
					So(w.Header().Get("Location"), ShouldEqual, fmt.Sprintf("/filters/%s/dimensions/geography/coverage?c=name-search&q=area#search--name", filterID))
				})

				Convey("And the status code should be 301", func() {
					So(w.Code, ShouldEqual, http.StatusMovedPermanently)
				})
			})
		})

		Convey("Given a valid parent search request", func() {
			stubFormData := url.Values{}
			stubFormData.Add("dimension", "geography")
			stubFormData.Add("coverage", "parent-search")
			stubFormData.Add("pq", "area")
			stubFormData.Add("is-search", "true")
			stubFormData.Add("larger-area", "country")
			stubFormData.Add("geog-id", "city")

			Convey("When the user is redirected to the get coverage screen", func() {
				const filterID = "1234"

				ff := NewFilterFlex(
					NewMockRenderClient(mockCtrl),
					NewMockFilterClient(mockCtrl),
					NewMockDatasetClient(mockCtrl),
					NewMockPopulationClient(mockCtrl),
					NewMockZebedeeClient(mockCtrl),
					cfg)
				w := runUpdateCoverage(filterID, "geography", stubFormData, ff.UpdateCoverage())

				Convey("Then the location header should match the get coverage screen with query persisted", func() {
					So(w.Header().Get("Location"), ShouldEqual, fmt.Sprintf("/filters/%s/dimensions/geography/coverage?c=parent-search&p=country&pq=area#search--parent", filterID))
				})

				Convey("And the status code should be 301", func() {
					So(w.Code, ShouldEqual, http.StatusMovedPermanently)
				})
			})
		})

		Convey("Given an invalid parent search request", func() {
			stubFormData := url.Values{}
			stubFormData.Add("dimension", "geography")
			stubFormData.Add("coverage", "parent-search")
			stubFormData.Add("pq", "area")
			stubFormData.Add("is-search", "true")
			stubFormData.Add("larger-area", "")
			stubFormData.Add("geog-id", "city")

			Convey("When the user is redirected to the get coverage screen", func() {
				const filterID = "1234"

				ff := NewFilterFlex(
					NewMockRenderClient(mockCtrl),
					NewMockFilterClient(mockCtrl),
					NewMockDatasetClient(mockCtrl),
					NewMockPopulationClient(mockCtrl),
					NewMockZebedeeClient(mockCtrl),
					cfg)
				w := runUpdateCoverage(filterID, "geography", stubFormData, ff.UpdateCoverage())

				Convey("Then the location header should match the get coverage screen with error parameter", func() {
					So(w.Header().Get("Location"), ShouldEqual, fmt.Sprintf("/filters/%s/dimensions/geography/coverage?c=parent-search&error=true", filterID))
				})

				Convey("And the status code should be 301", func() {
					So(w.Code, ShouldEqual, http.StatusMovedPermanently)
				})
			})
		})

		Convey("Given a valid continue request", func() {
			stubFormData := url.Values{}
			stubFormData.Add("dimension", "geography")
			stubFormData.Add("coverage", "name-search")
			stubFormData.Add("geog-id", "city")

			Convey("When the user makes the request", func() {
				const filterID = "1234"

				ff := NewFilterFlex(
					NewMockRenderClient(mockCtrl),
					NewMockFilterClient(mockCtrl),
					NewMockDatasetClient(mockCtrl),
					NewMockPopulationClient(mockCtrl),
					NewMockZebedeeClient(mockCtrl),
					cfg)
				w := runUpdateCoverage(filterID, "geography", stubFormData, ff.UpdateCoverage())

				Convey("Then the location header should match the review screen", func() {
					So(w.Header().Get("Location"), ShouldEqual, fmt.Sprintf("/filters/%s/dimensions", filterID))
				})

				Convey("And the status code should be 301", func() {
					So(w.Code, ShouldEqual, http.StatusMovedPermanently)
				})
			})
		})

		Convey("Given an invalid request", func() {
			Convey("When the request is missing the hidden required form values", func() {
				tests := map[string]url.Values{
					"Missing coverage":  {"larger-area": []string{"country"}, "dimension": []string{"geography"}, "geog-id": []string{"city"}},
					"Unknown coverage":  {"coverage": []string{"1234"}, "dimension": []string{"geography"}},
					"Missing dimension": {"coverage": []string{"default"}, "geog-id": []string{"city"}},
					"Missing geog-id":   {"coverage": []string{"default"}},
				}
				ff := NewFilterFlex(
					NewMockRenderClient(mockCtrl),
					NewMockFilterClient(mockCtrl),
					NewMockDatasetClient(mockCtrl),
					NewMockPopulationClient(mockCtrl),
					NewMockZebedeeClient(mockCtrl),
					cfg)

				for name, formData := range tests {
					Convey(name, func() {
						w := runUpdateCoverage("test", "test", formData, ff.UpdateCoverage())

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
