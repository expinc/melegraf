package config

type Config struct {
	Processors []ProcessorConfig `json:"processors"`
	Conveyors  []ConveyorConfig  `json:"conveyors"`
}
