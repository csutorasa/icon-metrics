package client

import (
	"github.com/csutorasa/icon-metrics/config"
	"github.com/csutorasa/icon-metrics/metrics"
)

// Room data holder.
type roomDescriptor struct {
	Id   string
	Name string
}

// Metrics session data holder.
type MetricsSession struct {
	sysId               string
	roomDescriptors     []roomDescriptor
	reportConfiguration *config.ReportConfiguration
}

// Creates a new session to report metrics.
func NewSession(sysId string, reportConfiguration *config.ReportConfiguration) *MetricsSession {
	return &MetricsSession{
		sysId:               sysId,
		roomDescriptors:     make([]roomDescriptor, 0),
		reportConfiguration: reportConfiguration,
	}
}

// Reports metrics based on device data.
func (session *MetricsSession) Report(values *DataPollResponse) {
	if len(session.roomDescriptors) == 0 {
		for id, thermostat := range values.Thermostats {
			if thermostat.Enabled == 0 {
				continue
			}
			session.roomDescriptors = append(session.roomDescriptors, roomDescriptor{Id: id, Name: thermostat.Name})
		}
	}

	metrics.ExternalTemperatureGauge.WithLabelValues(session.sysId).Set(values.ExternalTemperature)
	metrics.WaterTemperatureGauge.WithLabelValues(session.sysId).Set(values.WaterTemperature)
	for id, thermostat := range values.Thermostats {
		if thermostat.Enabled == 0 {
			continue
		}
		if thermostat.Live == 0 {
			if *session.reportConfiguration.RoomConnected {
				metrics.RoomConntectedGauge.WithLabelValues(session.sysId, id, thermostat.Name).Set(0)
			}
			session.removeRoom(id, thermostat.Name)
			continue
		}
		if *session.reportConfiguration.RoomConnected {
			metrics.RoomConntectedGauge.WithLabelValues(session.sysId, id, thermostat.Name).Set(1)
		}
		if *session.reportConfiguration.Temperature {
			metrics.RoomTemperatureGauge.WithLabelValues(session.sysId, id, thermostat.Name).Set(thermostat.Temperature)
		}
		if *session.reportConfiguration.DewTemperature {
			metrics.RoomDewTemperatureGauge.WithLabelValues(session.sysId, id, thermostat.Name).Set(thermostat.DewTemperature)
		}
		if *session.reportConfiguration.Relay {
			metrics.RelayGauge.WithLabelValues(session.sysId, id, thermostat.Name).Set(float64(thermostat.Relay))
		}
		if *session.reportConfiguration.Humidity {
			metrics.HumidityGauge.WithLabelValues(session.sysId, id, thermostat.Name).Set(thermostat.RelativeHumidity)
		}
		if *session.reportConfiguration.TargetTemperature {
			metrics.TargetTemperatureGauge.WithLabelValues(session.sysId, id, thermostat.Name).Set(thermostat.TargetTemperature())
		}
	}
}

// Resets all metrics.
func (session *MetricsSession) Reset() {
	metrics.WaterTemperatureGauge.DeleteLabelValues(session.sysId)
	metrics.ExternalTemperatureGauge.DeleteLabelValues(session.sysId)
	for _, roomDescriptor := range session.roomDescriptors {
		metrics.RoomConntectedGauge.DeleteLabelValues(session.sysId, roomDescriptor.Id, roomDescriptor.Name)
		session.removeRoom(roomDescriptor.Id, roomDescriptor.Name)
	}
	session.roomDescriptors = make([]roomDescriptor, 0)
}

// Removes metrics for a room
func (session *MetricsSession) removeRoom(id, name string) {
	metrics.RoomTemperatureGauge.DeleteLabelValues(session.sysId, id, name)
	metrics.RoomDewTemperatureGauge.DeleteLabelValues(session.sysId, id, name)
	metrics.RelayGauge.DeleteLabelValues(session.sysId, id, name)
	metrics.HumidityGauge.DeleteLabelValues(session.sysId, id, name)
	metrics.TargetTemperatureGauge.DeleteLabelValues(session.sysId, id, name)
}
