package utils

import (
	"net"
	"time"
	"github.com/google/logger"
)

type LoggingConnection struct {
	conn net.Conn
}

func (c *LoggingConnection) Read(b []byte) (n int, err error) {
	return c.conn.Read(b)
}

func (c *LoggingConnection) Write(b []byte) (n int, err error) {
	return c.conn.Write(b)
}

func (c *LoggingConnection) Close() error                       { return c.conn.Close() }
func (c *LoggingConnection) LocalAddr() net.Addr                { return c.conn.LocalAddr() }
func (c *LoggingConnection) RemoteAddr() net.Addr               { return c.conn.RemoteAddr() }
func (c *LoggingConnection) SetDeadline(t time.Time) error      { return c.conn.SetDeadline(t) }
func (c *LoggingConnection) SetReadDeadline(t time.Time) error  { return c.conn.SetReadDeadline(t) }
func (c *LoggingConnection) SetWriteDeadline(t time.Time) error { return c.conn.SetReadDeadline(t) }

type LoggingListener struct {
	Listener net.Listener
}

func (l *LoggingListener) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		logger.Errorf("Error while accepting %v\n", err)
	} else {
		logger.Infof("Accepted incoming connection from %v", conn.RemoteAddr())
	}

	return &LoggingConnection{conn: conn}, err
}

func (l *LoggingListener) Close() error   { return l.Listener.Close() }
func (l *LoggingListener) Addr() net.Addr { return l.Listener.Addr() }
