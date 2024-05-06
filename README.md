# Speedtest Exporter

This is a simple prometheus exporter for the Speedtest.net cli.

## Prerequisites

Before running the application the following things need to be done.

- Install the [Speedtest.net cli](https://www.speedtest.net/en/apps/cli)

## Prometheus example

```yaml
scrape_configs:
  - job_name: 'speedtest'
    metrics_path: /metrics
    scrape_interval: 15m
    scrape_timeout: 120s 
    static_configs:
      - targets: ['127.0.0.1:9798']
```

## Configuration

The following environment variables can be set to configure the application.

### PORT

The `PORT` environment variable sets the port on which the webserver starts. If non is set the default port is `9798`.

