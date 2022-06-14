package handlers

import (
	"context"
	"io"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/dimension"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-renderer/model"
)

// To mock interfaces in this file
//go:generate mockgen -source=clients.go -destination=mock_clients.go -package=handlers github.com/ONSdigital/dp-frontend-articles-controller/handlers

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
	UpdateDimensions(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, id, name, ifMatch string, dimension filter.Dimension) (dim filter.Dimension, eTag string, err error)
	SubmitFilter(ctx context.Context, userAuthToken, serviceAuthToken, downloadServiceToken, ifMatch string, sfr filter.SubmitFilterRequest) (resp *filter.SubmitFilterResponse, eTag string, err error)
}

// DatasetClient is an interface with methods required for a dataset client
type DatasetClient interface {
	GetOptions(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, id, edition, version, dimension string, q *dataset.QueryParams) (m dataset.Options, err error)
	GetVersionDimensions(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, id, edition, version string) (m dataset.VersionDimensions, err error)
}

// DimensionClient is an interface with methods required for a dimension client
type DimensionClient interface {
	GetAreas(ctx context.Context, input dimension.GetAreasInput) (dimension.GetAreasResponse, error)
}

type PopulationClient interface {
	GetPopulationAreaTypes(ctx context.Context, userAuthToken, serviceAuthToken, datasetID string) (population.GetAreaTypesResponse, error)
}
