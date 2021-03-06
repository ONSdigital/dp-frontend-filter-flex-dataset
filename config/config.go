package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dp-frontend-filter-flex-dataset
type Config struct {
	Debug                      bool          `envconfig:"DEBUG"`
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	APIRouterURL               string        `envconfig:"API_ROUTER_URL"`
	PatternLibraryAssetsPath   string        `envconfig:"PATTERN_LIBRARY_ASSETS_PATH"`
	SupportedLanguages         []string      `envconfig:"SUPPORTED_LANGUAGES"`
	SiteDomain                 string        `envconfig:"SITE_DOMAIN"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
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
		cfg.PatternLibraryAssetsPath = "//cdn.ons.gov.uk/dp-design-system/7129615"
	}

	return cfg, nil
}

func get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		Debug:                      false,
		BindAddr:                   "localhost:20100",
		APIRouterURL:               "http://localhost:23200/v1",
		SupportedLanguages:         []string{"en", "cy"},
		SiteDomain:                 "localhost",
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
	}

	return cfg, envconfig.Process("", cfg)
}
