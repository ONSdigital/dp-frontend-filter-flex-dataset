package routes

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
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
}

// Setup registers routes for the service
func Setup(ctx context.Context, r *mux.Router, cfg *config.Config, c Clients) {
	log.Info(ctx, "adding routes")
	r.StrictSlash(true).Path("/health").HandlerFunc(c.HealthCheckHandler)

	r.StrictSlash(true).Path("/filters/{filterID}/dimensions").Methods("GET").HandlerFunc(handlers.FilterFlexOverview(*cfg, c.Render, c.Filter))
	r.StrictSlash(true).Path("/filters/{filterID}/dimensions/{name}").Methods("GET").HandlerFunc(handlers.DimensionsSelector(*cfg, c.Render))
}
