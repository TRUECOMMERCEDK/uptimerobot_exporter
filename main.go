package main

import (
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"log/slog"
	"maragu.dev/env"
	"net/http"
	"os"
)

var (
	addr = flag.String("web.listen-address", ":9147", "Address on which to expose metrics and web interface.")
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	_ = env.Load()

	apiKey := env.GetStringOrDefault("UPTIMEROBOT_API_KEY", "")

	if apiKey == "" {
		logger.Error("UPTIMEROBOT_API_KEY is not set")
	}

	logger.Info("Uptimerobot exporter Starting", "binding_address", addr)

	prometheus.MustRegister(NewCollector(apiKey))
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
