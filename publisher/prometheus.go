package publisher

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type prometheusPublisher struct {
	server   *http.Server
}

func NewPrometherusPublisher(port int) Publisher {
	publisher := &prometheusPublisher{}
	publisher.server = &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        publisher,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return publisher
}

func (this *prometheusPublisher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/metrics" {
		promhttp.Handler().ServeHTTP(w, r)
	} else {
		w.WriteHeader(404)
	}
}

func (this *prometheusPublisher) Start() error {
	ln, err := net.Listen("tcp", this.server.Addr)
	if err != nil {
		return err
	}
	go func() {
		this.server.Serve(ln)
	}()
	return nil
}

func (this *prometheusPublisher) Stop(context context.Context) error {
	return this.server.Shutdown(context)
}

func (this *prometheusPublisher) Close() error {
	return this.server.Close()
}
