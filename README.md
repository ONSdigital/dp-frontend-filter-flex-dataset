# dp-frontend-filter-flex-dataset

[![GitHub release](https://img.shields.io/github/release/ONSdigital/dp-frontend-filter-flex-dataset.svg)](https://github.com/ONSdigital/dp-frontend-filter-flex-dataset/releases)

# Frontend service to host filter and flexing of datasets

Frontend service to host filter, flexing and rendering the templates for datasets

## Getting started

- Run `make debug`

## Dependencies

- No further dependencies other than those defined in `go.mod`

## Configuration

| Environment variable           | Default                     | Description                                                                                                                                           |
| ------------------------------ | --------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- |
| API_ROUTER_URL                 | <http://localhost:23200/v1> | The URL of the [dp-api-router](https://github.com/ONSdigital/dp-api-router)                                                                           |
| BIND_ADDR                      | :20100                      | The host and port to bind to                                                                                                                          |
| DEBUG                          | false                       | Enable debug mode                                                                                                                                     |
| DEFAULT_MAXIMUM_SEARCH_RESULTS | 50                          | Maximum paginated search results                                                                                                                      |
| ENABLE_MULTIVARIATE            | false                       | Enable 2021 [multivariate datasets](https://github.com/ONSdigital/dp-dataset-api/blob/5f9f4218b65aae4803809f4a876e9f72b9bf5305/models/dataset.go#L43) |
| GRACEFUL_SHUTDOWN_TIMEOUT      | 5s                          | The graceful shutdown timeout in seconds (`time.Duration` format)                                                                                     |
| HEALTHCHECK_CRITICAL_TIMEOUT   | 90s                         | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format)                                    |
| HEALTHCHECK_INTERVAL           | 30s                         | Time between self-healthchecks (`time.Duration` format)                                                                                               |
| PATTERN_LIBRARY_ASSETS_PATH    | ""                          | Pattern library location                                                                                                                              |
| SUPPORTED_LANGUAGES            | []string{"en", "cy"}        | Supported languages                                                                                                                                   |
| SITE_DOMAIN                    | localhost                   |                                                                                                                                                       |

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

## License

Copyright © 2022, Office for National Statistics (<https://www.ons.gov.uk>)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
