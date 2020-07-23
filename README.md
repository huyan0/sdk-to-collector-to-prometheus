# Prometheus Go Exporter Demo

This demo aims to check whether the Collector Prometheus Exporter successfully sends cumulative metircs retrieved from the Go SDK.

## Setup
[Download OpenTelemetry Collector](https://github.com/open-telemetry/opentelemetry-collector/releases)

[Download Prometheus](https://prometheus.io/download/)

## Running
### Run Prometheus
    prometheus --config.file prometheus.yml

This commmand starts a Prometheus instance that scrapes from `localhost:8800` once per second

#Run OT Collector:

    otel-col --config collector.config

This command starts a Collector instance that receives OTLP metrics and export the metric to Prometheus using the Prometheus exporter.

Generate Metric:

    go run main.go

This command starts generating metric events from two counter and a valuerecorder instument, and sends metrics using Go SDK's OTLP exporter once per second. 

## Checking Values
Access the Prometheus dashboard at `localhost:9090`. Select the graph option and `a_counter_1` or `a_counter_i` as the expression, then press Execute to see the updates. Note that the graph does not update automatically and that updates will not appear until the Prometheus instance says in the console that it is ready to begin retrieving data.
