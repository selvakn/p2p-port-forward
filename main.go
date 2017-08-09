package main

import (
	"fmt"
	"github.com/google/logger"
	"github.com/selvakn/libzt"
	"gopkg.in/alecthomas/kingpin.v2"
	"net"
	"os"
	"p2p-port-forward/utils"
)

const INTERNAL_ZT_PORT = 7878

func dialLocalService() (net.Conn, error) {
	return net.Dial("tcp", fmt.Sprintf("localhost:%s", *forwardPort))
}

func dialRemoteThroughTunnel(zt *libzt.ZT) func() (net.Conn, error) {
	return func() (net.Conn, error) {
		conn, err := zt.Connect6(*connectTo, INTERNAL_ZT_PORT)
		conn = (&utils.DataRateLoggingConnection{}).Init(conn)
		return conn, err
	}
}

var (
	network     = kingpin.Flag("network", "zerotier network id").Short('n').Default("8056c2e21c000001").String()
	forwardPort = kingpin.Flag("forward-port", "port to forward (in listen mode)").Short('f').Default("22").String()
	acceptPort  = kingpin.Flag("accept-port", "port to accept (in connect mode)").Short('a').Default("2222").String()

	connectTo = kingpin.Flag("connect-to", "server (zerotier) ip to connect").Short('c').String()
)

func main() {
	logger.Init("p2p-port-forward", false, false, os.Stdout)

	kingpin.Parse()

	zt := libzt.Init(*network, "./zt")

	logger.Infof("ipv4 = %v \n", zt.GetIPv4Address().String())
	logger.Infof("ipv6 = %v \n", zt.GetIPv6Address().String())

	if len(*connectTo) == 0 {
		logger.Info("Waiting for any client to connect")

		listener, _ := zt.Listen6(INTERNAL_ZT_PORT)
		loggingListener := &utils.LoggingListener{Listener: listener}
		dataRageLogginglistener := &utils.DataRateLoggingListener{Listener: loggingListener}

		go utils.Sync(dialLocalService, dataRageLogginglistener.Accept)

		<-utils.SetupCleanUpOnInterrupt(func() {
			listener.Close()
		})

	} else {
		ln, _ := net.Listen("tcp", fmt.Sprintf(":%s", *acceptPort))

		go utils.Sync(ln.Accept, dialRemoteThroughTunnel(zt))

		<-utils.SetupCleanUpOnInterrupt(func() {
			ln.Close()
		})

	}
}
