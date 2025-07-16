# uptimerobot_exporter

```console
0 - paused
1 - not checked yet
2 - up
8 - seems down
9 - down
```

## Installation
```console
sudo mkdir /opt/uptimerobot_exporter
cd /opt/uptimerobot_exporter
sudo tar -xvf uptimerobot_exporter_0.0.1_linux_amd64.tar.gz
sudo chmod 755 server
sudo chown prometheus:prometheus /opt/uptimerobot_exporter/*

sudo tee /etc/systemd/system/uptimerobot_exporter.service <<EOF
[Unit]
Description=UptimeRobot Exporter
Wants=network-online.target
After=network-online.target

[Service]
Type=simple
WorkingDirectory=/opt/uptimerobot_exporter
ExecStart=/opt/uptimerobot_exporter/server

Restart=always

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl enable --now uptimerobot_exporter.service
```

## Metrics exported
```console
# HELP uptimerobot_monitor_up Status of the UptimeRobot monitor
# TYPE uptimerobot_monitor_up gauge
uptimerobot_monitor_up{friendly_name="Google",id="1",type="1",url="https://www.google.com"} 2

Possible values for status:

0 = paused
1 = not checked yet
2 = up
8 = seems down
9 = down
```

## Prometheus configuration
```yaml
  - job_name: 'uptimerobot_exporter'
    static_configs:
      - targets:
          - localhost:9147
    relabel_configs:
      - source_labels: [__address__]
        regex: "([^:]+):\\d+"
        target_label: instance
```

## Filebeat configuration
```console
- type: journald
  enabled: true
  pipeline: filebeat
  id: service-uptimerobot-exporter
  include_matches.match:
    - _SYSTEMD_UNIT=uptimerobot_exporter.service
  fields:
    type: uptimerobot_exporter

  parsers:
    - ndjson:
      overwrite_keys: true
      add_error_key: true
      expand_keys: true
```
