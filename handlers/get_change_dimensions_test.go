package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/mocks"
	"github.com/ONSdigital/dp-renderer/helper"
	coreModel "github.com/ONSdigital/dp-renderer/model"
	gomock "github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetChangeDimensionsHandler(t *testing.T) {
	helper.InitialiseLocalisationsHelper(mocks.MockAssetFunction)
	mockCtrl := gomock.NewController(t)
	cfg := initialiseMockConfig()
	Convey("Change dimensions", t, func() {
		Convey("Given a valid page request", func() {
			mf := &filter.GetFilterResponse{
				Dataset: filter.Dataset{
					DatasetID: "test-dataset",
					Edition:   "2021",
					Version:   1,
				},
			}
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/filters/12345/dimensions/change", nil)

			Convey("When the filter is based off a multivariate dataset type", func() {
				md := dataset.DatasetDetails{
					Type: "multivariate",
				}

				mockRend := NewMockRenderClient(mockCtrl)
				mockRend.
					EXPECT().
					NewBasePageModel().
					Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain))
				mockRend.
					EXPECT().
					BuildPage(gomock.Any(), gomock.Any(), "dimensions")

				mockFc := NewMockFilterClient(mockCtrl)
				mockFc.
					EXPECT().
					GetFilter(gomock.Any(), gomock.Any()).
					Return(mf, nil)
				mockFc.
					EXPECT().
					GetDimensions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimensions{
						Items: []filter.Dimension{
							{
								Name:       "test dim",
								ID:         "td1",
								IsAreaType: new(bool),
							},
						},
					}, "", nil)
				mockFc.
					EXPECT().
					GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimension{
						Name:       "test dim",
						ID:         "td1",
						IsAreaType: new(bool),
					}, "", nil)

				mockDc := NewMockDatasetClient(mockCtrl)
				mockDc.
					EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(md, nil)

				mockPc := NewMockPopulationClient(mockCtrl)
				mockPc.
					EXPECT().
					GetDimensions(gomock.Any(), gomock.Any()).
					Return(population.GetDimensionsResponse{}, nil)

				mockZc := NewMockZebedeeClient(mockCtrl)
				mockZc.
					EXPECT().
					GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(zebedee.HomepageContent{}, nil)

				ff := NewFilterFlex(mockRend, mockFc, mockDc, mockPc, mockZc, cfg)
				router := mux.NewRouter()
				router.HandleFunc("/filters/12345/dimensions/change", ff.GetChangeDimensions())
				router.ServeHTTP(w, req)

				Convey("Then the status code should be 200", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
				})
			})

			Convey("When the filter is based off a multivariate dataset type and the user performs a search", func() {
				req := httptest.NewRequest(http.MethodGet, "/filters/12345/dimensions/change?q=test", nil)
				md := dataset.DatasetDetails{
					Type: "multivariate",
				}

				mockRend := NewMockRenderClient(mockCtrl)
				mockRend.
					EXPECT().
					NewBasePageModel().
					Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain))
				mockRend.
					EXPECT().
					BuildPage(gomock.Any(), gomock.Any(), "dimensions")

				mockFc := NewMockFilterClient(mockCtrl)
				mockFc.
					EXPECT().
					GetFilter(gomock.Any(), gomock.Any()).
					Return(mf, nil)
				mockFc.
					EXPECT().
					GetDimensions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimensions{
						Items: []filter.Dimension{
							{
								Name:       "test dim",
								ID:         "td1",
								IsAreaType: new(bool),
							},
						},
					}, "", nil)
				mockFc.
					EXPECT().
					GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimension{
						Name:       "test dim",
						ID:         "td1",
						IsAreaType: new(bool),
					}, "", nil)

				mockDc := NewMockDatasetClient(mockCtrl)
				mockDc.
					EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(md, nil)

				mockPc := NewMockPopulationClient(mockCtrl)
				mockPc.
					EXPECT().
					GetDimensions(gomock.Any(), gomock.Any()).
					Return(population.GetDimensionsResponse{}, nil).
					AnyTimes()

				mockZc := NewMockZebedeeClient(mockCtrl)
				mockZc.
					EXPECT().
					GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(zebedee.HomepageContent{}, nil)

				ff := NewFilterFlex(mockRend, mockFc, mockDc, mockPc, mockZc, cfg)
				router := mux.NewRouter()
				router.HandleFunc("/filters/12345/dimensions/change", ff.GetChangeDimensions())
				router.ServeHTTP(w, req)

				Convey("Then the status code should be 200", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
				})
			})

			Convey("When the filter is not based on a multivariate dataset type", func() {
				md := dataset.DatasetDetails{
					Type: "flex",
				}

				mockFc := NewMockFilterClient(mockCtrl)
				mockFc.
					EXPECT().
					GetFilter(gomock.Any(), filter.GetFilterInput{
						FilterID: "12345",
					}).
					Return(mf, nil)
				mockFc.
					EXPECT().
					GetDimensions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimensions{}, "", nil)

				mockDc := NewMockDatasetClient(mockCtrl)
				mockDc.
					EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(md, nil)

				mockZc := NewMockZebedeeClient(mockCtrl)
				mockZc.
					EXPECT().
					GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(zebedee.HomepageContent{}, nil)

				ff := NewFilterFlex(NewMockRenderClient(mockCtrl), mockFc, mockDc, NewMockPopulationClient(mockCtrl), mockZc, cfg)
				router := mux.NewRouter()
				router.HandleFunc("/filters/{filterID}/dimensions/change", ff.GetChangeDimensions())
				router.ServeHTTP(w, req)

				Convey("Then the user is redirected", func() {
					So(w.Code, ShouldEqual, http.StatusMovedPermanently)
				})

				Convey("And the location header is the overview page", func() {
					So(w.Header().Get("Location"), ShouldEqual, "/filters/12345/dimensions")
				})
			})

			Convey("When filter.GetDimensions api method responds with an error", func() {
				mockFc := NewMockFilterClient(mockCtrl)
				mockFc.
					EXPECT().
					GetDimensions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimensions{}, "", errors.New("Internal error"))

				ff := NewFilterFlex(NewMockRenderClient(mockCtrl), mockFc, NewMockDatasetClient(mockCtrl), NewMockPopulationClient(mockCtrl), NewMockZebedeeClient(mockCtrl), cfg)
				router := mux.NewRouter()
				router.HandleFunc("/filters/12345/dimensions/change", ff.GetChangeDimensions())
				router.ServeHTTP(w, req)

				Convey("Then the status code should be 500", func() {
					So(w.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})

			Convey("When zebedee.GetHomepageContent api method responds with an error", func() {
				md := dataset.DatasetDetails{
					Type: "multivariate",
				}

				mockRend := NewMockRenderClient(mockCtrl)
				mockRend.
					EXPECT().
					NewBasePageModel().
					Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain))
				mockRend.
					EXPECT().
					BuildPage(gomock.Any(), gomock.Any(), "dimensions")

				mockFc := NewMockFilterClient(mockCtrl)
				mockFc.
					EXPECT().
					GetFilter(gomock.Any(), gomock.Any()).
					Return(mf, nil)
				mockFc.
					EXPECT().
					GetDimensions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimensions{
						Items: []filter.Dimension{
							{
								Name:       "test dim",
								ID:         "td1",
								IsAreaType: new(bool),
							},
						},
					}, "", nil)
				mockFc.
					EXPECT().
					GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimension{
						Name:       "test dim",
						ID:         "td1",
						IsAreaType: new(bool),
					}, "", nil)

				mockDc := NewMockDatasetClient(mockCtrl)
				mockDc.
					EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(md, nil)

				mockPc := NewMockPopulationClient(mockCtrl)
				mockPc.
					EXPECT().
					GetDimensions(gomock.Any(), gomock.Any()).
					Return(population.GetDimensionsResponse{}, nil)

				mockZc := NewMockZebedeeClient(mockCtrl)
				mockZc.
					EXPECT().
					GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(zebedee.HomepageContent{}, errors.New("Internal error"))

				ff := NewFilterFlex(mockRend, mockFc, mockDc, mockPc, mockZc, cfg)
				router := mux.NewRouter()
				router.HandleFunc("/filters/12345/dimensions/change", ff.GetChangeDimensions())
				router.ServeHTTP(w, req)

				Convey("Then the status code should be 200", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
				})
			})

			Convey("When filter.GetFilter api method responds with an error", func() {
				md := dataset.DatasetDetails{
					Type: "multivariate",
				}

				mockFc := NewMockFilterClient(mockCtrl)
				mockFc.
					EXPECT().
					GetDimensions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimensions{
						Items: []filter.Dimension{
							{
								Name:       "test dim",
								ID:         "td1",
								IsAreaType: new(bool),
							},
						},
					}, "", nil)
				mockFc.
					EXPECT().
					GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimension{
						Name:       "test dim",
						ID:         "td1",
						IsAreaType: new(bool),
					}, "", nil)
				mockFc.
					EXPECT().
					GetFilter(gomock.Any(), gomock.Any()).
					Return(&filter.GetFilterResponse{}, errors.New("Internal error"))

				mockDc := NewMockDatasetClient(mockCtrl)
				mockDc.
					EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(md, nil)

				mockPc := NewMockPopulationClient(mockCtrl)
				mockPc.
					EXPECT().
					GetDimensions(gomock.Any(), gomock.Any()).
					Return(population.GetDimensionsResponse{}, nil)

				mockZc := NewMockZebedeeClient(mockCtrl)
				mockZc.
					EXPECT().
					GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(zebedee.HomepageContent{}, nil)

				ff := NewFilterFlex(NewMockRenderClient(mockCtrl), mockFc, mockDc, mockPc, mockZc, cfg)
				router := mux.NewRouter()
				router.HandleFunc("/filters/12345/dimensions/change", ff.GetChangeDimensions())
				router.ServeHTTP(w, req)

				Convey("Then the status code should be 500", func() {
					So(w.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})

			Convey("When the additional filter.GetDimension api method responds with an error", func() {
				md := dataset.DatasetDetails{
					Type: "multivariate",
				}

				mockFc := NewMockFilterClient(mockCtrl)
				mockFc.
					EXPECT().
					GetFilter(gomock.Any(), gomock.Any()).
					Return(mf, nil)
				mockFc.
					EXPECT().
					GetDimensions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimensions{
						Items: []filter.Dimension{{
							Name:       "test dim",
							ID:         "td1",
							IsAreaType: new(bool),
						}},
					}, "", nil)
				mockFc.
					EXPECT().
					GetDimension(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimension{}, "", errors.New("Internal error"))

				mockDc := NewMockDatasetClient(mockCtrl)
				mockDc.
					EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(md, nil)

				mockPc := NewMockPopulationClient(mockCtrl)
				mockPc.
					EXPECT().
					GetDimensions(gomock.Any(), gomock.Any()).
					Return(population.GetDimensionsResponse{}, nil)

				mockZc := NewMockZebedeeClient(mockCtrl)
				mockZc.
					EXPECT().
					GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(zebedee.HomepageContent{}, nil)

				ff := NewFilterFlex(NewMockRenderClient(mockCtrl), mockFc, mockDc, mockPc, mockZc, cfg)
				router := mux.NewRouter()
				router.HandleFunc("/filters/12345/dimensions/change", ff.GetChangeDimensions())
				router.ServeHTTP(w, req)

				Convey("Then the status code should be 500", func() {
					So(w.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})

			Convey("When dataset.Get api method responds with an error", func() {
				mockFc := NewMockFilterClient(mockCtrl)
				mockFc.
					EXPECT().
					GetFilter(gomock.Any(), gomock.Any()).
					Return(mf, nil)
				mockFc.
					EXPECT().
					GetDimensions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimensions{}, "", nil)

				mockDc := NewMockDatasetClient(mockCtrl)
				mockDc.
					EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dataset.DatasetDetails{}, errors.New("Internal error"))

				mockPc := NewMockPopulationClient(mockCtrl)
				mockPc.
					EXPECT().
					GetDimensions(gomock.Any(), gomock.Any()).
					Return(population.GetDimensionsResponse{}, nil)

				mockZc := NewMockZebedeeClient(mockCtrl)
				mockZc.
					EXPECT().
					GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(zebedee.HomepageContent{}, nil)

				ff := NewFilterFlex(NewMockRenderClient(mockCtrl), mockFc, mockDc, mockPc, mockZc, cfg)
				router := mux.NewRouter()
				router.HandleFunc("/filters/12345/dimensions/change", ff.GetChangeDimensions())
				router.ServeHTTP(w, req)

				Convey("Then the status code should be 500", func() {
					So(w.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})

			Convey("When population.GetDimensions api method responds with an error", func() {
				md := dataset.DatasetDetails{
					Type: "multivariate",
				}

				mockFc := NewMockFilterClient(mockCtrl)
				mockFc.
					EXPECT().
					GetFilter(gomock.Any(), gomock.Any()).
					Return(mf, nil)
				mockFc.
					EXPECT().
					GetDimensions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimensions{}, "", nil)

				mockDc := NewMockDatasetClient(mockCtrl)
				mockDc.
					EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(md, nil)

				mockPc := NewMockPopulationClient(mockCtrl)
				mockPc.
					EXPECT().
					GetDimensions(gomock.Any(), gomock.Any()).
					Return(population.GetDimensionsResponse{}, errors.New("Internal error"))

				mockZc := NewMockZebedeeClient(mockCtrl)
				mockZc.
					EXPECT().
					GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(zebedee.HomepageContent{}, nil)

				ff := NewFilterFlex(NewMockRenderClient(mockCtrl), mockFc, mockDc, mockPc, mockZc, cfg)
				router := mux.NewRouter()
				router.HandleFunc("/filters/12345/dimensions/change", ff.GetChangeDimensions())
				router.ServeHTTP(w, req)

				Convey("Then the status code should be 500", func() {
					So(w.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})

			Convey("When the additional population.GetDimensions api method call responds with an error", func() {
				req := httptest.NewRequest(http.MethodGet, "/filters/12345/dimensions/change?q=test", nil)
				md := dataset.DatasetDetails{
					Type: "multivariate",
				}

				mockFc := NewMockFilterClient(mockCtrl)
				mockFc.
					EXPECT().
					GetFilter(gomock.Any(), gomock.Any()).
					Return(mf, nil)
				mockFc.
					EXPECT().
					GetDimensions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(filter.Dimensions{}, "", nil)

				mockDc := NewMockDatasetClient(mockCtrl)
				mockDc.
					EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(md, nil)

				mockPc := NewMockPopulationClient(mockCtrl)
				mockPc.
					EXPECT().
					GetDimensions(gomock.Any(), gomock.Any()).
					Return(population.GetDimensionsResponse{}, nil)
				mockPc.
					EXPECT().
					GetDimensions(gomock.Any(), gomock.Any()).
					Return(population.GetDimensionsResponse{}, errors.New("Internal error"))

				mockZc := NewMockZebedeeClient(mockCtrl)
				mockZc.
					EXPECT().
					GetHomepageContent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(zebedee.HomepageContent{}, nil)

				ff := NewFilterFlex(NewMockRenderClient(mockCtrl), mockFc, mockDc, mockPc, mockZc, cfg)
				router := mux.NewRouter()
				router.HandleFunc("/filters/12345/dimensions/change", ff.GetChangeDimensions())
				router.ServeHTTP(w, req)

				Convey("Then the status code should be 500", func() {
					So(w.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})
		})
	})
}
