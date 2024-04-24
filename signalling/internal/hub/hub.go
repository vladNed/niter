package hub

import "github.com/indexone/signalling-server/internal/logging"

var (
	logger = logging.GetLogger(nil)

	// Singleton instance of the hub. This is the only instance of the hub that
	// will be used by the application.
	HubInstance = NewHub()
)

type BroadcastMessage struct {
	Message string
	Channel string
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan *BroadcastMessage
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan *BroadcastMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	logger.Info("Running hub instance.")
	for {
		select {
		case client := <-h.register:
			logger.Info("Registering new client")
			h.clients[client] = true
		case client := <-h.unregister:
			logger.Info("Unregister client")
			delete(h.clients, client)
		case message := <-h.broadcast:
			for client := range h.clients {
				hasChannel := false
				for _, channel := range client.channels {
					if channel == message.Channel {
						hasChannel = true
						break
					}
				}
				if !hasChannel {
					logger.Info("Client does not have channel: ", message.Channel)
					continue
				}
				client.WriteStream([]byte(message.Message))
			}
		}
	}
}
