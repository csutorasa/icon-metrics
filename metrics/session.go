package metrics

import (
	"github.com/csutorasa/icon-metrics/config"
	"github.com/csutorasa/icon-metrics/model"
)

// Metrics session data holder.
type MetricsSession interface {
	// Reports metrics based on device data.
	Report(values *model.DataPollResponse, reporter MetricsReporter)
	// Resets all metrics.
	Reset(reporter MetricsReporter)
}

// Room data holder.
type roomDescriptor struct {
	Id   string
	Name string
}

// Metrics session data holder.
type metricsSession struct {
	sysId               string
	roomDescriptors     []roomDescriptor
	reportConfiguration *config.ReportConfiguration
}

// Creates a new session to report metrics.
func NewSession(sysId string, reportConfiguration *config.ReportConfiguration) MetricsSession {
	return &metricsSession{
		sysId:               sysId,
		roomDescriptors:     make([]roomDescriptor, 0),
		reportConfiguration: reportConfiguration,
	}
}

// Reports metrics based on device data.
func (session *metricsSession) Report(values *model.DataPollResponse, reporter MetricsReporter) {
	if len(session.roomDescriptors) == 0 {
		for id, thermostat := range values.Thermostats {
			if thermostat.Enabled == 0 {
				continue
			}
			session.roomDescriptors = append(session.roomDescriptors, roomDescriptor{Id: id, Name: thermostat.Name})
		}
	}

	reporter.ExternalTemperature(session.sysId, values.ExternalTemperature)
	reporter.WaterTemperature(session.sysId, values.ExternalTemperature)
	for id, thermostat := range values.Thermostats {
		if thermostat.Enabled == 0 {
			continue
		}
		if thermostat.Live == 0 {
			reporter.RemoveRoom(session.sysId, id, thermostat.Name)
			reporter.RoomConnected(session.sysId, id, thermostat.Name, false)
			continue
		}
		if *session.reportConfiguration.RoomConnected {
			reporter.RoomConnected(session.sysId, id, thermostat.Name, true)
		}
		if *session.reportConfiguration.Temperature {
			reporter.RoomTemperature(session.sysId, id, thermostat.Name, thermostat.Temperature)
		}
		if *session.reportConfiguration.DewTemperature {
			reporter.RoomDewTemperature(session.sysId, id, thermostat.Name, thermostat.DewTemperature)
		}
		if *session.reportConfiguration.Relay {
			relay := false
			if thermostat.Relay > 0 {
				relay = true
			}
			reporter.RoomRelay(session.sysId, id, thermostat.Name, relay)
		}
		if *session.reportConfiguration.Humidity {
			reporter.RoomHumidity(session.sysId, id, thermostat.Name, thermostat.RelativeHumidity)
		}
		if *session.reportConfiguration.TargetTemperature {
			reporter.RoomTargetTemperature(session.sysId, id, thermostat.Name, thermostat.TargetTemperature())
		}
	}
}

// Resets all metrics.
func (session *metricsSession) Reset(reporter MetricsReporter) {
	for _, roomDescriptor := range session.roomDescriptors {
		reporter.RemoveRoom(session.sysId, roomDescriptor.Id, roomDescriptor.Name)
	}
	reporter.RemoveDevice(session.sysId)
	reporter.Connected(session.sysId, false)
	session.roomDescriptors = make([]roomDescriptor, 0)
}
