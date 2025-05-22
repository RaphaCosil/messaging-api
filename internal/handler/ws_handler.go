package handler

import (
    "net/http"
    "sync"
    "time"
    "log"
    "encoding/json"
    "github.com/gorilla/websocket"
    "github.com/RaphaCosil/messaging-api/internal/model"
    "github.com/RaphaCosil/messaging-api/internal/service"
)

type GenericMessage struct {
    Type      string      `json:"type"`
    Content   json.RawMessage `json:"content"`
    Timestamp time.Time   `json:"timestamp"`
}

type Client struct {
    conn     *websocket.Conn
    userId   uint
    chatId   uint
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
    messageService service.MessageService
}

func NewWebSocketHandler(hub *Hub, messageService service.MessageService) *WebSocketHandler {
    return &WebSocketHandler{
        Hub: hub,
        messageService: messageService,
    }
}

func (h *WebSocketHandler) HandleConnection(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }
    defer conn.Close()

    var username string
    var userId, chatId uint

    _, msg, err := conn.ReadMessage()
    if err != nil {
        return
    }

    var initialData struct {
        Username string `json:"username"`
        UserID   uint    `json:"user_id"`
        ChatID   uint    `json:"chat_id"`
    }

    err = json.Unmarshal(msg, &initialData)
    if err != nil {
        return
    }

    username = initialData.Username
    userId = initialData.UserID
    chatId = initialData.ChatID

    if username == "" || userId == 0 || chatId == 0 {
        conn.WriteMessage(websocket.TextMessage, []byte("Invalid initial data"))
        return
    }

    client := &Client{conn: conn, userId: userId, chatId: chatId, username: username}

    h.mu.Lock()
    if _, ok := h.Hub.clients[client]; ok {
        conn.WriteMessage(websocket.TextMessage, []byte("User already connected"))
        h.mu.Unlock()
        return
    }
    
    h.Hub.clients[client] = true
    h.mu.Unlock()

    content, _ := json.Marshal(username)
    h.Hub.broadcast <- GenericMessage{
        Type:      "join",
        Content:   content,
        Timestamp: time.Now(),
    }

    for {
        _, msgData, err := conn.ReadMessage()
        if err != nil {
            h.mu.Lock()
            delete(h.Hub.clients, client)
            h.mu.Unlock()
            return
        }

        log.Printf("Received message: %s", msgData)

        var msg GenericMessage
        err = json.Unmarshal(msgData, &msg)
        if err != nil {
            continue
        }

        go h.HandleMessage(msg)
    }
}

func (h *WebSocketHandler) HandleMessage(msg GenericMessage) {
    switch msg.Type {
    case "send-message":
        h.Hub.broadcast <- msg

        var message model.Message

        err := json.Unmarshal(msg.Content, &message)
        if err != nil {
            return
        }

        message.SentAt = time.Now()

        h.messageService.Create(message)

    case "update-message":
        h.Hub.broadcast <- msg
        var message model.Message

        err := json.Unmarshal(msg.Content, &message)
        if err != nil {
            log.Println("Error unmarshalling message:", err)
            return
        }

        message.SentAt = time.Now()

        h.messageService.Update(message.MessageID, message)

    case "delete-message":
        h.Hub.broadcast <- msg

        var messageId uint

        err := json.Unmarshal(msg.Content, &messageId)
        if err != nil {
            return
        }

        h.messageService.Delete(messageId)
    }
}