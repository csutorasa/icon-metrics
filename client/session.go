package client

import (
	"github.com/csutorasa/icon-metrics/metrics"
)

// Room data holder.
type roomDescriptor struct {
	Id   string
	Name string
}

// Metrics session data holder.
type MetricsSession struct {
	sysId           string
	roomDescriptors []roomDescriptor
}

// Creates a new session to report metrics.
func NewSession(sysId string) *MetricsSession {
	return &MetricsSession{
		sysId:           sysId,
		roomDescriptors: make([]roomDescriptor, 0),
	}
}

// Reports metrics based on device data.
func (session *MetricsSession) Report(values *DataPollResponse) {
	if len(session.roomDescriptors) == 0 {
		for id, thermostat := range values.Thermostats {
			if thermostat.Enabled == 0 || thermostat.Live == 0 {
				continue
			}
			session.roomDescriptors = append(session.roomDescriptors, roomDescriptor{Id: id, Name: thermostat.Name})
		}
	}

	metrics.ExternalTemperatureGauge.WithLabelValues(session.sysId).Set(values.ExternalTemperature)
	metrics.WaterTemperatureGauge.WithLabelValues(session.sysId).Set(values.WaterTemperature)
	for id, thermostat := range values.Thermostats {
		if thermostat.Enabled == 0 || thermostat.Live == 0 {
			continue
		}
		metrics.RoomTemperatureGauge.WithLabelValues(session.sysId, id, thermostat.Name).Set(thermostat.Temperature)
		metrics.RoomDewTemperatureGauge.WithLabelValues(session.sysId, id, thermostat.Name).Set(thermostat.DewTemperature)
		metrics.RelayGauge.WithLabelValues(session.sysId, id, thermostat.Name).Set(float64(thermostat.Relay))
		metrics.HumidityGauge.WithLabelValues(session.sysId, id, thermostat.Name).Set(thermostat.RelativeHumidity)
		metrics.TargetTemperatureGauge.WithLabelValues(session.sysId, id, thermostat.Name).Set(thermostat.TargetTemperature())
	}
}

// Resets all metrics.
func (session *MetricsSession) Reset() {
	metrics.WaterTemperatureGauge.DeleteLabelValues(session.sysId)
	metrics.ExternalTemperatureGauge.DeleteLabelValues(session.sysId)
	for _, roomDescriptor := range session.roomDescriptors {
		metrics.RoomTemperatureGauge.DeleteLabelValues(session.sysId, roomDescriptor.Id, roomDescriptor.Name)
		metrics.RoomDewTemperatureGauge.DeleteLabelValues(session.sysId, roomDescriptor.Id, roomDescriptor.Name)
		metrics.RelayGauge.DeleteLabelValues(session.sysId, roomDescriptor.Id, roomDescriptor.Name)
		metrics.HumidityGauge.DeleteLabelValues(session.sysId, roomDescriptor.Id, roomDescriptor.Name)
		metrics.TargetTemperatureGauge.DeleteLabelValues(session.sysId, roomDescriptor.Id, roomDescriptor.Name)
	}
	session.roomDescriptors = make([]roomDescriptor, 0)
}
