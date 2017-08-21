package utils

import (
	"io"
	"net"
)

func Sync(stream1 func() (net.Conn, error), stream2 func() (net.Conn, error), listenForNextConnection bool) {
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

		if listenForNextConnection {
			go sync(conn1, conn2)
		} else {
			sync(conn1, conn2)
		}
	}
}

func sync(source1 io.ReadWriteCloser, source2 io.ReadWriteCloser) {
	go func() {
		defer closeAll(source2, source1)

		_, err := io.Copy(source2, source1)
		if err != nil {
			log.Error(err)
		}
	}()
	_, err := io.Copy(source1, source2)
	if err != nil {
		log.Error(err)
	}
}

func closeAll(sources ...io.Closer) {
	log.Info("Closing all connections")
	for _, source := range sources {
		source.Close()
	}
}
