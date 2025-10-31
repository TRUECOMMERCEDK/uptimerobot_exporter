package collector

import (
	"net/http"
	"strconv"

	"log/slog"

	"github.com/prometheus/client_golang/prometheus"

	"uptimerobot_exporter/internal/uptimerobot"
)

// Collector implements the Prometheus collector interface for UptimeRobot.
type Collector struct {
	client     *uptimerobot.Client
	logger     *slog.Logger
	upDesc     *prometheus.Desc
	statusDesc *prometheus.Desc
}

// New creates a new UptimeRobot collector.
func New(apiKey string, httpClient *http.Client, logger *slog.Logger) prometheus.Collector {
	return &Collector{
		client: uptimerobot.NewClientWithHTTP(apiKey, httpClient),
		logger: logger,
		upDesc: prometheus.NewDesc(
			"uptimerobot_monitor_up",
			"Whether the UptimeRobot monitor is up (1) or down (0).",
			[]string{"id", "friendly_name", "url", "type"},
			nil,
		),
		statusDesc: prometheus.NewDesc(
			"uptimerobot_monitor_status",
			"Raw UptimeRobot monitor status code (2=up, 9=down, etc.).",
			[]string{"id", "friendly_name", "url", "type"},
			nil,
		),
	}
}

// Describe sends metric descriptors to Prometheus.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.upDesc
	ch <- c.statusDesc
}

// Collect fetches data from the UptimeRobot API and sends metrics.
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.logger.Debug("collecting monitor metrics from UptimeRobot")

	monitors, err := c.client.GetMonitors()
	if err != nil {
		c.logger.Error("failed to get monitors", "error", err)
		return
	}

	for _, m := range monitors {
		status := float64(m.Status)
		up := 0.0
		if m.Status == 2 {
			up = 1.0
		}

		c.logger.Debug("exporting monitor metric",
			"id", m.ID,
			"name", m.FriendlyName,
			"url", m.URL,
			"type", m.Type,
			"status", m.Status,
			"up", up,
		)

		ch <- prometheus.MustNewConstMetric(
			c.upDesc,
			prometheus.GaugeValue,
			up,
			strconv.Itoa(m.ID),
			m.FriendlyName,
			m.URL,
			strconv.Itoa(m.Type),
		)

		ch <- prometheus.MustNewConstMetric(
			c.statusDesc,
			prometheus.GaugeValue,
			status,
			strconv.Itoa(m.ID),
			m.FriendlyName,
			m.URL,
			strconv.Itoa(m.Type),
		)
	}

	c.logger.Debug("collection complete", "monitor_count", len(monitors))
}
