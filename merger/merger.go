package merger

import (
	"context"
	"io"
	"net/http"
	"sync"
	"time"

	prom "github.com/prometheus/client_model/go"
)

type Merger interface {
	Merge(w io.Writer) error
	AddSource(url string, labels []*prom.LabelPair)
}

type merger struct {
	mu      sync.Mutex
	timeout time.Duration
	client  *http.Client
	sources []*source
}

type source struct {
	url    string
	labels []*prom.LabelPair
}

func New(timeout time.Duration) Merger {
	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   false,
			DisableCompression:  false,
			MaxIdleConns:        1,
			MaxIdleConnsPerHost: 1,
			MaxConnsPerHost:     10,
			IdleConnTimeout:     5 * time.Minute,
		},
		Timeout: timeout,
	}
	return &merger{
		timeout: timeout,
		client:  client,
	}
}

// AddSource new source
func (m *merger) AddSource(url string, labels []*prom.LabelPair) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sources = append(m.sources, &source{url: url, labels: labels})
}

// Merge sources
func (m *merger) Merge(w io.Writer) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	return m.merge(ctx, w)
}
