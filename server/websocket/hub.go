package websocket

import (
	"encoding/json"
	"fmt"
	"log"
)

type Hub struct {
	clients map[*Client]bool

	// clients send messages this channel to be broadcast to all clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

type WSMsg struct {
	Event   string          `json:"event"`
	Payload json.RawMessage `json:"payload"`
}

type ConnectPayload struct {
	Id string `json:"id"`
}

type CommMsgPayload struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

func (h *Hub) handleMessage(msg []byte) {
	var envelope WSMsg
	if err := json.Unmarshal(msg, &envelope); err != nil {
		log.Println("Error parsing envelope:", err)
		return
	}

	switch envelope.Event {
	case "connection_opened":
		var p ConnectPayload
		if err := json.Unmarshal(envelope.Payload, &p); err != nil {
			log.Println("Bad connect payload:", err)
			return
		}
		fmt.Printf("Action: Client %s confirmed connection.\n", p.Id)

	case "communication":
		var p CommMsgPayload
		if err := json.Unmarshal(envelope.Payload, &p); err != nil {
			log.Println("Bad communication payload:", err)
			return
		}

		for client := range h.clients {
			select {
			case client.send <- msg:
			default:
				close(client.send)
				delete(h.clients, client)
			}
		}
		fmt.Printf("Action: Broadcast message from %s: %s\n", p.Id, p.Message)

	default:
		log.Printf("Unknown event type: %s\n", envelope.Event)
	}
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

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
			h.handleMessage(message)
		}
	}
}
