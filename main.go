package main

import (
	"github.com/selvakn/libzt"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"p2p-port-forward/client"
	"p2p-port-forward/logger"
	"p2p-port-forward/server"
	"p2p-port-forward/utils"
)

var (
	network     = kingpin.Flag("network", "zerotier network id").Short('n').Default("8056c2e21c000001").String()
	forwardPort = kingpin.Flag("forward-port", "port to forward (in listen mode)").Short('f').Default("22").String()
	acceptPort  = kingpin.Flag("accept-port", "port to accept (in connect mode)").Short('a').Default("2222").String()
	useUDP      = kingpin.Flag("use-udp", "UDP instead of TCP (TCP default)").Short('u').Default("false").Bool()

	connectTo = kingpin.Flag("connect-to", "server (zerotier) ip to connect").Short('c').String()
)

var log = logger.Logger

func main() {
	kingpin.Version("1.0.1")
	kingpin.Parse()

	zt := libzt.Init(*network, "./zt")

	log.Infof("ipv4 = %v ", zt.GetIPv4Address().String())
	log.Infof("ipv6 = %v ", zt.GetIPv6Address().String())

	var closableConn io.Closer

	if len(*connectTo) == 0 {
		forwarderServer := server.New(zt, *forwardPort, utils.GetIPProto(*useUDP))
		closableConn = forwarderServer.Listen()
	} else {
		forwarderClient := client.New(zt, *connectTo, *acceptPort, utils.GetIPProto(*useUDP))
		closableConn = forwarderClient.ListenAndSync()
	}

	<-utils.SetupCleanUpOnInterrupt(func() {
		if closableConn != nil {
			closableConn.Close()
		}
	})

}
