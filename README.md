# **Messaging API**

Chat API with support for multiple users, multiple rooms, and real-time WebSocket communication.

## Project Structure

```
.
├── cmd/server/main.go        # Entry point
├── internal/                 # Application core
│   ├── db/                   # Database connection
│   ├── handler/              # HTTP handlers
│   ├── model/                # Data models
│   ├── repository/           # Database access layer
│   ├── router/               # HTTP routes setup
│   └── service/              # Business logic layer
```

## API Routes

### **`/user`**

| Method | Route                      | Description           |
| ------ | -------------------------- | --------------------- |
| GET    | `/user`                    | List all users        |
| GET    | `/user/:id`                | Get user by ID        |
| POST   | `/user`                    | Create new user       |
| PUT    | `/user/:id`                | Update user           |
| DELETE | `/user/:id`                | Delete user           |
| GET    | `/user/username/:username` | Find user by username |
| GET    | `/user/chat/:chat_id`      | Find users by chat ID |

### **`/chat`**

| Method | Route                                      | Description                        |
| ------ | ------------------------------------------ | ---------------------------------- |
| GET    | `/chat`                                    | List all chats                     |
| GET    | `/chat/:id`                                | Get chat by ID                     |
| POST   | `/chat`                                    | Create new chat                    |
| PUT    | `/chat/:id`                                | Update chat                        |
| DELETE | `/chat/:id`                                | Delete chat                        |
| GET    | `/chat/user/:user_id`                      | Add user to chat                   |
| DELETE | `/chat/user/:user_id`                      | Remove user from chat              |
| GET    | `/chat/user/:user_id/chat/:chat_id/access` | Check if user has access to a chat |

### **`/ws`**

| Method | Route | Description                   |
| ------ | ----- | ----------------------------- |
| GET    | `/ws` | WebSocket connection endpoint |

## WebSocket Message Types

### **Connect to chat**

```json
{
  "type": "connect-chat",
  "content": {
    "customerd": 1,
    "chatId": 1,
    "username": "name"
  }
}
```

### **Send message**

```json
{
  "type": "send-message",
  "content": {
    "chatId": 1,
    "customerId": 1,
    "content": "Hi world."
  }
}
```

### **Update message**

```json
{
  "type": "update-message",
  "content": {
    "chatId": 1,
    "customerId": 1,
    "content": "Hello world!"
  }
}
```

### **Delete message**

```json
{
  "type": "delete-message",
  "content": 10
}
```

## WebSocket Handler

* Each chat has its own `Hub`.
* `clients` is a map of active WebSocket connections per chat.
* Messages are processed based on their `type`:

  * `join` and `leave` messages are automatically broadcasted.
  * `send-message`, `update-message`, and `delete-message` are broadcasted and persisted.
* Connection starts with a `connect-chat` message.

### Connection Flow

1. Client sends a message with `"type": "connect-chat"`.
2. Server creates or retrieves the `Hub` for the chat.
3. Client is registered in the `Hub`.
4. Server listens to incoming messages and forwards them to all clients in the same chat.

## Setup

1. Create a `.env` file with necessary environment variables.
2. Run the SQL script to create and populate the database tables.
3. (Optional) Import the Postman collection for testing.
4. Start the application:

```bash
go run cmd/server/main.go
```
