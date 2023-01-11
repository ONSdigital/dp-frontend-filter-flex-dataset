package model

import (
	coreModel "github.com/ONSdigital/dp-renderer/model"
)

// Selector represents page data for the Dimension selection screen
type Selector struct {
	coreModel.Page
	Dimensions       Dimension `json:"dimensions"`
	Selections       []Selection
	InitialSelection string
	IsAreaType       bool
	HasOptions       bool
	LeadText         string `json:"lead_text"`
	ErrorId          string `json:"error_id"`
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
}
