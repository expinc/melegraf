package processor

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/expinc/melegraf/config"
	"github.com/expinc/melegraf/conveyor"
	"github.com/expinc/melegraf/metric"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

// ProcessorRunner takes care of running processors
type ProcessorRunner interface {
	// Name returns the name of the processor
	Name() string

	// AddInput adds an input conveyor to the processor
	// The processor will be stopped and started again after this method is called
	AddInput(input *conveyor.Conveyor) error

	// RemoveInput removes an input conveyor from the processor
	// The processor will be stopped and started again after this method is called
	RemoveInput(name string) error

	// AddOutput adds an output conveyor to the processor
	// The processor will be stopped and started again after this method is called
	AddOutput(output *conveyor.Conveyor) error

	// RemoveOutput removes an output conveyor from the processor
	// The processor will be stopped and started again after this method is called
	RemoveOutput(name string) error

	// Start starts the processor
	Start() error

	// IsStarted returns true if the processor is started
	IsStarted() bool

	// Stop stops the processor
	Stop() error
}

// NewProcessorRunner creates a new processor runner
func NewProcessorRunner(procType string, cfg *config.ProcessorConfig) (ProcessorRunner, error) {
	proc, err := NewProcessor(procType, cfg)
	if err != nil {
		return nil, err
	}

	return &processorRunner{
		proc: proc,
	}, nil
}

type processorRunner struct {
	sync.Mutex

	proc      Processor
	inputs    []*conveyor.Conveyor
	outputs   []*conveyor.Conveyor
	isStarted bool
	stopChan  chan struct{}
}

var _ ProcessorRunner = (*processorRunner)(nil)

func (runner *processorRunner) Name() string {
	return runner.proc.Config().Name
}

func (runner *processorRunner) IsStarted() bool {
	runner.Lock()
	defer runner.Unlock()
	return runner.isStarted
}

func (runner *processorRunner) send(metrics []metric.Metric) {
	for _, output := range runner.outputs {
		for _, mt := range metrics {
			// Make a deep copy of the metric
			// This is necessary because the metric may be modified by the following processors
			mtCopy := mt.Copy()

			err := output.Put(mtCopy)
			if err != nil {
				logrus.Errorf("Processor \"%s\" failed to send metric to conveyor \"%s\": %v", runner.Name(), output.Name(), err)
			}
		}
	}
}

func (runner *processorRunner) startInternal() error {
	if runner.isStarted {
		logrus.Infof("Processor \"%s\" already started. Do nothing", runner.Name())
		return nil
	}

	err := runner.proc.Setup()
	if err != nil {
		return err
	}

	var crn *cron.Cron
	cronChan := make(chan struct{})
	if strings.TrimSpace(runner.proc.Config().CronSpec) != "" {
		crn = cron.New()
		crn.AddFunc(runner.proc.Config().CronSpec, func() { cronChan <- struct{}{} })
	}

	runner.stopChan = make(chan struct{})
	go func() {
		// select from stopChan, cronChan and inputs
		cases := make([]reflect.SelectCase, len(runner.inputs)+2)
		cases[0] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(runner.stopChan)}
		cases[1] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(cronChan)}
		for i, input := range runner.inputs {
			cases[i+2] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(input.GetChannel())}
		}

		if crn != nil {
			crn.Start()
		}

		stopped := false
		for {
			chosen, value, _ := reflect.Select(cases)

			switch chosen {
			case 0:
				logrus.Infof("Processor \"%s\" is being stopped", runner.Name())
				stopped = true

				close(runner.stopChan)
				runner.stopChan = nil

				if crn != nil {
					crn.Stop()
				}

				err2 := runner.proc.Close()
				if err2 != nil {
					logrus.Errorf("Processor \"%s\" failed to close: %v", runner.Name(), err2)
				}

				runner.isStarted = false
			case 1:
				out, err2 := runner.proc.OnCronTrigger()
				if err2 != nil {
					logrus.Errorf("Processor \"%s\" failed to process cron trigger: %v", runner.Name(), err2)
				} else {
					runner.send(out)
				}
			default:
				mt := value.Interface().(metric.Metric)
				out, err2 := runner.proc.OnReceive(mt)
				if err2 != nil {
					logrus.Errorf("Processor \"%s\" failed to process metric from conveyor \"%s\": %v", runner.Name(), runner.inputs[chosen-2].Name(), err2)
				} else {
					runner.send(out)
				}
			}

			if stopped {
				break
			}
		}
	}()

	runner.isStarted = true
	return nil
}

func (runner *processorRunner) Start() error {
	runner.Lock()
	defer runner.Unlock()
	return runner.startInternal()
}

func (runner *processorRunner) stopInternal() error {
	if !runner.isStarted {
		logrus.Infof("Processor \"%s\" already stopped. Do nothing", runner.Name())
		return nil
	}
	runner.stopChan <- struct{}{}
	return nil
}

func (runner *processorRunner) Stop() error {
	runner.Lock()
	defer runner.Unlock()
	return runner.stopInternal()
}

func (runner *processorRunner) AddInput(input *conveyor.Conveyor) error {
	runner.Lock()
	defer runner.Unlock()

	originallyStarted := runner.isStarted
	if originallyStarted {
		err := runner.stopInternal()
		if err != nil {
			return err
		}
	}

	for _, in := range runner.inputs {
		if in.Name() == input.Name() {
			return fmt.Errorf("input conveyor \"%s\" already exists", input.Name())
		}
	}

	runner.inputs = append(runner.inputs, input)
	input.OutputProcessorName = runner.Name()

	if originallyStarted {
		err := runner.startInternal()
		if err != nil {
			return err
		}
	}

	return nil
}

func (runner *processorRunner) RemoveInput(name string) error {
	runner.Lock()
	defer runner.Unlock()

	originallyStarted := runner.isStarted
	if originallyStarted {
		err := runner.stopInternal()
		if err != nil {
			return err
		}
	}

	removed := false
	for i, input := range runner.inputs {
		if input.Name() == name {
			runner.inputs = append(runner.inputs[:i], runner.inputs[i+1:]...)
			removed = true
			break
		}
	}

	if !removed {
		return fmt.Errorf("input conveyor \"%s\" not found", name)
	}

	if originallyStarted {
		err := runner.startInternal()
		if err != nil {
			return err
		}
	}

	return nil
}

func (runner *processorRunner) AddOutput(output *conveyor.Conveyor) error {
	runner.Lock()
	defer runner.Unlock()

	originallyStarted := runner.isStarted
	if originallyStarted {
		err := runner.stopInternal()
		if err != nil {
			return err
		}
	}

	for _, out := range runner.outputs {
		if out.Name() == output.Name() {
			return fmt.Errorf("output conveyor \"%s\" already exists", output.Name())
		}
	}

	runner.outputs = append(runner.outputs, output)
	output.InputProcessorName = runner.Name()

	if originallyStarted {
		err := runner.startInternal()
		if err != nil {
			return err
		}
	}

	return nil
}

func (runner *processorRunner) RemoveOutput(name string) error {
	runner.Lock()
	defer runner.Unlock()

	originallyStarted := runner.isStarted
	if originallyStarted {
		err := runner.stopInternal()
		if err != nil {
			return err
		}
	}

	removed := false
	for i, output := range runner.outputs {
		if output.Name() == name {
			runner.outputs = append(runner.outputs[:i], runner.outputs[i+1:]...)
			removed = true
			break
		}
	}

	if !removed {
		return fmt.Errorf("output conveyor \"%s\" not found", name)
	}

	if originallyStarted {
		err := runner.startInternal()
		if err != nil {
			return err
		}
	}

	return nil
}
