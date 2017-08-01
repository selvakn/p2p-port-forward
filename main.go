package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"
	"unsafe"
	"github.com/op/go-logging"
	"io"
)

/*
#cgo CFLAGS: -I ./libzt/include
#cgo darwin LDFLAGS: -L ./libzt/darwin/ -lzt -lstdc++

#include "libzt.h"
#include <stdlib.h>
#include <stdio.h>
#include <unistd.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <string.h>
#include <netdb.h>
*/
import "C"

const NETWORK_ID = "8056c2e21c000001"
const PORT = 50718 // 7878
const BUF_SIZE = 2800

var log = logging.MustGetLogger("tunneler")
var totalBytesSent = 0
var totalBytesReceived = 0

func setupCleanUpOnInterrupt() chan bool {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	cleanupDone := make(chan bool)

	go func() {
		for range signalChan {
			log.Info("\nReceived an interrupt, shutting dow.\n")
			cleanupDone <- true
		}
	}()
	return cleanupDone
}

func getOtherIP() string {
	if len(os.Args) >= 2 {
		return os.Args[1]
	} else {
		return ""
	}
}

func validate(value C.int, message string) {
	if value < 0 {
		log.Info(message)
		os.Exit(1)
	}
}

func bindAndListen() (C.int) {
	sockfd := C.zts_socket(syscall.AF_INET6, syscall.SOCK_STREAM, 0)
	validate(sockfd, "Error in opening socket")

	serverSocket := syscall.RawSockaddrInet6{Flowinfo: 0, Family: syscall.AF_INET6, Port: PORT}
	retVal := C.zts_bind(sockfd, (*C.struct_sockaddr)(unsafe.Pointer(&serverSocket)), C.sizeof_struct_sockaddr_in6)
	validate(retVal, "ERROR on binding")
	log.Debugf("Bind Complete")

	C.zts_listen(sockfd, 1)
	log.Debugf("Listening")

	return sockfd
}

func accept(sockfd C.int) C.int {
	clientSocket := syscall.RawSockaddrInet6{}
	clientSocketLength := C.sizeof_struct_sockaddr_in6
	connSockfd := C.zts_accept(sockfd, (*C.struct_sockaddr)(unsafe.Pointer(&clientSocket)), (*C.socklen_t)(unsafe.Pointer(&clientSocketLength)))

	validate(connSockfd, "ERROR on accept")

	log.Info("Accepted incoming connection from client")

	return connSockfd
}

func parseIPV6(ipString string) [16]byte {
	ip := net.ParseIP(ipString)
	var arr [16]byte
	copy(arr[:], ip)
	return arr
}

func initZT() {
	C.zts_simple_start(C.CString("./zt"), C.CString(NETWORK_ID))

	ipv4Address := make([]byte, C.ZT_MAX_IPADDR_LEN)
	ipv6Address := make([]byte, C.ZT_MAX_IPADDR_LEN)

	C.zts_get_ipv4_address(C.CString(NETWORK_ID), (*C.char)(unsafe.Pointer(&ipv4Address[0])), C.ZT_MAX_IPADDR_LEN)
	log.Infof("ipv4 = %s \n", string(ipv4Address[:C.ZT_MAX_IPADDR_LEN]))

	C.zts_get_ipv6_address(C.CString(NETWORK_ID), (*C.char)(unsafe.Pointer(&ipv6Address[0])), C.ZT_MAX_IPADDR_LEN)
	log.Infof("ipv6 = %s \n", string(ipv6Address[:C.ZT_MAX_IPADDR_LEN]))
}

func connectToOther() C.int {
	arr := parseIPV6(getOtherIP())

	clientSocket := syscall.RawSockaddrInet6{Flowinfo: 0, Family: syscall.AF_INET6, Port: PORT, Addr: arr}

	sockfd := C.zts_socket(syscall.AF_INET6, syscall.SOCK_STREAM, 0)
	validate(sockfd, "Error in opening socket")

	retVal := C.zts_connect(sockfd, (*C.struct_sockaddr)(unsafe.Pointer(&clientSocket)), C.sizeof_struct_sockaddr_in6)
	validate(retVal, "Error in connect client")

	return sockfd
}

func validateErr(err error, message string) bool {
	if err != nil {
		log.Infof("%s: %v\n", message, err)
	}
	return err != nil
}

func bridge(readWriteCloser io.ReadWriteCloser, sockfd C.int) {

	buffer1 := make([]byte, BUF_SIZE)
	go func() {
		for {
			rlen, err := readWriteCloser.Read(buffer1)
			if err == io.EOF || validateErr(err, "Error reading from stream") {
				readWriteCloser.Close()
				break
			}

			wlen, writeErr := syscall.Write((int)(sockfd), buffer1[:rlen])
			if validateErr(writeErr, "Error writing to zt") {
				break
			}

			totalBytesSent += wlen
			log.Debugf("Total sent so far: %d\n", totalBytesSent)
		}
	}()

	buffer2 := make([]byte, BUF_SIZE)
	go func() {
		for {
			rlen, err := syscall.Read((int)(sockfd), buffer2)

			if rlen == 0 || validateErr(err, "Error reading from zt") {
				break
			}

			wlen, writeErr := readWriteCloser.Write(buffer2[:rlen])
			if validateErr(writeErr, "Error writing to stream") {
				break
			}

			totalBytesReceived += wlen
			log.Debugf("Total received so far: %d\n", totalBytesReceived)
		}
	}()
}

func main() {
	initZT()
	defer C.zts_stop()

	if len(getOtherIP()) == 0 {
		sockfd := bindAndListen()
		defer C.zts_close(sockfd)

		//go func() {
		//	for {
				newSockfd := accept(sockfd)

				conn, _ := net.Dial("tcp", "localhost:22")
				bridge(conn, newSockfd)
			//}
		//}()
	} else {
		sockfd := connectToOther()
		defer C.zts_close(sockfd)

		ln, _ := net.Listen("tcp", ":2222")
		defer ln.Close()

		conn, _ := ln.Accept()
		bridge(conn, sockfd)
	}

	<-setupCleanUpOnInterrupt()
}
