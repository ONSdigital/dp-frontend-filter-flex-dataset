package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/config"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/service"
	"github.com/ONSdigital/log.go/v2/log"

	dpotelgo "github.com/ONSdigital/dp-otel-go"
)

func main() {
	log.Namespace = "dp-frontend-filter-flex-dataset"
	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Fatal(ctx, "application unexpectedly failed", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func run(ctx context.Context) error {
	// Create error channel for os signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Create service initialiser and an error channel for service errors
	svcList := service.NewServiceList(&service.Init{})
	svcErrors := make(chan error, 1)

	// Read config
	cfg, err := config.Get()
	if err != nil {
		return fmt.Errorf("unable to retrieve service configuration: %w", err)
	}

	log.Info(ctx, "got service configuration", log.Data{"config": cfg})

	// Set up Open Telemetry
	otelConfig := dpotelgo.Config{
		OtelServiceName:          cfg.OTServiceName,
		OtelExporterOtlpEndpoint: cfg.OTExporterOTLPEndpoint,
	}

	otelShutdown, oErr := dpotelgo.SetupOTelSDK(ctx, otelConfig)
	if oErr != nil {
		return fmt.Errorf("error setting up OpenTelemetry - hint: ensure OTEL_EXPORTER_OTLP_ENDPOINT is set. %w", oErr)
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	// Run service
	svc := service.New()
	if err := svc.Init(ctx, cfg, svcList); err != nil {
		return fmt.Errorf("failed to initialise service: %w", err)
	}
	svc.Run(ctx, svcErrors)

	// Blocks until an os interrupt or a fatal error occurs
	select {
	case err := <-svcErrors:
		log.Error(ctx, "service error received", err)
	case sig := <-signals:
		log.Info(ctx, "os signal received", log.Data{"signal": sig})
	}

	return svc.Close(ctx)
}
