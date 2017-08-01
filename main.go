package main

import (
	"net"
	"./forwarder"
	"./listener"
	"./utils"
	"./libzt"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("util")

const PORT = 50718 // 7878
const NETWORK_ID = "8056c2e21c000001"

func main() {
	libzt.SimpleStart("./zt", NETWORK_ID)

	log.Infof("ipv4 = %s \n", libzt.GetIpv4Address(NETWORK_ID))
	log.Infof("ipv6 = %s \n", libzt.GetIpv6Address(NETWORK_ID))

	if len(forwarder.GetOtherIP()) == 0 {
		ztListener, _ := libzt.Listen6(PORT)

		go func() {
			for {
				conn, _ := ztListener.Accept()
				go listener.HandleIncoming(conn)
			}
		}()

		<-utils.SetupCleanUpOnInterrupt(func() {
			ztListener.Close()
		})

	} else {
		ln, _ := net.Listen("tcp", ":2222")

		go func() {
			for {
				conn, err := ln.Accept()
				if err == nil {
					go forwarder.HandleOutgoing(conn, PORT)
				}
			}
		}()

		<-utils.SetupCleanUpOnInterrupt(func() {
			ln.Close()
		})

	}

}
