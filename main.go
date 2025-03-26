package main

import (
	"log/slog"
	"maragu.dev/env"
	"net"
	"net/http"
	"os"
	"strconv"
	"uptimerobot_exporter/config"
	"uptimerobot_exporter/prober"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	_ = env.Load()

	host := env.GetStringOrDefault("HOST", "0.0.0.0")
	port := env.GetIntOrDefault("PORT", 9147)

	address := net.JoinHostPort(host, strconv.Itoa(port))

	envConfig := config.Config{
		UptimeKey:  env.GetStringOrDefault("API_KEY", ""),
		UptimeHost: env.GetStringOrDefault("API_HOST", "api.uptimerobot.com"),
	}

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		prober.Handler(w, r, envConfig, logger)
	})

	logger.Info("UptimeRobot Exporter Starting", "binding_address", address)

	if err := http.ListenAndServe(address, nil); err != nil {
		logger.Error("UptimeRobot Exporter Start failed", "binding_address", address)
	}

}
