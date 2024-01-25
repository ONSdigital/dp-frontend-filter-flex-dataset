package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dp-frontend-filter-flex-dataset
type Config struct {
	APIRouterURL                string        `envconfig:"API_ROUTER_URL"`
	BindAddr                    string        `envconfig:"BIND_ADDR"`
	Debug                       bool          `envconfig:"DEBUG"`
	DefaultMaximumSearchResults int           `envconfig:"DEFAULT_MAXIMUM_SEARCH_RESULTS"`
	EnableMultivariate          bool          `envconfig:"ENABLE_MULTIVARIATE"`
	GracefulShutdownTimeout     time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval         time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout  time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	OTBatchTimeout              time.Duration `encconfig:"OTEL_BATCH_TIMEOUT"`
	OTServiceName               string        `envconfig:"OTEL_SERVICE_NAME"`
	OTExporterOTLPEndpoint      string        `envconfig:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	PatternLibraryAssetsPath    string        `envconfig:"PATTERN_LIBRARY_ASSETS_PATH"`
	SiteDomain                  string        `envconfig:"SITE_DOMAIN"`
	SupportedLanguages          []string      `envconfig:"SUPPORTED_LANGUAGES"`
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	cfg, err := get()
	if err != nil {
		return nil, err
	}

	if cfg.Debug {
		cfg.PatternLibraryAssetsPath = "http://localhost:9002/dist/assets"
	} else {
		cfg.PatternLibraryAssetsPath = "//cdn.ons.gov.uk/dp-design-system/e0a75c3"
	}

	return cfg, nil
}

func get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		APIRouterURL:                "http://localhost:23200/v1",
		BindAddr:                    "localhost:20100",
		Debug:                       false,
		DefaultMaximumSearchResults: 50,
		EnableMultivariate:          false,
		GracefulShutdownTimeout:     5 * time.Second,
		HealthCheckInterval:         30 * time.Second,
		HealthCheckCriticalTimeout:  90 * time.Second,
		OTBatchTimeout:              5 * time.Second,
		OTExporterOTLPEndpoint:      "localhost:4317",
		OTServiceName:               "dp-frontend-filter-flex-dataset",
		SiteDomain:                  "localhost",
		SupportedLanguages:          []string{"en", "cy"},
	}

	return cfg, envconfig.Process("", cfg)
}
