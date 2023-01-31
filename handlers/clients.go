package handlers

import (
	"context"
	"io"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-renderer/model"
)

// To mock interfaces in this file
//go:generate mockgen -source=clients.go -destination=mock_clients.go -package=handlers github.com/ONSdigital/dp-frontend-filter-flex-dataset/handlers

// ClientError is an interface that can be used to retrieve the status code if a client has errored
type ClientError interface {
	Error() string
	Code() int
}

// RenderClient is an interface with methods for rendering a template
type RenderClient interface {
	BuildPage(w io.Writer, pageModel interface{}, templateName string)
	NewBasePageModel() model.Page
}

// FilterClient is an interface with the methods required for a filter client
type FilterClient interface {
	GetFilter(ctx context.Context, input filter.GetFilterInput) (f *filter.GetFilterResponse, err error)
	GetJobState(ctx context.Context, userAuthToken, serviceAuthToken, downloadServiceToken, collectionID, filterID string) (f filter.Model, eTag string, err error)
	GetDimension(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, filterID, name string) (dim filter.Dimension, eTag string, err error)
	GetDimensions(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, filterID string, q *filter.QueryParams) (dims filter.Dimensions, eTag string, err error)
	GetDimensionOptions(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, filterID, name string, q *filter.QueryParams) (opts filter.DimensionOptions, eTag string, err error)
	UpdateDimensions(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, id, name, ifMatch string, dimension filter.Dimension) (dim filter.Dimension, eTag string, err error)
	AddDimensionValue(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, filterID, name, value, ifMatch string) (eTag string, err error)
	RemoveDimensionValue(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, filterID, name, value, ifMatch string) (eTag string, err error)
	DeleteDimensionOptions(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, filterID, name string) (eTag string, err error)
	SubmitFilter(ctx context.Context, userAuthToken, serviceAuthToken, downloadServiceToken, ifMatch string, sfr filter.SubmitFilterRequest) (resp *filter.SubmitFilterResponse, eTag string, err error)
	AddFlexDimension(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, id, name string, options []string, isAreaType bool, ifMatch string) (eTag string, err error)
	RemoveDimension(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, filterID, name, ifMatch string) (eTag string, err error)
}

// DatasetClient is an interface with methods required for a dataset client
type DatasetClient interface {
	Get(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, datasetID string) (m dataset.DatasetDetails, err error)
	GetOptions(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, id, edition, version, dimension string, q *dataset.QueryParams) (m dataset.Options, err error)
	GetVersion(ctx context.Context, userAuthToken, serviceAuthToken, downloadServiceAuthToken, collectionID, datasetID, edition, version string) (m dataset.Version, err error)
	GetVersionDimensions(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, id, edition, version string) (m dataset.VersionDimensions, err error)
}

// PopulationClient is an interface with methods required for a population client
type PopulationClient interface {
	GetAreaTypes(ctx context.Context, input population.GetAreaTypesInput) (population.GetAreaTypesResponse, error)
	GetAreas(ctx context.Context, input population.GetAreasInput) (population.GetAreasResponse, error)
	GetAreaTypeParents(ctx context.Context, input population.GetAreaTypeParentsInput) (population.GetAreaTypeParentsResponse, error)
	GetArea(ctx context.Context, input population.GetAreaInput) (population.GetAreaResponse, error)
	GetCategorisations(ctx context.Context, input population.GetCategorisationsInput) (population.GetCategorisationsResponse, error)
	GetDimensions(ctx context.Context, input population.GetDimensionsInput) (population.GetDimensionsResponse, error)
	GetDimensionsDescription(ctx context.Context, input population.GetDimensionsDescriptionInput) (population.GetDimensionsResponse, error)
	GetParentAreaCount(ctx context.Context, input population.GetParentAreaCountInput) (int, error)
}

// ZebedeeClient is an interface with methods required for the zebedee client
type ZebedeeClient interface {
	GetHomepageContent(ctx context.Context, userAccessToken, collectionID, lang, path string) (m zebedee.HomepageContent, err error)
}
