package hub

import (
	"bytes"
	"time"

	"github.com/gorilla/websocket"
	"github.com/indexone/signalling-server/internal/cache"
)

var (
	writeWait = 1 * time.Second
	newline   = []byte{'\n'}
	space     = []byte{' '}
)

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	channels []string
	state    State
}

func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		hub:   hub,
		conn:  conn,
		state: New,
	}
}

func (c *Client) Register() {
	c.hub.register <- c
}

func (c *Client) Unregister() {
	cache.MemcacheInstance.ClearData(c)
	c.conn.Close()
	c.hub.unregister <- c
}

func (c *Client) handleChannelSubscribe(message []byte) {
	channelSubscribeRequest, err := parseChannelRequest(message)
	if err != nil {
		logger.Error("Error parsing channel request: ", err)
		return
	}
	c.channels = channelSubscribeRequest.Channels
	logger.Info("Client subscribed to channels: ", channelSubscribeRequest.Channels)

	if ContainsChannel(c.channels, MarketplaceChannel) {
		for _, value := range cache.MemcacheInstance.All() {
			c.WriteStream(*value)
		}
	}

	c.state = Registered
}

func (c *Client) handleMessage(message []byte) {
	message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
	messageRequest, err := parseMessageRequest(message)
	if err != nil {
		logger.Error("Error parsing message request: ", err)
		return
	}
	switch messageRequest := messageRequest.(type) {
	case *CreateOfferRequest:
		cache.MemcacheInstance.Set(messageRequest.OfferID, c, message)
		broadcastMessage := BroadcastMessage{
			Channel: MarketplaceChannel,
			Message: string(message),
		}
		c.hub.broadcast <- &broadcastMessage
		logger.Info("New offer created: ", messageRequest.OfferID)
	case *AnswerOfferRequest:
		if !cache.MemcacheInstance.Contains(messageRequest.OfferID) {
			response := MessageResponse{Status: WS_BAD_REQUEST_STATUS, Details: "Offer not found"}
			responsePayload, err := parseMessageResponse(response)
			if err != nil {
				logger.Error("Error parsing message response: ", err)
				break
			}
			c.WriteStream(responsePayload)
			logger.Error("Offer not found")
			return
		}
		client, ok := cache.MemcacheInstance.Get(messageRequest.OfferID)
		if !ok {
			logger.Error("Error casting client from cache")
			return
		}
		cl := client.(*Client)
		cl.WriteStream(message)
	default:
		logger.Error("Invalid message type")
	}
}

func (c *Client) ReadStream() {
	defer c.Unregister()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Error("Error reading message from websocket: ", err)
			}
			break
		}

		if c.state == New {
			c.handleChannelSubscribe(message)
			continue
		}

		c.handleMessage(message)
	}
}

func (c *Client) WriteStream(message []byte) {
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))

	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return
	}
	w.Write(message)

	if err := w.Close(); err != nil {
		return
	}
}
