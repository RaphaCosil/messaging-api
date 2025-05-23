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

type InitialData struct {
    Username string `json:"username"`
    UserID   uint    `json:"user_id"`
    ChatID   uint    `json:"chat_id"`
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
    }
    return hub
}


func (h *WebSocketHandler) HandleConnection(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Error upgrading connection:", err)
        return
    }
    defer conn.Close()

    _, msg, err := conn.ReadMessage()
    if err != nil {
        log.Println("Error reading initial message:", err)
        return
    }

    client, _, err := h.HandleLogin(msg, conn)
    if err != nil || client == nil {
        return
    }

    for {
        _, msgData, err := conn.ReadMessage()
        if err != nil {
            log.Println("Error reading message:", err)
            h.HandleDisconnect(client)
            return
        }

        log.Printf("Received message: %s", msgData)

        var msg GenericMessage
        err = json.Unmarshal(msgData, &msg)
        if err != nil {
            log.Println("Error unmarshalling message:", err)
            continue
        }

        go h.HandleMessage(msg)
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
            Timestamp: time.Now(),
        }
    } else {
        h.mu.Unlock()
        log.Println("Client not found in hub")
    }
}

func (h *WebSocketHandler) HandleLogin(msg []byte, conn *websocket.Conn) (*Client, *Hub, error) {
    var initialData InitialData
    err := json.Unmarshal(msg, &initialData)
    if err != nil {
        log.Println("Error unmarshalling initial data:", err)
        return nil, nil, err
    }
    if initialData.Username == "" || initialData.UserID == 0 || initialData.ChatID == 0 {
        log.Println("Invalid initial data")
        return nil, nil, err
    }

    client := &Client{
        conn:     conn,
        userId:   initialData.UserID,
        chatId:   initialData.ChatID,
        username: initialData.Username,
    }

    hub := h.GetOrCreateHub(client.chatId)

    h.mu.Lock()
    defer h.mu.Unlock()


    if _, ok := hub.clients[client]; ok {
        conn.WriteMessage(websocket.TextMessage, []byte("User  already connected"))
        return nil, nil, nil
    }

    hub.clients[client] = true

    log.Printf("Client %s connected to chatId: %d", client.username, client.chatId)

    content, _ := json.Marshal(client.username)
    hub.broadcast <- GenericMessage{
        Type:      "join",
        Content:   content,
        Timestamp: time.Now(),
    }

    return client, hub, nil
}


func (h *WebSocketHandler) HandleMessage(msg GenericMessage) {
    switch msg.Type {
    case "send-message":
        var message model.Message

        err := json.Unmarshal(msg.Content, &message)
        if err != nil {
            return
        }
        h.hubs[message.ChatID].broadcast <- msg

        message.SentAt = time.Now()

        h.messageService.Create(message)

    case "update-message":
        var message model.Message

        err := json.Unmarshal(msg.Content, &message)
        if err != nil {
            log.Println("Error unmarshalling message:", err)
            return
        }
        h.hubs[message.ChatID].broadcast <- msg

        message.SentAt = time.Now()

        h.messageService.Update(message.MessageID, message)

    case "delete-message":
        var message model.Message

        err := json.Unmarshal(msg.Content, &message)
        if err != nil {
            return
        }
        h.hubs[message.ChatID].broadcast <- msg

        h.messageService.Delete(message.MessageID, message.CustomerID)
    }
}