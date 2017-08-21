package utils

import (
	"net"
	"time"
)

type DuplexUDPConnection struct {
	UDPConn  *net.UDPConn
	rUDPAddr *net.UDPAddr
}

func (c *DuplexUDPConnection) Read(b []byte) (n int, err error) {
	len, rUDPAddr, err := c.UDPConn.ReadFromUDP(b)
	c.rUDPAddr = rUDPAddr
	return len, err
}

func (c *DuplexUDPConnection) Write(b []byte) (n int, err error) {
	return c.UDPConn.WriteToUDP(b, c.rUDPAddr)
}

func (c *DuplexUDPConnection) Close() error                       { return c.UDPConn.Close() }
func (c *DuplexUDPConnection) LocalAddr() net.Addr                { return c.UDPConn.LocalAddr() }
func (c *DuplexUDPConnection) RemoteAddr() net.Addr               { return c.UDPConn.RemoteAddr() }
func (c *DuplexUDPConnection) SetDeadline(t time.Time) error      { return c.UDPConn.SetDeadline(t) }
func (c *DuplexUDPConnection) SetReadDeadline(t time.Time) error  { return c.UDPConn.SetReadDeadline(t) }
func (c *DuplexUDPConnection) SetWriteDeadline(t time.Time) error { return c.UDPConn.SetReadDeadline(t) }
