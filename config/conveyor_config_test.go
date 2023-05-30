package config

import (
	"encoding/json"
	"testing"
)

func TestValidateConveyorConfig_Succeed(t *testing.T) {
	configStr := `
	{
		"name": "cpu2host",
		"size": 10,
		"input": "cpu_usage_collector",
		"output": "hostname_modifier"
	}
	`

	var config ConveyorConfig
	err := json.Unmarshal([]byte(configStr), &config)
	if err != nil {
		t.Error(err)
	}

	err = config.Validate()
	if err != nil {
		t.Error(err)
	}
}

func TestValidateConveyorConfig_Fail_SpaceName(t *testing.T) {
	configStr := `
	{
		"name": " \r\n\t",
		"size": 10,
		"input": "cpu_usage_collector",
		"output": "hostname_modifier"
	}
	`

	var config ConveyorConfig
	err := json.Unmarshal([]byte(configStr), &config)
	if err != nil {
		t.Error(err)
	}

	err = config.Validate()
	if err == nil {
		t.Errorf("Config with space name should be invalid")
	}
}

func TestValidateConveyorConfig_Fail_InvalidSize(t *testing.T) {
	configStr := `
	{
		"name": "cpu2host",
		"size": 0,
		"input": "cpu_usage_collector",
		"output": "hostname_modifier"
	}
	`

	var config ConveyorConfig
	err := json.Unmarshal([]byte(configStr), &config)
	if err != nil {
		t.Error(err)
	}

	err = config.Validate()
	if err == nil {
		t.Errorf("Config with zero size should be invalid")
	}
}

func TestValidateConveyorConfig_Fail_SpaceInput(t *testing.T) {
	configStr := `
	{
		"name": "cpu2host",
		"size": 10,
		"input": " \r\n\t",
		"output": "hostname_modifier"
	}
	`

	var config ConveyorConfig
	err := json.Unmarshal([]byte(configStr), &config)
	if err != nil {
		t.Error(err)
	}

	err = config.Validate()
	if err == nil {
		t.Errorf("Config with space input should be invalid")
	}
}

func TestValidateConveyorConfig_Fail_SpaceOutput(t *testing.T) {
	configStr := `
	{
		"name": "cpu2host",
		"size": 10,
		"input": "cpu_usage_collector",
		"output": " \r\n\t"
	}
	`

	var config ConveyorConfig
	err := json.Unmarshal([]byte(configStr), &config)
	if err != nil {
		t.Error(err)
	}

	err = config.Validate()
	if err == nil {
		t.Errorf("Config with space output should be invalid")
	}
}
