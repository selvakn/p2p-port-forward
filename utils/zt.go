package utils

import "unsafe"

/*
#cgo CFLAGS: -I ./../libzt/include
#cgo darwin LDFLAGS: -L ${SRCDIR}/../libzt/darwin/ -lzt -lstdc++
#cgo linux LDFLAGS: -L ${SRCDIR}/../libzt/linux/ -lzt -lstdc++

#include "libzt.h"
#include <netdb.h>
*/
import "C"

const NETWORK_ID = "8056c2e21c000001"

func InitZT() {
	C.zts_simple_start(C.CString("./zt"), C.CString(NETWORK_ID))

	ipv4Address := make([]byte, C.ZT_MAX_IPADDR_LEN)
	ipv6Address := make([]byte, C.ZT_MAX_IPADDR_LEN)

	C.zts_get_ipv4_address(C.CString(NETWORK_ID), (*C.char)(unsafe.Pointer(&ipv4Address[0])), C.ZT_MAX_IPADDR_LEN)
	log.Infof("ipv4 = %s \n", string(ipv4Address[:C.ZT_MAX_IPADDR_LEN]))

	C.zts_get_ipv6_address(C.CString(NETWORK_ID), (*C.char)(unsafe.Pointer(&ipv6Address[0])), C.ZT_MAX_IPADDR_LEN)
	log.Infof("ipv6 = %s \n", string(ipv6Address[:C.ZT_MAX_IPADDR_LEN]))
}
