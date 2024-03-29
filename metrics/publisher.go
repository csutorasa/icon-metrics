package metrics

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// HTTP server
type PrometheusPublisher interface {
	io.Closer
	// Starts to listen and serve.
	Start() error
	// Stops serving and listening.
	Stop(context context.Context) error
}

// HTTP server
type prometheusPublisher struct {
	server *http.Server
}

// Creates a new server with the given port
func NewPrometheusPublisher(port int) PrometheusPublisher {
	publisher := &prometheusPublisher{}
	promhttpHandler := promhttp.Handler()
	mux := http.NewServeMux()
	mux.Handle("GET /metrics", promhttpHandler)
	mux.HandleFunc("GET /status", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	publisher.server = &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return publisher
}

// Starts to listen and serve.
func (publisher *prometheusPublisher) Start() error {
	ln, err := net.Listen("tcp", publisher.server.Addr)
	if err != nil {
		return err
	}
	go func() {
		publisher.server.Serve(ln)
	}()
	return nil
}

// Stops serving and listening.
func (publisher *prometheusPublisher) Stop(context context.Context) error {
	return publisher.server.Shutdown(context)
}

// Cleans up resources.
func (publisher *prometheusPublisher) Close() error {
	return publisher.server.Close()
}
