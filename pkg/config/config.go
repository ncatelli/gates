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
	return fmt.Sprintf("gate type is invalid: [%s]", e.gate)
}

// ErrInvalidOutputType represents an invalid output being passed in the `OUTPUT_TYPE`
// env.
type ErrInvalidOutputType struct {
	output string
}

func (e *ErrInvalidOutputType) Error() string {
	return fmt.Sprintf("output type is invalid: [%s]", e.output)
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
	OutputTy   string `env:"OUTPUT_TYPE"`
	OutputAddr string `env:"OUTPUT_ADDR" envDefault:"127.0.0.1"`
}

// New initializes a Config, attempting to parse parames from Envs.
func New() (Config, error) {
	c := Config{}

	if err := env.Parse(&c); err != nil {
		return c, err
	}

	valid := false
	switch c.GateTy {
	case "and":
		valid = true
	case "not":
		valid = true
	default:
		valid = false
	}

	if !valid {
		return c, &ErrInvalidGateType{
			gate: c.GateTy,
		}
	}

	switch c.GateTy {
	case "gate":
		valid = true
	default:
		valid = false
	}

	if !valid {
		return c, &ErrInvalidOutputType{
			output: c.OutputTy,
		}
	}

	return c, nil
}
