//go:build js && wasm
// +build js,wasm

package main

import (
	"encoding/json"
	"syscall/js"

	"github.com/indexone/niter/core/bitcoin"
	"github.com/indexone/niter/core/config"
	"github.com/indexone/niter/core/discovery"
	"github.com/indexone/niter/core/discovery/schemas"
	"github.com/indexone/niter/core/logging"
	"github.com/indexone/niter/core/p2p"
	"github.com/indexone/niter/core/p2p/protocol"
	txs "github.com/indexone/niter/core/transactions"
	"github.com/indexone/niter/core/utils"
)

var (
	ConfigSet     bool
	logger        logging.Logger = logging.NewLogger(logging.INFO)
	wsClient      *discovery.WSClient
	peer          *p2p.Peer
	eventsChannel chan protocol.PeerEvents = make(chan protocol.PeerEvents)
	msgChannel    chan schemas.Message     = make(chan schemas.Message)
	btcWallet     *bitcoin.Wallet

	// Broadcast transaction
	txPool        *txs.TxPool = txs.NewTxPool()
	txPoolChannel chan txs.Tx = make(chan txs.Tx)
	txPoolSignal  chan uint8  = make(chan uint8)
)

const VERSION = "0.1.0"

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
	client, err := discovery.NewWSClient(msgChannel)
	if err != nil {
		return js.Global().Get("Error").New("Error creating WS client: " + err.Error())
	}
	wsClient = client
	go client.Start()
	return nil
}

// Start webRTC connection
func startPeerClient() interface{} {
	newPeer, err := p2p.NewPeer(eventsChannel, msgChannel)
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

			if len(args) == 0 {
				reject.Invoke(js.Global().Get("Error").New("No arguments provided"))
				resolve.Invoke(js.Undefined())
				return
			}
			// Parse the offer details from the app input
			offerJson := args[0].String()
			var offerDetails schemas.OfferDetails
			if err := json.Unmarshal([]byte(offerJson), &offerDetails); err != nil {
				reject.Invoke(js.Global().Get("Error").New("Error parsing offer details: " + err.Error()))
				resolve.Invoke(js.Undefined())
				return
			}

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
			event := <-eventsChannel
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
			offerId := utils.Hash([]byte(encodedSDP))[:6]
			peer.ActiveOfferId = offerId
			offerPayload := schemas.OfferMessage{
				Type:         localSess.Type.String(),
				OfferID:      offerId,
				OfferSDP:     encodedSDP,
				OfferDetails: offerDetails,
			}
			err = wsClient.Write(offerPayload)
			if err != nil {
				reject.Invoke(js.Global().Get("Error").New("Error sending offer: " + err.Error()))
				resolve.Invoke(js.Undefined())
				return
			}
			peer.State = protocol.PeerNegotiating
			resolve.Invoke(js.ValueOf(offerId))
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
			peer.ActiveOfferId = args[0].String()
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
			event := <-eventsChannel
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

// Wallet methods
func initWallet(this js.Value, args []js.Value) interface{} {
	handler := js.FuncOf(func(this js.Value, inputs []js.Value) interface{} {
		resolve := inputs[0]
		reject := inputs[1]

		go func() {
			if len(args) == 1 {
				wif := args[0].String()
				chainParams := config.Config.GetChainParams()
				wallet, err := bitcoin.LoadWallet(wif, chainParams)
				if err != nil {
					reject.Invoke(js.Global().Get("Error").New("Could not load BTC wallet"))
					resolve.Invoke(js.Undefined())
					return
				}
				btcWallet = wallet
				resolve.Invoke(js.ValueOf(wif))
				return
			}

			chainParams := config.Config.GetChainParams()
			wallet, err := bitcoin.GenerateWallet(chainParams)
			if err != nil {
				reject.Invoke(js.Global().Get("Error").New("Could not generate BTC wallet"))
				resolve.Invoke(js.Undefined())
				return
			}

			btcWallet = wallet
			wif, err := btcWallet.WIF()
			if err != nil {
				reject.Invoke(js.Global().Get("Error").New("Could not generate WIF"))
				resolve.Invoke(js.Undefined())
				return
			}
			resolve.Invoke(js.ValueOf(wif))
		}()

		return nil
	})

	return js.Global().Get("Promise").New(handler)
}

func getWalletAddress(this js.Value, args []js.Value) interface{} {
	if btcWallet == nil {
		return js.Global().Get("Error").New("Wallet not initialized")
	}

	address, err := btcWallet.Address().Serialize()
	if err != nil {
		return js.Global().Get("Error").New("Could not fetch segwit address")
	}
	return js.ValueOf(address)
}

func getPendingBroadcast(this js.Value, args []js.Value) interface{} {
	handler := js.FuncOf(func(this js.Value, inputs []js.Value) interface{} {
		resolve := inputs[0]
		reject := inputs[1]

		go func() {
			pendingTx := txPool.Next()
			if pendingTx == nil {
				resolve.Invoke(js.ValueOf(nil))
				reject.Invoke(js.Undefined())
				return
			}

			txJson := pendingTx.Serialize()
			resolve.Invoke(js.ValueOf(txJson))
			reject.Invoke(js.Undefined())
		}()

		return nil
	})

	return js.Global().Get("Promise").New(handler)
}

func markPendingBroadcast(this js.Value, args []js.Value) interface{} {
	if len(args) == 0 {
		return js.Global().Get("Error").New("No arguments provided")
	}

	txQueueId := args[0].Int()
	txPool.Mark()
	txPoolSignal <- uint8(txQueueId)

	return js.ValueOf(true)
}

func main() {
	go txs.RunTxPoolHandler(txPool, txPoolChannel)
	jsGlobal := js.Global()

	// System
	jsGlobal.Set("wasmVersion", VERSION)
	jsGlobal.Set("wasmInit", js.FuncOf(initialize))

	// P2P
	jsGlobal.Set("wasmCreateOffer", js.FuncOf(createOffer))
	jsGlobal.Set("wasmCreateAnswer", js.FuncOf(createAnswer))
	jsGlobal.Set("wasmPollOffers", js.FuncOf(pollOffers))
	jsGlobal.Set("wasmGetPeerState", js.FuncOf(getPeerState))
	jsGlobal.Set("wasmSendData", js.FuncOf(wasmSendData))
	jsGlobal.Set("wasmPollExchangeData", js.FuncOf(pollExchangeData))

	// Wallet
	jsGlobal.Set("wasmInitWallet", js.FuncOf(initWallet))
	jsGlobal.Set("wasmGetWalletAddress", js.FuncOf(getWalletAddress))

	// Events
	jsGlobal.Set("wasmGetPendingBroadcast", js.FuncOf(getPendingBroadcast))
	jsGlobal.Set("wasmMarkPendingBroadcast", js.FuncOf(markPendingBroadcast))

	// This is a blocking call to keep the wasm running.
	<-make(chan struct{})
}
