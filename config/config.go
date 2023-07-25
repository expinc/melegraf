package config

import "encoding/json"

type Config interface {
	Validate() error
}

type CustomConfig interface {
	Config
	json.Unmarshaler
}
