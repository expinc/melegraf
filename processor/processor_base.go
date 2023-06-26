package processor

import (
	"fmt"
	"strings"
	"sync"

	"github.com/expinc/melegraf/config"
	"github.com/expinc/melegraf/conveyor"
	"github.com/expinc/melegraf/metric"
	"github.com/robfig/cron/v3"
)

// ProcessorBase is a base struct for all processors
// It implements ProcessorCommonMixin interface
// It is intended to be embedded into concrete processors
// It embeds sync.Mutex so that only one goroutine can modify the processor at a time
// However, it is caller's responsibility to lock the processor before calling any of its mutator methods
type ProcessorBase struct {
	sync.Mutex

	config    config.ProcessorConfig
	inputs    []*conveyor.Conveyor
	outputs   []*conveyor.Conveyor
	isStarted bool
	cron      *cron.Cron
	cronChan  chan struct{}
}

// var _ Processor = (*ProcessorBase)(nil)

func (proc *ProcessorBase) Name() string {
	return proc.config.Name
}

func (proc *ProcessorBase) AddInput(input *conveyor.Conveyor) error {
	for _, in := range proc.inputs {
		if in.Name() == input.Name() {
			return fmt.Errorf("input conveyor \"%s\" already exists", input.Name())
		}
	}

	proc.inputs = append(proc.inputs, input)
	input.OutputProcessorName = proc.Name()
	return nil
}

func (proc *ProcessorBase) RemoveInput(name string) error {
	for i, input := range proc.inputs {
		if input.Name() == name {
			proc.inputs = append(proc.inputs[:i], proc.inputs[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("input conveyor \"%s\" not found", name)
}

func (proc *ProcessorBase) AddOutput(output *conveyor.Conveyor) error {
	for _, out := range proc.outputs {
		if out.Name() == output.Name() {
			return fmt.Errorf("output conveyor \"%s\" already exists", output.Name())
		}
	}

	proc.outputs = append(proc.outputs, output)
	output.InputProcessorName = proc.Name()
	return nil
}

func (proc *ProcessorBase) RemoveOutput(name string) error {
	for i, output := range proc.outputs {
		if output.Name() == name {
			proc.outputs = append(proc.outputs[:i], proc.outputs[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("output conveyor \"%s\" not found", name)
}

func (proc *ProcessorBase) Start() error {
	if proc.isStarted {
		return fmt.Errorf("processor \"%s\" already started", proc.Name())
	}
	proc.isStarted = true

	// Launch cron job
	proc.cronChan = make(chan struct{})
	if strings.TrimSpace(proc.config.CronSpec) != "" {
		proc.cron = cron.New()
		proc.cron.AddFunc(proc.config.CronSpec, func() { proc.cronChan <- struct{}{} })
	}

	return nil
}

func (proc *ProcessorBase) IsStarted() bool {
	return proc.isStarted
}

func (proc *ProcessorBase) Stop() error {
	if !proc.isStarted {
		return fmt.Errorf("processor \"%s\" not started", proc.Name())
	}
	proc.isStarted = false

	// Stop cron job
	if proc.cron != nil {
		proc.cron.Stop()
	}

	return nil
}

func (proc *ProcessorBase) Send(mt metric.Metric) error {
	for _, output := range proc.outputs {
		// Make a deep copy of the metric
		// This is necessary because the metric may be modified by the following processorss
		mtCopy := mt.Copy()

		err := output.Put(&mtCopy)
		if err != nil {
			return err
		}
	}

	return nil
}
