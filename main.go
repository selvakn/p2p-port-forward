package main

import (
	"net"
	"p2p-port-forward/forwarder"
	"p2p-port-forward/utils"
	"p2p-port-forward/libzt"
	"fmt"
	"p2p-port-forward/listener"
	"github.com/google/logger"
	"os"
)

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
	logger.Init("p2p-port-forward", false, false, os.Stdout)

	zt := libzt.Init(NETWORK_ID, "./zt")

	logger.Infof("ipv4 = %v \n", zt.GetIPv4Address().String())
	logger.Infof("ipv6 = %v \n", zt.GetIPv6Address().String())

	if len(forwarder.GetOtherIP()) == 0 {
		ztListener, _ := zt.Listen6(PORT)
		loggingListener := utils.LoggingListener{Listener: ztListener}

		go utils.Sync(dialLocalService, loggingListener.Accept)

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
