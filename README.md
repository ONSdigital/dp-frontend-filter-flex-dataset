# dp-frontend-filter-flex-dataset

Frontend service to host filter and flexing of datasets

## Getting started

* Run `make debug`

## Dependencies

* No further dependencies other than those defined in `go.mod`

## Configuration

| Environment variable         | Default   | Description
| ---------------------------- | --------- | -----------
| DEBUG                        | false     | Enable debug mode
| BIND_ADDR                    | :20100    | The host and port to bind to
| API_ROUTER_URL               | <http://localhost:23200/v1> | The URL of the [dp-api-router](https://github.com/ONSdigital/dp-api-router)
| PATTERN_LIBRARY_ASSETS_PATH  | ""        | Pattern library location
| SUPPORTED_LANGUAGES          | []string{"en", "cy"}   | Supported languages
| SITE_DOMAIN                  | localhost |
| GRACEFUL_SHUTDOWN_TIMEOUT    | 5s        | The graceful shutdown timeout in seconds (`time.Duration` format)
| HEALTHCHECK_INTERVAL         | 30s       | Time between self-healthchecks (`time.Duration` format)
| HEALTHCHECK_CRITICAL_TIMEOUT | 90s       | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format)

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

## License

Copyright © 2021, Office for National Statistics (<https://www.ons.gov.uk>)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
