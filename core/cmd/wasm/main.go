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
	"github.com/indexone/niter/core/p2p/protocol"
	"github.com/indexone/niter/core/utils"
)

var (
	ConfigSet     bool
	logger        logging.Logger = logging.NewLogger(logging.INFO)
	wsClient      *discovery.WSClient
	peer          *p2p.Peer
	eventsChannel chan protocol.PeerEvents = make(chan protocol.PeerEvents)
)

const VERSION = "0.1.25"

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
	client, err := discovery.NewWSClient(peer)
	if err != nil {
		return js.Global().Get("Error").New("Error creating WS client: " + err.Error())
	}
	wsClient = client
	go client.Start()
	return nil
}

// Start webRTC connection
func startPeerClient() interface{} {
	newPeer, err := p2p.NewPeer(eventsChannel)
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

			// Start the peer client
			err = startPeerClient()
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

			logger.Debug("Module initialized completed")
			resolve.Invoke(js.Undefined())
		}()

		return nil
	})

	return js.Global().Get("Promise").New(handler)
}

// Checks if the peer is initialized. If not, it returns an error.
func isPeerInitialized() interface{} {
	if peer == nil {
		return js.Global().Get("Error").New("Peer not initialized")
	}
	return nil
}

// The create offer function, starts the local connection for the web rtc node
// and creates a data channel. With this the node is ready to create an offer
// and send it to the signalling server to be forwarded to the other peer.
//
// This way the node becomes an initiator node.
func createOffer(this js.Value, args []js.Value) interface{} {
	if err := isPeerInitialized(); err != nil {
		return err
	}

	handler := js.FuncOf(func(this js.Value, inputs []js.Value) interface{} {
		resolve := inputs[0]
		reject := inputs[1]
		go func() {
			err := peer.StartInitiator()
			if err != nil {
				reject.Invoke(js.Global().Get("Error").New("Error starting initiator: " + err.Error()))
				resolve.Invoke(js.Undefined())
				return
			}
			_, err = peer.CreateOffer()
			if err != nil {
				reject.Invoke(js.Global().Get("Error").New("Error creating offer: " + err.Error()))
				resolve.Invoke(js.Undefined())
				return
			}
			event := <- eventsChannel
			if event != protocol.InitiatorICECandidate {
				reject.Invoke(js.Global().Get("Error").New("Error gathering ICE candidates"))
				resolve.Invoke(js.Undefined())
				return
			}

			localSess := peer.LocalConnection.LocalDescription()
			if localSess == nil {
				reject.Invoke(js.Global().Get("Error").New("Error unmarshalling local description: "))
				resolve.Invoke(js.Undefined())
				return
			}
			encodedSDP, err := utils.EncodeSDP(localSess)
			if err != nil {
				reject.Invoke(js.Global().Get("Error").New("Error encoding local SDP: " + err.Error()))
				resolve.Invoke(js.Undefined())
				return
			}
			offerId := utils.Hash([]byte(encodedSDP))
			offerPayload := schemas.OfferMessage{
				Type:     localSess.Type.String(),
				OfferID:  offerId[:6],
				OfferSDP: encodedSDP,
			}
			err = wsClient.Write(offerPayload)
			if err != nil {
				reject.Invoke(js.Global().Get("Error").New("Error sending offer: " + err.Error()))
				resolve.Invoke(js.Undefined())
				return
			}
			peer.State = protocol.PeerNegotiating
			resolve.Invoke(js.ValueOf(offerId[:6]))
		}()

		return nil
	})
	return js.Global().Get("Promise").New(handler)
}

