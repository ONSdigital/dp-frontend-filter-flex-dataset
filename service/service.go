package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/dimension"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/health"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/assets"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/config"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/routes"
	render "github.com/ONSdigital/dp-renderer/v2"
	"github.com/ONSdigital/dp-renderer/v2/middleware/renderror"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

var (
	// BuildTime represents the time in which the service was built
	BuildTime string
	// GitCommit represents the commit (SHA-1) hash of the service that is running
	GitCommit string
	// Version represents the version of the service that is running
	Version string
)

// Service contains the healthcheck, server and serviceList for the controller
type Service struct {
	Config             *config.Config
	HealthCheck        HealthChecker
	Server             HTTPServer
	ServiceList        *ExternalServiceList
	routerHealthClient *health.Client
}

// New creates a new service
func New() *Service {
	return &Service{}
}

// Init initialises all the service dependencies, including healthcheck with checkers, api and middleware
func (svc *Service) Init(ctx context.Context, cfg *config.Config, serviceList *ExternalServiceList) (err error) {
	log.Info(ctx, "initialising service")

	svc.Config = cfg
	svc.ServiceList = serviceList

	// Get health client for api router
	svc.routerHealthClient = serviceList.GetHealthClient("api-router", cfg.APIRouterURL)

	dimensionClient, err := dimension.NewWithHealthClient(svc.routerHealthClient)
	if err != nil {
		return fmt.Errorf("failed to create dimensions API client: %w", err)
	}

	populationClient, err := population.NewWithHealthClient(svc.routerHealthClient)
	if err != nil {
		return fmt.Errorf("failed to create population API client: %w", err)
	}

	// Initialise clients
	clients := routes.Clients{
		Render:     render.NewWithDefaultClient(assets.Asset, assets.AssetNames, cfg.PatternLibraryAssetsPath, cfg.SiteDomain),
		Filter:     filter.NewWithHealthClient(svc.routerHealthClient),
		Dataset:    dataset.NewWithHealthClient(svc.routerHealthClient),
		Dimension:  dimensionClient,
		Population: populationClient,
		Zebedee:    zebedee.NewWithHealthClient(svc.routerHealthClient),
	}

	// Get healthcheck with checkers
	if svc.HealthCheck, err = serviceList.GetHealthCheck(cfg, BuildTime, GitCommit, Version); err != nil {
		return fmt.Errorf("failed to create health check: %w", err)
	}
	if err = svc.registerCheckers(ctx, clients); err != nil {
		return fmt.Errorf("failed to register checkers: %w", err)
	}
	clients.HealthCheckHandler = svc.HealthCheck.Handler

	// Initialise router
	r := mux.NewRouter()
	middleware := []alice.Constructor{
		renderror.Handler(clients.Render),
	}
	newAlice := alice.New(middleware...).Then(r)
	routes.Setup(ctx, r, cfg, clients)
	svc.Server = serviceList.GetHTTPServer(cfg.BindAddr, newAlice)

	return nil
}

// Run starts an initialised service
func (svc *Service) Run(ctx context.Context, svcErrors chan error) {
	log.Info(ctx, "Starting service", log.Data{"config": svc.Config})

	// Start healthcheck
	svc.HealthCheck.Start(ctx)

	// Start HTTP server
	log.Info(ctx, "Starting server")
	go func() {
		if err := svc.Server.ListenAndServe(); err != nil {
			log.Fatal(ctx, "failed to start http listen and serve", err)
			svcErrors <- err
		}
	}()
}

// Close gracefully shuts the service down in the required order, with timeout
func (svc *Service) Close(ctx context.Context) error {
	log.Info(ctx, "commencing graceful shutdown")
	ctx, cancel := context.WithTimeout(ctx, svc.Config.GracefulShutdownTimeout)
	hasShutdownError := false

	go func() {
		defer cancel()

		// stop healthcheck, as it depends on everything else
		log.Info(ctx, "stop health checkers")
		svc.HealthCheck.Stop()

		// TODO: close any backing services here, e.g. client connections to databases

		// stop any incoming requests
		if err := svc.Server.Shutdown(ctx); err != nil {
			log.Error(ctx, "failed to shutdown http server", err)
			hasShutdownError = true
		}
	}()

	// wait for shutdown success (via cancel) or failure (timeout)
	<-ctx.Done()

	// timeout expired
	if ctx.Err() == context.DeadlineExceeded {
		log.Error(ctx, "shutdown timed out", ctx.Err())
		return ctx.Err()
	}

	// other error
	if hasShutdownError {
		err := errors.New("failed to shutdown gracefully")
		log.Error(ctx, "failed to shutdown gracefully ", err)
		return err
	}

	log.Info(ctx, "graceful shutdown was successful")
	return nil
}

func (svc *Service) registerCheckers(ctx context.Context, c routes.Clients) error {
	hasErrors := false

	if err := svc.HealthCheck.AddCheck("API router", svc.routerHealthClient.Checker); err != nil {
		hasErrors = true
		log.Error(ctx, "failed to add API router checker", err)
	}

	if hasErrors {
		return errors.New("Error(s) registering checkers for healthcheck")
	}

	return nil
}
