package metrics

import (
	"context"
	"errors"
	"time"

	"github.com/csutorasa/icon-metrics/internal/config"
	"github.com/csutorasa/icon-metrics/pkg/client"
	"github.com/csutorasa/icon-metrics/pkg/model"
)

type iconMetricsClient struct {
	client  client.IconClient
	metrics *Metrics
	config  *config.ReportConfiguration
}

func NewIconMetricsClient(c client.IconClient, m *Metrics, config *config.ReportConfiguration) client.IconClient {
	return &iconMetricsClient{
		client:  c,
		metrics: m,
		config:  config,
	}
}

func (c *iconMetricsClient) SysId() string {
	return c.client.SysId()
}

func (c *iconMetricsClient) Close() error {
	c.metrics.RemoveSystem(c.client.SysId())
	return c.client.Close()
}

func (c *iconMetricsClient) ReadValues() (*model.DataPollResponse, error) {
	t := NewTimer()
	values, err := c.client.ReadValues()
	c.httpClientRequest("read_values", err, t.End())
	return values, err
}

func (c *iconMetricsClient) ReadValuesContext(ctx context.Context) (*model.DataPollResponse, error) {
	t := NewTimer()
	values, err := c.client.ReadValuesContext(ctx)
	c.httpClientRequest("read_values", err, t.End())
	return values, err
}

func (c *iconMetricsClient) SetThermostatSettings(tab int, thermosSettings model.ThermostatSettings) error {
	t := NewTimer()
	err := c.client.SetThermostatSettings(tab, thermosSettings)
	c.httpClientRequest("set_thermostat_settings", err, t.End())
	return err
}

func (c *iconMetricsClient) SetThermostatSettingsContext(ctx context.Context, tab int, thermosSettings model.ThermostatSettings) error {
	t := NewTimer()
	err := c.client.SetThermostatSettingsContext(ctx, tab, thermosSettings)
	c.httpClientRequest("set_thermostat_settings", err, t.End())
	return err
}

func (c *iconMetricsClient) SetGeneralSettings(tab int, generalSettings *model.GeneralSettings) error {
	t := NewTimer()
	err := c.client.SetGeneralSettings(tab, generalSettings)
	c.httpClientRequest("set_general_settings", err, t.End())
	return err
}

func (c *iconMetricsClient) SetGeneralSettingsContext(ctx context.Context, tab int, generalSettings *model.GeneralSettings) error {
	t := NewTimer()
	err := c.client.SetGeneralSettingsContext(ctx, tab, generalSettings)
	c.httpClientRequest("set_general_settings", err, t.End())
	return err
}

func (c *iconMetricsClient) httpClientRequest(endpointName string, err error, d time.Duration) {
	if !*c.config.HttpClient {
		return
	}
	if err == nil {
		c.metrics.Http.Observe(c.SysId(), endpointName, 200, d)
		return
	}
	errHttpStatus := &client.ErrHttpStatus{}
	if errors.As(err, &errHttpStatus) {
		c.metrics.Http.Observe(c.SysId(), endpointName, errHttpStatus.StatusCode, d)
		return
	}
	errHttpBodyUnmarshal := &client.ErrHttpBodyUnmarshal{}
	if errors.As(err, &errHttpBodyUnmarshal) {
		c.metrics.Http.Observe(c.SysId(), endpointName, 200, d)
		return
	}
	c.metrics.Http.Observe(c.SysId(), endpointName, 0, d)
}
