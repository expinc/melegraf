package config

type MelegrafConfig struct {
	Processors []ProcessorConfig `json:"processors"`
	Conveyors  []ConveyorConfig  `json:"conveyors"`
}

var _ Config = (*MelegrafConfig)(nil)

func (config *MelegrafConfig) Validate() error {
	for _, processor := range config.Processors {
		if err := processor.Validate(); err != nil {
			return err
		}
	}

	for _, conveyor := range config.Conveyors {
		if err := conveyor.Validate(); err != nil {
			return err
		}
	}

	return nil
}
