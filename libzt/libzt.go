package libzt

/*
#cgo CFLAGS: -I ./include
#cgo darwin LDFLAGS: -L ${SRCDIR}/darwin/ -lzt -lstdc++
#cgo linux LDFLAGS: -L ${SRCDIR}/linux/ -lzt -lstdc++

#include "libzt.h"
#include <netdb.h>
*/
import "C"
import "unsafe"

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