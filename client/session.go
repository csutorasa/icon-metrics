package client

import (
	"github.com/csutorasa/icon-metrics/metrics"
)

type MetricsSession interface {
	Report(values *DataPollResponse)
	Reset()
}

type roomDescriptor struct {
	Id   string
	Name string
}

type simpleSession struct {
	sysId           string
	roomDescriptors []roomDescriptor
}

func NewSession(sysId string) MetricsSession {
	return &simpleSession{
		sysId:           sysId,
		roomDescriptors: make([]roomDescriptor, 0),
	}
}

func (this *simpleSession) Report(values *DataPollResponse) {
	if len(this.roomDescriptors) == 0 {
		for id, thermostat := range values.Thermostats {
			if thermostat.Enabled == 0 {
				continue
			}
			this.roomDescriptors = append(this.roomDescriptors, roomDescriptor{Id: id, Name: thermostat.Name})
		}
	}

	metrics.ExternalTemperatureGauge.WithLabelValues(this.sysId).Set(values.ExternalTemperature)
	metrics.WaterTemperatureGauge.WithLabelValues(this.sysId).Set(values.WaterTemperature)
	for id, thermostat := range values.Thermostats {
		if thermostat.Enabled == 0 {
			continue
		}
		metrics.RoomTemperatureGauge.WithLabelValues(this.sysId, id, thermostat.Name).Set(thermostat.Temperature)
		metrics.RoomDewTemperatureGauge.WithLabelValues(this.sysId, id, thermostat.Name).Set(thermostat.DewTemperature)
		metrics.RelayGauge.WithLabelValues(this.sysId, id, thermostat.Name).Set(float64(thermostat.Relay))
		metrics.HumidityGauge.WithLabelValues(this.sysId, id, thermostat.Name).Set(thermostat.RelativeHumidity)
		metrics.TargetTemperatureGauge.WithLabelValues(this.sysId, id, thermostat.Name).Set(thermostat.TargetTemperature())
	}
}

func (this *simpleSession) Reset() {
	metrics.WaterTemperatureGauge.DeleteLabelValues(this.sysId)
	for _, roomDescriptor := range this.roomDescriptors {
		metrics.RoomTemperatureGauge.DeleteLabelValues(this.sysId, roomDescriptor.Id, roomDescriptor.Name)
		metrics.RelayGauge.DeleteLabelValues(this.sysId, roomDescriptor.Id, roomDescriptor.Name)
		metrics.HumidityGauge.DeleteLabelValues(this.sysId, roomDescriptor.Id, roomDescriptor.Name)
		metrics.TargetTemperatureGauge.DeleteLabelValues(this.sysId, roomDescriptor.Id, roomDescriptor.Name)
	}
	this.roomDescriptors = make([]roomDescriptor, 0)
}
