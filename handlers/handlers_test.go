package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/dimension"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/config"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/model"
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
					DimensionsSelector(mockRend, mockFilter, NewMockDimensionClient(mockCtrl)),
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
				DimensionsSelector(NewMockRenderClient(mockCtrl), mockFilter, NewMockDimensionClient(mockCtrl)),
			)

			Convey("Then the status code should be 404", func() {
				So(w.Code, ShouldEqual, http.StatusNotFound)
			})
		})

		Convey("Given a dimension param which is not an area type", func() {
			const dimensionName = "siblings"

			stubDimension := filter.Dimension{
				Name:       dimensionName,
				IsAreaType: toBoolPtr(false),
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
					DimensionsSelector(mockRend, mockFilter, NewMockDimensionClient(mockCtrl)),
				)

				Convey("And the status code should be 200", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
				})
			})

			// This can be removed once we start using the name/ID.
			Convey("Then the dimensions API should be queried using the display name", func() {
				mockFilter := NewMockFilterClient(mockCtrl)
				mockFilter.
					EXPECT().
					GetJobState(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Model{}, "", nil).
					AnyTimes()
				mockFilter.
					EXPECT().
					GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), "Siblings").
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
					DimensionsSelector(mockRend, mockFilter, NewMockDimensionClient(mockCtrl)),
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
					DimensionsSelector(NewMockRenderClient(mockCtrl), mockFilter, NewMockDimensionClient(mockCtrl)),
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
				IsAreaType: toBoolPtr(true),
			}

			Convey("When area types are returned", func() {
				Convey("Then the page should contain a list of area type selections", func() {
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

					mockDimension := NewMockDimensionClient(mockCtrl)
					mockDimension.EXPECT().
						GetAreaTypes(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(
							dimension.GetAreaTypesResponse{
								AreaTypes: []dimension.AreaType{{
									ID:         dimensionName,
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
										Value:      dimensionName,
										Label:      dimensionLabel,
										TotalCount: 1,
									},
								},
							},
							"selector",
						)

					w := runDimensionsSelector(dimensionName, DimensionsSelector(mockRend, mockFilter, mockDimension))

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

					mockDimension := NewMockDimensionClient(mockCtrl)
					mockDimension.EXPECT().
						GetAreaTypes(gomock.Any(), gomock.Any(), gomock.Any(), cantabularID).
						Return(dimension.GetAreaTypesResponse{}, nil)

					mockRend := NewMockRenderClient(mockCtrl)
					mockRend.EXPECT().
						NewBasePageModel().
						Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain)).
						AnyTimes()

					mockRend.
						EXPECT().
						BuildPage(gomock.Any(), gomock.Any(), "selector").
						AnyTimes()

					w := runDimensionsSelector(dimensionName, DimensionsSelector(mockRend, mockFilter, mockDimension))

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

				mockDimension := NewMockDimensionClient(mockCtrl)
				mockDimension.EXPECT().
					GetAreaTypes(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dimension.GetAreaTypesResponse{}, errors.New("oh no"))

				mockRend := NewMockRenderClient(mockCtrl)
				mockRend.EXPECT().
					NewBasePageModel().
					Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain)).
					AnyTimes()

				w := runDimensionsSelector(dimensionName, DimensionsSelector(mockRend, mockFilter, mockDimension))

				Convey("Then the status code should be 500", func() {
					So(w.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})
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

func toBoolPtr(val bool) *bool {
	return &val
}
