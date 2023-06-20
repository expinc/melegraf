package config

import (
	"errors"
	"strings"

	"github.com/expinc/melegraf/globals"
)

type ProcessorConfig struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	CronSpec string `json:"cronSpec"`
	Params   Config `json:"params"`
}

var _ Config = (*ProcessorConfig)(nil)

func (config *ProcessorConfig) Validate() error {
	if strings.TrimSpace(config.Name) == "" {
		return errors.New("name is required")
	}

	if strings.TrimSpace(config.Type) == "" {
		return errors.New("type is required")
	}

	if _, err := globals.CronParser.Parse(strings.TrimSpace(config.CronSpec)); err != nil {
		return err
	}

	if nil != config.Params {
		if err := config.Params.Validate(); err != nil {
			return err
		}
	}

	return nil
}
