package processors

import (
	"fmt"

	"github.com/expinc/melegraf/config"
	"github.com/expinc/melegraf/metric"
	"github.com/expinc/melegraf/processor"
)

const (
	ProcessorTypeDummy = "dummy"
)

func init() {
	processor.RegisterProcessorConstructor(ProcessorTypeDummy, NewDummyProcessor)
}

type dummyProcessor struct {
	cfg           *config.ProcessorConfig
	hasSetup      bool
	countReceived int
	countSent     int
}

var _ processor.Processor = (*dummyProcessor)(nil)

// NewDummyProcessor creates a new dummy processor
func NewDummyProcessor(cfg *config.ProcessorConfig) (processor.Processor, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}

	_, ok := cfg.Params.(*DummyConfig)
	if !ok {
		return nil, fmt.Errorf("invalid params type for dummy processor: %T", cfg.Params)
	}

	return &dummyProcessor{
		cfg: cfg,
	}, nil
}

func (proc *dummyProcessor) Config() *config.ProcessorConfig {
	return proc.cfg
}

func (proc *dummyProcessor) Setup() error {
	proc.hasSetup = true
	return nil
}

func (proc *dummyProcessor) Close() error {
	proc.hasSetup = false
	return nil
}

func (proc *dummyProcessor) OnReceive(mt metric.Metric) ([]metric.Metric, error) {
	proc.countReceived++
	mt.Tags = append(mt.Tags, metric.Tag{Key: "dummy", Value: "received"})
	proc.countSent++
	return []metric.Metric{mt}, nil
}

func (proc *dummyProcessor) OnCronTrigger() ([]metric.Metric, error) {
	mt := metric.Metric{
		Name: "dummy",
		Tags: []metric.Tag{
			{Key: "dummy", Value: "cron"},
		},
		Fields: []metric.Field{
			{Key: "dummy", Value: "cron"},
		},
	}
	proc.countSent++
	return []metric.Metric{mt}, nil
}
