// Main application package.
package main

import (
	"context"
	"errors"
	"log"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/csutorasa/icon-metrics/internal/args"
	"github.com/csutorasa/icon-metrics/internal/config"
	"github.com/csutorasa/icon-metrics/internal/metrics"
	"github.com/csutorasa/icon-metrics/internal/metrics/prometheus"
	"github.com/csutorasa/icon-metrics/internal/server"
	"github.com/csutorasa/icon-metrics/internal/signal"
	"github.com/csutorasa/icon-metrics/pkg/client"
)

// Main logger instance
var logger *log.Logger = log.Default()

func main() {
	t := metrics.NewTimer()
	args := args.ParseArgs()
	logger.Printf("Arguments - loaded in %s %s", t.End(), args)
	c := readConfigFile(args.Config)
	m := prometheus.RegisterMetrics()
	s := startServer(c.Port)
	defer closeServer(s)

	clients := createClients(c.Devices, m)
	if len(clients) == 0 {
		logger.Printf("System - no clients could be initialized, exiting")
		return
	}

	var wg sync.WaitGroup
	cancelFuncs := []context.CancelFunc{}
	for client, device := range clients {
		cancel := reportValuesLoop(client, device, m, &wg)
		cancelFuncs = append(cancelFuncs, cancel)
	}
	s.SetStatus(server.ServerStatusOK)
	interruptHandler(s, cancelFuncs)
	wg.Wait()
}

// Returns configuration from file.
func readConfigFile(configPath string) *config.Configuration {
	t := metrics.NewTimer()
	logger.Printf("Configuration - reading from %s", configPath)
	c, err := config.ReadConfigFile(configPath)
	if err != nil {
		logger.Panicf("Configuration - failed to load caused by %s", err.Error())
	}
	logger.Printf("Configuration - loaded in %s %d devices", t.End(), len(c.Devices))
	return c
}

func startServer(port int) server.IconMetricsServer {
	t := metrics.NewTimer()
	logger.Printf("Server - starting on port %d", port)
	s := server.NewIconMetricsServer(port)
	err := s.Start()
	if err != nil {
		logger.Panicf("Server - failed to start caused by %s", err.Error())
	}
	logger.Printf("Server - started in %s", t.End())
	return s
}

func closeServer(s server.IconMetricsServer) {
	t := metrics.NewTimer()
	logger.Printf("Server - stopping")
	s.Close()
	logger.Printf("Server - stopped in %s", t.End())
}

// Handles OS signals for shutdown.
func interruptHandler(s server.IconMetricsServer, cancelFuncs []context.CancelFunc) {
	c := signal.InterruptChannel()
	go func() {
		for {
			switch <-c {
			case signal.InterruptTypeGraceful:
				s.SetStatus(server.ServerStatusStopping)
				logger.Printf("System - graceful shutdown initiated")
				for _, c := range cancelFuncs {
					c()
				}
			case signal.InterruptTypeForced:
				s.SetStatus(server.ServerStatusStopping)
				logger.Printf("System - forced shutdown initiated")
				os.Exit(0)
			}
		}
	}()
}

func createClients(devices []*config.IconConfiguration, m *metrics.Metrics) map[client.IconClient]*config.IconConfiguration {
	clients := map[client.IconClient]*config.IconConfiguration{}
	for _, device := range devices {
		client, err := newIconClient(device, m)
		if err != nil {
			logger.Printf("Client %s - failed to create client for %s", device.SysId, device.Url)
			continue
		}
		clients[client] = device
	}
	return clients
}

func newIconClient(device *config.IconConfiguration, m *metrics.Metrics) (client.IconClient, error) {
	u, err := url.Parse(device.Url)
	if err != nil {
		return nil, err
	}
	c, err := client.NewIconClient(u, device.SysId, device.Password)
	if err != nil {
		return nil, err
	}
	return metrics.NewIconMetricsClient(c, m, device.Report), nil
}

func closeIconClient(client client.IconClient) {
	t := metrics.NewTimer()
	logger.Printf("Client %s - disonnecting", client.SysId())
	err := client.Close()
	if err != nil {
		logger.Printf("Client %s - failed to disonnect caused by %s", client.SysId(), err.Error())
	} else {
		logger.Printf("Client %s - disconnected in %s", client.SysId(), t.End())
	}
}

func reportValuesLoop(client client.IconClient, device *config.IconConfiguration, m *metrics.Metrics, wg *sync.WaitGroup) context.CancelFunc {
	delay := time.Duration(device.Delay) * time.Second
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer closeIconClient(client)
		reportValues(ctx, client, delay, m, device.Report)
	}()
	return cancel
}

// Main loop for handling a single iCON device.
func reportValues(ctx context.Context, c client.IconClient, d time.Duration, m *metrics.Metrics, reportConfig *config.ReportConfiguration) {
	reportedRooms := map[metrics.RoomDescriptor]bool{}
	m.Connected.Set(c.SysId(), 0)
	for {
		values, err := c.ReadValuesContext(ctx)
		if err != nil {
			logger.Printf("Client %s - failed to read values caused by %s", c.SysId(), err.Error())
			m.RemoveSystem(c.SysId())
			value := sleep(ctx, d)
			if value > 0 {
				break
			}
			continue
		}
		m.Connected.Set(c.SysId(), 1)
		metrics.Report(values, c.SysId(), m, reportConfig, &reportedRooms)
		value := sleep(ctx, d)
		if value > 0 {
			break
		}
	}
}

// Sleeps for the duration, which can be interrupted.
func sleep(ctx context.Context, d time.Duration) int {
	c, cancel := context.WithTimeout(ctx, d)
	defer cancel()
	<-c.Done()
	if errors.Is(c.Err(), context.DeadlineExceeded) {
		return 0
	}
	return 1
}
