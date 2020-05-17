# prometheus-exporter-merger

Merges Prometheus metrics from multiple sources.

## But Why?!

Sometimes you need to scrape Prometheus metrics from multiple containers in a single pod,
but you can't do this using annotations: [prometheus/prometheus#3756](https://github.com/prometheus/prometheus/issues/3756).

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
      keyX: valueX
      keyY: Y
  - url: http://127.0.0.1:8082/metrics
    labels:
      key2: Z
```

Another way to pass configuration by setting environment variables:

```bash
export LISTEN=":8080"
export SCRAPE_TIMEOUT="20s"
export URL_1=http://127.0.0.1:801/api/v1/metrics/prometheus,keyX:valueX,keyY:Y
export URL_2=http://0.0.0.0:7070/api/v1/metrics/prometheus,key2:Z
```

## Kubernetes

The prometheus-exporter-merger is supposed to run as a sidecar.

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
        env:
        - name: LISTEN
          value: :8080
        - name: SCRAPE_TIMEOUT
          value: 20s
        - name: URL_COMMON
          value: http://127.0.0.1:8081/api/v1/metrics/prometheus,type:common
        - name: URL_AUDIT
          value: http://127.0.0.1:8082/api/v1/metrics/prometheus,type:audit
...
```
