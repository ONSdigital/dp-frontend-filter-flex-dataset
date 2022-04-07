package handlers

import (
	"context"
	"io"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/dimension"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
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
	GetJobState(ctx context.Context, userAuthToken, serviceAuthToken, downloadServiceToken, collectionID, filterID string) (f filter.Model, eTag string, err error)
	UpdateFlexBlueprint(ctx context.Context, userAuthToken, serviceAuthToken, downloadServiceToken, collectionID string, m filter.Model, doSubmit bool, populationType string, ifMatch string) (filter.Model, string, error)
}

// DatasetClient is an interface with methods required for a dataset client
type DatasetClient interface {
	GetOptions(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, id, edition, version, dimension string, q *dataset.QueryParams) (m dataset.Options, err error)
}

// DimensionClient is an interface with methods required for a dimension client
type DimensionClient interface {
	GetAreaTypes(ctx context.Context, userAuthToken, serviceAuthToken, datasetID string) (dimension.GetAreaTypesResponse, error)
}
