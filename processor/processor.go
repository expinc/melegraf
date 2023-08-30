package processor

import (
	"fmt"

	"github.com/expinc/melegraf/config"
	"github.com/expinc/melegraf/metric"
)

// Processor contains all methods that must be implemented by concrete processors
type Processor interface {
	// Config returns the configuration of the processor
	// The returned value must not be modified
	Config() *config.ProcessorConfig

	// Setup is called once when the processor is started
	// It is typically used to establish connections to external systems
	// or to open necessary files
	Setup() error

	// Close is called once when the processor is stopped
	// It is typically used to close established connections
	// or to close opened files
	Close() error

	// OnReceive is called when a metric is received from an input conveyor
	// It returns a slice of metrics that should be sent to output conveyors
	OnReceive(mt metric.Metric) (out []metric.Metric, err error)

	// OnCronTrigger is called when a cron trigger is fired
	// It returns a slice of metrics that should be sent to output conveyors
	OnCronTrigger() (out []metric.Metric, err error)
}

// ProcessorConstructor is a function that creates a new processor
type ProcessorConstructor func(cfg *config.ProcessorConfig) (Processor, error)

var (
	procType2Constructor = map[string]ProcessorConstructor{}
)

// RegisterProcessorConstructor registers a processor constructor of a given type
func RegisterProcessorConstructor(procType string, constructor ProcessorConstructor) {
	procType2Constructor[procType] = constructor
}

// NewProcessor creates a new processor
func NewProcessor(procType string, cfg *config.ProcessorConfig) (Processor, error) {
	constructor, ok := procType2Constructor[procType]
	if !ok {
		return nil, fmt.Errorf("invalid processor type: %s", procType)
	}

	proc, err := constructor(cfg)
	if err != nil {
		return nil, err
	}

	return proc, nil
}
