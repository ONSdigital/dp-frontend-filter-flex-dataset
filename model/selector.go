package model

import (
	coreModel "github.com/ONSdigital/dp-renderer/model"
)

type Selector struct {
	coreModel.Page
	Dimensions Dimension `json:"dimensions"`
}
