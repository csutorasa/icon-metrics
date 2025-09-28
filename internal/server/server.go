package server

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
type IconMetricsServer interface {
	io.Closer
	// Starts to listen and serve.
	Start() error
	// Stops serving and listening.
	Stop(context context.Context) error
}

// HTTP server
type iconMetricsServer struct {
	server *http.Server
}

// Creates a new server with the given port
func NewIconMetricsServer(port int) IconMetricsServer {
	server := &iconMetricsServer{}
	mux := http.NewServeMux()
	mux.Handle("GET /metrics", promhttp.Handler())
	mux.HandleFunc("GET /status", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	server.server = &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return server
}

// Starts to listen and serve.
func (publisher *iconMetricsServer) Start() error {
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
func (publisher *iconMetricsServer) Stop(context context.Context) error {
	return publisher.server.Shutdown(context)
}

// Cleans up resources.
func (publisher *iconMetricsServer) Close() error {
	return publisher.server.Close()
}
