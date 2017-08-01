package listener

import (
	"syscall"
	"unsafe"
	"github.com/op/go-logging"
	"./../utils"
	"net"
)

/*
#cgo CFLAGS: -I ./../libzt/include
#cgo darwin LDFLAGS: -L ${SRCDIR}/../libzt/darwin/ -lzt -lstdc++
#cgo linux LDFLAGS: -L ${SRCDIR}/../libzt/linux/ -lzt -lstdc++

#include "libzt.h"
#include <netdb.h>
*/
import "C"


var log = logging.MustGetLogger("tunneler")

func BindAndListen(port uint16) (int) {
	sockfd := C.zts_socket(syscall.AF_INET6, syscall.SOCK_STREAM, 0)
	utils.Validate((int)(sockfd), "Error in opening socket")

	serverSocket := syscall.RawSockaddrInet6{Flowinfo: 0, Family: syscall.AF_INET6, Port: port}
	retVal := C.zts_bind(sockfd, (*C.struct_sockaddr)(unsafe.Pointer(&serverSocket)), C.sizeof_struct_sockaddr_in6)
	utils.Validate((int)(retVal), "ERROR on binding")
	log.Debugf("Bind Complete")

	C.zts_listen(sockfd, 1)
	log.Debugf("Listening")

	return (int)(sockfd)
}

func Accept(sockfd int) int {
	clientSocket := syscall.RawSockaddrInet6{}
	clientSocketLength := C.sizeof_struct_sockaddr_in6
	connSockfd := C.zts_accept((C.int)(sockfd), (*C.struct_sockaddr)(unsafe.Pointer(&clientSocket)), (*C.socklen_t)(unsafe.Pointer(&clientSocketLength)))

	utils.Validate((int)(connSockfd), "ERROR on accept")

	log.Info("Accepted incoming connection from client")

	return (int)(connSockfd)
}

func HandleIncoming(sockfd int) {
	//defer C.zts_close(sockfd)

	conn, _ := net.Dial("tcp", "localhost:22")
	defer conn.Close()

	go utils.ReceiveFrom(conn, sockfd)
	utils.SendTo(sockfd, conn)
}

