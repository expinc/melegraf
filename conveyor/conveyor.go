package conveyor

import (
	"errors"
	"sync"

	"github.com/expinc/melegraf/metric"
)

type Conveyor struct {
	name                string
	channel             chan *metric.Metric
	closeOnce           sync.Once
	InputProcessorName  string
	OutputProcessorName string
}

func NewConveyor(name string, size int) *Conveyor {
	return &Conveyor{
		name:    name,
		channel: make(chan *metric.Metric, size),
	}
}

func (conveyor *Conveyor) Name() string {
	return conveyor.name
}

func (conveyor *Conveyor) Put(metric *metric.Metric) (err error) {
	// recover panic and return error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("conveyor is disposed")
		}
	}()

	select {
	case conveyor.channel <- metric:
		err = nil
	default:
		err = errors.New("conveyor is full")
	}

	return
}

func (conveyor *Conveyor) GetChannel() <-chan *metric.Metric {
	return conveyor.channel
}

func (conveyor *Conveyor) Dispose() {
	conveyor.closeOnce.Do(func() {
		close(conveyor.channel)
	})
}
