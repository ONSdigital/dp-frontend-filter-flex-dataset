package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/helpers"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/mocks"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/model"
	"github.com/ONSdigital/dp-renderer/helper"
	coreModel "github.com/ONSdigital/dp-renderer/model"
	gomock "github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDimensionsHandler(t *testing.T) {
	helper.InitialiseLocalisationsHelper(mocks.MockAssetFunction)
	mockCtrl := gomock.NewController(t)
	cfg := initialiseMockConfig()

	Convey("Dimensions Selector", t, func() {
		Convey("Given a valid dimension param for a filter", func() {
			Convey("Then the page title contains the dimension name", func() {
				const dimensionName = "Number Of Siblings"

				mockFilter := NewMockFilterClient(mockCtrl)
				mockFilter.
					EXPECT().
					GetJobState(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Model{}, "", nil)
				mockFilter.
					EXPECT().
					GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimension{Name: dimensionName}, "", nil).
					AnyTimes()

				mockRend := NewMockRenderClient(mockCtrl)
				mockRend.
					EXPECT().
					NewBasePageModel().
					Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain)).
					AnyTimes()

				mockRend.
					EXPECT().
					BuildPage(gomock.Any(), pageHasTitle{dimensionName}, gomock.Any())

				w := runDimensionsSelector(
					"number+of+siblings",
					DimensionsSelector(mockRend, mockFilter, NewMockPopulationClient(mockCtrl)),
				)

				Convey("And the status code should be 200", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
				})
			})
		})

		Convey("Given a dimension param which is missing from a filter", func() {
			mockFilter := NewMockFilterClient(mockCtrl)
			mockFilter.
				EXPECT().
				GetJobState(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(filter.Model{}, "", nil) // No filter dimensions
			mockFilter.
				EXPECT().
				GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(filter.Dimension{}, "", &filter.ErrInvalidFilterAPIResponse{ExpectedCode: http.StatusOK, ActualCode: http.StatusNotFound}).
				AnyTimes()

			w := runDimensionsSelector(
				"city",
				DimensionsSelector(NewMockRenderClient(mockCtrl), mockFilter, NewMockPopulationClient(mockCtrl)),
			)

			Convey("Then the status code should be 404", func() {
				So(w.Code, ShouldEqual, http.StatusNotFound)
			})
		})

		Convey("Given a dimension param which is not an area type", func() {
			const dimensionName = "siblings"

			stubDimension := filter.Dimension{
				Name:       dimensionName,
				IsAreaType: helpers.ToBoolPtr(false),
			}

			// This will change, but represents the current non-area-type behaviour.
			Convey("Then the page should contain no selections", func() {
				mockFilter := NewMockFilterClient(mockCtrl)
				mockFilter.
					EXPECT().
					GetJobState(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Model{}, "", nil).
					AnyTimes()
				mockFilter.
					EXPECT().
					GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(stubDimension, "", nil).
					AnyTimes()

				mockRend := NewMockRenderClient(mockCtrl)
				mockRend.
					EXPECT().
					NewBasePageModel().
					Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain)).
					AnyTimes()

				mockRend.
					EXPECT().
					// Assert that there are no selections passed to BuildPage
					BuildPage(gomock.Any(), pageMatchesSelections{}, gomock.Any())

				w := runDimensionsSelector(
					dimensionName,
					DimensionsSelector(mockRend, mockFilter, NewMockPopulationClient(mockCtrl)),
				)

				Convey("And the status code should be 200", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
				})
			})

			Convey("Then the page should have the area type bool set to false", func() {
				mockFilter := NewMockFilterClient(mockCtrl)
				mockFilter.
					EXPECT().
					GetJobState(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Model{}, "", nil).
					AnyTimes()
				mockFilter.
					EXPECT().
					GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(stubDimension, "", nil).
					AnyTimes()

				mockRend := NewMockRenderClient(mockCtrl)
				mockRend.
					EXPECT().
					NewBasePageModel().
					Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain)).
					AnyTimes()

				mockRend.
					EXPECT().
					// Assert that the area type boolean is false
					BuildPage(gomock.Any(), pageIsAreaType{false}, gomock.Any())

				w := runDimensionsSelector(
					dimensionName,
					DimensionsSelector(mockRend, mockFilter, NewMockPopulationClient(mockCtrl)),
				)

				Convey("And the status code should be 200", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
				})
			})

			Convey("Then the dimensions API should be queried using the dimension name", func() {
				mockFilter := NewMockFilterClient(mockCtrl)
				mockFilter.
					EXPECT().
					GetJobState(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Model{}, "", nil).
					AnyTimes()
				mockFilter.
					EXPECT().
					GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), "siblings").
					Return(stubDimension, "", nil).
					AnyTimes()

				mockRend := NewMockRenderClient(mockCtrl)
				mockRend.
					EXPECT().
					NewBasePageModel().
					Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain)).
					AnyTimes()

				mockRend.
					EXPECT().
					BuildPage(gomock.Any(), gomock.Any(), gomock.Any())

				w := runDimensionsSelector(
					dimensionName,
					DimensionsSelector(mockRend, mockFilter, NewMockPopulationClient(mockCtrl)),
				)

				Convey("And the status code should be 200", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
				})
			})

			Convey("When the filter API responds with an error", func() {
				mockFilter := NewMockFilterClient(mockCtrl)
				mockFilter.EXPECT().
					GetJobState(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Model{}, "", errors.New("oh no")).
					AnyTimes()

				w := runDimensionsSelector(
					dimensionName,
					DimensionsSelector(NewMockRenderClient(mockCtrl), mockFilter, NewMockPopulationClient(mockCtrl)),
				)

				Convey("Then the status code should be 500", func() {
					So(w.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})
		})

		Convey("Given an area type", func() {
			const dimensionName = "city"

			stubAreaTypeDimension := filter.Dimension{
				Name:       dimensionName,
				IsAreaType: helpers.ToBoolPtr(true),
			}

			Convey("When area types are returned", func() {
				Convey("Then the page should contain a list of area type selections", func() {
					const dimensionID = "city"
					const dimensionLabel = "City"

					mockFilter := NewMockFilterClient(mockCtrl)
					mockFilter.EXPECT().
						GetJobState(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(filter.Model{}, "", nil).
						AnyTimes()
					mockFilter.
						EXPECT().
						GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(stubAreaTypeDimension, "", nil).
						AnyTimes()

					mockPc := NewMockPopulationClient(mockCtrl)
					mockPc.EXPECT().
						GetPopulationAreaTypes(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(
							population.GetAreaTypesResponse{
								AreaTypes: []population.AreaType{{
									ID:         dimensionID,
									Label:      dimensionLabel,
									TotalCount: 1,
								}},
							},
							nil,
						).
						AnyTimes()

					mockRend := NewMockRenderClient(mockCtrl)
					mockRend.EXPECT().
						NewBasePageModel().
						Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain)).
						AnyTimes()

					// Validate page data contains selections
					mockRend.EXPECT().
						BuildPage(
							gomock.Any(),
							pageMatchesSelections{
								selections: []model.Selection{
									{
										Value:      dimensionID,
										Label:      dimensionLabel,
										TotalCount: 1,
									},
								},
							},
							"selector",
						)

					w := runDimensionsSelector(dimensionName, DimensionsSelector(mockRend, mockFilter, mockPc))

					Convey("And the status code should be 200", func() {
						So(w.Code, ShouldEqual, http.StatusOK)
					})
				})

				Convey("Then the dimensions API client should request area types using the cantabular ID", func() {
					const cantabularID = "cantabular"

					mockFilter := NewMockFilterClient(mockCtrl)
					mockFilter.EXPECT().
						GetJobState(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(filter.Model{PopulationType: cantabularID}, "", nil).
						AnyTimes()
					mockFilter.
						EXPECT().
						GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(stubAreaTypeDimension, "", nil).
						AnyTimes()

					mockPc := NewMockPopulationClient(mockCtrl)
					mockPc.EXPECT().
						GetPopulationAreaTypes(gomock.Any(), gomock.Any(), gomock.Any(), cantabularID).
						Return(population.GetAreaTypesResponse{}, nil)

					mockRend := NewMockRenderClient(mockCtrl)
					mockRend.EXPECT().
						NewBasePageModel().
						Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain)).
						AnyTimes()

					mockRend.
						EXPECT().
						BuildPage(gomock.Any(), gomock.Any(), "selector").
						AnyTimes()

					w := runDimensionsSelector(dimensionName, DimensionsSelector(mockRend, mockFilter, mockPc))

					Convey("And the status code should be 200", func() {
						So(w.Code, ShouldEqual, http.StatusOK)
					})
				})

				Convey("Then the page should have the area type bool set to true", func() {
					mockFilter := NewMockFilterClient(mockCtrl)
					mockFilter.EXPECT().
						GetJobState(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(filter.Model{}, "", nil).
						AnyTimes()
					mockFilter.
						EXPECT().
						GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(stubAreaTypeDimension, "", nil).
						AnyTimes()

					mockPc := NewMockPopulationClient(mockCtrl)
					mockPc.EXPECT().
						GetPopulationAreaTypes(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(population.GetAreaTypesResponse{}, nil)

					mockRend := NewMockRenderClient(mockCtrl)
					mockRend.EXPECT().
						NewBasePageModel().
						Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain)).
						AnyTimes()

					mockRend.
						EXPECT().
						// Assert that the area type boolean is true
						BuildPage(gomock.Any(), pageIsAreaType{true}, gomock.Any())

					w := runDimensionsSelector(dimensionName, DimensionsSelector(mockRend, mockFilter, mockPc))

					Convey("And the status code should be 200", func() {
						So(w.Code, ShouldEqual, http.StatusOK)
					})
				})
			})

			Convey("Given a truthy error query param", func() {
				req := httptest.NewRequest(http.MethodGet, "/filters/1234/dimensions/city?error=true", nil)

				Convey("Then the page should contain a populated error", func() {
					mockFilter := NewMockFilterClient(mockCtrl)
					mockFilter.
						EXPECT().
						GetJobState(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(filter.Model{}, "", nil)
					mockFilter.
						EXPECT().
						GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(filter.Dimension{IsAreaType: helpers.ToBoolPtr(true)}, "", nil).
						AnyTimes()

					mockPc := NewMockPopulationClient(mockCtrl)
					mockPc.EXPECT().
						GetPopulationAreaTypes(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(population.GetAreaTypesResponse{}, nil).
						AnyTimes()

					mockRend := NewMockRenderClient(mockCtrl)
					mockRend.
						EXPECT().
						NewBasePageModel().
						Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain)).
						AnyTimes()

					mockRend.
						EXPECT().
						// Confirm the page contains an error
						BuildPage(gomock.Any(), pageHasError{true}, gomock.Any())

					selector := DimensionsSelector(mockRend, mockFilter, mockPc)

					w := httptest.NewRecorder()
					router := mux.NewRouter()
					router.HandleFunc("/filters/{filterID}/dimensions/{name}", selector)
					router.ServeHTTP(w, req)

					Convey("And the status code should be 200", func() {
						So(w.Code, ShouldEqual, http.StatusOK)
					})
				})
			})

			Convey("When the dimension API responds with an error", func() {
				mockFilter := NewMockFilterClient(mockCtrl)
				mockFilter.EXPECT().
					GetJobState(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Model{}, "", nil)
				mockFilter.
					EXPECT().
					GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(stubAreaTypeDimension, "", nil).
					AnyTimes()

				mockPc := NewMockPopulationClient(mockCtrl)
				mockPc.EXPECT().
					GetPopulationAreaTypes(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(population.GetAreaTypesResponse{}, errors.New("oh no"))

				mockRend := NewMockRenderClient(mockCtrl)
				mockRend.EXPECT().
					NewBasePageModel().
					Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain)).
					AnyTimes()

				w := runDimensionsSelector(dimensionName, DimensionsSelector(mockRend, mockFilter, mockPc))

				Convey("Then the status code should be 500", func() {
					So(w.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})
		})
	})
}

