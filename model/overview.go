package model

import (
	coreModel "github.com/ONSdigital/dp-renderer/model"
)

// Overview represents the data to display the overview page
type Overview struct {
	coreModel.Page
	FilterID              string      `json:"filter_id"`
	Panel                 Panel       `json:"panel"`
	Dimensions            []Dimension `json:"dimensions"`
	IsMultivariate        bool        `json:"is_multivariate"`
	HasSDC                bool        `json:"has_sdc"`
	ImproveResults        coreModel.Collapsible
	DimensionDescriptions coreModel.Collapsible
}
