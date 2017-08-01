package main

import (
	"net"
	"./forwarder"
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
				ztConn, err := ztListener.Accept()
				if err != nil {
					log.Error(err)
					return
				}
				conn, _ := net.Dial("tcp", "localhost:22")
				if err != nil {
					log.Error(err)
					return
				}

				go utils.Sync(ztConn, conn)
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
				if err != nil {
					log.Error(err)
					return
				}
				ztConn, err := libzt.Connect6(forwarder.GetOtherIP(), PORT)
				if err != nil {
					log.Error(err)
					return
				}

				go utils.Sync(conn, ztConn)
			}
		}()

		<-utils.SetupCleanUpOnInterrupt(func() {
			ln.Close()
		})

	}

}