// The create answer function should be called only when the node is convinced
// about an offer and wants to connect to the initiator node. Usually the
// initiator node should have created a data channel and sent an offer SDP to
// the signalling server.
//
// The node will get the offer SDP from the signalling server and set it as the
// remote description. Then it will create an answer and send it back to the
// initiator node, thus rendering the node a responder node.
func createAnswer(this js.Value, args []js.Value) interface{} {
	if err := isPeerInitialized(); err != nil {
		return err
	}

	handler := js.FuncOf(func(this js.Value, inputs []js.Value) interface{} {
		resolve := inputs[0]
		reject := inputs[1]
		go func() {

			err := peer.StartResponder()
			if err != nil {
				reject.Invoke(js.Global().Get("Error").New("Error starting responder: " + err.Error()))
				resolve.Invoke(js.Undefined())
				return
			}
			offerData, ok := discovery.Cache.GetOffer(args[0].String())
			if !ok {
				reject.Invoke(js.Global().Get("Error").New("Offer not found"))
				resolve.Invoke(js.Undefined())
				return
			}
			offerSDP := offerData["OfferSDP"].(string)
			err = peer.SetOffer(offerSDP)
			if err != nil {
				reject.Invoke(js.Global().Get("Error").New("Error setting offer: " + err.Error()))
				resolve.Invoke(js.Undefined())
				return
			}
			_, err = peer.CreateAnswer()
			if err != nil {
				reject.Invoke(js.Global().Get("Error").New("Error creating answer: " + err.Error()))
				resolve.Invoke(js.Undefined())
				return
			}
			event := <- eventsChannel
			if event != protocol.ResponderICECandidate {
				reject.Invoke(js.Global().Get("Error").New("Error gathering ICE candidates"))
				resolve.Invoke(js.Undefined())
				return
			}
			localSession := peer.LocalConnection.LocalDescription()
			if localSession == nil {
				reject.Invoke(js.Global().Get("Error").New("Error unmarshalling local description"))
				resolve.Invoke(js.Undefined())
				return
			}
			encodedAnswer, err := utils.EncodeSDP(localSession)
			if err != nil {
				reject.Invoke(js.Global().Get("Error").New("Error encoding answer: " + err.Error()))
				resolve.Invoke(js.Undefined())
				return
			}
			answerPayload := schemas.AnswerMessage{
				Type:      localSession.Type.String(),
				OfferID:   args[0].String(),
				AnswerSDP: encodedAnswer,
			}
			err = wsClient.Write(answerPayload)
			if err != nil {
				reject.Invoke(js.Global().Get("Error").New("Error sending answer: " + err.Error()))
				resolve.Invoke(js.Undefined())
				return
			}
			logger.Debug("Answer sent")
			resolve.Invoke(js.Undefined())
		}()

		return nil
	})

	return js.Global().Get("Promise").New(handler)
}

func pollOffers(this js.Value, args []js.Value) interface{} {
	handler := js.FuncOf(func(this js.Value, inputs []js.Value) interface{} {
		resolve := inputs[0]
		reject := inputs[1]
		go func() {
			offers := discovery.Cache.AllOffers()

			resolve.Invoke(js.ValueOf(offers))
			reject.Invoke(js.Undefined())
		}()

		return nil
	})

	return js.Global().Get("Promise").New(handler)
}

func getPeerState(this js.Value, args []js.Value) interface{} {
	if err := isPeerInitialized(); err != nil {
		return err
	}

	remotePeer := ""
	if peer.RemotePeer != nil {
		remotePeer = peer.RemotePeer.Id
	}
	payload := map[string]interface{}{
		"id":         peer.Id(),
		"state":      peer.State.String(),
		"remotePeer": remotePeer,
	}

	return js.ValueOf(payload)
}

func wasmSendData(this js.Value, args []js.Value) interface{} {
	if err := isPeerInitialized(); err != nil {
		return err
	}
	data := args[0].String()
	err := peer.SendData([]byte(data))
	if err != nil {
		return js.Global().Get("Error").New("Error sending data: " + err.Error())
	}
	return js.Undefined()
}

func pollExchangeData(this js.Value, args []js.Value) interface{} {
	if err := isPeerInitialized(); err != nil {
		return err
	}

	handler := js.FuncOf(func(this js.Value, inputs []js.Value) interface{} {
		resolve := inputs[0]
		go func() {
			data := peer.ExchangeData
			payload := make([]interface{}, 0)
			for _, d := range data {
				payload = append(payload, map[string]interface{}{
					"side":      d.Side,
					"data":      d.Data,
					"timestamp": d.Timestamp,
				})
			}
			resolve.Invoke(js.ValueOf(payload))
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
	jsGlobal.Set("wasmCreateAnswer", js.FuncOf(createAnswer))
	jsGlobal.Set("wasmPollOffers", js.FuncOf(pollOffers))
	jsGlobal.Set("wasmGetPeerState", js.FuncOf(getPeerState))
	jsGlobal.Set("wasmSendData", js.FuncOf(wasmSendData))
	jsGlobal.Set("wasmPollExchangeData", js.FuncOf(pollExchangeData))

	// This is a blocking call to keep the wasm running.
	<-make(chan struct{})
}
