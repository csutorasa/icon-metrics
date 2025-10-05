package metrics

import (
	"github.com/csutorasa/icon-metrics/internal/config"
	"github.com/csutorasa/icon-metrics/pkg/model"
)

// Room data holder.
type RoomDescriptor struct {
	Id   string
	Name string
}

// Reports metrics based on device data.
func Report(values *model.DataPollResponse, sysId string, m *Metrics, config *config.ReportConfiguration, reportedRooms *map[RoomDescriptor]bool) {
	for room := range *reportedRooms {
		found := false
		for id, thermostat := range values.Thermostats {
			if room.Id == id && room.Name == thermostat.Name && thermostat.Enabled > 0 {
				found = true
				break
			}
		}
		if !found {
			m.RemoveRoom(sysId, room.Id, room.Name)
		}
	}

	setSystemValue(config.ExternalTemperature, m.ExternalTemperature, sysId, values.ExternalTemperature)
	setSystemValue(config.Heating, m.Heating, sysId, boolToFloat64(values.HeatingCooling == model.Heating))
	setSystemValue(config.Eco, m.Eco, sysId, boolToFloat64(values.ComfortEco == model.Eco))

	newReportedRooms := map[RoomDescriptor]bool{}
	for id, thermostat := range values.Thermostats {
		if thermostat.Enabled == 0 {
			m.RemoveRoom(sysId, id, thermostat.Name)
			continue
		}
		newReportedRooms[RoomDescriptor{Id: id, Name: thermostat.Name}] = true
		if thermostat.Live == 0 {
			setRoomValue(config.RoomConnected, m.RoomConnected, sysId, id, thermostat.Name, 0)
			m.RoomDewTemperature.Remove(sysId, id, thermostat.Name)
			m.RoomHumidity.Remove(sysId, id, thermostat.Name)
			m.RoomRelay.Remove(sysId, id, thermostat.Name)
			m.RoomTargetTemperature.Remove(sysId, id, thermostat.Name)
			m.RoomTemperature.Remove(sysId, id, thermostat.Name)
			continue
		}
		setRoomValue(config.RoomConnected, m.RoomConnected, sysId, id, thermostat.Name, 1)
		setRoomValue(config.Temperature, m.RoomTemperature, sysId, id, thermostat.Name, thermostat.Temperature)
		setRoomValue(config.DewTemperature, m.RoomDewTemperature, sysId, id, thermostat.Name, thermostat.DewTemperature)
		setRoomValue(config.Relay, m.RoomRelay, sysId, id, thermostat.Name, intToFloat64(thermostat.Relay))
		setRoomValue(config.Humidity, m.RoomHumidity, sysId, id, thermostat.Name, thermostat.RelativeHumidity)
		setRoomValue(config.TargetTemperature, m.RoomTargetTemperature, sysId, id, thermostat.Name, thermostat.TargetTemperature())
	}
	*reportedRooms = newReportedRooms
}

func setSystemValue(enabled *bool, m SystemMetric, sysId string, f float64) {
	if *enabled {
		m.Set(sysId, f)
	} else {
		m.Remove(sysId)
	}
}

func setRoomValue(enabled *bool, m RoomMetric, sysId string, id string, name string, f float64) {
	if *enabled {
		m.Set(sysId, id, name, f)
	} else {
		m.Remove(sysId, id, name)
	}
}
