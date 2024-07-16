// Metrics data holding and manipulation.
package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricsReporter interface {
	// Registers application uptime metric
	Uptime()
	SystemMetricsReporter
	RoomMetricsReporter
	HttpMetricsReporter
}

// Registers application uptime metric
func (r *metricsReporter) Uptime() {
	startTime := time.Now()
	promauto.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "uptime",
		Help: "Uptime for the service",
	}, func() float64 { return float64(time.Duration(time.Since(startTime).Milliseconds())) })
}

// Device related required parameters
var genericParameters = []string{"sysId"}

type SystemMetricsReporter interface {
	// Reports if the device connection is active or not.
	Connected(sysId string, connected bool)
	// Reports the water temperature.
	WaterTemperature(sysId string, temperature float64)
	// Reports the external temperature.
	ExternalTemperature(sysId string, temperature float64)
	// Reports if the controller is set to heating or cooling.
	Heating(sysId string, heating bool)
	// Reports if the controller is set to eco or normal mode.
	Eco(sysId string, eco bool)
	// Removes device from reporting.
	RemoveDevice(sysId string)
}

type systemMetricsReporter struct {
	connectedGauge           *prometheus.GaugeVec
	waterTemperatureGauge    *prometheus.GaugeVec
	externalTemperatureGauge *prometheus.GaugeVec
	heatingGauge             *prometheus.GaugeVec
	ecoGauge                 *prometheus.GaugeVec
}

func newSystemPrometheusReporter() SystemMetricsReporter {
	return &systemMetricsReporter{
		connectedGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "icon_controller_connected",
			Help: "For each controller, reports 1 if the controller is ready to be read, 0 otherwise.",
		}, genericParameters),
		waterTemperatureGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "icon_water_temperature",
			Help: "For each controller, reports the cooling or heating water temperature",
		}, genericParameters),
		externalTemperatureGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "icon_external_temperature",
			Help: "For each controller, reports external temperature",
		}, genericParameters),
		heatingGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "icon_heating",
			Help: "For each controller, reports 1 if the controller is set to heating mode, 0 otherwise",
		}, genericParameters),
		ecoGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "icon_eco",
			Help: "For each controller, reports 1 if the controller is in economy mode, 0 otherwise",
		}, genericParameters),
	}
}

func (r *systemMetricsReporter) Connected(sysId string, connected bool) {
	gauge := r.connectedGauge.WithLabelValues(sysId)
	if connected {
		gauge.Set(1)
	} else {
		gauge.Set(0)
	}
}

func (r *systemMetricsReporter) WaterTemperature(sysId string, temperature float64) {
	r.waterTemperatureGauge.WithLabelValues(sysId).Set(temperature)
}

func (r *systemMetricsReporter) ExternalTemperature(sysId string, temperature float64) {
	r.externalTemperatureGauge.WithLabelValues(sysId).Set(temperature)
}

func (r *systemMetricsReporter) Heating(sysId string, heating bool) {
	gauge := r.heatingGauge.WithLabelValues(sysId)
	if heating {
		gauge.Set(1)
	} else {
		gauge.Set(0)
	}
}

func (r *systemMetricsReporter) Eco(sysId string, eco bool) {
	gauge := r.ecoGauge.WithLabelValues(sysId)
	if eco {
		gauge.Set(1)
	} else {
		gauge.Set(0)
	}
}

func (r *systemMetricsReporter) RemoveDevice(sysId string) {
	r.connectedGauge.DeleteLabelValues(sysId)
	r.waterTemperatureGauge.DeleteLabelValues(sysId)
	r.externalTemperatureGauge.DeleteLabelValues(sysId)
	r.heatingGauge.DeleteLabelValues(sysId)
	r.ecoGauge.DeleteLabelValues(sysId)
}

// Room related required parameters
var roomParameters = append(genericParameters, "id", "room")

type RoomMetricsReporter interface {
	// Reports if the device connection is active or not.
	RoomConnected(sysId string, id string, room string, connected bool)
	// Reports the room temperature.
	RoomTemperature(sysId string, id string, room string, temperature float64)
	// Reports the room dew temperature.
	RoomDewTemperature(sysId string, id string, room string, temperature float64)
	// Reports the room target temperature.
	RoomTargetTemperature(sysId string, id string, room string, temperature float64)
	// Reports the relative humidity.
	RoomHumidity(sysId string, id string, room string, humidity float64)
	// Reports the relay state.
	RoomRelay(sysId string, id string, room string, connected bool)
	// Removes a room from reporting.
	RemoveRoom(sysId string, id string, room string)
}

