package model

import coreModel "github.com/ONSdigital/dp-renderer/model"

// ChangeDimensions represents the data to display a ChangeDimensions page
type ChangeDimensions struct {
	coreModel.Page
	Output SearchOutput `json:"output"`
	Search SearchField  `json:"search"`
}
