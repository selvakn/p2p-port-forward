package libzt

import (
	"net"
	"time"
	"syscall"
)

type Connection struct {
	fd         int
}

func (c *Connection) Read(b []byte) (n int, err error) {
	return syscall.Read(c.fd, b)
}

func (c *Connection) Write(b []byte) (n int, err error) {
	return syscall.Write(c.fd, b)
}

func (c *Connection) Close() error {
	return syscall.Close(c.fd)
}


func (c *Connection) LocalAddr() net.Addr {
	return nil // TODO: Implement
}

func (c *Connection) RemoteAddr() net.Addr {
	return nil // TODO: Implement

}
func (c *Connection) SetDeadline(time.Time) error      { return nil }
func (c *Connection) SetReadDeadline(time.Time) error  { return nil }
func (c *Connection) SetWriteDeadline(time.Time) error { return nil }
