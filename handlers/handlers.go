package handlers

import "github.com/ONSdigital/dp-frontend-filter-flex-dataset/config"

// FilterFlex represents the handlers for filtering and flexing
type FilterFlex struct {
	Render                      RenderClient
	FilterClient                FilterClient
	DatasetClient               DatasetClient
	PopulationClient            PopulationClient
	ZebedeeClient               ZebedeeClient
	EnableMultivariate          bool
	DefaultMaximumSearchResults int
}

// NewFilterFlex creates a new instance of FilterFlex
func NewFilterFlex(rc RenderClient, fc FilterClient, dc DatasetClient, pc PopulationClient, zc ZebedeeClient, cfg *config.Config) *FilterFlex {
	return &FilterFlex{
		Render:                      rc,
		FilterClient:                fc,
		DatasetClient:               dc,
		PopulationClient:            pc,
		ZebedeeClient:               zc,
		EnableMultivariate:          cfg.EnableMultivariate,
		DefaultMaximumSearchResults: cfg.DefaultMaximumSearchResults,
	}
}
