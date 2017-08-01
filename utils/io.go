package utils

import (
	"io"
	"syscall"
)

import "C"

var totalBytesSent = 0
var totalBytesReceived = 0
const BUF_SIZE = 2800


func ReceiveFrom(toWriter io.Writer, fromSockFd int) {
	buffer := make([]byte, BUF_SIZE)
	for {
		rlen, err := syscall.Read(fromSockFd, buffer)

		if rlen == 0 || ValidateErr(err, "Error reading from zt") {
			break
		}

		wlen, writeErr := toWriter.Write(buffer[:rlen])
		if ValidateErr(writeErr, "Error writing to stream") {
			break
		}

		totalBytesReceived += wlen
		log.Debugf("Total received so far: %d\n", totalBytesReceived)
	}
}

func SendTo(toSockfd int, fromReader io.Reader) {
	buffer := make([]byte, BUF_SIZE)
	for {
		rlen, err := fromReader.Read(buffer)
		if err == io.EOF || ValidateErr(err, "Error reading from stream") {
			log.Info("Stream conn closed")
			break
		}

		wlen, writeErr := syscall.Write(toSockfd, buffer[:rlen])
		if ValidateErr(writeErr, "Error writing to zt") {
			break
		}

		totalBytesSent += wlen
		log.Debugf("Total sent so far: %d\n", totalBytesSent)
	}
}

