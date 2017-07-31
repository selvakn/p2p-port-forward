package main

import (
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"unsafe"

	"github.com/songgao/water"
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

func bindAndListen() (int, int) {
	sockfd := C.zts_socket(syscall.AF_INET6, syscall.SOCK_STREAM, 0)
	validate(sockfd, "Error in opening socket")

	serverSocket := syscall.RawSockaddrInet6{Flowinfo: 0, Family: syscall.AF_INET6, Port: PORT}
	retVal := C.zts_bind(sockfd, (*C.struct_sockaddr)(unsafe.Pointer(&serverSocket)), C.sizeof_struct_sockaddr_in6)
	validate(retVal, "ERROR on binding")
	log.Debugf("Bind Complete")

	C.zts_listen(sockfd, 1)
	log.Debugf("Listening")

	clientSocket := syscall.RawSockaddrInet6{}
	clientSocketLength := C.sizeof_struct_sockaddr_in6
	connSockfd := C.zts_accept(sockfd, (*C.struct_sockaddr)(unsafe.Pointer(&clientSocket)), (*C.socklen_t)(unsafe.Pointer(&clientSocketLength)))
	validate(connSockfd, "ERROR on accept")

	log.Info("Accepted incoming connection from client")

	return int(sockfd), int(connSockfd)
}

func parseIPV6(ipString string) [16]byte {
	ip := net.ParseIP(ipString)
	var arr [16]byte
	copy(arr[:], ip)
	return arr
}

func ifconfig(args ...string) {
	cmd := exec.Command("ifconfig", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if nil != err {
		log.Infof("Error running command:", err)
	}
}

func setupTun(initater bool) *water.Interface {
	iface, _ := water.New(water.Config{
		DeviceType: water.TUN,
	})

	log.Infof("Interface Name: %s\n", iface.Name())

	if initater {
		ifconfig(iface.Name(), "10.1.0.10", "10.1.0.20", "up")
	} else {
		ifconfig(iface.Name(), "10.1.0.20", "10.1.0.10", "up")
	}

	return iface
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

func connectToOther() int {
	arr := parseIPV6(getOtherIP())

	clientSocket := syscall.RawSockaddrInet6{Flowinfo: 0, Family: syscall.AF_INET6, Port: PORT, Addr: arr}

	sockfd := C.zts_socket(syscall.AF_INET6, syscall.SOCK_STREAM, 0)
	validate(sockfd, "Error in opening socket")

	retVal := C.zts_connect(sockfd, (*C.struct_sockaddr)(unsafe.Pointer(&clientSocket)), C.sizeof_struct_sockaddr_in6)
	validate(retVal, "Error in connect client")

	return (int)(sockfd)
}

func validateErr(err error, message string) bool {
	if err != nil {
		log.Infof("%s: %v\n", message, err)
	}
	return err != nil
}

func bridge(readWriteCloser io.ReadWriter, sockfd int, oneWay bool, read bool) {

	if !oneWay || read {
		buffer1 := make([]byte, BUF_SIZE)
		totalRead := 0
		totalSent := 0

		go func() {
			for {
				rlen, err := readWriteCloser.Read(buffer1)
				if validateErr(err, "Error reading from tun") {
					break
				}

				totalRead += rlen
				log.Debugf("Total read so far: %d\n", totalRead)

				wlen, writeErr := syscall.Write(sockfd, buffer1[:rlen])
				if validateErr(writeErr, "Error writing to zt") {
					break
				}

				totalSent += wlen
				log.Debugf("Total sent so far: %d\n", totalSent)
			}
		}()
	}

	if !oneWay || !read {
		buffer2 := make([]byte, BUF_SIZE)
		totalReceived := 0
		totalSaved := 0

		go func() {
			for {
				rlen, err := syscall.Read(sockfd, buffer2)
				if validateErr(err, "Error reading from zt") {
					break
				}

				totalReceived += rlen
				log.Debugf("Total received so far: %d\n", totalReceived)

				wlen, writeErr := readWriteCloser.Write(buffer2[:rlen])
				if validateErr(writeErr, "Error writing to tun") {
					break
				}

				totalSaved += wlen
				log.Debugf("Total saved so far: %d\n", totalSaved)

			}
		}()
	}
}

func main() {
	initZT()
	defer C.zts_stop()

	if len(getOtherIP()) == 0 {
		tunInterfaceStream := setupTun(true).ReadWriteCloser
		//tunInterfaceStream, _ := os.Create("/tmp/test")


		sockfd, finalSockfd := bindAndListen()

		defer syscall.Close(sockfd)
		defer syscall.Close(finalSockfd)

		bridge(tunInterfaceStream, finalSockfd, false, false)
		defer tunInterfaceStream.Close()
	} else {
		tunInterfaceStream := setupTun(false).ReadWriteCloser
		//tunInterfaceStream, _ := os.Open("/Users/selva/repos/libzt/examples/cpp/main")

		sockfd := connectToOther()
		defer syscall.Close(sockfd)

		bridge(tunInterfaceStream, sockfd, false, true)
		defer tunInterfaceStream.Close()
	}

	<-setupCleanUpOnInterrupt()
}
