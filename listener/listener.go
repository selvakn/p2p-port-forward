package listener

import (
	"syscall"
	"github.com/op/go-logging"
	"../utils"
	"../libzt"
	"net"
)

var log = logging.MustGetLogger("tunneler")

func BindAndListen(port uint16) (int) {
	sockfd := libzt.Socket(syscall.AF_INET6, syscall.SOCK_STREAM, 0)
	utils.Validate(sockfd, "Error in opening socket")

	serverSocket := syscall.RawSockaddrInet6{Flowinfo: 0, Family: syscall.AF_INET6, Port: port}
	retVal := libzt.Bind6(sockfd, serverSocket)
	utils.Validate(retVal, "ERROR on binding")
	log.Debugf("Bind Complete")

	libzt.Listen(sockfd, 1)
	log.Debugf("Listening")

	return sockfd
}

func Accept(sockfd int) int {
	connSockfd, _, _ := libzt.Accept6(sockfd)

	utils.Validate(connSockfd, "ERROR on accept")

	log.Info("Accepted incoming connection from client")

	return connSockfd
}

func HandleIncoming(sockfd int) {
	//defer libzt.Close(sockfd)

	conn, _ := net.Dial("tcp", "localhost:22")
	defer conn.Close()

	go utils.ReceiveFrom(conn, sockfd)
	utils.SendTo(sockfd, conn)
}
