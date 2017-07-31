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

func reverse(numbers []C.int) []C.int {
	for i := 0; i < len(numbers)/2; i++ {
		j := len(numbers) - i - 1
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
	return numbers
}

func setupCleanUpOnInterrupt(fdsToClose []C.int) chan bool {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	cleanupDone := make(chan bool)

	go func() {
		for range signalChan {
			log.Info("\nReceived an interrupt, shutting dow.\n")
			log.Infof("Going to close :%v", fdsToClose)
			for _, fd := range reverse(fdsToClose) {
				C.zts_close(fd)
			}
			log.Infof("Closing Complete");
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

func bindAndListen(onAccept func(newSockfd C.int)) (C.int) {
	sockfd := C.zts_socket(syscall.AF_INET6, syscall.SOCK_STREAM, 0)
	validate(sockfd, "Error in opening socket")

	serverSocket := syscall.RawSockaddrInet6{Flowinfo: 0, Family: syscall.AF_INET6, Port: PORT}
	retVal := C.zts_bind(sockfd, (*C.struct_sockaddr)(unsafe.Pointer(&serverSocket)), C.sizeof_struct_sockaddr_in6)
	validate(retVal, "ERROR on binding")
	log.Debugf("Bind Complete")

	C.zts_listen(sockfd, 1)
	log.Debugf("Listening")

	//go func() {
	//	for {
			clientSocket := syscall.RawSockaddrInet6{}
			clientSocketLength := C.sizeof_struct_sockaddr_in6
			connSockfd := C.zts_accept(sockfd, (*C.struct_sockaddr)(unsafe.Pointer(&clientSocket)), (*C.socklen_t)(unsafe.Pointer(&clientSocketLength)))

			validate(connSockfd, "ERROR on accept")

			log.Info("Accepted incoming connection from client")

			onAccept(connSockfd)
		//}
	//}()

	return sockfd
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

func bridge(readWriteCloser io.ReadWriteCloser, sockfd int) {

	buffer1 := make([]byte, BUF_SIZE)
	go func() {
		for {
			rlen, err := readWriteCloser.Read(buffer1)
			if err == io.EOF || validateErr(err, "Error reading from stream") {
				readWriteCloser.Close()
				break
			}

			wlen, writeErr := syscall.Write(sockfd, buffer1[:rlen])
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
			rlen, err := syscall.Read(sockfd, buffer2)
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

	var allFDs []C.int

	if len(getOtherIP()) == 0 {
		sockfd := bindAndListen(func(finalSockfd C.int) {
			conn, _ := net.Dial("tcp", "localhost:22")

			bridge(conn, (int)(finalSockfd))
			allFDs = append(allFDs, finalSockfd)
		})
		defer C.zts_close(sockfd)
	} else {
		sockfd := connectToOther()
		defer C.zts_close(sockfd)

		ln, _ := net.Listen("tcp", ":2222")
		defer ln.Close()

		conn, _ := ln.Accept()
		bridge(conn, (int)(sockfd))
	}

	<-setupCleanUpOnInterrupt(allFDs)
}
