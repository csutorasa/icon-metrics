package client

import (
	"errors"
	"fmt"
	"strings"
)

type HC int
type CE int

const (
	Heating HC = 0
	Cooling HC = 1
)

const (
	Comfort CE = 0
	Eco     CE = 1
)

type DataPollResponse struct {
	SysId                       string         `json:"SYSID"`
	SERVICE                     int            `json:"SERVICE"`
	Version                     string         `json:"VER"`
	HeatingCooling              HC             `json:"HC"`
	ComfortEco                  CE             `json:"CE"`
	ON                          int            `json:"ON"`
	ExternalTemperature         float64        `json:"ETEMP"`
	WaterTemperature            float64        `json:"WTEMP"`
	Pump                        int            `json:"PUMP"`
	Error                       int            `json:"ERR"`
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

func (this *DataPollResponse) TargetTemperature() float64 {
	if this.HeatingCooling == Heating {
		if this.ComfortEco == Comfort {
			return this.HeatingTargetTemperature
		} else {
			return this.EcoHeatingTargetTemperature
		}
	} else {
		if this.ComfortEco == Comfort {
			return this.CoolingTargetTemperature
		} else {
			return this.EcoCoolingTargetTemperature
		}
	}
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
	ComfortEco                  CE      `json:"CE"`
	HeatingCooling              HC      `json:"HC"`
	OpenWindowInput             int     `json:"DI"`
	HeatingTargetTemperature    float64 `json:"XAH"`
	CoolingTargetTemperature    float64 `json:"XAC"`
	EcoHeatingTargetTemperature float64 `json:"ECOH"`
	EcoCoolingTargetTemperature float64 `json:"ECOC"`
	ParentalLock                int     `json:"PL"`
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

func (this *DP) TargetTemperature() float64 {
	if this.HeatingCooling == Heating {
		if this.ComfortEco == Comfort {
			return this.HeatingTargetTemperature
		} else {
			return this.EcoHeatingTargetTemperature
		}
	} else {
		if this.ComfortEco == Comfort {
			return this.CoolingTargetTemperature
		} else {
			return this.EcoCoolingTargetTemperature
		}
	}
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
	Result  string         `json:"result"`
	Refresh bool           `json:"refresh"`
	Errors  map[string]any `json:"errors"`
}

func (this *actionResponse) IsSuccess() bool {
	return this.Result == "success"
}

func (this *actionResponse) CreateError() error {
	var sb strings.Builder
	for error, data := range this.Errors {
		sb.WriteString(error)
		sb.WriteString(" ")
		sb.WriteString(fmt.Sprintf("%v", data))
	}
	return errors.New(sb.String())
}
