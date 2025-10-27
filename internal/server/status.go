package server

import "net/http"

// HTTP server status
type ServerStatus string

const (
	ServerStatusStarting ServerStatus = "starting"
	ServerStatusOK       ServerStatus = "OK"
	ServerStatusStopping ServerStatus = "stopping"
)

// Sets the status of the server.
func (s *iconMetricsServer) SetStatus(status ServerStatus) {
	s.status = status
}

func (s *iconMetricsServer) statusHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(s.status))
}
