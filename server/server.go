package server

import (
	"io"
	"github.com/selvakn/libzt"
	"net"
	"p2p-port-forward/constants"
	"p2p-port-forward/utils"
	"fmt"
	"p2p-port-forward/logger"
)

var log = logger.Logger

type Server struct {
	zt      *libzt.ZT
	port    string
	ipProto utils.IPProto
}

func New(zt *libzt.ZT, port string, proto utils.IPProto) Server {
	return Server{zt: zt, ipProto: proto, port: port}
}

func (s *Server) Listen() io.Closer {
	log.Info("Waiting for any client to connect")

	listener, _ := s.zt.Listen6(constants.INTERNAL_ZT_PORT)
	loggingListener := &utils.LoggingListener{Listener: listener}
	dataRageLogginglistener := &utils.DataRateLoggingListener{Listener: loggingListener}

	go utils.Sync(s.dialLocalService(), dataRageLogginglistener.Accept, true)
	return dataRageLogginglistener
}

func (s *Server) dialLocalService() func() (net.Conn, error) {
	return func() (net.Conn, error) {
		return net.Dial(s.ipProto.GetName(), fmt.Sprintf("localhost:%s", s.port))
	}
}
