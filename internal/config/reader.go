package config

import (
	"errors"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

// Returns the config that is read from the file.
func ReadConfig(r io.Reader) (*Configuration, error) {
	decoder := yaml.NewDecoder(r)
	config := &Configuration{}
	err := decoder.Decode(config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse yaml: %w", err)
	}
	err = validateConfig(config)
	if err != nil {
		return config, fmt.Errorf("invalid config: %w", err)
	}
	setDefaults(config)
	return config, nil
}

// Returns the config that is read from the file.
func ReadConfigFile(filepath string) (*Configuration, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	defer f.Close()
	return ReadConfig(f)
}

// Scans the config for invalid settings.
func validateConfig(config *Configuration) error {
	if len(config.Devices) == 0 {
		return errors.New("there are no devices to monitor")
	}
	for i, device := range config.Devices {
		if device.SysId == "" {
			return fmt.Errorf("device config at %d position is missing sysid", i)
		}
		if device.Url == "" {
			return fmt.Errorf("device config at %d position is missing url", i)
		}
	}
	return nil
}
