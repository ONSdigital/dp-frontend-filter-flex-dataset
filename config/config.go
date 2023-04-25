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
	EnableCustomSort            bool          `envconfig:"ENABLE_CUSTOM_SORT"`
	EnableMultivariate          bool          `envconfig:"ENABLE_MULTIVARIATE"`
	GracefulShutdownTimeout     time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval         time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout  time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
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
		cfg.PatternLibraryAssetsPath = "//cdn.ons.gov.uk/dp-design-system/0c107a5"
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
		EnableCustomSort:            false,
		EnableMultivariate:          false,
		GracefulShutdownTimeout:     5 * time.Second,
		HealthCheckInterval:         30 * time.Second,
		HealthCheckCriticalTimeout:  90 * time.Second,
		SiteDomain:                  "localhost",
		SupportedLanguages:          []string{"en", "cy"},
	}

	return cfg, envconfig.Process("", cfg)
}
