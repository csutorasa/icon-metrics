package prometheus

import (
	"strconv"
	"time"

	"github.com/csutorasa/icon-metrics/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var genericParameters = []string{"sysId"}

type PrometheusSystemMetric struct {
	gauge *prometheus.GaugeVec
}

func NewPrometheusSystemMetric(opts prometheus.GaugeOpts) *PrometheusSystemMetric {
	return &PrometheusSystemMetric{
		gauge: promauto.NewGaugeVec(opts, genericParameters),
	}
}

func (m *PrometheusSystemMetric) Set(sysId string, f float64) {
	m.gauge.WithLabelValues(sysId).Set(f)
}

func (m *PrometheusSystemMetric) Remove(sysId string) {
	m.gauge.DeleteLabelValues(sysId)
}

var roomParameters = append(genericParameters, "id", "room")

type PrometheusRoomMetric struct {
	gauge *prometheus.GaugeVec
}

func NewPrometheusRoomMetric(opts prometheus.GaugeOpts) *PrometheusRoomMetric {
	return &PrometheusRoomMetric{
		gauge: promauto.NewGaugeVec(opts, roomParameters),
	}
}

func (m *PrometheusRoomMetric) Set(sysId string, roomId string, roomName string, f float64) {
	m.gauge.WithLabelValues(sysId, roomId, roomName).Set(f)
}

func (m *PrometheusRoomMetric) RemoveSystem(sysId string) {
	m.gauge.DeletePartialMatch(prometheus.Labels{"sysId": sysId})
}

func (m *PrometheusRoomMetric) Remove(sysId string, roomId string, roomName string) {
	m.gauge.DeleteLabelValues(sysId, roomId, roomName)
}

var httpParameters = append(genericParameters, "name", "response")

type PrometheushttpMetric struct {
	summary *prometheus.SummaryVec
}

func NewPrometheusHttpMetric(opts prometheus.SummaryOpts) *PrometheushttpMetric {
	return &PrometheushttpMetric{
		summary: promauto.NewSummaryVec(opts, httpParameters),
	}
}

func (m *PrometheushttpMetric) Observe(sysId string, name string, statusCode int, duration time.Duration) {
	m.summary.WithLabelValues(sysId, name, strconv.Itoa(statusCode)).Observe(duration.Seconds())
}

func (m *PrometheushttpMetric) RemoveSystem(sysId string) {
	m.summary.DeletePartialMatch(prometheus.Labels{"sysId": sysId})
}

// Registers application metrics
func RegisterMetrics() *metrics.Metrics {
	promauto.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "uptime",
		Help: "Uptime for the service",
	}, func() float64 { return float64(metrics.Uptime().Seconds()) })
	return &metrics.Metrics{
		Connected: NewPrometheusSystemMetric(prometheus.GaugeOpts{
			Name: "icon_controller_connected",
			Help: "For each controller, reports 1 if the controller is ready to be read, 0 otherwise.",
		}),
		WaterTemperature: NewPrometheusSystemMetric(prometheus.GaugeOpts{
			Name: "icon_water_temperature",
			Help: "For each controller, reports the cooling or heating water temperature",
		}),
		ExternalTemperature: NewPrometheusSystemMetric(prometheus.GaugeOpts{
			Name: "icon_external_temperature",
			Help: "For each controller, reports external temperature",
		}),
		Heating: NewPrometheusSystemMetric(prometheus.GaugeOpts{
			Name: "icon_heating",
			Help: "For each controller, reports 1 if the controller is set to heating mode, 0 otherwise",
		}),
		Eco: NewPrometheusSystemMetric(prometheus.GaugeOpts{
			Name: "icon_eco",
			Help: "For each controller, reports 1 if the controller is in economy mode, 0 otherwise",
		}),
		RoomConnected: NewPrometheusRoomMetric(prometheus.GaugeOpts{
			Name: "icon_room_connected",
			Help: "For each room, reports 1 if the room is connected to the controller, 0 otherwise",
		}),
		RoomTemperature: NewPrometheusRoomMetric(prometheus.GaugeOpts{
			Name: "icon_temperature",
			Help: "For each room, reports the room temperature",
		}),
		RoomDewTemperature: NewPrometheusRoomMetric(prometheus.GaugeOpts{
			Name: "icon_dew_temperature",
			Help: "For each room, reports the room dew temperature",
		}),
		RoomTargetTemperature: NewPrometheusRoomMetric(prometheus.GaugeOpts{
			Name: "icon_target_temperature",
			Help: "For each room, reports the target temperature",
		}),
		RoomHumidity: NewPrometheusRoomMetric(prometheus.GaugeOpts{
			Name: "icon_humidity",
			Help: "For each room, reports the relative humidity",
		}),
		RoomRelay: NewPrometheusRoomMetric(prometheus.GaugeOpts{
			Name: "icon_relay_on",
			Help: "For each room, reports 1 if the relay is open, 0 otherwise",
		}),
		Http: NewPrometheusHttpMetric(prometheus.SummaryOpts{
			Name: "icon_http_client_seconds",
			Help: "iCon HTTP client requests",
		}),
	}
}
