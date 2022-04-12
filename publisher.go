package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// HTTP server
type Publisher struct {
	server *http.Server
}

// Creates a new server with the given port
func NewPrometherusPublisher(port int) *Publisher {
	publisher := &Publisher{}
	publisher.server = &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        publisher,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return publisher
}

// Main logic of the server.
func (publisher *Publisher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/metrics" {
		if r.Method == http.MethodGet {
			promhttp.Handler().ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	} else if r.URL.Path == "/status" {
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// Starts to listen and serve.
func (publisher *Publisher) Start() error {
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
func (publisher *Publisher) Stop(context context.Context) error {
	return publisher.server.Shutdown(context)
}

// Cleans up resources.
func (publisher *Publisher) Close() error {
	return publisher.server.Close()
}
