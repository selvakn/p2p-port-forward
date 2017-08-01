package listener

import (
	"net"
	"io"
)

func HandleIncoming(ztConn net.Conn) {
	conn, _ := net.Dial("tcp", "localhost:22")
	defer conn.Close()

	go io.Copy(conn, ztConn)
	io.Copy(ztConn, conn)
}
