package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dp-frontend-filter-flex-dataset
type Config struct {
	Debug                       bool          `envconfig:"DEBUG"`
	EnableMultivariate          bool          `envconfig:"ENABLE_MULTIVARIATE"`
	EnableCustomSort            bool          `envconfig:"ENABLE_CUSTOM_SORT"`
	BindAddr                    string        `envconfig:"BIND_ADDR"`
	APIRouterURL                string        `envconfig:"API_ROUTER_URL"`
	DefaultMaximumSearchResults int           `envconfig:"DEFAULT_MAXIMUM_SEARCH_RESULTS"`
	PatternLibraryAssetsPath    string        `envconfig:"PATTERN_LIBRARY_ASSETS_PATH"`
	SupportedLanguages          []string      `envconfig:"SUPPORTED_LANGUAGES"`
	SiteDomain                  string        `envconfig:"SITE_DOMAIN"`
	GracefulShutdownTimeout     time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval         time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout  time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
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
		cfg.PatternLibraryAssetsPath = "//cdn.ons.gov.uk/dp-design-system/0c107a5"
	}

	return cfg, nil
}

func get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		Debug:                       false,
		EnableMultivariate:          false,
		EnableCustomSort:            false,
		BindAddr:                    "localhost:20100",
		APIRouterURL:                "http://localhost:23200/v1",
		DefaultMaximumSearchResults: 50,
		SupportedLanguages:          []string{"en", "cy"},
		SiteDomain:                  "localhost",
		GracefulShutdownTimeout:     5 * time.Second,
		HealthCheckInterval:         30 * time.Second,
		HealthCheckCriticalTimeout:  90 * time.Second,
	}

	return cfg, envconfig.Process("", cfg)
}
