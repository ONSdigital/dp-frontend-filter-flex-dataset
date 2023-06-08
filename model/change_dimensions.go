package model

import coreModel "github.com/ONSdigital/dp-renderer/v2/model"

// ChangeDimensions represents the data to display a ChangeDimensions page
type ChangeDimensions struct {
	coreModel.Page
	Output           SearchOutput `json:"output"`
	SearchOutput     SearchOutput `json:"search_output"`
	Search           SearchField  `json:"search"`
	FormAction       string       `json:"form_action"`
	Panel            Panel        `json:"panel"`
	HasSDC           bool         `json:"has_sdc"`
	MaxVariableError bool         `json:"max_variable_error"`
	ImproveResults   coreModel.Collapsible
}
