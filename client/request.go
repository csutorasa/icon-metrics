package client

import (
	"fmt"
)

type ThermostatSettings map[int]*ThermostatSetting

func (this ThermostatSettings) ToValues(tab int) map[string][]string {
	data := make(map[string][]string)
	for signal, thermosSetting := range this {
		for key, value := range thermosSetting.ToValues(tab, signal) {
			data[key] = value
		}
	}
	data["tab"] = []string{fmt.Sprintf("%d", tab)}
	data["form"] = []string{"thermos_data"}
	return data
}

type ThermostatSetting struct {
	HeatingCooling              bool
	Installed                   bool
	EcoCoolingTargetTemperature float64
	EcoHeatingTargetTemperature float64
	CoolingTargetTemperature    float64
	HeatingTargetTemperature    float64
	Cef                         bool
	Cec                         bool
	Name                        string
	ManualRange                 float64
	RegBHeating                 float64
	RegBCooling                 float64
}

func (this *ThermostatSetting) ToValues(tab int, signal int) map[string][]string {
	id := fmt.Sprintf("%d_%d", tab, signal)
	data := map[string][]string{
		("cooling@" + id): {fmt.Sprintf("%.1f", this.CoolingTargetTemperature)},
		("heating@" + id): {fmt.Sprintf("%.1f", this.HeatingTargetTemperature)},
		("ecoc@" + id):    {fmt.Sprintf("%.1f", this.EcoCoolingTargetTemperature)},
		("ecoh@" + id):    {fmt.Sprintf("%.1f", this.EcoHeatingTargetTemperature)},
		("name@" + id):    {this.Name},
		("lim@" + id):     {fmt.Sprintf("%.1f", this.ManualRange)},
		("dxh@" + id):     {fmt.Sprintf("%.1f", this.RegBHeating)},
		("dxc@" + id):     {fmt.Sprintf("%.1f", this.RegBCooling)},
	}
	if this.HeatingCooling {
		data["hc@"+id] = []string{"on"}
	}
	if this.Installed {
		data["installed@"+id] = []string{"on"}
	}
	if this.Cef {
		data["cef@"+id] = []string{"on"}
	}
	if this.Cec {
		data["cec@"+id] = []string{"on"}
	}
	return data
}

type GeneralSettings struct {
	ComfortEcoMode              string
	ComfortEcoTab               int
	ComfortEcoSignal            int
	HeatingCoolingMode          string
	HeatingCoolingTab           int
	HeatingCoolingSignal        int
	HeatingTargetTemperature    int
	CoolingTargetTemperature    int
	EcoHeatingTargetTemperature int
	EcoCoolingTargetTemperature int
}

func (this *GeneralSettings) ToValues(tab int) map[string][]string {
	id := fmt.Sprintf("%d", tab)
	return map[string][]string{
		"func@ce_0":   {this.ComfortEcoMode},
		"icon@ce_0":   {fmt.Sprintf("%d", this.ComfortEcoTab)},
		"signal@ce_0": {fmt.Sprintf("%d", this.ComfortEcoSignal)},
		"func@hc_0":   {this.HeatingCoolingMode},
		"icon@hc_0":   {fmt.Sprintf("%d", this.HeatingCoolingTab)},
		"signal@hc_0": {fmt.Sprintf("%d", this.HeatingCoolingSignal)},
		"xah":         {fmt.Sprintf("%d", this.HeatingTargetTemperature)},
		"xac":         {fmt.Sprintf("%d", this.CoolingTargetTemperature)},
		"ecoh":        {fmt.Sprintf("%d", this.EcoHeatingTargetTemperature)},
		"ecoc":        {fmt.Sprintf("%d", this.EcoCoolingTargetTemperature)},
		"tab":         {id},
		"form":        {"general"},
	}
}
