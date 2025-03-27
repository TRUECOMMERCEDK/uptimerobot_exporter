# uptimerobot_exporter





## Installation
```console
sudo useradd --no-create-home --shell /bin/false uptimerobotexporter
sudo mkdir /opt/uptimerobot_exporter
cd /opt/uptimerobot_exporter
sudo tar -xvf uptimerobot_exporter_0.0.1_linux_amd64.tar.gz
sudo chmod 755 uptimeexporterserver
sudo chown uptimerobotexporter:uptimerobotexporter /opt/uptimerobot_exporter/*
sudo ln -s /opt/uptimerobot_exporter/uptimeexporterserver /usr/local/bin/uptimeexporterserver

sudo tee /etc/systemd/system/uptimerobot_exporter.service <<EOF
[Unit]
Description=UptimeRobot Exporter
Wants=network-online.target
After=network-online.target

[Service]
Type=simple
WorkingDirectory=/opt/uptimerobot_exporter
ExecStart=/usr/local/bin/uptimeexporterserver

Restart=always

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl enable --now uptimerobot_exporter.service
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
