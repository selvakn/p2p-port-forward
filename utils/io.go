package utils

import (
	"io"
	"net"
)

func Sync(stream1 func() (net.Conn, error), stream2 func() (net.Conn, error)) {
	for {
		conn1, err := stream1()
		if err != nil {
			log.Error(err)
			return
		}
		conn2, err := stream2()

		if err != nil {
			log.Error(err)
			return
		}

		go sync(conn1, conn2)
	}
}

func sync(source1 io.ReadWriteCloser, source2 io.ReadWriteCloser) {
	go func() {
		defer source2.Close()
		defer source1.Close()

		io.Copy(source2, source1)
	}()
	io.Copy(source1, source2)
}
