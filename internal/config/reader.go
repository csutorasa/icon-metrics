package config

import (
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
	err = config.validateConfig()
	if err != nil {
		return config, fmt.Errorf("invalid config: %w", err)
	}
	config.setDefaults()
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
