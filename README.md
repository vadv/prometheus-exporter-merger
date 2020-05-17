# prometheus-exporter-merger

Merges Prometheus metrics from multiple sources.

## But Why?!

> [prometheus/prometheus#3756](https://github.com/prometheus/prometheus/issues/3756)

To start the exporter:

```
prometheus-exporter-merger --config config.yaml
```

Config example:

```yaml
listen: :8080
scrap_timeout: 20s
sources:
  - url: http://127.0.0.1:8081/metrics
    labels:
      key1: value1
  - url: http://127.0.0.1:8082/metrics
    labels:
      key2: value2
```