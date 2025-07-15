# 💬 Lowkey Chatroom

A real-time web-based chatroom built using **Golang**, **WebSockets**, and **OpenSSL**, with a clean, mobile-friendly UI inspired by iMessage. This app supports multi-user messaging across rooms, WebSocket-based communication, user login, and session management.

![screenshot](./screenshot.png)

---

## ⚙️ Features

- 🧠 Built in **Go** using the `net/http` and `gorilla/websocket` packages
- 🔐 Secure login using **OpenSSL**-generated OTPs (JWT/session-less)
- 🔁 Concurrent message handling using **Goroutines**
- 🪢 **WebSocket-based** real-time messaging
- 👥 Multiple chatrooms (rooms can be joined dynamically)
- 📲 Responsive chat UI with iOS-style layout
- 🌙 Dark mode aesthetics
- 🚪 Logout/login flow with UI control

---

## 🧪 Technologies Used

| Layer        | Tech/Library         |
|--------------|----------------------|
| Backend      | Go, Gorilla WebSocket |
| Concurrency  | Goroutines, Channels |
| Security     | OTPs, OpenSSL        |
| Frontend     | HTML, CSS (iOS style), Vanilla JS |
| Deployment   | HTTPS-ready, WebSocket TLS support |

---

## 📁 Project Structure

```
chatapp/
│
├── main.go                  # Go WebSocket + HTTP server
├── manager.go              # Room, login, and socket frontend
├── cert/
│   ├── cert.pem             # SSL certificate (OpenSSL)
│   └── key.pem              # SSL key (OpenSSL)
├── go.mod
├── go.sum
├── gencert.sh               # Script to generate SSL certs
├── .gitignore                # Git ignore file
├── otp.go
├── event.go
├── client.go
├── screenshot.png
├── screenshots/             # Directory for screenshots
│   ├── login.png            # Login screenshot
│   └── chat.png             # Chat view screenshot
├── frontend/
│   └── index.html           # Frontend HTML file
└── README.md
```

---

## 🚀 Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/thatquietkid/websockets-project-golang.git
cd websockets-project-golang
```

### 2. Generate SSL Certificate (if not already)

```bash
bash gencert.sh
```

> Use localhost (127.0.0.1) as Common Name when prompted.

### 3. Run the Server

```bash
go run main.go
```

> The server starts on `https://localhost:8080`.

---

## 🧪 Sample Users

> You can mock login with any username/password combo. OTP or JWT-based login expected in production.
```
Current user: admin
Password: admin1234
```

---

## 🧵 How It Works

- User logs in via a login form
- Server issues a one-time token (OTP) or JWT (via `/login`)
- A WebSocket is established at `/ws?otp=...`
- All messages are routed through WebSocket events (`send message`, `new_message`, `change_room`)
- The frontend handles dynamic rendering and displays messages in real-time

---

## Future Improvements

- Persistent storage (messages are in-memory only for now!)
- Real authentication backend (mock login for now!)
- Rate limiting

---

## 📸 Screenshots

| Login | Chat View |
|-------|-----------|
| ![](./screenshots/login.png) | ![](./screenshots/chat.png) |

---

## 🤝 Contributions

Feel free to fork the repo, file issues, and submit PRs. All feedback is welcome!

---

## 📜 License

MIT License © 2025 [Nitin Chauhan](LICENSE)