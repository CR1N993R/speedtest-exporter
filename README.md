# Speedtest Exporter

Speedtest Exporter is a straightforward Prometheus exporter written in Go specifically for the Speedtest.net CLI.

## Prerequisites

Before running the application, make sure to complete the following steps:

- **Install the [Speedtest.net CLI](https://www.speedtest.net/en/apps/cli)**

## Prometheus Configuration Example

Here's an example configuration snippet for Prometheus:

```yaml
scrape_configs:
  - job_name: 'speedtest'
    metrics_path: /metrics
    scrape_interval: 15m
    scrape_timeout: 120s 
    static_configs:
      - targets: ['127.0.0.1:9798']
```

## Metrics

The exporter exposes the following metrics obtained from Speedtest:

- `speedtest_download`: Download speed (B/s)
- `speedtest_upload`: Upload speed (B/s)
- `speedtest_jitter`: Jitter (ms)
- `speedtest_ping`: Ping (ms)

## Configuration Options

You can configure the application using the following environment variables:

### `PORT`

- **Description:** Sets the port for the webserver to listen on.
- **Default:** `9798` if not specified.

## Development

To start the application:

```shell
make start
```

To build the application to an executable:

```shell
make build
```