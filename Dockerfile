FROM golang:1.14-alpine3.11 as builder

RUN apk add --no-cache git make gcc libc-dev

WORKDIR /github.com/vadv/prometheus-exporter-merger
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build --ldflags "-s -w -linkmode external -extldflags -static" --tags netcgo -o /prometheus-exporter-merger

FROM scratch
USER nobody

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder  /prometheus-exporter-merger /prometheus-exporter-merger

EXPOSE 8080
CMD ["/prometheus-exporter-merger", "--config", "/config/prometheus-exporter-merger.yaml"]
