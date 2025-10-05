package server

import (
	"net/http"
	"time"
)

func (s *iconMetricsServer) metricsHandlerFunc(w http.ResponseWriter, r *http.Request) {
	switch s.status {
	case ServerStatusStarting:
		w.Header().Add("Retry-After", time.Now().Add(time.Duration(3*time.Second)).UTC().Format(http.TimeFormat))
		w.WriteHeader(http.StatusServiceUnavailable)
	case ServerStatusStopping:
		w.WriteHeader(http.StatusServiceUnavailable)
	default:
		s.metricsHandler.ServeHTTP(w, r)
	}
}
