package handler

import (
    "net/http"
    "sync"
    "time"
    "github.com/gorilla/websocket"
)

type GenericMessage struct {
    Type      string      `json:"type"`
    Content   interface{} `json:"content"`
    Timestamp time.Time   `json:"timestamp"`
}

type Client struct {
    conn     *websocket.Conn
    username string
}

type Hub struct {
    clients   map[*Client]bool
    broadcast chan GenericMessage
    mu        sync.Mutex
}

func NewHub() *Hub {
    return &Hub{
        clients:   make(map[*Client]bool),
        broadcast: make(chan GenericMessage),
    }
}

func (h *Hub) Run() {
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

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

type WebSocketHandler struct {
    Hub *Hub
    mu  sync.Mutex
}

func (h *WebSocketHandler) HandleConnection(w http.ResponseWriter, r *http.Request) {
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

    client := &Client{conn: conn, username: username}
    h.mu.Lock()
    h.Hub.clients[client] = true
    h.mu.Unlock()

    h.Hub.broadcast <- GenericMessage{
        Type:      "join",
        Content:   username,
        Timestamp: time.Now(),
    }
}

func (h *WebSocketHandler) HandleMessage(msg GenericMessage) {
    h.mu.Lock()
    defer h.mu.Unlock()

    for client := range h.Hub.clients {
        err := client.conn.WriteJSON(msg)
        if err != nil {
            client.conn.Close()
            delete(h.Hub.clients, client)
        }
    }
}

func NewWebSocketHandler(hub *Hub) *WebSocketHandler {
    return &WebSocketHandler{
        Hub: hub,
    }
}