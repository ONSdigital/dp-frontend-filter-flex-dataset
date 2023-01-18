package routes

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/dimension"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/config"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/handlers"
	render "github.com/ONSdigital/dp-renderer"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// Clients - struct containing all the clients for the controller
type Clients struct {
	HealthCheckHandler func(w http.ResponseWriter, req *http.Request)
	Render             *render.Render
	Filter             *filter.Client
	Dataset            *dataset.Client
	Dimension          *dimension.Client
	Population         *population.Client
	Zebedee            *zebedee.Client
}

// Setup registers routes for the service
func Setup(ctx context.Context, r *mux.Router, cfg *config.Config, c Clients) {
	log.Info(ctx, "adding routes")

	ff := handlers.NewFilterFlex(c.Render, c.Filter, c.Dataset, c.Population, c.Zebedee, cfg)

	r.StrictSlash(true).Path("/health").HandlerFunc(c.HealthCheckHandler)

	r.StrictSlash(true).Path("/filters/{filterID}/submit").Methods("POST").HandlerFunc(ff.Submit())

	r.StrictSlash(true).Path("/filters/{filterID}/dimensions").Methods("GET").HandlerFunc(ff.FilterFlexOverview())
	if cfg.EnableMultivariate {
		r.StrictSlash(true).Path("/filters/{filterID}/dimensions/change").Methods("GET").HandlerFunc(ff.GetChangeDimensions())
		r.StrictSlash(true).Path("/filters/{filterID}/dimensions/change").Methods("POST").HandlerFunc(ff.PostChangeDimensions())
	}
	r.StrictSlash(true).Path("/filters/{filterID}/dimensions/{name}").Methods("GET").HandlerFunc(ff.DimensionSelector())
	r.StrictSlash(true).Path("/filters/{filterID}/dimensions/{name}").Methods("POST").HandlerFunc(ff.ChangeDimension())

	r.StrictSlash(true).Path("/filters/{filterID}/dimensions/geography/coverage").Methods("GET").HandlerFunc(ff.GetCoverage())
	r.StrictSlash(true).Path("/filters/{filterID}/dimensions/geography/coverage").Methods("POST").HandlerFunc(ff.UpdateCoverage())
}
