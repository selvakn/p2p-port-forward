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
	return dataRateCounter{
		rateCounter: ratecounter.NewRateCounter(10 * time.Second),
	}
}

type dataRateCounter struct {
	rateCounter *ratecounter.RateCounter
	//rateUpdater chan uint
}

func (c dataRateCounter) GetDataRate() datasize.ByteSize {
	return (datasize.ByteSize)(c.rateCounter.Rate()/10) * datasize.B
}

func (c dataRateCounter) CaptureEvent(rate int) {
	c.rateCounter.Incr(int64(rate))
	//return (datasize.ByteSize)(c.rateCounter.Rate()/10) * datasize.B
}
