package main

import (
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"maragu.dev/env"
)

var (
	addr = flag.String("web.listen-address", ":9149", "Address on which to expose metrics and web interface.")
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
