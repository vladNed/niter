//go:build js && wasm
// +build js,wasm

package main

import (
	"encoding/json"
	"syscall/js"

	"github.com/indexone/niter/core/config"
	"github.com/indexone/niter/core/discovery"
	"github.com/indexone/niter/core/discovery/schemas"
	"github.com/indexone/niter/core/logging"
	"github.com/indexone/niter/core/p2p"
	"github.com/google/uuid"
)

var (
	ConfigSet bool
	logger    logging.Logger = logging.NewLogger(logging.INFO)
	wsClient  *discovery.WSClient
	peer      *p2p.Peer
)

const VERSION = "0.0.2"

// Sets the config within the wasm file. This can only be done at initialization.
// The input of the function is a JSON string representing the config.
func setConfig(inputs []js.Value) interface{} {
	if ConfigSet {
		return js.Global().Get("Error").New("Module already initialized")
	}

	configJson := inputs[0].String()
	var cfg config.DefaultConfig
	if err := json.Unmarshal([]byte(configJson), &cfg); err != nil {
		return js.Global().Get("Error").New("Error parsing config: " + err.Error())
	}

	config.Config = &cfg
	ConfigSet = true
	logger = logging.NewLogger(cfg.LogLevel)

	return nil
}

// Start the websocket client
func startWSClient() interface{} {
	if !ConfigSet && wsClient == nil {
		return js.Global().Get("Error").New("Config not set")
	}
	client, err := discovery.NewWSClient()
	if err != nil {
		return js.Global().Get("Error").New("Error creating WS client: " + err.Error())
	}
	wsClient = client
	go client.Start()
	return nil
}

// Start webRTC connection
func startPeerClient() interface{} {
	newPeer, err := p2p.NewPeer()
	if err != nil {
		return js.Global().Get("Error").New("Error creating peer: " + err.Error())
	}
	peer = newPeer
	return nil
}

// Initializes the module. It return a promise that resolves when the module is initialized.
func initialize(this js.Value, args []js.Value) interface{} {
	handler := js.FuncOf(func(this js.Value, inputs []js.Value) interface{} {
		resolve := inputs[0]
		reject := inputs[1]
		go func() {
			logger.Debug("Initializing module...")
			if len(args) == 0 {
				err := js.Global().Get("Error").New("No arguments provided")
				reject.Invoke(err)
				resolve.Invoke(js.Undefined())
				return
			}

			// Sets the config for the module
			err := setConfig(args)
			if err != nil {
				reject.Invoke(err)
				resolve.Invoke(js.Undefined())
				return
			}

			// Start the websocket client
			err = startWSClient()
			if err != nil {
				reject.Invoke(err)
				resolve.Invoke(js.Undefined())
				return
			}

			// Start the peer client
			err = startPeerClient()
			if err != nil {
				reject.Invoke(err)
				resolve.Invoke(js.Undefined())
				return
			}

			logger.Debug("Module initialized completed")
			resolve.Invoke(js.Undefined())
		}()

		return nil
	})

	return js.Global().Get("Promise").New(handler)
}

func isPeerInitialized() interface{} {
	if peer == nil {
		return js.Global().Get("Error").New("Peer not initialized")
	}
	return nil
}

func createOffer(this js.Value, args []js.Value) interface{} {
	if err := isPeerInitialized(); err != nil {
		return err
	}


	handler := js.FuncOf(func(this js.Value, inputs []js.Value) interface{} {
		resolve := inputs[0]
		reject := inputs[1]
		go func() {
			logger.Debug("Creating new offer")
			offer, err := peer.CreateOffer()
			if err != nil {
				reject.Invoke(js.Global().Get("Error").New("Error creating offer: " + err.Error()))
				resolve.Invoke(js.Undefined())
				return
			}

			// TODO: Hash and make id
			offerId := uuid.New().String()
			offerPayload := schemas.OfferMessage{
				Type: offer.Type.String(),
				OfferID: offerId,
				OfferDescription: offer,
			}

			offerJson, err := json.Marshal(offerPayload)
			err = wsClient.Write(offerJson)
			if err != nil {
				reject.Invoke(js.Global().Get("Error").New("Error sending offer: " + err.Error()))
				resolve.Invoke(js.Undefined())
				return
			}
			logger.Debug("Offer sent")
			resolve.Invoke(js.ValueOf(offerId))
		}()

		return nil
	})
	return js.Global().Get("Promise").New(handler)
}

func main() {
	jsGlobal := js.Global()
	jsGlobal.Set("wasmVersion", VERSION)
	jsGlobal.Set("wasmInit", js.FuncOf(initialize))
	jsGlobal.Set("wasmCreateOffer", js.FuncOf(createOffer))

	// This is a blocking call to keep the wasm running.
	<-make(chan struct{})
}
