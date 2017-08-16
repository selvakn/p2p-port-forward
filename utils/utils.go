package utils

import (
	"os"
	"os/signal"
	"p2p-port-forward/logger"
)

var log = logger.Logger

func SetupCleanUpOnInterrupt(callback func()) chan bool {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	cleanupDone := make(chan bool)

	go func() {
		for range signalChan {
			log.Info("\nReceived an interrupt, shutting dow.")
			callback()

			cleanupDone <- true
		}
	}()
	return cleanupDone
}

func Validate(value int, message string) {
	if value < 0 {
		log.Error(message)
		os.Exit(1)
	}
}
