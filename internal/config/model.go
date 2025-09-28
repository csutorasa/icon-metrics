package config

// Configuration root
type Configuration struct {
	Port    int                  `yaml:"port"`
	Devices []*IconConfiguration `yaml:"devices"`
}

// iCON device configuration
type IconConfiguration struct {
	Url      string               `yaml:"url"`
	SysId    string               `yaml:"sysid"`
	Password string               `yaml:"password"`
	Delay    int                  `yaml:"delay"`
	Report   *ReportConfiguration `yaml:"report"`
}

// iCON device report configuration
type ReportConfiguration struct {
	// metrics.RoomConntectedGauge
	ControllerConnected *bool `yaml:"controllerConnected"`
	// metrics.HttpGauge
	HttpClient *bool `yaml:"httpClient"`
	// metrics.WaterTemperatureGauge
	WaterTemperature *bool `yaml:"waterTemperature"`
	// metrics.ExternalTemperatureGauge
	ExternalTemperature *bool `yaml:"externalTemperature"`
	// metrics.HeatingGauge
	Heating *bool `yaml:"heating"`
	// metrics.EcoGauge
	Eco *bool `yaml:"eco"`
	// metrics.RoomConntectedGauge
	RoomConnected *bool `yaml:"roomConnected"`
	// metrics.RoomTemperatureGauge
	Temperature *bool `yaml:"temperature"`
	// metrics.RoomDewTemperatureGauge
	DewTemperature *bool `yaml:"dewTemperature"`
	// metrics.RelayGauge
	Relay *bool `yaml:"relay"`
	// metrics.HumidityGauge
	Humidity *bool `yaml:"humidity"`
	// metrics.TargetTemperatureGauge
	TargetTemperature *bool `yaml:"targetTemperature"`
}