type roomMetricsReporter struct {
	roomConntectedGauge        *prometheus.GaugeVec
	roomTemperatureGauge       *prometheus.GaugeVec
	roomDewTemperatureGauge    *prometheus.GaugeVec
	roomTargetTemperatureGauge *prometheus.GaugeVec
	roomHumidityGauge          *prometheus.GaugeVec
	roomRelayGauge             *prometheus.GaugeVec
}

func newRoomMetricsReporter() RoomMetricsReporter {
	return &roomMetricsReporter{
		roomConntectedGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "icon_room_connected",
			Help: "For each room, reports 1 if the room is connected to the controller, 0 otherwise",
		}, roomParameters),
		roomTemperatureGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "icon_temperature",
			Help: "For each room, reports the room temperature",
		}, roomParameters),
		roomDewTemperatureGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "icon_dew_temperature",
			Help: "For each room, reports the room dew temperature",
		}, roomParameters),
		roomTargetTemperatureGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "icon_target_temperature",
			Help: "For each room, reports the target temperature",
		}, roomParameters),
		roomHumidityGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "icon_humidity",
			Help: "For each room, reports the relative humidity",
		}, roomParameters),
		roomRelayGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "icon_relay_on",
			Help: "For each room, reports 1 if the relay is open, 0 otherwise",
		}, roomParameters),
	}
}

func (r *roomMetricsReporter) RoomConnected(sysId string, id string, room string, connected bool) {
	gauge := r.roomConntectedGauge.WithLabelValues(sysId, id, room)
	if connected {
		gauge.Set(1)
	} else {
		gauge.Set(0)
	}
}

func (r *roomMetricsReporter) RoomTemperature(sysId string, id string, room string, temperature float64) {
	r.roomTemperatureGauge.WithLabelValues(sysId, id, room).Set(temperature)
}

func (r *roomMetricsReporter) RoomDewTemperature(sysId string, id string, room string, temperature float64) {
	r.roomDewTemperatureGauge.WithLabelValues(sysId, id, room).Set(temperature)
}

func (r *roomMetricsReporter) RoomTargetTemperature(sysId string, id string, room string, temperature float64) {
	r.roomTargetTemperatureGauge.WithLabelValues(sysId, id, room).Set(temperature)
}

func (r *roomMetricsReporter) RoomHumidity(sysId string, id string, room string, humidity float64) {
	r.roomHumidityGauge.WithLabelValues(sysId, id, room).Set(humidity)
}

func (r *roomMetricsReporter) RoomRelay(sysId string, id string, room string, connected bool) {
	gauge := r.roomRelayGauge.WithLabelValues(sysId, id, room)
	if connected {
		gauge.Set(1)
	} else {
		gauge.Set(0)
	}
}

func (r *roomMetricsReporter) RemoveRoom(sysId string, id string, room string) {
	r.roomConntectedGauge.DeleteLabelValues(sysId, id, room)
	r.roomTemperatureGauge.DeleteLabelValues(sysId, id, room)
	r.roomDewTemperatureGauge.DeleteLabelValues(sysId, id, room)
	r.roomRelayGauge.DeleteLabelValues(sysId, id, room)
	r.roomHumidityGauge.DeleteLabelValues(sysId, id, room)
	r.roomTargetTemperatureGauge.DeleteLabelValues(sysId, id, room)
}

// HTTP response related required parameters
var httpParameters = append(genericParameters, "name", "response")

type HttpMetricsReporter interface {
	// Reports the HTTP response status code along with the duration.
	HttpClientRequest(sysId string, name string, statusCode int, duration time.Duration)
}

type httpMetricsReporter struct {
	httpSummary *prometheus.SummaryVec
}

func newHttpMetricsReporter() HttpMetricsReporter {
	return &httpMetricsReporter{
		httpSummary: promauto.NewSummaryVec(prometheus.SummaryOpts{
			Name: "icon_http_client_seconds",
			Help: "iCon HTTP client requests",
		}, httpParameters),
	}
}

func (r *httpMetricsReporter) HttpClientRequest(sysId string, name string, statusCode int, duration time.Duration) {
	r.httpSummary.WithLabelValues(sysId, name, strconv.Itoa(statusCode)).Observe(duration.Seconds())
}

type metricsReporter struct {
	HttpMetricsReporter
	RoomMetricsReporter
	SystemMetricsReporter
}

func NewPrometheusReporter() MetricsReporter {
	return &metricsReporter{
		SystemMetricsReporter: newSystemPrometheusReporter(),
		RoomMetricsReporter:   newRoomMetricsReporter(),
		HttpMetricsReporter:   newHttpMetricsReporter(),
	}
}