func runDimensionsSelector(dimension string, selector func(http.ResponseWriter, *http.Request)) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", fmt.Sprintf("/filters/1234/dimensions/%s", dimension), nil)

	router := mux.NewRouter()
	router.HandleFunc("/filters/{filterID}/dimensions/{name}", selector)

	router.ServeHTTP(w, req)

	return w
}

// pageMatchesSelections is a gomock matcher that confirms a selection page
// contains the correct selections (i.e. radio buttons).
type pageMatchesSelections struct {
	selections []model.Selection
}

func (c pageMatchesSelections) Matches(x interface{}) bool {
	page, ok := x.(model.Selector)
	if !ok {
		return false
	}

	return reflect.DeepEqual(c.selections, page.Selections)
}

func (c pageMatchesSelections) String() string {
	return fmt.Sprintf("is equal to %+v", c.selections)
}

// pageMatchesSelections is a gomock matcher that confirms a selection page
// has the correct page title.
type pageHasTitle struct {
	title string
}

func (p pageHasTitle) Matches(x interface{}) bool {
	page, ok := x.(model.Selector)
	if !ok {
		return false
	}

	return p.title == page.Page.Metadata.Title
}

func (p pageHasTitle) String() string {
	return fmt.Sprintf("title is equal to \"%s\"", p.title)
}

// pageIsAreaType is a gomock matcher that confirms a selection page
// `IsAreaType` boolean is set to the expected value.
type pageIsAreaType struct {
	expected bool
}

func (c pageIsAreaType) Matches(x interface{}) bool {
	page, ok := x.(model.Selector)
	if !ok {
		return false
	}

	return page.IsAreaType == c.expected
}

func (c pageIsAreaType) String() string {
	return fmt.Sprintf("is equal to %+v", c.expected)
}

// pageHasError is a gomock matcher that confirms a selection page
// has a populated error.
type pageHasError struct {
	expected bool
}

func (p pageHasError) Matches(x interface{}) bool {
	page, ok := x.(model.Selector)
	if !ok {
		return false
	}

	if p.expected {
		return page.Error.Title != ""
	}

	return page.Error.Title == ""
}

func (p pageHasError) String() string {
	return fmt.Sprintf("is equal to %+v", p.expected)
}
