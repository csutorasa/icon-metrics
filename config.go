package main

import (
	"errors"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Configuration root
type Configuration struct {
	Port    int                  `yaml:"port"`
	Devices []*IconConfiguration `yaml:"devices"`
}

// iCON device configuration
type IconConfiguration struct {
	Url      string `yaml:"url"`
	SysId    string `yaml:"sysid"`
	Password string `yaml:"password"`
	Delay    int    `yaml:"delay"`
}

// Returns the config that is read from the file.
func ReadConfig(filepath string) (*Configuration, error) {
	config := &Configuration{}
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}
	err = validateConfig(config)
	if err != nil {
		return config, err
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
			device.Delay = 60
		}
	}
	return nil
}
