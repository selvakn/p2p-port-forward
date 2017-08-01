package forwarder

import (
	"os"
	"net"
	"syscall"
	"unsafe"
	"../utils"
)

/*
#cgo CFLAGS: -I ./../libzt/include
#cgo darwin LDFLAGS: -L ${SRCDIR}/../libzt/darwin/ -lzt -lstdc++
#cgo linux LDFLAGS: -L ${SRCDIR}/../libzt/linux/ -lzt -lstdc++

#include "libzt.h"
#include <netdb.h>
*/
import "C"


func GetOtherIP() string {
	if len(os.Args) >= 2 {
		return os.Args[1]
	} else {
		return ""
	}
}

func HandleOutgoing(conn net.Conn, port uint16) {
	defer conn.Close()

	sockfd := connectToOther(port)
	//defer C.zts_close(sockfd)

	go utils.ReceiveFrom(conn, sockfd)
	utils.SendTo(sockfd, conn)
}

func connectToOther(port uint16) int {
	arr := parseIPV6(GetOtherIP())

	clientSocket := syscall.RawSockaddrInet6{Flowinfo: 0, Family: syscall.AF_INET6, Port: port, Addr: arr}

	sockfd := C.zts_socket(syscall.AF_INET6, syscall.SOCK_STREAM, 0)
	utils.Validate((int)(sockfd), "Error in opening socket")

	retVal := C.zts_connect(sockfd, (*C.struct_sockaddr)(unsafe.Pointer(&clientSocket)), C.sizeof_struct_sockaddr_in6)
	utils.Validate((int)(retVal), "Error in connect client")

	return (int)(sockfd)
}

func parseIPV6(ipString string) [16]byte {
	ip := net.ParseIP(ipString)
	var arr [16]byte
	copy(arr[:], ip)
	return arr
}

