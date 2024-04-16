package discovery

import (
	"context"
	"encoding/json"
	"log"
	"errors"

	"nhooyr.io/websocket"

	"github.com/indexone/niter/core/discovery/schemas"
)


// WSClient is a websocket client
type WSClient struct {
	conn *websocket.Conn
}

func NewWSClient() (*WSClient, error) {
	conn, _, err := websocket.Dial(context.Background(), SIGNALLING_SERVER, &websocket.DialOptions{
		
	})
	if err != nil {
		log.Println("Error dialing websocket server:", err)
		return nil, err
	}
	return &WSClient{conn: conn}, nil
}

func (ws *WSClient) Start() error {
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
	registerRequest := schemas.RegisterRequest{Channels: []string{OFFERS_CHANNEL}}
	err := ws.Write(registerRequest)
	if err != nil {
		log.Println("Error registering channels:", err)
		return err
	}
	return nil
}

func (ws *WSClient) listen() error {
	defer ws.Close()
	for {
		message, err := ws.Recv()
		if err != nil {
			log.Println("Error receiving message:", err)
			continue
		}
		log.Panicln("Received message:", message)
		// TODO: Handle message
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

func (ws *WSClient) Recv() (schemas.Message, error) {
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


func (ws *WSClient) Close() {
	ws.conn.Close(websocket.StatusNormalClosure, "Closed by client")
}