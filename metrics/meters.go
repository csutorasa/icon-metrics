package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var genericParameters = []string{"sysId"}
var httpParameters = append(genericParameters, "name", "response")
var roomParameters = append(genericParameters, "id", "room")

var ConntectedGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "icon_controller_connected",
	Help: "Read values",
}, genericParameters)
var HttpGauge = promauto.NewSummaryVec(prometheus.SummaryOpts{
	Name: "icon_http_client_seconds",
	Help: "iCon HTTP client requests",
}, httpParameters)

var startTime time.Time

func RegisterUptime() {
	startTime = time.Now()
	promauto.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "uptime",
		Help: "Uptime for the service",
	}, func() float64 { return float64(time.Duration(time.Since(startTime).Milliseconds())) })
}

var WaterTemperatureGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "icon_water_temperature",
	Help: "Water temperature",
}, genericParameters)
var ExternalTemperatureGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "icon_external_temperature",
	Help: "External temperature",
}, genericParameters)
var RoomTemperatureGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "icon_temperature",
	Help: "Room temperature",
}, roomParameters)
var RoomDewTemperatureGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "icon_dew_temperature",
	Help: "Room temperature",
}, roomParameters)
var RelayGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "icon_relay_on",
	Help: "Relay on flag",
}, roomParameters)
var HumidityGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "icon_humidity",
	Help: "Relative humidity",
}, roomParameters)
var TargetTemperatureGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "icon_target_temperature",
	Help: "Target temperature",
}, roomParameters)
