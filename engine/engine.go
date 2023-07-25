package engine

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/expinc/melegraf/config"
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
	cfg config.MelegrafConfig
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
			logrus.Info("Got signal, exiting gracefully...")
			end = true
		case <-reloadChan:
			logrus.Info("Got reload signal, reloading...")
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
