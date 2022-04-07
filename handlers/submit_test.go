package handlers

// import (
// 	"errors"
// 	"strings"
// 	"testing"

// 	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
// 	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
// 	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/config"
// 	gomock "github.com/golang/mock/gomock"
// 	. "github.com/smartystreets/goconvey/convey"
// )

// func TestSubmitHandler(t *testing.T) {
// 	mockCtrl := gomock.NewController(t)
// 	//cfg := initialiseMockConfig()
// 	ctx := gomock.Any()
// 	var mockServiceAuthToken, mockDownloadToken, mockUserAuthToken, mockCollectionID, mockFilterID string

// 	Convey("test submit handler", t, func() {
// 		// mockCfg := config.Config{EnableCensusPages: true}
// 		// mockVersions := dataset.VersionsList{
// 		// 	Items: []dataset.Version{
// 		// 		{}, // deliberately empty
// 		// 		{
// 		// 			Dimensions: []dataset.VersionDimension{
// 		// 				{
// 		// 					Name: "aggregate",
// 		// 				},
// 		// 				{
// 		// 					Name: "time",
// 		// 				},
// 		// 			},
// 		// 		},
// 		// 	},

// 		Convey("test Submit handler, starts a filter-outputs job and redirects", func() {
// 			mockClient := NewMockFilterClient(mockCtrl)
// 			mockFm := filter.Model{}
// 			mockClient.EXPECT().GetJobState(ctx, mockUserAuthToken, mockServiceAuthToken, mockDownloadToken, mockCollectionID, mockFilterID).Return()
// 			mockClient.EXPECT().UpdateFlexBlueprint(ctx, mockUserAuthToken, mockServiceAuthToken, mockDownloadToken, mockCollectionID, mockFm, true, "", "").Return("12345", "testETag", nil)

// 			body := strings.NewReader("dimension=aggregate")
// 			w := testResponse(301, body, "/datasets/1234/editions/2021/versions/1/filter-flex", mockClient, mockDatasetClient, true, mockCfg)

// 			location := w.Header().Get("Location")
// 			So(location, ShouldNotBeEmpty)

// 			So(location, ShouldEqual, "/filters/12345/dimensions/aggregate")
// 		})

// 		// Convey("test post route fails if config is false", func() {
// 		// 	mockCfg := config.Config{EnableCensusPages: false}
// 		// 	mockDatasetClient := NewMockDatasetClient(mockCtrl)
// 		// 	mockFilterClient := NewMockFilterClient(mockCtrl)
// 		// 	body := strings.NewReader("")

// 		// 	testResponse(500, body, "/datasets/1234/editions/2021/versions/1/filter-flex", mockFilterClient, mockDatasetClient, true, mockCfg)
// 		// })

// 		// Convey("test CreateFilterFlexID returns 500 if unable to create a blueprint on filter api", func() {
// 		// 	mockDatasetClient := NewMockDatasetClient(mockCtrl)
// 		// 	mockDatasetClient.EXPECT().GetVersion(ctx, userAuthToken, serviceAuthToken, "", collectionID, "1234", "2021", "1").Return(mockVersions.Items[0], nil)
// 		// 	mockDatasetClient.EXPECT().Get(ctx, userAuthToken, serviceAuthToken, collectionID, "1234").Return(dataset.DatasetDetails{IsBasedOn: &dataset.IsBasedOn{}}, nil)
// 		// 	mockFilterClient := NewMockFilterClient(mockCtrl)
// 		// 	mockFilterClient.EXPECT().CreateFlexibleBlueprint(ctx, userAuthToken, serviceAuthToken, "", collectionID, "1234", "2021", "1", gomock.Any(), "").Return("", "", errors.New("unable to create filter blueprint"))
// 		// 	body := strings.NewReader("")

// 		// 	testResponse(500, body, "/datasets/1234/editions/2021/versions/1/filter-flex", mockFilterClient, mockDatasetClient, true, mockCfg)
// 		// })
// 	})
// }
