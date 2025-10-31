package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"uptimerobot_exporter/internal/collector"
	"uptimerobot_exporter/internal/logging"
)

var (
	listenAddr    = flag.String("web.listen-address", ":9149", "Address on which to expose metrics and web interface.")
	apiKey        = flag.String("api-key", "", "UptimeRobot API key for authentication.")
	tlsSkipVerify = flag.Bool("tls-skip-verify", false, "Skip TLS certificate verification for UptimeRobot API requests.")

	flagLogFormat = flag.String("log.format", "text", "Log format: json or text.")
	flagLogLevel  = flag.String("log.level", "info", "Log level: debug, info, warn, error.")
)

var release = "dev"

func main() {
	flag.Parse()

	logger := logging.NewWithOptions(*flagLogFormat, *flagLogLevel)
	slog.SetDefault(logger)

	if *apiKey == "" {
		logger.Error("missing required --api-key flag, exiting")
		os.Exit(1)
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: *tlsSkipVerify}, // #nosec G402 — intentional flag
		},
	}

	logger.Info("Starting uptimeRobot exporter",
		slog.String("version", release),
		slog.String("addr", *listenAddr),
		slog.Bool("tls-skip-verify", *tlsSkipVerify),
		slog.String("log_format", *flagLogFormat),
		slog.String("log_level", *flagLogLevel))

	// Register Prometheus collector
	prometheus.MustRegister(collector.New(*apiKey, httpClient, logger))

	// Expose metrics and health endpoints
	http.Handle("GET /metrics", promhttp.Handler())
	http.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"version": release,
		})
		if err != nil {
			return
		}
	})

	logger.Info("listening for scrape requests", "addr", *listenAddr)
	if err := http.ListenAndServe(*listenAddr, nil); err != nil {
		logger.Error("http server error", "error", err)
		os.Exit(1)
	}
}
