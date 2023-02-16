package model

import (
	"fmt"
)

// Experimental!
type ThermostatSettings map[int]*ThermostatSetting

// Experimental!
func (settings ThermostatSettings) ToValues(tab int) map[string][]string {
	data := make(map[string][]string)
	for signal, thermosSetting := range settings {
		for key, value := range thermosSetting.ToValues(tab, signal) {
			data[key] = value
		}
	}
	data["tab"] = []string{fmt.Sprintf("%d", tab)}
	data["form"] = []string{"thermos_data"}
	return data
}

// Experimental!
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

// Experimental!
func (setting *ThermostatSetting) ToValues(tab int, signal int) map[string][]string {
	id := fmt.Sprintf("%d_%d", tab, signal)
	data := map[string][]string{
		("cooling@" + id): {fmt.Sprintf("%.1f", setting.CoolingTargetTemperature)},
		("heating@" + id): {fmt.Sprintf("%.1f", setting.HeatingTargetTemperature)},
		("ecoc@" + id):    {fmt.Sprintf("%.1f", setting.EcoCoolingTargetTemperature)},
		("ecoh@" + id):    {fmt.Sprintf("%.1f", setting.EcoHeatingTargetTemperature)},
		("name@" + id):    {setting.Name},
		("lim@" + id):     {fmt.Sprintf("%.1f", setting.ManualRange)},
		("dxh@" + id):     {fmt.Sprintf("%.1f", setting.RegBHeating)},
		("dxc@" + id):     {fmt.Sprintf("%.1f", setting.RegBCooling)},
	}
	if setting.HeatingCooling {
		data["hc@"+id] = []string{"on"}
	}
	if setting.Installed {
		data["installed@"+id] = []string{"on"}
	}
	if setting.Cef {
		data["cef@"+id] = []string{"on"}
	}
	if setting.Cec {
		data["cec@"+id] = []string{"on"}
	}
	return data
}

// Experimental!
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

// Experimental!
func (settings *GeneralSettings) ToValues(tab int) map[string][]string {
	id := fmt.Sprintf("%d", tab)
	return map[string][]string{
		"func@ce_0":   {settings.ComfortEcoMode},
		"icon@ce_0":   {fmt.Sprintf("%d", settings.ComfortEcoTab)},
		"signal@ce_0": {fmt.Sprintf("%d", settings.ComfortEcoSignal)},
		"func@hc_0":   {settings.HeatingCoolingMode},
		"icon@hc_0":   {fmt.Sprintf("%d", settings.HeatingCoolingTab)},
		"signal@hc_0": {fmt.Sprintf("%d", settings.HeatingCoolingSignal)},
		"xah":         {fmt.Sprintf("%d", settings.HeatingTargetTemperature)},
		"xac":         {fmt.Sprintf("%d", settings.CoolingTargetTemperature)},
		"ecoh":        {fmt.Sprintf("%d", settings.EcoHeatingTargetTemperature)},
		"ecoc":        {fmt.Sprintf("%d", settings.EcoCoolingTargetTemperature)},
		"tab":         {id},
		"form":        {"general"},
	}
}
