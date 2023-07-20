package engine

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/expinc/melegraf/config"
	"github.com/spf13/viper"
)

type Engine interface {
	Run() error
}

type engine struct {
	cfg config.MelegrafConfig
}

func NewEngine() Engine {
	return &engine{}
}

func (eg *engine) loadConfigFromFile(file string) error {
	// TODO
	return nil
}

func (eg *engine) Run() error {
	var err error
	end := false
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	reloadChan := make(chan struct{}, 1)
	for {
		err = eg.cleanup()
		if err != nil {
			break
		}

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
			break
		}

		select {
		case <-sigchan:
			// TODO: log
			end = true
		case <-reloadChan:
			// TODO: log
			end = false
		}

		if end {
			break
		}
	}

	return err
}

func (eg *engine) cleanup() error {
	// TODO
	return nil
}

func (eg *engine) bootstrap() error {
	// TODO
	return nil
}
