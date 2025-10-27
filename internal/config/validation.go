package config

import (
	"errors"
	"fmt"
	"net/url"
)

// Scans the config for invalid settings.
func (config *Configuration) validateConfig() error {
	if len(config.Devices) == 0 {
		return errors.New("there are no devices to monitor")
	}
	for i, device := range config.Devices {
		if device.SysId == "" {
			return fmt.Errorf("device config at %d position is missing sysid", i)
		}
		if device.Url == "" {
			return fmt.Errorf("device config %s is missing url", device.SysId)
		}
		u, err := url.Parse(device.Url)
		if err != nil {
			return fmt.Errorf("device config %s has an invalid url \"%s\"", device.SysId, device.Url)
		}
		if u.Scheme != "http" && u.Scheme != "https" {
			return fmt.Errorf("device config %s has an unsupported scheme \"%s\"", device.SysId, u.Scheme)
		}
	}
	return nil
}
