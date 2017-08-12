package utils

import (
	"github.com/paulbellamy/ratecounter"
	"github.com/c2h5oh/datasize"
	"time"
)

var rateUpdateInterval = 10 * time.Second

type DataRateCounter interface {
	GetDataRate() datasize.ByteSize
	CaptureEvent(rate int)
}

func NewRateCounter() DataRateCounter {
	counter := dataRateCounter{
		rateCounter:     ratecounter.NewRateCounter(rateUpdateInterval),
		noActivityTimer: time.NewTimer(time.Second),
	}
	go counter.updateOnNoActivity()

	return counter
}

type dataRateCounter struct {
	rateCounter     *ratecounter.RateCounter
	noActivityTimer *time.Timer
}

func (c dataRateCounter) GetDataRate() datasize.ByteSize {
	return (datasize.ByteSize)(c.rateCounter.Rate()/10) * datasize.B
}

func (c dataRateCounter) CaptureEvent(rate int) {
	c.rateCounter.Incr(int64(rate))
	c.noActivityTimer.Reset(rateUpdateInterval)
}

func (c dataRateCounter) updateOnNoActivity() {
	for {
		<-c.noActivityTimer.C
		c.rateCounter.Incr(0)
		c.noActivityTimer.Reset(rateUpdateInterval)
	}
}
