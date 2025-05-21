package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/RaphaCosil/messaging-api/internal/model"
)

func NewHub() *model.Hub {
	return &model.Hub{
		clients:  make(map[*model.Client]bool),
		broadcast: make(chan model.GenericMessage),
	}
}

func (h *model.Hub) Run() {
	for {
		msg := <-h.broadcast
		h.mu.Lock()
		for client := range h.clients {
			err := client.conn.WriteJSON(msg)
			if err != nil {
				client.conn.Close()
				delete(h.clients, client)
			}
		}
		h.mu.Unlock()
	}
}