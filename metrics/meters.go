// Metrics data holding and manipulation.
package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Device related required parameters
var genericParameters = []string{"sysId"}

// HTTP response related required parameters
var httpParameters = append(genericParameters, "name", "response")

// Room related required parameters
var roomParameters = append(genericParameters, "id", "room")

// Reports if the device connection is active or not.
var ConntectedGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "icon_controller_connected",
	Help: "Read values",
}, genericParameters)

// Reports the HTTP response status code along with the duration.
var HttpGauge = promauto.NewSummaryVec(prometheus.SummaryOpts{
	Name: "icon_http_client_seconds",
	Help: "iCon HTTP client requests",
}, httpParameters)

// Application start time
var startTime time.Time

// Registers application uptime metric
func RegisterUptime() {
	startTime = time.Now()
	promauto.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "uptime",
		Help: "Uptime for the service",
	}, func() float64 { return float64(time.Duration(time.Since(startTime).Milliseconds())) })
}

// Reports the water temperature.
var WaterTemperatureGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "icon_water_temperature",
	Help: "Water temperature",
}, genericParameters)

// Reports the external temperature.
var ExternalTemperatureGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "icon_external_temperature",
	Help: "External temperature",
}, genericParameters)

// Reports the room temperature.
var RoomTemperatureGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "icon_temperature",
	Help: "Room temperature",
}, roomParameters)
var RoomDewTemperatureGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "icon_dew_temperature",
	Help: "Room temperature",
}, roomParameters)

// Reports the relay state.
var RelayGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "icon_relay_on",
	Help: "Relay on flag",
}, roomParameters)

// Reports the relative humidity.
var HumidityGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "icon_humidity",
	Help: "Relative humidity",
}, roomParameters)

// Reports the target temperature.
var TargetTemperatureGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "icon_target_temperature",
	Help: "Target temperature",
}, roomParameters)
