package engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/expinc/melegraf/config"
	"github.com/expinc/melegraf/conveyor"
	"github.com/expinc/melegraf/processor"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	})
}

type Engine interface {
	Run() error
}

type engine struct {
	cfg       config.MelegrafConfig
	procs     map[string]processor.ProcessorRunner
	conveyors map[string]*conveyor.Conveyor
}

func NewEngine() Engine {
	return &engine{}
}

func (eg *engine) loadConfigFromFile(file string) error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	eg.cfg = config.MelegrafConfig{}
	err = json.Unmarshal(content, &eg.cfg)
	if err != nil {
		return err
	}

	return nil
}

func (eg *engine) Run() error {
	logrus.Info("Starting melegraf...")
	var err error
	end := false
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	reloadChan := make(chan struct{}, 1)
	for {
		err = eg.loadConfigFromFile(viper.ConfigFileUsed())
		if err != nil {
			break
		}

		err = eg.cfg.Validate()
		if err != nil {
			break
		}

		err = eg.bootstrap()
		if err != nil {
			eg.cleanup()
			break
		}

		select {
		case <-sigchan:
			logrus.Info("Got signal, exiting gracefully...")
			end = true
		case <-reloadChan:
			logrus.Info("Got reload signal, reloading...")
			end = false
		}

		err = eg.cleanup()
		if err != nil {
			break
		}

		if end {
			break
		}
	}

	return err
}

func (eg *engine) cleanup() error {
	logrus.Info("Cleaning up...")
	failed := false

	// Stop all processors
	if eg.procs != nil {
		for _, proc := range eg.procs {
			err := proc.Stop()
			if err != nil {
				logrus.Errorf("Error stopping processor %s: %s", proc.Name(), err)
				failed = true
			}
		}
	}

	// Remove all conveyors
	if eg.conveyors != nil {
		for _, conv := range eg.conveyors {
			if eg.procs != nil {
				if strings.TrimSpace(conv.InputProcessorName) != "" {
					proc, ok := eg.procs[conv.InputProcessorName]
					if ok {
						proc.RemoveInput(conv.Name())
					} else {
						logrus.Errorf("processor %s not found", conv.InputProcessorName)
						failed = true
					}
				}

				if strings.TrimSpace(conv.OutputProcessorName) != "" {
					proc, ok := eg.procs[conv.OutputProcessorName]
					if ok {
						proc.RemoveOutput(conv.Name())
					} else {
						logrus.Errorf("processor %s not found", conv.OutputProcessorName)
						failed = true
					}
				}
			}

			conv.Dispose()
		}

		eg.conveyors = nil
	}

	// Remove all processors
	eg.procs = nil

	if failed {
		return fmt.Errorf("cleanup failed")
	}

	return nil
}

func (eg *engine) bootstrap() error {
	logrus.Info("Bootstrapping...")

	// Create all processors
	eg.procs = make(map[string]processor.ProcessorRunner)
	for _, procCfg := range eg.cfg.Processors {
		proc, err := processor.NewProcessorRunner(procCfg.Type, procCfg)
		if err != nil {
			return err
		}

		eg.procs[procCfg.Name] = proc
	}

	// Create all conveyors
	eg.conveyors = make(map[string]*conveyor.Conveyor)
	for _, convCfg := range eg.cfg.Conveyors {
		conv := conveyor.NewConveyor(convCfg.Name, int(convCfg.Size))
		eg.conveyors[convCfg.Name] = conv

		inputProc, ok := eg.procs[convCfg.Input]
		if !ok {
			return fmt.Errorf("processor \"%s\" not found", convCfg.Input)
		}
		inputProc.AddOutput(conv)

		outputProc, ok := eg.procs[convCfg.Output]
		if !ok {
			return fmt.Errorf("processor \"%s\" not found", convCfg.Output)
		}
		outputProc.AddInput(conv)
	}

	// Start all processors
	for _, proc := range eg.procs {
		err := proc.Start()
		if err != nil {
			return err
		}
	}

	return nil
}
