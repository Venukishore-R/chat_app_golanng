# Go Chat Room Application

A simple concurrent chat application in Go that allows clients to dynamically join, send messages, and leave the chat room using RESTful HTTP APIs. The core architecture utilizes goroutines and channels to handle multiple clients concurrently and efficiently.

---

## 🧠 Features

- Clients can join and leave a central chat room
- Messages are broadcast to all connected clients
- Fully concurrent design using goroutines and channels
- RESTful API to interact with the chat system
- Message timeout support to avoid blocking

---

## 📦 Project Structure

```
CHAT_APP/
├── cmd/
│   └── main.go                    # Entry point of the application
├── endpoints/
│   └── endpoints.go               # HTTP route definitions
├── handlers/
│   └── handlers.go                # Business logic for each endpoint
├── internal/app/services/
│   ├── client/
│   │   └── client.go              # Client definition and logic
│   └── server/
│       └── server.go              # ChatRoom (server) implementation
├── go.mod                         # Module definition
└── go.sum                         # Dependency checksum
```

---

## ⚙️ Tech Stack

- **Language**: Go (Golang)
- **Concurrency**: Goroutines, Channels
- **Web Server**: net/http

---

## 🚀 Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) (v1.16 or higher)

### Run the App

```bash
go run main.go
```

By default, the server will start at `http://localhost:8080`.

---

## 🔌 API Endpoints

### ➕ Join Chat

- **Endpoint:** `/join?id=<client_id>`
- **Method:** `GET`
- **Description:** Adds a new client to the chat room.

### 💬 Send Message

- **Endpoint:** `/send?id=<client_id>&message=<message>`
- **Method:** `POST` or `GET`
- **Description:** Sends a message from the specified client to the chat room.

### 📥 Get Messages

- **Endpoint:** `/messages?id=<client_id>`
- **Method:** `GET`
- **Description:** Fetches broadcast messages for the specified client.
- **Note:** Supports timeout to avoid indefinite blocking.

### ❌ Leave Chat

- **Endpoint:** `/leave?id=<client_id>`
- **Method:** `GET`
- **Description:** Removes the client from the chat room.

---

## 🧵 Concurrency & Thread-Safety

- All client operations are processed through Go channels.
- Goroutines manage each client’s message delivery.
- The central `ChatRoom` ensures safe access to shared state (like the list of connected clients).
- Proper cleanup is performed when a client leaves.

---

## 🧪 Example Flow

1. Client joins the chat:

   ```
   GET /join?id=user1
   ```

2. Client sends a message:

   ```
   GET /send?id=user1&message=HelloWorld
   ```

3. Client retrieves messages:

   ```
   GET /messages?id=user1
   ```

4. Client leaves the chat:
   ```
   GET /leave?id=user1
   ```

---

## 🎯 Bonus Features

- ✅ `/messages` endpoint supports timeout for 10's to prevent blocking
- ✅ Clients that leave no longer receive broadcast messages

---

## 🤝 Contributing

Feel free to fork this project and submit pull requests for improvements, bug fixes, or additional features!

---

## 📄 License

This project is licensed under the [MIT License](LICENSE).

---

## 🧍 Author

**Venukishore Ramasamy**  
[GitHub Profile](https://github.com/Venukishore-R)
