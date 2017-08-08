package utils

import (
	"io"
	"net"
	"github.com/google/logger"
)

func Sync(stream1 func() (net.Conn, error), stream2 func() (net.Conn, error)) {
	for {
		conn1, err := stream1()
		if err != nil {
			logger.Error(err)
			return
		}
		conn2, err := stream2()

		if err != nil {
			logger.Error(err)
			return
		}

		go sync(conn1, conn2)
	}
}

func sync(source1 io.ReadWriteCloser, source2 io.ReadWriteCloser) {
	// FixMe: Cauing seg fault with the zt connections on close ¯\_(ツ)_/¯

	//defer source1.Close()
	//defer source2.Close()

	go io.Copy(source1, source2)
	io.Copy(source2, source1)
}
