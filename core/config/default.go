package config

import (
	"github.com/indexone/niter/core/logging"
	"github.com/pion/webrtc/v4"
)

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

func GetICEConfiguration() webrtc.Configuration {
	return webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}
}
