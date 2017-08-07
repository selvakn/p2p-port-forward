package libzt

/*
#cgo CFLAGS: -I ./include
#cgo darwin LDFLAGS: -L ${SRCDIR}/darwin/ -lzt -lstdc++ -lm -std=c++11
#cgo linux LDFLAGS: -L ${SRCDIR}/linux/ -lzt -lstdc++ -lm -std=c++11

#include "libzt.h"
#include <netdb.h>
*/
import "C"
import (
	"unsafe"
	"syscall"
	"net"
	"errors"
	"encoding/binary"
)

const ZT_MAX_IPADDR_LEN = C.ZT_MAX_IPADDR_LEN

type ZT struct {
	id       string
	homePath string
}

func Init(id string, homePath string) (*ZT) {
	zt := &ZT{id: id, homePath: homePath }
	C.zts_simple_start(C.CString(homePath), C.CString(id))
	return zt
}

func (zt *ZT) GetIPv4Address() string {
	address := make([]byte, ZT_MAX_IPADDR_LEN)
	C.zts_get_ipv4_address(C.CString(zt.id), (*C.char)(unsafe.Pointer(&address[0])), C.ZT_MAX_IPADDR_LEN)
	return string(address)
}

func (zt *ZT) GetIPv6Address() string {
	address := make([]byte, ZT_MAX_IPADDR_LEN)
	C.zts_get_ipv6_address(C.CString(zt.id), (*C.char)(unsafe.Pointer(&address[0])), C.ZT_MAX_IPADDR_LEN)
	return string(address)
}

func (zt *ZT) Listen6(port uint16) (net.Listener, error) {
	fd := socket(syscall.AF_INET6, syscall.SOCK_STREAM, 0)
	if fd < 0 {
		return nil, errors.New("Error in opening socket")
	}

	serverSocket := syscall.RawSockaddrInet6{Flowinfo: 0, Family: syscall.AF_INET6, Port: htonl(port)}
	retVal := bind6(fd, serverSocket)
	if retVal < 0 {
		return nil, errors.New("ERROR on binding")
	}

	retVal = listen(fd, 1)
	if retVal < 0 {
		return nil, errors.New("ERROR listening")
	}

	return &TCP6Listener{fd: fd}, nil
}

func (zt *ZT) Connect6(ip string, port uint16) (net.Conn, error) {
	clientSocket := syscall.RawSockaddrInet6{Flowinfo: 0, Family: syscall.AF_INET6, Port: htonl(port), Addr: parseIPV6(ip)}

	fd := socket(syscall.AF_INET6, syscall.SOCK_STREAM, 0)
	if fd < 0 {
		return nil, errors.New("Error in opening socket")
	}

	retVal := (int)(C.zts_connect(cint(fd), (*C.struct_sockaddr)(unsafe.Pointer(&clientSocket)), syscall.SizeofSockaddrInet6))
	if retVal < 0 {
		return nil, errors.New("Unable to connect")
	}

	conn := &Connection{
		fd: fd,
	}
	return conn, nil
}

func htonl(number uint16) uint16 {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, number)
	return *(*uint16)(unsafe.Pointer(&bytes[0]))
}

func close(fd int) int {
	return (int)(C.zts_close(cint(fd)))
}

func socket(family int, socketType int, protocol int) int {
	return (int)(C.zts_socket(cint(family), cint(socketType), cint(protocol)))
}

func listen(fd int, backlog int) int {
	return (int)(C.zts_listen(cint(fd), cint(backlog)))
}

func bind6(fd int, sockerAddr syscall.RawSockaddrInet6) int {
	return (int)(C.zts_bind(cint(fd), (*C.struct_sockaddr)(unsafe.Pointer(&sockerAddr)), syscall.SizeofSockaddrInet6))
}

func accept6(fd int) (int, syscall.RawSockaddrInet6) {
	socketAddr := syscall.RawSockaddrInet6{}
	socketLength := syscall.SizeofSockaddrInet6
	return (int)(C.zts_accept(cint(fd), (*C.struct_sockaddr)(unsafe.Pointer(&socketAddr)), (*C.socklen_t)(unsafe.Pointer(&socketLength)))), socketAddr
}

func cint(value int) C.int {
	return (C.int)(value)
}

func parseIPV6(ipString string) [16]byte {
	ip := net.ParseIP(ipString)
	var arr [16]byte
	copy(arr[:], ip)
	return arr
}
