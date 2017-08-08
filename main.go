package main

import (
	"net"
	"p2p-port-forward/forwarder"
	"p2p-port-forward/utils"
	"github.com/selvakn/libzt"
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
		conn, err := zt.Connect6(forwarder.GetOtherIP(), PORT)
		conn = (&utils.DataRateLoggingConnection{}).Init(conn)
		return conn, err
	}
}

func main() {
	logger.Init("p2p-port-forward", false, false, os.Stdout)

	zt := libzt.Init(NETWORK_ID, "./zt")

	logger.Infof("ipv4 = %v \n", zt.GetIPv4Address().String())
	logger.Infof("ipv6 = %v \n", zt.GetIPv6Address().String())

	if len(forwarder.GetOtherIP()) == 0 {
		listener, _ := zt.Listen6(PORT)
		loggingListener := &utils.LoggingListener{Listener: listener}
		dataRageLogginglistener := &utils.DataRateLoggingListener{Listener: loggingListener}

		go utils.Sync(dialLocalService, dataRageLogginglistener.Accept)

		<-utils.SetupCleanUpOnInterrupt(func() {
			listener.Close()
		})

	} else {
		ln, _ := net.Listen("tcp", fmt.Sprintf(":%s", forwarder.LocalPortToForward()))

		go utils.Sync(ln.Accept, dialRemoteThroughTunnel(zt))

		<-utils.SetupCleanUpOnInterrupt(func() {
			ln.Close()
		})

	}
}
