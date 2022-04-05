# iCON-metrics

This is an command line application which reads data from [NGBS iCON smart home control systems](https://www.ngbsh.hu/en/icon.html).

## Read data

Data is read from the control system(s).
To define the control systems to be scraped, [config file](config.yml) needs to be updated.

```yaml
devices:
  - url: http://192.168.1.10 # device address
    sysid: '123123123123' # device ID (printed on the controller)
    password: '123123123123' # password (same as sysid if empty)
    delay: 60 # delay in seconds between reads (default 60)
  - url: http://192.168.1.11 # device address
    sysid: '321321321321' # device ID (printed on the controller)
```

## Install on linux

Download the application from [GitHub releases](https://github.com/csutorasa/icon-metrics/releases) and unzip it.

```bash
wget https://github.com/csutorasa/icon-metrics/releases/download/1.0.0/icon-metrics-linux-amd64.zip
unzip icon-metrics-linux-amd64.zip
```

Create systemd service.

```bash
touch /etc/systemd/system/icon-metrics.service
```

```ini
[Unit]
Description=iCON metrics publisher

[Install]
WantedBy=multi-user.target

[Service]
Type=simple
ExecStart=/path/to/icon-metrics
WorkingDirectory=/path/to
Restart=always
```

[Automatic install script](linux_installer.sh)

```bash
curl -s https://raw.githubusercontent.com/csutorasa/icon-metrics/master/linux_installer.sh | bash -s amd64
```

## Metrics

Metrics are hosted in [prometheus](https://prometheus.io/) format.
The http server port can be configured in the [config file](config.yml).

```yaml
port: 8080
```

You need to extend the prometheus config (prometheus.yml) to scrape this application.

```yaml
scrape_configs:
  # Other scrape configs can be here

  - job_name: 'icon-metrics'
    static_configs:
    - targets: ['localhost:8080']
```

Available metrics:

| Metric                    | Type    | Description                                         |
| ------------------------- | ------- | --------------------------------------------------- |
| uptime                    | gauge   | uptime in milliseconds                              |
| icon_controller_connected | gauge   | 1 if the controller is ready to be read 0 otherwise |
| icon_http_client_seconds  | summary | icon HTTP request durations in seconds              |
| icon_external_temperature | gauge   | external temperature                                |
| icon_water_temperature    | gauge   | water temperature                                   |
| icon_temperature          | gauge   | room temperature                                    |
| icon_relay_on             | gauge   | 1 if room relay is open 0 otherwise                 |
| icon_humidity             | gauge   | room humidity                                       |
| icon_target_temperature   | gauge   | room target temperature                             |
| icon_dew_temperature      | gauge   | room dew temperature                                |

## Grafana dashboard

![grafana_image](https://user-images.githubusercontent.com/6968192/112884441-ce609c00-90cf-11eb-86e5-9dce7dab8e2a.png)
