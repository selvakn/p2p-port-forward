package utils

import (
	"fmt"
	"github.com/c2h5oh/datasize"
	"github.com/gosuri/uilive"
	"net"
	"time"
)

type TransferRate struct {
	up   datasize.ByteSize
	down datasize.ByteSize
}

type DataRateLoggingConnection struct {
	conn            net.Conn
	writer          *uilive.Writer
	upRateCounter   DataRateCounter
	downRateCounter DataRateCounter
}

func (c *DataRateLoggingConnection) Init(conn net.Conn) *DataRateLoggingConnection {
	c.conn = conn
	c.writer = uilive.New()
	c.writer.Start()
	c.upRateCounter = NewRateCounter()
	c.downRateCounter = NewRateCounter()

	go c.updateStats()

	return c
}

func (c *DataRateLoggingConnection) Read(b []byte) (n int, err error) {
	len, err := c.conn.Read(b)
	c.downRateCounter.CaptureEvent(len)

	return len, err
}
func (c *DataRateLoggingConnection) Write(b []byte) (n int, err error) {
	len, err := c.conn.Write(b)
	c.upRateCounter.CaptureEvent(len)

	return len, err
}

func (c *DataRateLoggingConnection) Close() error {
	c.writer.Stop()
	return c.conn.Close()
}

func (c *DataRateLoggingConnection) LocalAddr() net.Addr           { return c.conn.LocalAddr() }
func (c *DataRateLoggingConnection) RemoteAddr() net.Addr          { return c.conn.RemoteAddr() }
func (c *DataRateLoggingConnection) SetDeadline(t time.Time) error { return c.conn.SetDeadline(t) }
func (c *DataRateLoggingConnection) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}
func (c *DataRateLoggingConnection) SetWriteDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *DataRateLoggingConnection) getTransferRate() TransferRate {
	return TransferRate{up: c.upRateCounter.GetDataRate(), down: c.downRateCounter.GetDataRate()}
}

func (c *DataRateLoggingConnection) updateStats() {
	throttle := time.Tick(time.Second)
	for {
		<-throttle
		rate := c.getTransferRate()
		fmt.Fprintf(c.writer, "Up: %v/s, Down: %v/s              			\n", rate.up.HR(), rate.down.HR())
	}
}

type DataRateLoggingListener struct {
	Listener net.Listener
}

func (l *DataRateLoggingListener) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	return (&DataRateLoggingConnection{}).Init(conn), err
}

func (l DataRateLoggingListener) Close() error   { return l.Listener.Close() }
func (l DataRateLoggingListener) Addr() net.Addr { return l.Listener.Addr() }
