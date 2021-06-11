package chat

import (
	"GoStore/pkg/models"
	uuid "github.com/satori/go.uuid"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	uuid uuid.UUID
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	//f *os.File
	history []models.Message
}

/*type Message struct {
	UUID    string
	Message []byte
}*/

func NewHub() *Hub {
	/*hub := &Hub{
		uuid: uuid.NewV1(),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}

	err := errors.New("da")
	hub.f, err = os.OpenFile("./" + hub.uuid.String() + ".json",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	return hub*/
	return &Hub{
		uuid:       uuid.NewV1(),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		history:    make([]models.Message, 0),
	}
}

func (h *Hub) Clients() map[*Client]bool {
	return h.clients
}

func (h *Hub) History() []models.Message {
	return h.history
}

func (h *Hub) GetUUID() uuid.UUID {
	return h.uuid
}

/*func (h *Hub) CloseFile() {
	h.f.Close()
}*/

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
