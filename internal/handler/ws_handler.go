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

type InitialData struct {
    Username  string `json:"username"`
    CustomerID uint   `json:"customerId"`
    ChatID   uint    `json:"chatId"`
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
    hubs           map[uint]*Hub
    mu             sync.Mutex
    messageService service.MessageService
}

func NewWebSocketHandler(hubs map[uint]*Hub, messageService service.MessageService) *WebSocketHandler {
    return &WebSocketHandler{
        hubs:           hubs,
        messageService: messageService,
    }
}

func (h *WebSocketHandler) GetOrCreateHub(chatId uint) *Hub {
    h.mu.Lock()
    defer h.mu.Unlock()

    hub, exists := h.hubs[chatId]
    if !exists {
        hub = NewHub()
        h.hubs[chatId] = hub
        go hub.Run()
        log.Printf("Created new hub for chatId: %d", chatId)
    } else {
        log.Printf("Using existing hub for chatId: %d", chatId)
    }

    return hub
}


func (h *WebSocketHandler) HandleConnection(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        h.HandleError(err, "upgrade", nil)
        return
    }
    defer conn.Close()

    _, msg, err := conn.ReadMessage()
    if err != nil {
        h.HandleError(err, "read", nil)
        return
    }

    var genericMessage GenericMessage
    err = json.Unmarshal(msg, &genericMessage)
    if err != nil {
        log.Println("Error unmarshalling message:", err)
        conn.WriteMessage(websocket.TextMessage, []byte("Invalid message format"))
        return
    }

    client, _, err := h.HandleLogin(genericMessage, conn)
    if err != nil || client == nil {
        return
    }

    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            log.Println("Error reading message:", err)
            h.HandleDisconnect(client)
            return
        }

        log.Printf("Received message: %s", msg)

        var genericMessage GenericMessage
        err = json.Unmarshal(msg, &genericMessage)
        if err != nil {
            log.Println("Error unmarshalling message:", err)
            continue
        }

        go h.HandleMessage(genericMessage)
    }
}

func (h *WebSocketHandler) HandleDisconnect(client *Client) {
    h.mu.Lock()
    if _, ok := h.hubs[client.chatId].clients[client]; ok {
        delete(h.hubs[client.chatId].clients, client)
        h.mu.Unlock()
        content, _ := json.Marshal(client.username)
        h.hubs[client.chatId].broadcast <- GenericMessage{
            Type:      "leave",
            Content:   content,
        }
    } else {
        h.mu.Unlock()
        log.Println("Client not found in hub")
    }
}

func (h *WebSocketHandler) HandleLogin(msg GenericMessage, conn *websocket.Conn) (*Client, *Hub, error) {
    var initialData InitialData

    err := json.Unmarshal(msg.Content, &initialData)
    if err != nil {
        h.HandleError(err, "unmarshal", nil)
        log.Println("Error unmarshalling initial data:", err)
        return nil, nil, err
    }
    if initialData.Username == "" || initialData.CustomerID == 0 || initialData.ChatID == 0 {
        h.HandleError(err, "validation", nil)
        return nil, nil, err
    }

    client := &Client{
        conn:     conn,
        userId:   initialData.CustomerID,
        chatId:   initialData.ChatID,
        username: initialData.Username,
    }

    hub := h.GetOrCreateHub(client.chatId)

    h.mu.Lock()
    defer h.mu.Unlock()


    if _, ok := hub.clients[client]; ok {
        conn.WriteMessage(websocket.TextMessage, []byte("User  already connected"))
        h.HandleError(err, "already connected", client)
        return nil, nil, nil
    }

    hub.clients[client] = true

    content, _ := json.Marshal(client.username)
    hub.broadcast <- GenericMessage{
        Type:      "join",
        Content:   content,
    }

    return client, hub, nil
}


func (h *WebSocketHandler) HandleMessage(msg GenericMessage) {
    switch msg.Type {
    case "send-message":
        var message model.Message

        err := json.Unmarshal(msg.Content, &message)
        if err != nil {
            h.HandleError(err, "unmarshal", nil)
            return
        }
        h.hubs[message.ChatID].broadcast <- msg

        message.SentAt = time.Now()

        h.messageService.Create(message)

    case "update-message":
        var message model.Message

        err := json.Unmarshal(msg.Content, &message)
        if err != nil {
            h.HandleError(err, "unmarshal", nil)
            return
        }
        h.hubs[message.ChatID].broadcast <- msg

        message.SentAt = time.Now()

        h.messageService.Update(message.MessageID, message)

    case "delete-message":
        var message model.Message

        err := json.Unmarshal(msg.Content, &message)
        if err != nil {
            h.HandleError(err, "unmarshal", nil)
            return
        }
        h.hubs[message.ChatID].broadcast <- msg

        h.messageService.Delete(message.MessageID, message.CustomerID)
    }
}

func (h *WebSocketHandler) HandleError(err error, errorType string, client *Client) {
    if err != nil {
        log.Printf("Error [%s]: %v", errorType, err)
        if client != nil {
            conn := client.conn
            if conn != nil {
                conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
            }
        }
    }
}
