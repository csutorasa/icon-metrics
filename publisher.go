package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Publisher struct {
	server *http.Server
}

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

func (this *Publisher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func (this *Publisher) Start() error {
	ln, err := net.Listen("tcp", this.server.Addr)
	if err != nil {
		return err
	}
	go func() {
		this.server.Serve(ln)
	}()
	return nil
}

func (this *Publisher) Stop(context context.Context) error {
	return this.server.Shutdown(context)
}

func (this *Publisher) Close() error {
	return this.server.Close()
}
