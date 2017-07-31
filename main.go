package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"unsafe"

	"github.com/songgao/water"
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
const BUF_SIZE = 2000

func setupCleanUpOnInterrupt() chan bool {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	cleanupDone := make(chan bool)

	go func() {
		for range signalChan {
			fmt.Println("\nReceived an interrupt, shutting dow.\n")

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
		fmt.Println(message)
		os.Exit(1)
	}
}

func bindAndListen(sockfd C.int) int {
	serverSocket := syscall.RawSockaddrInet6{Flowinfo: 0, Family: syscall.AF_INET6, Port: PORT}
	retVal := C.zts_bind(sockfd, (*C.struct_sockaddr)(unsafe.Pointer(&serverSocket)), C.sizeof_struct_sockaddr_in6)
	validate(retVal, "ERROR on binding")
	fmt.Println("Bind Complete")

	C.zts_listen(sockfd, 1)
	fmt.Println("Listening")

	clientSocket := syscall.RawSockaddrInet6{}
	clientSocketLength := C.sizeof_struct_sockaddr_in6
	newSockfd := C.zts_accept(sockfd, (*C.struct_sockaddr)(unsafe.Pointer(&clientSocket)), (*C.socklen_t)(unsafe.Pointer(&clientSocketLength)))
	validate(newSockfd, "ERROR on accept")
	fmt.Println("Accepted")

	clientIpAddress := make([]byte, C.ZT_MAX_IPADDR_LEN)
	C.inet_ntop(syscall.AF_INET6, unsafe.Pointer(&clientSocket.Addr), (*C.char)(unsafe.Pointer(&clientIpAddress[0])), C.ZT_MAX_IPADDR_LEN)
	fmt.Printf("Incoming connection from client having IPv6 address: %s\n", string(clientIpAddress[:C.ZT_MAX_IPADDR_LEN]))

	return int(newSockfd)
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
		log.Fatalln("Error running command:", err)
	}
}

func setupTun(initater bool) *water.Interface {
	iface, _ := water.New(water.Config{
		DeviceType: water.TUN,
	})

	log.Printf("Interface Name: %s\n", iface.Name())

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
	log.Printf("ipv4 = %s \n", string(ipv4Address[:C.ZT_MAX_IPADDR_LEN]))

	C.zts_get_ipv6_address(C.CString(NETWORK_ID), (*C.char)(unsafe.Pointer(&ipv6Address[0])), C.ZT_MAX_IPADDR_LEN)
	log.Printf("ipv6 = %s \n", string(ipv6Address[:C.ZT_MAX_IPADDR_LEN]))
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

func validateErr(err error, message string) {
	if err != nil {
		log.Println(message)
	}
}

func bridge(readWriteCloser io.ReadWriteCloser, sockfd int) {
	buffer1 := make([]byte, BUF_SIZE)
	go func() {
		for {
			plen, err := readWriteCloser.Read(buffer1)
			validateErr(err, "Error reading from tun")

			_, writeErr := syscall.Write(sockfd, buffer1[:plen])
			validateErr(writeErr, "Error writing to zt")
		}
	}()

	buffer2 := make([]byte, BUF_SIZE)

	go func() {
		for {
			plen, err := syscall.Read(sockfd, buffer2)
			validateErr(err, "Error reading from zt")

			_, writeErr := readWriteCloser.Write(buffer2[:plen])
			validateErr(writeErr, "Error writing to tun")

		}
	}()
}

func main() {
	initZT()
	defer C.zts_stop()

	sockfd := C.zts_socket(syscall.AF_INET6, syscall.SOCK_STREAM, 0)
	validate(sockfd, "Error in opening socket")
	defer C.zts_close(sockfd)

	var finalSockfd int
	var iface *water.Interface

	if len(getOtherIP()) == 0 {
		iface = setupTun(true)
		finalSockfd = bindAndListen(sockfd)
	} else {
		iface = setupTun(false)
		finalSockfd = connectToOther()
	}

	bridge(iface.ReadWriteCloser, finalSockfd)
	defer iface.ReadWriteCloser.Close()


	<-setupCleanUpOnInterrupt()
}
