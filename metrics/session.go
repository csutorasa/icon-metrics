package metrics

import (
	"time"

	"github.com/csutorasa/icon-metrics/config"
	"github.com/csutorasa/icon-metrics/model"
)

// Metrics session data holder.
type MetricsSession interface {
	// Reports connected metric.
	Connected(connected bool)
	// Reports metrics based on device data.
	Report(values *model.DataPollResponse)
	// Reports HTTP metrics.
	HttpClientRequest(endpointName string, statusCode int, duration time.Duration)
	// Resets all metrics.
	Reset()
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
	reporter            MetricsReporter
}

// Creates a new session to report metrics.
func NewSession(sysId string, reportConfiguration *config.ReportConfiguration, reporter MetricsReporter) MetricsSession {
	return &metricsSession{
		sysId:               sysId,
		roomDescriptors:     make([]roomDescriptor, 0),
		reportConfiguration: reportConfiguration,
		reporter:            reporter,
	}
}

// Reports connected metric.
func (session *metricsSession) Connected(connected bool) {
	if *session.reportConfiguration.ControllerConnected {
		session.reporter.Connected(session.sysId, connected)
	}
}

// Reports metrics based on device data.
func (session *metricsSession) Report(values *model.DataPollResponse) {
	if len(session.roomDescriptors) == 0 {
		for id, thermostat := range values.Thermostats {
			if thermostat.Enabled == 0 {
				continue
			}
			session.roomDescriptors = append(session.roomDescriptors, roomDescriptor{Id: id, Name: thermostat.Name})
		}
	}

	if *session.reportConfiguration.ExternalTemperature {
		session.reporter.ExternalTemperature(session.sysId, values.ExternalTemperature)
	}
	if *session.reportConfiguration.WaterTemperature {
		session.reporter.WaterTemperature(session.sysId, values.WaterTemperature)
	}
	if *session.reportConfiguration.Heating {
		session.reporter.Heating(session.sysId, values.HeatingCooling == model.Heating)
	}
	if *session.reportConfiguration.Eco {
		session.reporter.Eco(session.sysId, values.ComfortEco == model.Eco)
	}
	for id, thermostat := range values.Thermostats {
		if thermostat.Enabled == 0 {
			continue
		}
		if thermostat.Live == 0 {
			session.reporter.RemoveRoom(session.sysId, id, thermostat.Name)
			session.reporter.RoomConnected(session.sysId, id, thermostat.Name, false)
			continue
		}
		if *session.reportConfiguration.RoomConnected {
			session.reporter.RoomConnected(session.sysId, id, thermostat.Name, true)
		}
		if *session.reportConfiguration.Temperature {
			session.reporter.RoomTemperature(session.sysId, id, thermostat.Name, thermostat.Temperature)
		}
		if *session.reportConfiguration.DewTemperature {
			session.reporter.RoomDewTemperature(session.sysId, id, thermostat.Name, thermostat.DewTemperature)
		}
		if *session.reportConfiguration.Relay {
			relay := false
			if thermostat.Relay > 0 {
				relay = true
			}
			session.reporter.RoomRelay(session.sysId, id, thermostat.Name, relay)
		}
		if *session.reportConfiguration.Humidity {
			session.reporter.RoomHumidity(session.sysId, id, thermostat.Name, thermostat.RelativeHumidity)
		}
		if *session.reportConfiguration.TargetTemperature {
			session.reporter.RoomTargetTemperature(session.sysId, id, thermostat.Name, thermostat.TargetTemperature())
		}
	}
}

// Reports HTTP metrics.
func (session *metricsSession) HttpClientRequest(endpointName string, statusCode int, duration time.Duration) {
	if *session.reportConfiguration.HttpClient {
		session.reporter.HttpClientRequest(session.sysId, endpointName, statusCode, duration)
	}
}

// Resets all metrics.
func (session *metricsSession) Reset() {
	for _, roomDescriptor := range session.roomDescriptors {
		session.reporter.RemoveRoom(session.sysId, roomDescriptor.Id, roomDescriptor.Name)
	}
	session.reporter.RemoveDevice(session.sysId)
	session.reporter.Connected(session.sysId, false)
	session.roomDescriptors = make([]roomDescriptor, 0)
}
