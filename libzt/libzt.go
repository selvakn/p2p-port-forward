package libzt

/*
#cgo CFLAGS: -I ./include
#cgo darwin LDFLAGS: -L ${SRCDIR}/darwin/ -lzt -lstdc++
#cgo linux LDFLAGS: -L ${SRCDIR}/linux/ -lzt -lstdc++

#include "libzt.h"
#include <netdb.h>
*/
import "C"
import (
	"unsafe"
	"syscall"
	"net"
	"errors"
)

const ZT_MAX_IPADDR_LEN = C.ZT_MAX_IPADDR_LEN

func SimpleStart(homePath, networkId string) {
	C.zts_simple_start(C.CString(homePath), C.CString(networkId))
}

func GetIpv4Address(networkId string) string {
	address := make([]byte, ZT_MAX_IPADDR_LEN)
	C.zts_get_ipv4_address(C.CString(networkId), (*C.char)(unsafe.Pointer(&address[0])), C.ZT_MAX_IPADDR_LEN)
	return string(address)
}

func GetIpv6Address(networkId string) string {
	address := make([]byte, ZT_MAX_IPADDR_LEN)
	C.zts_get_ipv6_address(C.CString(networkId), (*C.char)(unsafe.Pointer(&address[0])), C.ZT_MAX_IPADDR_LEN)
	return string(address)
}

// TODO: Return err as second value

func Close(fd int) int {
	return (int)(C.zts_close(cint(fd)))
}

func Listen6(port uint16) (net.Listener, error) {
	fd := Socket(syscall.AF_INET6, syscall.SOCK_STREAM, 0)
	if fd < 0 {
		return nil, errors.New("Error in opening socket")
	}

	serverSocket := syscall.RawSockaddrInet6{Flowinfo: 0, Family: syscall.AF_INET6, Port: port}
	retVal := Bind6(fd, serverSocket)
	if retVal < 0 {
		return nil, errors.New("ERROR on binding")
	}

	retVal = Listen(fd, 1)
	if retVal < 0 {
		return nil, errors.New("ERROR listening")
	}

	return &TCP6Listener{fd: fd}, nil
}

func Socket(family int, socketType int, protocol int) int {
	return (int)(C.zts_socket(cint(family), cint(socketType), cint(protocol)))
}

func Connect6(fd int, socketAddr syscall.RawSockaddrInet6) int {
	return (int)(C.zts_connect(cint(fd), (*C.struct_sockaddr)(unsafe.Pointer(&socketAddr)), syscall.SizeofSockaddrInet6))
}

func Bind6(fd int, sockerAddr syscall.RawSockaddrInet6) int {
	return (int)(C.zts_bind(cint(fd), (*C.struct_sockaddr)(unsafe.Pointer(&sockerAddr)), syscall.SizeofSockaddrInet6))
}

func Accept6(fd int) (int, syscall.RawSockaddrInet6) {
	socketAddr := syscall.RawSockaddrInet6{}
	socketLength := syscall.SizeofSockaddrInet6
	return (int)(C.zts_accept(cint(fd), (*C.struct_sockaddr)(unsafe.Pointer(&socketAddr)), (*C.socklen_t)(unsafe.Pointer(&socketLength)))), socketAddr
}

func Listen(fd int, backlog int) int {
	return (int)(C.zts_listen(cint(fd), cint(backlog)))
}

func cint(value int) C.int {
	return (C.int)(value)
}
