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
	conveyors map[string]conveyor.Conveyor
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
	// Create all processors

	// Create all conveyors

	// Start all processors

	// TODO
	return nil
}
