# prometheus-exporter-merger

Merges Prometheus metrics from multiple sources.

## But Why?!

> [prometheus/prometheus#3756](https://github.com/prometheus/prometheus/issues/3756)

To start the exporter:

```
prometheus-exporter-merger --config /config/prometheus-exporter-merger.yaml
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

## Kubernetes

The prometheus-exporter-merger is supposed to run as a sidecar.
By default, config must be available in the container by the path: `/config/prometheus-exporter-merger.yaml`.

```yaml
...
  template:
    metadata:
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "8080"
...
    spec:
      containers:
...
      - name: prometheus-exporter-merger
        image: vadv/prometheus-exporter-merger
        volumeMounts:
        - name: config
          mountPath: /config
...
```
