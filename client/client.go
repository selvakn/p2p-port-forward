package client

import (
	"fmt"
	"github.com/selvakn/libzt"
	"io"
	"net"
	"p2p-port-forward/constants"
	"p2p-port-forward/logger"
	"p2p-port-forward/utils"
)

var log = logger.Logger

type Client struct {
	zt           *libzt.ZT
	networkProto utils.IPProto
	connectTo    string
	port         string
}

func New(zt *libzt.ZT, ipToConnectTo string, port string, networkProto utils.IPProto) Client {
	return Client{zt: zt, networkProto: networkProto, port: port, connectTo: ipToConnectTo}
}

func (c *Client) ListenAndSync() io.Closer {
	if c.networkProto == utils.UDP {
		return c.listenAndSyncUDP()
	} else {
		return c.listenAndSyncTCP()
	}
}

func (c *Client) listenAndSyncUDP() io.Closer {
	go utils.Sync(c.listenUDP(), c.dialRemoteThroughTunnel(), false)
	return nil
}

func (c *Client) listenAndSyncTCP() io.Closer {
	ln, _ := net.Listen(c.networkProto.GetName(), fmt.Sprintf(":%s", c.port))
	go utils.Sync(ln.Accept, c.dialRemoteThroughTunnel(), true)
	return ln
}

func (c *Client) dialRemoteThroughTunnel() func() (net.Conn, error) {
	return func() (net.Conn, error) {
		log.Infof("Attempting a remote connection")
		conn, err := c.zt.Connect6(c.connectTo, constants.INTERNAL_ZT_PORT)
		conn = (&utils.DataRateLoggingConnection{}).Init(conn)
		return conn, err
	}
}

func (c *Client) listenUDP() func() (net.Conn, error) {
	fmt.Println("return Listening....")
	return func() (net.Conn, error) {
		addr, _ := net.ResolveUDPAddr(c.networkProto.GetName(), fmt.Sprintf(":%s", c.port))
		fmt.Println("Listening....")
		conn, err := net.ListenUDP(c.networkProto.GetName(), addr)
		return &utils.DuplexUDPConnection{UDPConn: conn}, err

	}
}
