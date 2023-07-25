package config

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/expinc/melegraf/globals"
)

type CustomConfigConstructor func() CustomConfig

var (
	type2Params = map[string]CustomConfigConstructor{}
)

type ProcessorConfig struct {
	Name     string       `json:"name"`
	Type     string       `json:"type"`
	CronSpec string       `json:"cronSpec"`
	Params   CustomConfig `json:"params"`
}

var _ CustomConfig = (*ProcessorConfig)(nil)

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

	if config.Params != nil {
		if err := config.Params.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (config *ProcessorConfig) UnmarshalJSON(data []byte) error {
	aux := &struct {
		Name     string          `json:"name"`
		Type     string          `json:"type"`
		CronSpec string          `json:"cronSpec"`
		Params   json.RawMessage `json:"params"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	config.Name = aux.Name
	config.Type = aux.Type
	config.CronSpec = aux.CronSpec

	if aux.Params != nil {
		paramConstructor, ok := type2Params[config.Type]
		if !ok {
			return errors.New("invalid processor type")
		}
		params := paramConstructor()
		if err := json.Unmarshal(aux.Params, &params); err != nil {
			return err
		}
		config.Params = params
	}

	return nil
}
