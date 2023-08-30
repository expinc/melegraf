package processors

import (
	"encoding/json"

	"github.com/expinc/melegraf/config"
)

type DummyConfig struct {
	PropValue string `json:"propValue"`
}

var _ config.CustomConfig = (*DummyConfig)(nil)

func NewDummyConfig() config.CustomConfig {
	return &DummyConfig{}
}

func init() {
	config.RegisterCustomConfigConstructor(ProcessorTypeDummy, NewDummyConfig)
}

func (cfg *DummyConfig) Validate() error {
	return nil
}

func (cfg *DummyConfig) UnmarshalJSON(data []byte) error {
	aux := &struct {
		PropValue string `json:"propValue"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	cfg.PropValue = aux.PropValue
	return nil
}
