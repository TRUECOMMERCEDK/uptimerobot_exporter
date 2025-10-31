# uptimerobot_exporter

Prometheus exporter for [UptimeRobot](https://uptimerobot.com) monitor statuses.

The exporter queries the UptimeRobot REST API and exposes monitor state and status metrics for Prometheus scraping.  

---

## Features

* Exposes the current status of all UptimeRobot monitors
* Exports both raw status codes and simple up/down values
* Secure HTTPS access with optional `--tls-skip-verify`
* Structured JSON or text logging (`--log.format`)
* Configurable log level (`--log.level`)
* Lightweight single-binary design with no dependencies

---

## Getting Started

### Prerequisites

You need:
* Go 1.24 or later
* An [UptimeRobot API key](https://uptimerobot.com/dashboard.php#mySettings)

### Run

```bash
go build -o uptimerobot_exporter ./cmd/uptimerobot_exporter
./uptimerobot_exporter \
  --api-key="xxxxxx" \
  --web.listen-address=":9149"
```

Open your browser and visit:

```bash
http://localhost:9149/metrics
```

You should see Prometheus metrics like:

```bash
# HELP uptimerobot_monitor_up Whether the UptimeRobot monitor is up (1) or down (0).
# TYPE uptimerobot_monitor_up gauge
uptimerobot_monitor_up{id="778899",friendly_name="Example Site",url="https://example.com",type="1"} 1

# HELP uptimerobot_monitor_status Raw UptimeRobot monitor status code (2=up, 9=down, etc.).
# TYPE uptimerobot_monitor_status gauge
uptimerobot_monitor_status{id="778899",friendly_name="Example Site",url="https://example.com",type="1"} 2
```

### Command-line Flags

| Flag                   | Default      | Description                                     |
| ---------------------- | ------------ | ----------------------------------------------- |
| `--web.listen-address` | `:9149`      | Address and port on which to expose metrics     |
| `--api-key`            | *(required)* | UptimeRobot API key for authentication          |
| `--tls-skip-verify`    | `false`      | Skip TLS certificate verification for API calls |
| `--log.format`         | `json`       | Log format: `json` or `text`                    |
| `--log.level`          | `info`       | Log level: `debug`, `info`, `warn`, or `error`  |

### Endpoints

| Path       | Description                            |
| ---------- | -------------------------------------- |
| `/metrics` | Prometheus metrics endpoint            |
| `/healthz` | Simple health check returning `200 OK` |

### Metrics

| Metric                       | Type  | Description                                              |
| ---------------------------- | ----- | -------------------------------------------------------- |
| `uptimerobot_monitor_up`     | Gauge | `1` if the monitor is up, `0` if down                    |
| `uptimerobot_monitor_status` | Gauge | Raw UptimeRobot monitor status code (2=up, 9=down, etc.) |


```console
0 - paused
1 - not checked yet
2 - up
8 - seems down
9 - down
```
