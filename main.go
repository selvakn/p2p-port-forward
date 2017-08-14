package main

import (
	"fmt"
	"github.com/google/logger"
	"github.com/selvakn/libzt"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"net"
	"os"
	"p2p-port-forward/client"
	"p2p-port-forward/constants"
	"p2p-port-forward/utils"
)

func dialLocalService(networkProto utils.IPProto) func() (net.Conn, error) {
	return func() (net.Conn, error) {
		return net.Dial(networkProto.GetName(), fmt.Sprintf("localhost:%s", *forwardPort))
	}
}

var (
	network     = kingpin.Flag("network", "zerotier network id").Short('n').Default("8056c2e21c000001").String()
	forwardPort = kingpin.Flag("forward-port", "port to forward (in listen mode)").Short('f').Default("22").String()
	acceptPort  = kingpin.Flag("accept-port", "port to accept (in connect mode)").Short('a').Default("2222").String()
	useUDP      = kingpin.Flag("use-udp", "UDP instead of TCP (TCP default)").Short('u').Default("false").Bool()

	connectTo = kingpin.Flag("connect-to", "server (zerotier) ip to connect").Short('c').String()
)

func main() {
	logger.Init("p2p-port-forward", false, false, os.Stdout)

	kingpin.Version("1.0.1")
	kingpin.Parse()

	zt := libzt.Init(*network, "./zt")

	logger.Infof("ipv4 = %v \n", zt.GetIPv4Address().String())
	logger.Infof("ipv6 = %v \n", zt.GetIPv6Address().String())

	var closableConn io.Closer

	if len(*connectTo) == 0 {
		logger.Info("Waiting for any client to connect")

		listener, _ := zt.Listen6(constants.INTERNAL_ZT_PORT)
		loggingListener := &utils.LoggingListener{Listener: listener}
		dataRageLogginglistener := &utils.DataRateLoggingListener{Listener: loggingListener}

		go utils.Sync(dialLocalService(utils.GetIPProto(*useUDP)), dataRageLogginglistener.Accept)

		closableConn = listener
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
