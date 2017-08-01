package utils

import (
	"os"
	"os/signal"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("util")

func SetupCleanUpOnInterrupt(callback func()) chan bool {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	cleanupDone := make(chan bool)

	go func() {
		for range signalChan {
			log.Info("\nReceived an interrupt, shutting dow.\n")
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

func ValidateErr(err error, message string) bool {
	if err != nil {
		log.Infof("%s: %v\n", message, err)
	}
	return err != nil
}