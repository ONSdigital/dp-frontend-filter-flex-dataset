package model

import (
	coreModel "github.com/ONSdigital/dp-renderer/model"
)

// Coverage represents the data to display the coverage page
type Coverage struct {
	coreModel.Page
	Geography     string              `json:"geography"`
	Dimension     string              `json:"dimension"`
	HasNoResults  bool                `json:"has_no_results"`
	Search        string              `json:"search"`
	DisplaySearch bool                `json:"display_search"`
	SearchResults []SelectableElement `json:"search_results"`
	Options       []SelectableElement `json:"options"`
	ParentSelect  []SelectableElement `json:"parent_select"`
	NameSearch    SearchField         `json:"name_search"`
	ParentSearch  SearchField         `json:"parent_search"`
}

/* SelectableElement represents the data required for a selectable element.
Text is the human readable label.
Value is the value sent to the server.
IsSelected is a boolean representing whether the element is selected.
IsDisabled is a boolean representing whether the element is disabled */
type SelectableElement struct {
	Text       string `json:"text"`
	Value      string `json:"value"`
	IsSelected bool   `json:"is_selected"`
	IsDisabled bool   `json:"is_disabled"`
}

// SearchField represents the data required to populate the search input partial
type SearchField struct {
	Value    string `json:"value"`
	Name     string `json:"name"`
	ID       string `json:"id"`
	Language string `json:"language"`
}
