package main

import (
	"flag"
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
	"github.com/csutorasa/icon-metrics/publisher"
)

var logger *log.Logger = log.Default()
var startTime time.Time = time.Now()

func main() {
	configPath := parseArgs()

	c, err := readConfig(configPath)
	if err != nil {
		logger.Panicf("Failed to load configuration caused by %s", err.Error())
	}

	metrics.RegisterUptime()

	logger.Printf("Starting prometheus server on port %d", c.Port)
	start := metrics.NewTimer()
	p := publisher.NewPrometherusPublisher(c.Port)

	err = p.Start()
	if err != nil {
		logger.Panicf("Failed to start prometheus server on port %d caused by %s", c.Port, err.Error())
	}
	defer func() {
		start := metrics.NewTimer()
		logger.Printf("Stopping prometheus server on port %d", c.Port)
		p.Close()
		logger.Printf("Successfully stopped Prometheus server on port %d under %s", c.Port, start.End().String())
	}()
	logger.Printf("Successfully started Prometheus server on port %d under %s", c.Port, start.End().String())

	var wg sync.WaitGroup
	channels := make([]chan int, 0)
	for _, device := range c.Devices {
		client, err := client.NewIconClient(device.Url, device.SysId, device.Password)
		if err != nil {
			logger.Printf("Failed to create client for device %s @ %s", device.SysId, device.Url)
			continue
		}
		delay := time.Duration(device.Delay) * time.Second
		ch := make(chan int, 0)
		channels = append(channels, ch)
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				start := metrics.NewTimer()
				logger.Printf("Disonnecting from %s", client.SysId())
				err := client.Close()
				if err != nil {
					logger.Printf("Failed to disonnect from %s caused by %s", client.SysId(), err.Error())
				} else {
					logger.Printf("Successfully disconnected from %s under %s", client.SysId(), start.End().String())
				}
			}()
			reportValues(client, ch, delay)
		}()
	}
	if len(channels) != 0 {
		interruptHandler(channels)
		wg.Wait()
	}
}

func parseArgs() string {
	configPath := flag.String("config", "", "Configuration file url")
	flag.Parse()
	if *configPath == "" {
		dir := filepath.Dir(os.Args[0])
		return filepath.Join(dir, "config.yml")
	}
	return *configPath
}

func readConfig(configPath string) (*config.Configuration, error) {
	logger.Printf("Loading configuration from %s", configPath)
	c, err := config.ReadConfig(configPath)
	if err != nil {
		return nil, err
	}
	logger.Printf("Configuration is loaded from %s", configPath)
	return c, nil
}

func interruptHandler(channels []chan int) {
	c := make(chan os.Signal)
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
				break
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

func reportValues(c client.IconClient, trigger chan int, d time.Duration) {
	session := client.NewSession(c.SysId())
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

func sleep(trigger chan int, d time.Duration) int {
	go func() {
		time.Sleep(d)
		trigger <- 0
	}()
	return <-trigger
}
