package model

import "github.com/ONSdigital/dp-api-clients-go/v2/filter"

// Dimension represents the data for a single dimension
type Dimension struct {
	Options        []string `json:"options"`
	IsTruncated    bool     `json:"is_truncated"`
	TruncateLink   string   `json:"truncate_link"`
	OptionsCount   int      `json:"options_count"`
	Name           string   `json:"name"`
	ID             string   `json:"id"`
	URI            string   `json:"uri"`
	IsGeography    bool     `json:"is_geography"`
	HasCategories  bool     `json:"has_categories"`
	HasChange      bool     `json:"has_change"`
	FeedbackAPIURL string   `json:"feedback_api_url"`
}

// FilterDimension represents a DTO for filter.Dimension with the additional OptionsCount field
type FilterDimension struct {
	filter.Dimension
	OptionsCount        int
	CategorisationCount int
}
