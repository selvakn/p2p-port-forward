package forwarder

import (
	"os"
	"net"
	"../libzt"
	"io"
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

	ztConn, _ := libzt.Connect6(GetOtherIP(), port)
	defer conn.Close()

	go io.Copy(conn, ztConn)
	io.Copy(ztConn, conn)
}

