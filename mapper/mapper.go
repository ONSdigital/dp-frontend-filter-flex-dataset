package mapper

import (
	"context"

	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/config"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/model"
	coreModel "github.com/ONSdigital/dp-renderer/model"
)

func CreateFilterFlexOverview(ctx context.Context, basePage coreModel.Page, cfg config.Config) model.Overview {
	m := model.Overview{
		Page: basePage,
	}
	m.BetaBannerEnabled = true
	m.Type = "filter-flex-overview"
	m.Metadata.Title = "Filter flex overview"

	return m
}
