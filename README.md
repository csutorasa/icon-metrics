# Icon-metrics

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
| icon_water_temperature    | gauge   | water temperature                                   |
| icon_temperature          | gauge   | room temperature                                    |
| icon_relay_on             | gauge   | 1 if room relay is open 0 otherwise                 |
| icon_humidity             | gauge   | room humidity                                       |
| icon_target_temperature   | gauge   | room target temperature                             |

## Grafana dashboard

![grafana_image](https://user-images.githubusercontent.com/6968192/112884128-6dd15f00-90cf-11eb-9b1e-072caac8391b.png)

