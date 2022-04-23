# iCON-metrics

This is an command line application which reads data from [NGBS iCON smart home control systems](https://www.ngbsh.hu/en/icon.html).
The data is processed and is exposed to be scraped from [prometheus](https://prometheus.io/).
Data from prometheus can be display in different ways, but [grafana](https://grafana.com/) is recommended.

## Read data

Data is read from the control system(s).
To define the control systems to be scraped, [config file](config.yml) needs to be updated.

```yaml
port: 8080 # port to run on (defaults to 80)
devices:
  - url: http://192.168.1.10 # device address
    sysid: '123123123123' # device ID (printed on the controller)
    password: '123123123123' # password (defaults to sysid if empty)
    delay: 15 # delay in seconds between reads (defaults to 60)
  - url: http://192.168.1.11 # device address
    sysid: '321321321321' # device ID (printed on the controller)
```

Config.yml validation can be done via the [schema](config.schema.json).
For further configuration options use the [schema](config.schema.json) to explore and validate your config file.

## Install on linux

Download the application from [GitHub releases](https://github.com/csutorasa/icon-metrics/releases) and unzip it.

```bash
wget https://github.com/csutorasa/icon-metrics/releases/download/1.2.0/icon-metrics-linux-amd64.zip
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
StandardOutput=/var/log/icon-metrics/log.log
StandardError=/var/log/icon-metrics/error.log
Restart=always
```

[Automatic install script](linux_installer.sh)

```bash
curl -s https://raw.githubusercontent.com/csutorasa/icon-metrics/master/linux_installer.sh | sudo bash -s amd64
```

## Docker image

Create your config.yml and run the command below

```bash
docker run -d -v /path/to/your/config.yml:/app/config.yml -p8080:8080 csutorasa/icon-metrics:latest
```

If you do not want use 8080 port then use `-p${YOUR_PORT}:8080`.

## Metrics

### Prometheus scraper

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

### Metrics reporting

Most metrics can be disabled from the configuaration separately for each device in the [config file](config.yml).

Available metrics:

| Metric                    | Type    | Description                                          | Enable configuration flag |
| ------------------------- | ------- | ---------------------------------------------------- | ------------------------- |
| uptime                    | gauge   | uptime in milliseconds                               | N/A                       |
| icon_controller_connected | gauge   | 1 if the controller is ready to be read 0 otherwise  | N/A                       |
| icon_room_connected       | gauge   | 1 if the room is connected to the device 0 otherwise | roomConnected             |
| icon_http_client_seconds  | summary | icon HTTP request durations in seconds               | httpClient                |
| icon_external_temperature | gauge   | external temperature                                 | externalTemperature       |
| icon_water_temperature    | gauge   | water temperature                                    | waterTemperature          |
| icon_temperature          | gauge   | room temperature                                     | temperature               |
| icon_relay_on             | gauge   | 1 if room relay is open 0 otherwise                  | relay                     |
| icon_humidity             | gauge   | room humidity                                        | humidity                  |
| icon_target_temperature   | gauge   | room target temperature                              | targetTemperature         |
| icon_dew_temperature      | gauge   | room dew temperature                                 | dewTemperature            |

## Grafana dashboard

This data is designed to be displayed in a [grafana dashboard](https://grafana.com/docs/grafana/latest/dashboards/).
An [example dashboard](grafana.json) is available to be [imported](https://grafana.com/docs/grafana/latest/dashboards/export-import/).

![grafana_image](https://user-images.githubusercontent.com/6968192/164945203-c4c97804-d30a-498e-9c6f-9cbee4633f12.png)
