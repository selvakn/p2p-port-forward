package forwarder

import (
	"os"
	"net"
	"syscall"
	"../utils"
	"../libzt"
)


func GetOtherIP() string {
	if len(os.Args) >= 2 {
		return os.Args[1]
	} else {
		return ""
	}
}

func HandleOutgoing(conn net.Conn, port uint16) {
	defer conn.Close()

	sockfd := connectToOther(port)
	//defer libzt.Close(sockfd)

	go utils.ReceiveFrom(conn, sockfd)
	utils.SendTo(sockfd, conn)
}

func connectToOther(port uint16) int {
	arr := parseIPV6(GetOtherIP())

	clientSocket := syscall.RawSockaddrInet6{Flowinfo: 0, Family: syscall.AF_INET6, Port: port, Addr: arr}

	sockfd := libzt.Socket(syscall.AF_INET6, syscall.SOCK_STREAM, 0)
	utils.Validate(sockfd, "Error in opening socket")

	retVal := libzt.Connect6(sockfd, clientSocket)
	utils.Validate(retVal, "Error in connect client")

	return sockfd
}

func parseIPV6(ipString string) [16]byte {
	ip := net.ParseIP(ipString)
	var arr [16]byte
	copy(arr[:], ip)
	return arr
}

