package model

import (
	coreModel "github.com/ONSdigital/dp-renderer/v2/model"
)

// Selector represents page data for the Dimension selection screen
type Selector struct {
	coreModel.Page
	Dimensions        Dimension `json:"dimensions"`
	Selections        []Selection
	InitialSelection  string
	IsAreaType        bool
	LeadText          string `json:"lead_text"`
	ErrorId           string `json:"error_id"`
	Panel             Panel  `json:"panel"`
	EnableFeedbackAPI bool   `json:"enable_feedback_api"`
	FeedbackAPIURL    string `json:"feedback_api_url"`
}

// Selection represents a dimension selection (e.g. an Area-type of City)
type Selection struct {
	Value           string
	Label           string
	Description     string
	Categories      []string
	CategoriesCount int
	TotalCount      int
	IsTruncated     bool
	TruncateLink    string
	IsSuggested     bool
}
