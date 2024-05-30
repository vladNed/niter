package discovery

import (
	"context"
	"encoding/json"
	"errors"

	"nhooyr.io/websocket"

	"github.com/indexone/niter/core/config"
	"github.com/indexone/niter/core/discovery/schemas"
	"github.com/indexone/niter/core/logging"
	"github.com/indexone/niter/core/p2p"
)

var logger = logging.NewLogger(config.Config.LogLevel)

// WSClient is a websocket client
type WSClient struct {
	conn    *websocket.Conn
	rtcPeer *p2p.Peer
}

func NewWSClient(rtcPeer *p2p.Peer) (*WSClient, error) {
	conn, _, err := websocket.Dial(context.Background(), SIGNALLING_SERVER, nil)
	if err != nil {
		logger.Error("Error connecting to signalling server:", err.Error())
		return nil, err
	}

	logger.Debug("Connected to signalling server")
	return &WSClient{conn: conn, rtcPeer: rtcPeer}, nil
}

func (ws *WSClient) Start() error {
	logger.Debug("Starting WS client")
	err := ws.registerChannels()
	if err != nil {
		return err
	}
	err = ws.listen()
	if err != nil {
		return err
	}
	return nil
}

func (ws *WSClient) registerChannels() error {
	registerRequest := schemas.RegisterRequest{Channels: []string{OFFERS_CHANNEL, MARKETPLACE_CHANNEL}}
	err := ws.Write(registerRequest)
	if err != nil {
		logger.Error("Error registering channels:", err.Error())
		return err
	}
	return nil
}

func (ws *WSClient) listen() error {
	defer ws.Close()
	for {
		msg, err := ws.recv()
		if err != nil {
			logger.Error("Error receiving message:", err.Error())
			return nil
		}
		logger.Debug("Received message")
		if shouldExit := ws.handleRecvMessages(msg); shouldExit {
			return nil
		}
	}
}

func (ws *WSClient) Write(payload interface{}) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	err = ws.conn.Write(context.Background(), websocket.MessageText, payloadBytes)
	if err != nil {
		return err
	}
	return nil
}

// recv receives a message from the WebSocket connection.
//
// It reads the message type and payload from the connection,
// and then parses the payload into a schemas.Message object.
//
// If there is an error reading or parsing the message, an error is returned.
// Otherwise, the parsed message is returned.
func (ws *WSClient) recv() (schemas.Message, error) {
	msgType, payload, err := ws.conn.Read(context.Background())
	if err != nil {
		return nil, err
	}
	if msgType != websocket.MessageText {
		return nil, errors.New("invalid message type")
	}

	message, err := schemas.ParseReceivedMessage(payload)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (ws *WSClient) handleRecvMessages(msg schemas.Message) bool {
	shouldExit := false
	switch msg := msg.(type) {
	case *schemas.OfferMessage:
		logger.Debug("Received offer message")
		Cache.AddOffer(*msg)
	case *schemas.AnswerMessage:
		logger.Debug("Received answer message")
		ws.rtcPeer.SetOffer(msg.AnswerSDP)
		shouldExit = true
	default:
		logger.Warn("Unknown message type")
	}

	return shouldExit
}

func (ws *WSClient) Close() {
	ws.conn.Close(websocket.StatusNormalClosure, "Closed by client")
}
