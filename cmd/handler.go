package main

import (
	"log"
	"net/http"

	"github.com/vadv/prometheus-exporter-merger/merger"
)

type handler struct {
	m merger.Merger
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/healthz":
		w.WriteHeader(http.StatusOK)
	default:
		err := h.m.Merge(w)
		if err != nil {
			log.Printf("[ERROR] %s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
