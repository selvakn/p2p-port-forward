package main

import (
	"net"
	"p2p-port-forward/forwarder"
	"p2p-port-forward/utils"
	"p2p-port-forward/libzt"
	"github.com/op/go-logging"
	"fmt"
	"p2p-port-forward/listener"
)

var log = logging.MustGetLogger("util")

const PORT = 7878
const NETWORK_ID = "8056c2e21c000001"

func dialLocalService() (net.Conn, error) {
	return net.Dial("tcp", fmt.Sprintf("localhost:%s", listener.LocalPortToListen()))
}

func dialRemoteThroughTunnel(zt *libzt.ZT) func() (net.Conn, error) {
	return func() (net.Conn, error) {
		return zt.Connect6(forwarder.GetOtherIP(), PORT)
	}
}

func main() {
	zt := libzt.Init(NETWORK_ID, "./zt")

	log.Infof("ipv4 = %s \n", zt.GetIPv4Address())
	log.Infof("ipv6 = %s \n", zt.GetIPv6Address())

	if len(forwarder.GetOtherIP()) == 0 {
		ztListener, _ := zt.Listen6(PORT)
		go utils.Sync(dialLocalService, ztListener.Accept)

		<-utils.SetupCleanUpOnInterrupt(func() {
			ztListener.Close()
		})

	} else {
		ln, _ := net.Listen("tcp", fmt.Sprintf(":%s", forwarder.LocalPortToForward()))

		go utils.Sync(ln.Accept, dialRemoteThroughTunnel(zt))

		<-utils.SetupCleanUpOnInterrupt(func() {
			ln.Close()
		})

	}
}
