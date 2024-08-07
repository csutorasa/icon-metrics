// Main application package.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/csutorasa/icon-metrics/client"
	"github.com/csutorasa/icon-metrics/config"
	"github.com/csutorasa/icon-metrics/metrics"
)

// Main logger instance
var logger *log.Logger = log.Default()

func main() {
	configPath := parseArgs()

	c, err := readConfig(configPath)
	if err != nil {
		logger.Panicf("Failed to load configuration caused by %s", err.Error())
	}

	reporter := metrics.NewPrometheusReporter()
	reporter.Uptime()

	logger.Printf("Starting http server on port %d", c.Port)
	start := metrics.NewTimer()
	p := metrics.NewPrometheusPublisher(c.Port)

	err = p.Start()
	if err != nil {
		logger.Panicf("Failed to start http server on port %d caused by %s", c.Port, err.Error())
	}
	defer func() {
		start := metrics.NewTimer()
		logger.Printf("Stopping http server on port %d", c.Port)
		p.Close()
		logger.Printf("Successfully stopped http server on port %d under %s", c.Port, start.End().String())
	}()
	logger.Printf("Successfully started http server on port %d under %s", c.Port, start.End().String())

	var wg sync.WaitGroup
	channels := make([]chan int, 0)
	for _, device := range c.Devices {
		reportConfig := device.Report
		session := metrics.NewSession(device.SysId, reportConfig, reporter)
		client, err := client.NewIconClient(device.Url, device.SysId, device.Password, session)
		if err != nil {
			logger.Printf("Failed to create client for device %s @ %s", device.SysId, device.Url)
			continue
		}
		delay := time.Duration(device.Delay) * time.Second
		ch := make(chan int)
		channels = append(channels, ch)
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				start := metrics.NewTimer()
				logger.Printf("Disonnecting from %s", client.SysId())
				err := client.Close()
				reporter.RemoveDevice(client.SysId())
				if err != nil {
					logger.Printf("Failed to disonnect from %s caused by %s", client.SysId(), err.Error())
				} else {
					logger.Printf("Successfully disconnected from %s under %s", client.SysId(), start.End().String())
				}
			}()
			reportValues(client, ch, delay, session)
		}()
	}
	if len(channels) != 0 {
		interruptHandler(channels)
		wg.Wait()
	}
}

// Parses configuration file path from command line options
func parseArgs() string {
	configPath := flag.String("config", "", "Configuration file url")
	flag.Parse()
	if *configPath == "" {
		dir := filepath.Dir(os.Args[0])
		return filepath.Join(dir, "config.yml")
	}
	return *configPath
}

// Returns configuration from file.
func readConfig(configPath string) (*config.Configuration, error) {
	logger.Printf("Loading configuration from %s", configPath)
	c, err := config.ReadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	logger.Printf("Configuration is loaded from %s", configPath)
	return c, nil
}

// Handles OS signals for shutdown.
func interruptHandler(channels []chan int) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
	closing := false
	go func() {
		for {
			switch <-c {
			case syscall.SIGINT:
				if !closing {
					closing = true
					logger.Printf("SIGINT received, graceful shutdown initiated")
					for _, c := range channels {
						ch := c

						go func() {
							ch <- 1
						}()
					}
				} else {
					logger.Printf("SIGINT received again, force shutdown initiated")
					os.Exit(0)
				}
			case syscall.SIGTERM:
				logger.Printf("SIGTERM received, force shutdown initiated")
				os.Exit(0)
			case syscall.SIGABRT:
				logger.Printf("SIGABRT received, force shutdown initiated")
				os.Exit(0)
			}
		}
	}()
}

// Main loop for handling a single iCON device.
func reportValues(c client.IconClient, trigger chan int, d time.Duration, session metrics.MetricsSession) {
	session.Connected(false)
	for {
		if !c.IsLoggedIn() {
			logger.Printf("Connecting to %s", c.SysId())
			err := c.Login()
			if err != nil {
				logger.Printf("Failed to connect to %s caused by %s", c.SysId(), err.Error())
				session.Reset()
				value := sleep(trigger, d)
				if value > 0 {
					break
				}
				continue
			}
			logger.Printf("Connected to %s", c.SysId())
			session.Connected(true)
		}
		values, err := c.ReadValues()
		if err != nil {
			logger.Printf("Failed to read values from %s caused by %s", c.SysId(), err.Error())
			session.Reset()
			value := sleep(trigger, d)
			if value > 0 {
				break
			}
			continue
		}
		session.Report(values)
		value := sleep(trigger, d)
		if value > 0 {
			break
		}
	}
}

// Sleeps for the duration, which can be interrupted.
func sleep(trigger chan int, d time.Duration) int {
	go func() {
		time.Sleep(d)
		trigger <- 0
	}()
	return <-trigger
}
