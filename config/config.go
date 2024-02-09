package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Configuration root
type Configuration struct {
	Port    int                  `yaml:"port"`
	Devices []*IconConfiguration `yaml:"devices"`
}

// iCON device configuration
type IconConfiguration struct {
	Url      string               `yaml:"url"`
	SysId    string               `yaml:"sysid"`
	Password string               `yaml:"password"`
	Delay    int                  `yaml:"delay"`
	Report   *ReportConfiguration `yaml:"report"`
}

// iCON device report configuration
type ReportConfiguration struct {
	// metrics.HttpGauge
	HttpClient *bool `yaml:"httpClient"`
	// metrics.WaterTemperatureGauge
	WaterTemperature *bool `yaml:"waterTemperature"`
	// metrics.ExternalTemperatureGauge
	ExternalTemperature *bool `yaml:"externalTemperature"`
	// metrics.RoomConntectedGauge
	RoomConnected *bool `yaml:"roomConnected"`
	// metrics.RoomTemperatureGauge
	Temperature *bool `yaml:"temperature"`
	// metrics.RoomDewTemperatureGauge
	DewTemperature *bool `yaml:"dewTemperature"`
	// metrics.RelayGauge
	Relay *bool `yaml:"relay"`
	// metrics.HumidityGauge
	Humidity *bool `yaml:"humidity"`
	// metrics.TargetTemperatureGauge
	TargetTemperature *bool `yaml:"targetTemperature"`
}

// Returns the config that is read from the file.
func ReadConfig(filepath string) (*Configuration, error) {
	config := &Configuration{}
	data, err := os.ReadFile(filepath)
	if err != nil {
		return config, fmt.Errorf("failed to read config: %w", err)
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("failed to parse yaml: %w", err)
	}
	err = validateConfig(config)
	if err != nil {
		return config, fmt.Errorf("invalid config: %w", err)
	}
	return config, nil
}

// Scans the config for invalid settings.
func validateConfig(config *Configuration) error {
	if config.Port == 0 {
		config.Port = 80
	}
	if config.Devices == nil || len(config.Devices) == 0 {
		return errors.New("there are no devices to monitor")
	}
	for i, device := range config.Devices {
		if device.SysId == "" {
			return fmt.Errorf("device config at %d position is missing sysid", i)
		}
		if device.Url == "" {
			return fmt.Errorf("device config at %d position is missing url", i)
		}
		if device.Password == "" {
			device.Password = device.SysId
		}
		if device.Delay == 0 {
			device.Delay = 15
		}
		if device.Report == nil {
			defaultReport := ReportConfiguration{
				HttpClient:          enabled(),
				WaterTemperature:    enabled(),
				ExternalTemperature: enabled(),
				RoomConnected:       enabled(),
				Temperature:         enabled(),
				DewTemperature:      enabled(),
				Relay:               enabled(),
				Humidity:            enabled(),
				TargetTemperature:   enabled(),
			}
			device.Report = &defaultReport
		} else {
			if device.Report.HttpClient == nil {
				device.Report.HttpClient = enabled()
			}
			if device.Report.WaterTemperature == nil {
				device.Report.WaterTemperature = enabled()
			}
			if device.Report.ExternalTemperature == nil {
				device.Report.ExternalTemperature = enabled()
			}
			if device.Report.RoomConnected == nil {
				device.Report.RoomConnected = enabled()
			}
			if device.Report.Temperature == nil {
				device.Report.Temperature = enabled()
			}
			if device.Report.DewTemperature == nil {
				device.Report.DewTemperature = enabled()
			}
			if device.Report.Relay == nil {
				device.Report.Relay = enabled()
			}
			if device.Report.Humidity == nil {
				device.Report.Humidity = enabled()
			}
			if device.Report.TargetTemperature == nil {
				device.Report.TargetTemperature = enabled()
			}
		}
	}
	return nil
}

func enabled() *bool {
	b := true
	return &b
}
