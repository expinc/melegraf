package config

import (
	"encoding/json"
	"testing"
)

func TestValidateProcessorConfig_Succeed(t *testing.T) {
	configStr := `
	{
		"name": "cpu_usage_collector",
		"type": "cpu_usage_collector",
		"cronSpec": "@every 1s"
	}
	`

	var config ProcessorConfigBase
	err := json.Unmarshal([]byte(configStr), &config)
	if err != nil {
		t.Error(err)
	}

	err = config.ValidateBase()
	if err != nil {
		t.Error(err)
	}
}

func TestValidateProcessorConfig_Fail_SpaceName(t *testing.T) {
	configStr := `
	{
		"name": " \r\n\t",
		"type": "cpu_usage_collector",
		"cronSpec": "@every 1s"
	}
	`

	var config ProcessorConfigBase
	err := json.Unmarshal([]byte(configStr), &config)
	if err != nil {
		t.Error(err)
	}

	err = config.ValidateBase()
	if err == nil {
		t.Errorf("Config with space name should be invalid")
	}
}

func TestValidateProcessorConfig_Fail_SpaceType(t *testing.T) {
	configStr := `
	{
		"name": "cpu_usage_collector",
		"type": " \r\n\t",
		"cronSpec": "@every 1s"
	}
	`

	var config ProcessorConfigBase
	err := json.Unmarshal([]byte(configStr), &config)
	if err != nil {
		t.Error(err)
	}

	err = config.ValidateBase()
	if err == nil {
		t.Errorf("Config with space type should be invalid")
	}
}

func TestValidateProcessorConfig_Fail_InvalidCron(t *testing.T) {
	configStr := `
	{
		"name": "cpu_usage_collector",
		"type": "cpu_usage_collector",
		"cronSpec": "invalid"
	}
	`

	var config ProcessorConfigBase
	err := json.Unmarshal([]byte(configStr), &config)
	if err != nil {
		t.Error(err)
	}

	err = config.ValidateBase()
	if err == nil {
		t.Errorf("Config with invalid cron should be invalid")
	}
}
