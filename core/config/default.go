package config

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/indexone/niter/core/logging"
	"github.com/pion/webrtc/v4"
)

var Config = NewDefaultConfig()

type DefaultConfig struct {
	LogLevel            logging.LogLevel `json:"logLevel"`
	SignallingServerURL string           `json:"signallingServerUrl"`
	Network             string           `json:"network"`
}

func (c *DefaultConfig) GetChainParams() *chaincfg.Params {
	switch c.Network {
	case "mainnet":
		return &chaincfg.MainNetParams
	case "testnet":
		return &chaincfg.TestNet3Params
	case "regtest":
		return &chaincfg.RegressionNetParams
	default:
		return &chaincfg.MainNetParams
	}
}
func NewDefaultConfig() *DefaultConfig {
	return &DefaultConfig{
		LogLevel:            logging.DEBUG,
		SignallingServerURL: "",
		Network:             "mainnet",
	}
}

func GetICEConfiguration() webrtc.Configuration {
	return webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{
					"stun:stun.l.google.com:19302",
					"stun:stun1.l.google.com:19302",
					"stun:stun2.l.google.com:19302",
				},
			},
		},
	}
}
