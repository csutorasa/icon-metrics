package config_test

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/csutorasa/icon-metrics/internal/config"
)

func TestReadConfig(t *testing.T) {
	buf := bytes.NewBuffer([]byte(`devices:
  - url: http://192.168.1.10 # device address
    sysid: '123123123123' # device ID (printed on the controller)
`))
	config, err := config.ReadConfig(buf)
	if err != nil {
		t.Fatalf("error is not expected but got %s", err.Error())
	}
	if config.Port != 80 {
		t.Logf("Expected default port 80 but got %d", config.Port)
		t.Fail()
	}
	if len(config.Devices) != 1 {
		t.Logf("Expected 1 devices but got %d", len(config.Devices))
		t.Fail()
	} else {
		device := config.Devices[0]
		if device.Url != "http://192.168.1.10" {
			t.Logf("Expected url http://192.168.1.10 but got %s", device.Url)
			t.Fail()
		}
		if device.SysId != "123123123123" {
			t.Logf("Expected sysId 123123123123 but got %s", device.SysId)
			t.Fail()
		}
		if device.Password != "123123123123" {
			t.Logf("Expected password 123123123123 but got %s", device.Password)
			t.Fail()
		}
		if device.Delay != 15 {
			t.Logf("Expected delay 15 but got %d", device.Delay)
			t.Fail()
		}
		if device.Report == nil {
			t.Logf("Expected report but got nil")
			t.Fail()
		} else {
			if *device.Report.ControllerConnected != true {
				t.Logf("Expected ControllerConnected true but got %t", *device.Report.ControllerConnected)
				t.Fail()
			}
		}
	}
}

func TestReadConfigEmptyConfig(t *testing.T) {
	buf := bytes.NewBuffer([]byte(""))
	_, err := config.ReadConfig(buf)
	if err == nil {
		t.Fatalf("error is expected but got none")
	}
	if !errors.Is(err, io.EOF) {
		t.Fatalf("EOF is expected but got %s", err.Error())
	}
}

func TestReadConfigInvalidConfig(t *testing.T) {
	buf := bytes.NewBuffer([]byte("invalid"))
	_, err := config.ReadConfig(buf)
	if err == nil {
		t.Fatalf("error is expected but got none")
	}
	if !strings.Contains(err.Error(), "cannot unmarshal") {
		t.Fatalf("unmarshal error is expected but got %s", err.Error())
	}
}
