package main

import (
	"net"
	"syscall"
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
		sockfd := listener.BindAndListen(PORT)

		go func() {
			for {
				newSockfd := listener.Accept(sockfd)
				go listener.HandleIncoming(newSockfd)
			}
		}()

		<-utils.SetupCleanUpOnInterrupt(func() {
			syscall.Close((int)(sockfd))
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
