package processor

import (
	"github.com/expinc/melegraf/conveyor"
	"github.com/expinc/melegraf/metric"
)

// ProcessorCommonMixin is a common interface for all processors
type ProcessorCommonMixin interface {
	Name() string
	AddInput(input *conveyor.Conveyor) error
	RemoveInput(name string) error
	AddOutput(output *conveyor.Conveyor) error
	RemoveOutput(name string) error
	Start() error
	IsStarted() bool
	Stop() error
	Send(mt metric.Metric) error
}

// ConcreteProcessor is an interface for concrete processors
type ConcreteProcessor interface {
	// Setup is called once when the processor is started
	// It is typically used to establish connections to external systems
	// or to open necessary files
	Setup() error

	// Close is called once when the processor is stopped
	// It is typically used to close established connections
	// or to close opened files
	Close() error

	// OnReceive is called when a metric is received from an input conveyor
	OnReceive() error

	// OnCronTrigger is called when a cron trigger is fired
	OnCronTrigger() error
}

// Processor is an interface for all processors
type Processor interface {
	ProcessorCommonMixin
	ConcreteProcessor
}
