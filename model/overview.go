package model

import (
	coreModel "github.com/ONSdigital/dp-renderer/model"
)

type Overview struct {
	coreModel.Page
	FilterID   string      `json:"filter_id"`
	Dimensions []Dimension `json:"dimensions"`
}
