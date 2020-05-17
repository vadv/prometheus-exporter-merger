package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	prom "github.com/prometheus/client_model/go"
	"github.com/vadv/prometheus-exporter-merger/merger"
)

func main() {

	var (
		configPath = flag.String("config", "config.yaml", "Path to config")
	)
	flag.Parse()

	c, err := parseConfig(*configPath)
	if err != nil {
		panic(err)
	}

	m := merger.New(c.Timeout)
	for _, s := range c.Sources {
		var labels []*prom.LabelPair
		for k, v := range s.Labels {
			k, v := k, v
			labels = append(labels, &prom.LabelPair{Name: &k, Value: &v})
		}
		m.AddSource(s.Url, labels)
	}

	srv := &http.Server{Addr: c.Listen, Handler: &handler{m: m}}
	go srv.ListenAndServe()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		panic(err)
	}
}
