package signal

import (
	"os"
	"os/signal"
	"syscall"
)

type InterruptType int

const (
	InterruptTypeNone InterruptType = iota
	InterruptTypeGraceful
	InterruptTypeForced
)

func InterruptChannel() <-chan InterruptType {
	interrupt := make(chan os.Signal, 1)
	output := make(chan InterruptType)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
	closing := false
	go func() {
		for {
			switch <-interrupt {
			case syscall.SIGINT:
				if !closing {
					closing = true
					go func() {
						output <- InterruptTypeGraceful
					}()
				} else {
					output <- InterruptTypeForced
					os.Exit(0)
				}
			case syscall.SIGTERM, syscall.SIGABRT:
				output <- InterruptTypeForced
				os.Exit(0)
			}
		}
	}()
	return output
}
