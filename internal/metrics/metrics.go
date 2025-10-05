// Metrics data holding and manipulation.
package metrics

import (
	"time"
)

type SystemMetric interface {
	Set(sysId string, f float64)
	Remove(sysId string)
}

type RoomMetric interface {
	Set(sysId string, roomId string, roomName string, f float64)
	RemoveSystem(sysId string)
	Remove(sysId string, roomId string, roomName string)
}

type HttpMetric interface {
	Observe(sysId string, name string, statusCode int, duration time.Duration)
	RemoveSystem(sysId string)
}

type Metrics struct {
	Connected             SystemMetric
	WaterTemperature      SystemMetric
	ExternalTemperature   SystemMetric
	Heating               SystemMetric
	Eco                   SystemMetric
	RoomConnected         RoomMetric
	RoomTemperature       RoomMetric
	RoomDewTemperature    RoomMetric
	RoomTargetTemperature RoomMetric
	RoomHumidity          RoomMetric
	RoomRelay             RoomMetric
	Http                  HttpMetric
}

func (m *Metrics) RemoveSystem(sysId string) {
	m.RoomConnected.RemoveSystem(sysId)
	m.RoomDewTemperature.RemoveSystem(sysId)
	m.RoomHumidity.RemoveSystem(sysId)
	m.RoomRelay.RemoveSystem(sysId)
	m.RoomTargetTemperature.RemoveSystem(sysId)
	m.RoomTemperature.RemoveSystem(sysId)
	m.Connected.Remove(sysId)
	m.ExternalTemperature.Remove(sysId)
	m.Heating.Remove(sysId)
	m.Eco.Remove(sysId)
	m.Http.RemoveSystem(sysId)
}

func (m *Metrics) RemoveRoom(sysId string, roomId string, roomName string) {
	m.RoomConnected.Remove(sysId, roomId, roomName)
	m.RoomDewTemperature.Remove(sysId, roomId, roomName)
	m.RoomHumidity.Remove(sysId, roomId, roomName)
	m.RoomRelay.Remove(sysId, roomId, roomName)
	m.RoomTargetTemperature.Remove(sysId, roomId, roomName)
	m.RoomTemperature.Remove(sysId, roomId, roomName)
}

func boolToFloat64(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

func intToFloat64(i int) float64 {
	if i > 0 {
		return 1
	}
	return 0
}
