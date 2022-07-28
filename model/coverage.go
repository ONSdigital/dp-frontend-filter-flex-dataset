package model

import (
	coreModel "github.com/ONSdigital/dp-renderer/model"
)

// Coverage represents the data to display the coverage page
type Coverage struct {
	coreModel.Page
	Geography     string         `json:"geography"`
	Dimension     string         `json:"dimension"`
	HasNoResults  bool           `json:"has_no_results"`
	Search        string         `json:"search"`
	DisplaySearch bool           `json:"display_search"`
	SearchResults []SearchResult `json:"search_results"`
	Options       []Option       `json:"options"`
}

// SearchResult represents the data required to display a search result
type SearchResult struct {
	Label      string `json:"label"`
	ID         string `json:"id"`
	IsSelected bool   `json:"is_selected"`
}

// Option represents the data required to display an option
type Option struct {
	Label string `json:"label"`
	ID    string `json:"id"`
}
