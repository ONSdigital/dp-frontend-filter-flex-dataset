package model

import (
	coreModel "github.com/ONSdigital/dp-renderer/v2/model"
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
	SetParent          string              `json:"set_parent"`
	FeedbackAPIURL     string              `json:"feedback_api_url"`
}
