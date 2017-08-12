package utils

import (
	"github.com/paulbellamy/ratecounter"
	"github.com/c2h5oh/datasize"
	"time"
)

type DataRateCounter interface {
	GetDataRate() datasize.ByteSize
	CaptureEvent(rate int)
}

func NewRateCounter() DataRateCounter {
	counter := dataRateCounter{
		rateCounter:     ratecounter.NewRateCounter(10 * time.Second),
		noActivityTimer: time.NewTimer(time.Second),
	}
	go counter.updateOnNoActivity()

	return counter
}

type dataRateCounter struct {
	rateCounter *ratecounter.RateCounter
	noActivityTimer *time.Timer
}

func (c dataRateCounter) GetDataRate() datasize.ByteSize {
	return (datasize.ByteSize)(c.rateCounter.Rate()/10) * datasize.B
}

func (c dataRateCounter) CaptureEvent(rate int) {
	c.rateCounter.Incr(int64(rate))
	c.noActivityTimer.Reset(time.Second)
}

func (c dataRateCounter) updateOnNoActivity() {
	for {
		<-c.noActivityTimer.C
		c.rateCounter.Incr(0)
	}
}
