package main

import (
	"net"
	"p2p-port-forward/utils"
	"github.com/selvakn/libzt"
	"fmt"
	"github.com/google/logger"
	"os"
	"gopkg.in/alecthomas/kingpin.v2"
)

const INTERNAL_ZT_PORT = 7878

func dialLocalService() (net.Conn, error) {
	return net.Dial("tcp", fmt.Sprintf("localhost:%s", *forwardPort))
}

func dialRemoteThroughTunnel(zt *libzt.ZT) func() (net.Conn, error) {
	return func() (net.Conn, error) {
		conn, err := zt.Connect6(*serverIp, INTERNAL_ZT_PORT)
		conn = (&utils.DataRateLoggingConnection{}).Init(conn)
		return conn, err
	}
}

var (
	network     = kingpin.Flag("network", "zerotier network id").Short('n').Default("8056c2e21c000001").String()
	mode        = kingpin.Flag("listen-mode", "listen mode").Short('l').Default("false").Bool()
	forwardPort = kingpin.Flag("forward-port", "port to forward (in listen mode)").Short('f').Default("22").String()
	acceptPort  = kingpin.Flag("accept-port", "port to accept (in connect mode)").Short('a').Default("2222").String()

	serverIp = kingpin.Flag("server", "server (zerotier) ip").Short('s').String()
)

func main() {
	logger.Init("p2p-port-forward", false, false, os.Stdout)

	kingpin.Parse()

	zt := libzt.Init(*network, "./zt")

	logger.Infof("ipv4 = %v \n", zt.GetIPv4Address().String())
	logger.Infof("ipv6 = %v \n", zt.GetIPv6Address().String())

	if *mode {
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
