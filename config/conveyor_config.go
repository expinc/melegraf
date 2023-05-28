package config

import (
	"errors"
	"strings"
)

type ConveyorConfig struct {
	Name   string `json:"name"`
	Size   uint   `json:"size"`
	Input  string `json:"input"`
	Output string `json:"output"`
}

func (config *ConveyorConfig) Validate() error {
	if strings.TrimSpace(config.Name) == "" {
		return errors.New("name is required")
	}

	if config.Size == 0 {
		return errors.New("size must be a positive number")
	}

	if strings.TrimSpace(config.Input) == "" {
		return errors.New("input is required")
	}

	if strings.TrimSpace(config.Output) == "" {
		return errors.New("output is required")
	}

	return nil
}
