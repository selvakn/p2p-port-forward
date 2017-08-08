package libzt

import (
	"net"
	"errors"
	"syscall"
	"fmt"
)

type TCP6Listener struct {
	fd        int
	localIP   net.IP
	localPort uint16
}

func (l *TCP6Listener) Accept() (net.Conn, error) {
	acceptedFd, sockAddr := accept6(l.fd)
	if acceptedFd < 0 {
		return nil, errors.New("Unable to accept new connection")
	}

	conn := &Connection{
		fd:         acceptedFd,

		localIP:    l.localIP,
		localPort:  l.localPort,

		remoteIp:   net.IP(sockAddr.Addr[:]),
		remotePort: sockAddr.Port,
	}
	return conn, nil
}

func (l *TCP6Listener) Close() error {
	return syscall.Close(l.fd)
}

func (l *TCP6Listener) Addr() net.Addr {
	addr, _ := net.ResolveTCPAddr("tcp6", fmt.Sprintf("[%s]:%d", l.localIP.String(), l.localPort))
	return addr
}
