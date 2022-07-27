package model

import (
	coreModel "github.com/ONSdigital/dp-renderer/model"
)

// Coverage represents the data to display the coverage page
type Coverage struct {
	coreModel.Page
	Geography     string         `json:"geography"`
	IsSearch      bool           `json:"is_search"`
	Search        string         `json:"search"`
	SearchResults []SearchResult `json:"search_results"`
	AreasAdded    []string       `json:"areas_added"`
}

// SearchResult represents the data required to display a search result
type SearchResult struct {
	Label      string `json:"label"`
	ID         string `json:"id"`
	IsSelected bool   `json:"is_selected"`
}
