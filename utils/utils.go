package utils

import (
	"github.com/google/logger"
	"os"
	"os/signal"
)

func SetupCleanUpOnInterrupt(callback func()) chan bool {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	cleanupDone := make(chan bool)

	go func() {
		for range signalChan {
			logger.Info("\nReceived an interrupt, shutting dow.\n")
			callback()

			cleanupDone <- true
		}
	}()
	return cleanupDone
}

func Validate(value int, message string) {
	if value < 0 {
		logger.Error(message)
		os.Exit(1)
	}
}
