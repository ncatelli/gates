package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

// ErrInvalidGateType represents an invalid gate being passed in the `GATE_TYPE`
// env.
type ErrInvalidGateType struct {
	gate string
}

func (e *ErrInvalidGateType) Error() string {
	return fmt.Sprintf("gate type is invalid: [%v]", e.gate)
}

// ErrUndefinedConfig represents a configuration hasn't been specified.
type ErrUndefinedConfig struct{}

func (e *ErrUndefinedConfig) Error() string {
	return "configuration has not been specified"
}

// Config stores configuration parameters for interacting with the server at a
// global level.
type Config struct {
	ListenAddr string `env:"ListenAddr" envDefault:"0.0.0.0:8080"`
	GateTy     string `env:"GATE_TYPE"`
	OutputAddr string `env:"OUTPUT_ADDR"`
}

// New initializes a Config, attempting to parse parames from Envs.
func New() (Config, error) {
	c := Config{}

	if err := env.Parse(&c); err != nil {
		return c, err
	}

	valid := false
	switch c.GateTy {
	case "not":
		valid = true
	default:
		valid = false
	}

	if valid {
		return c, nil
	} else {
		return c, &ErrInvalidGateType{
			gate: c.GateTy,
		}
	}
}
