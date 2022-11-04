package model

import (
	coreModel "github.com/ONSdigital/dp-renderer/model"
)

// Coverage represents the data to display the coverage page
type Coverage struct {
	coreModel.Page
	Geography          string              `json:"geography"`
	Dimension          string              `json:"dimension"`
	GeographyID        string              `json:"geography_id"`
	ParentSelect       []SelectableElement `json:"parent_select"`
	NameSearch         SearchField         `json:"name_search"`
	ParentSearch       SearchField         `json:"parent_search"`
	CoverageType       string              `json:"coverage_type"`
	NameSearchOutput   SearchOutput        `json:"name_search_output"`
	ParentSearchOutput SearchOutput        `json:"parent_search_output"`
	IsSelectParents    bool                `json:"is_select_parents"`
	OptionType         string              `json:"option_type"`
}

/* SearchOutput represents the presentable data required to display search output section
HasNoResults is a bool which displays messaging if there are no search results
SearchResults is an array of search results
Options is an array of previously added options
Language is the user set language */
type SearchOutput struct {
	HasNoResults  bool                `json:"has_no_results"`
	SearchResults []SelectableElement `json:"search_results"`
	Options       []SelectableElement `json:"options"`
	Language      string              `json:"language"`
	coreModel.Pagination
}

/* SelectableElement represents the data required for a selectable element.
Text is the human readable label.
Value is the value sent to the server.
Name is the name attribute.
IsSelected is a boolean representing whether the element is selected.
IsDisabled is a boolean representing whether the element is disabled */
type SelectableElement struct {
	Text       string `json:"text"`
	Value      string `json:"value"`
	Name       string `json:"name"`
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
