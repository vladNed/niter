package config

import "github.com/indexone/niter/core/logging"

var Config = NewDefaultConfig()

type DefaultConfig struct {
	LogLevel            logging.LogLevel `json:"logLevel"`
	SignallingServerURL string           `json:"signallingServerUrl"`
}

func NewDefaultConfig() *DefaultConfig {
	return &DefaultConfig{
		LogLevel:            logging.DEBUG,
		SignallingServerURL: "",
	}
}