package websocket
import (
    "net/http"
    "github.com/gorilla/websocket"
	"github.com/RaphaCosil/messaging-api/internal/model"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketHandler struct {
    Hub *model.Hub
}

func (h *WebSocketHandler) HandleConnection(w http.ResponseWriter, r *http.Request){
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	var username string
	
	_, msg, err := conn.ReadMessage()
	if err != nil {
		return
	}

	username = string(msg)

	client := &Client{conn: conn, username, username}
	h.mu.Lock()
	h.clients[client] = true
	h.mu.Unlock()


	h.broadcast <- model.Message{
		0,

	}

}

func (h *WebSocketHandler) HandleMessage(msg model.GenericMessage) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for client := range h.clients {
		err := client.conn.WriteJSON(msg)
		if err != nil {
			client.conn.Close()
			delete(h.clients, client)
		}
	}
}