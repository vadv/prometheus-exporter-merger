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
	mu           sync.Mutex
	scrapTimeout time.Duration
	client       *http.Client
	sources      []*source
}

type source struct {
	url    string
	labels []*prom.LabelPair
}

func New(scrapTimeout time.Duration) Merger {
	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   false,
			DisableCompression:  false,
			MaxIdleConns:        1,
			MaxIdleConnsPerHost: 1,
			MaxConnsPerHost:     10,
			IdleConnTimeout:     5 * time.Minute,
		},
		Timeout: scrapTimeout,
	}
	return &merger{
		scrapTimeout: scrapTimeout,
		client:       client,
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
	ctx, cancel := context.WithTimeout(context.Background(), m.scrapTimeout)
	defer cancel()
	return m.merge(ctx, w)
}
