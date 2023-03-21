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
	ShowGetDataButton     bool        `json:"show_get_data_button"`
	DisableGetDataButton  bool        `json:"disable_get_data_button"`
	HasSDC                bool        `json:"has_sdc"`
	MaxVariableError      bool        `json:"max_variable_error"`
	ImproveResults        coreModel.Collapsible
	DimensionDescriptions coreModel.Collapsible
}
