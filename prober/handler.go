package prober

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"time"
	"uptimerobot_exporter/config"
)

type Monitors struct {
	Monitors []struct {
		ID              int    `json:"id"`
		FriendlyName    string `json:"friendly_name"`
		URL             string `json:"url"`
		Type            int    `json:"type"`
		SubType         string `json:"sub_type"`
		KeywordType     int    `json:"keyword_type"`
		KeywordCaseType int    `json:"keyword_case_type"`
		KeywordValue    string `json:"keyword_value"`
		Port            string `json:"port"`
		Interval        int    `json:"interval"`
		Timeout         int    `json:"timeout"`
		Status          int    `json:"status"`
		CreateDatetime  int    `json:"create_datetime"`
	} `json:"monitors"`
}

func Handler(w http.ResponseWriter, r *http.Request, c config.Config, logger *slog.Logger) {

	MonitorStatusGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "monitor_status",
			Help:      "Displays monitor status",
			Namespace: "uptimerobot",
		},
		[]string{
			"id",
			"friendly_name",
		},
	)

	registry := prometheus.NewRegistry()
	registry.MustRegister(MonitorStatusGauge)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			DialContext: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			DisableKeepAlives:     true,
			MaxIdleConnsPerHost:   -1,
		},
	}

	addr := fmt.Sprintf("%s", c.UptimeHost)
	req, err := http.NewRequest("POST", "https://"+addr+"/v2/getMonitors?format=json&api_key="+c.UptimeKey, nil)

	req.Close = true

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error getting data", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("http status code is not 200", err)
		return
	}

	bodyText, err := io.ReadAll(resp.Body)

	var monitors Monitors

	err = json.Unmarshal(bodyText, &monitors)
	if err != nil {
		fmt.Println(err)

	}

	for _, v := range monitors.Monitors {
		MonitorStatusGauge.WithLabelValues(strconv.Itoa(v.ID), v.FriendlyName).Set(float64(v.Status))
	}

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
	return

}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
