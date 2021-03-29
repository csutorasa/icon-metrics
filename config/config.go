package config

import (
	"errors"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Port    int                  `yaml:"port"`
	Devices []*IconConfiguration `yaml:"devices"`
}

type IconConfiguration struct {
	Url      string `yaml:"url"`
	SysId    string `yaml:"sysid"`
	Password string `yaml:"password"`
	Delay    int    `yaml:"delay"`
}

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

func validateConfig(config *Configuration) error {
	if config.Port == 0 {
		config.Port = 80
	}
	if config.Devices == nil || len(config.Devices) == 0 {
		return errors.New("There are no devices to monitor")
	}
	for i, device := range config.Devices {
		if device.SysId == "" {
			return errors.New(fmt.Sprintf("Device config at %d position is missing sysid", i))
		}
		if device.Url == "" {
			return errors.New(fmt.Sprintf("Device config at %d position is missing url", i))
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
