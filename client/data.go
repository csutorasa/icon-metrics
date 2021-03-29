package client

type DataPollResponse struct {
	SysId                       string         `json:"SYSID"`
	SERVICE                     int            `json:"SERVICE"`
	Version                     string         `json:"VER"`
	Heating                     int            `json:"HC"`
	CE                          int            `json:"CE"`
	ON                          int            `json:"ON"`
	ETEMP                       float64        `json:"ETEMP"`
	WaterTemperature            float64        `json:"WTEMP"`
	Pump                        int            `json:"PUMP"`
	ERR                         int            `json:"ERR"`
	OverheatWarning             int            `json:"OVERHEAT"`
	FrostWarning                int            `json:"WFROST"`
	HeatingTargetTemperature    float64        `json:"XAH"`
	CoolingTargetTemperature    float64        `json:"XAC"`
	EcoHeatingTargetTemperature float64        `json:"ECOH"`
	EcoCoolingTargetTemperature float64        `json:"ECOC"`
	SIG                         int            `json:"SIG"`
	SW                          int            `json:"SW"`
	Email                       string         `json:"EMAIL"`
	Timezone                    string         `json:"TZ"`
	Thermostats                 map[string]*DP `json:"DP"`
	TPR                         TPR            `json:"TPR"`
}

type DP struct {
	Enabled                     int     `json:"ON"`
	IHC                         int     `json:"IHC"`
	Live                        int     `json:"LIVE"`
	Temperature                 float64 `json:"TEMP"`
	RelativeHumidity            float64 `json:"RH"`
	DewTemperature              float64 `json:"DEW"`
	ManualRange                 float64 `json:"LIM"`
	DWP                         int     `json:"DWP"`
	FrostWarning                int     `json:"FROST"`
	CE                          int     `json:"CE"`
	HC                          int     `json:"HC"`
	DI                          int     `json:"DI"`
	HeatingTargetTemperature    float64 `json:"XAH"`
	CoolingTargetTemperature    float64 `json:"XAC"`
	EcoHeatingTargetTemperature float64 `json:"ECOH"`
	EcoCoolingTargetTemperature float64 `json:"ECOC"`
	PL                          int     `json:"PL"`
	CEF                         int     `json:"CEF"`
	CEC                         int     `json:"CEC"`
	RegBHeating                 int     `json:"DXH"`
	RegBCooling                 int     `json:"DXC"`
	Relay                       int     `json:"OUT"`
	WP                          int     `json:"WP"`
	MV                          int     `json:"MV"`
	TPR                         int     `json:"TPR"`
	Name                        string  `json:"NAME"`
}

type TPR struct {
	Heat TPRData `json:"HEAT"`
	Cool TPRData `json:"COOL"`
}

type TPRData struct {
}

const success = "success"
const failure = "failure"

type actionResponse struct {
	Result  string `json:"result"`
	Refresh bool   `json:"refresh"`
}
