package config

import (
	"errors"
	"strings"

	"github.com/expinc/melegraf/globals"
)

type ProcessorConfig interface {
	Validate() error
}

type ProcessorConfigBase struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	CronSpec string `json:"cronSpec"`
}

func (config *ProcessorConfigBase) ValidateBase() error {
	if strings.TrimSpace(config.Name) == "" {
		return errors.New("name is required")
	}

	if strings.TrimSpace(config.Type) == "" {
		return errors.New("type is required")
	}

	if _, err := globals.CronParser.Parse(strings.TrimSpace(config.CronSpec)); err != nil {
		return err
	}

	return nil
}
