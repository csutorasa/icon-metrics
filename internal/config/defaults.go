package config

// Scans the config for invalid settings.
func setDefaults(config *Configuration) {
	setDefault(&config.Port, 80)
	for _, device := range config.Devices {
		setDefault(&device.Password, device.SysId)
		setDefault(&device.Delay, 15)
		device.Report = getDefault(device.Report, ReportConfiguration{})
		device.Report.ControllerConnected = getDefault(device.Report.ControllerConnected, true)
		device.Report.HttpClient = getDefault(device.Report.HttpClient, true)
		device.Report.WaterTemperature = getDefault(device.Report.WaterTemperature, true)
		device.Report.ExternalTemperature = getDefault(device.Report.ExternalTemperature, true)
		device.Report.Heating = getDefault(device.Report.Heating, true)
		device.Report.Eco = getDefault(device.Report.Eco, true)
		device.Report.RoomConnected = getDefault(device.Report.RoomConnected, true)
		device.Report.Temperature = getDefault(device.Report.Temperature, true)
		device.Report.DewTemperature = getDefault(device.Report.DewTemperature, true)
		device.Report.Relay = getDefault(device.Report.Relay, true)
		device.Report.Humidity = getDefault(device.Report.Humidity, true)
		device.Report.TargetTemperature = getDefault(device.Report.TargetTemperature, true)
	}
}

func setDefault[T comparable](value *T, defValue T) {
	var def T
	if *value == def {
		*value = defValue
	}
}

func getDefault[T any](value *T, defValue T) *T {
	if value == nil {
		t := new(T)
		*t = defValue
		return t
	}
	return value
}
