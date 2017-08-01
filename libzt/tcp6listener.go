package libzt

import (
	"net"
	"errors"
	"syscall"
)

type TCP6Listener struct {
	fd int
}

func (l *TCP6Listener) Accept() (net.Conn, error) {
	acceptedFd, _ := accept6(l.fd)
	if acceptedFd < 0 {
		return nil, errors.New("Unable to accept new connection")
	}

	conn := &Connection{
		fd: acceptedFd,
	}
	return conn, nil
}

func (l *TCP6Listener) Close() error {
	return syscall.Close(l.fd)
}

func (l *TCP6Listener) Addr() net.Addr {
	return nil //TODO: Implement
}


