package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
	"encoding/json"

	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		CheckOrigin: checkOrigin,
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
	}
)

type Manager struct {
	clients ClientList
	sync.RWMutex
	otps RetentionMap
	handlers map[string]EventHandler
}

func NewManager(ctx context.Context) *Manager {
	m := &Manager{
		clients: make(ClientList),
		handlers: make(map[string]EventHandler),
		otps: NewRetentionMap(ctx, 5 * time.Second),
	}
	m.setupEventHandlers()
	return m
}

func (m *Manager) setupEventHandlers() {
	m.handlers[EventSendMessage] = SendMessage
	m.handlers[EventChangeRoom] = ChatRoomHandler
}

func ChatRoomHandler(event Event, c *Client) error {
	var changeRoomEvent ChangeRoomEvent
	if err := json.Unmarshal(event.Payload, &changeRoomEvent); err != nil {
		return fmt.Errorf("error unmarshaling event payload: %w", err)
	}

	c.chatroom = changeRoomEvent.Name
	return nil
}

func SendMessage(event Event, c *Client) error {
	var chatevent SendMessageEvent
	if err := json.Unmarshal(event.Payload, &chatevent); err != nil {
		return fmt.Errorf("error unmarshaling event payload: %w", err)
	}

	var broadMessage NewMessageEvent
	broadMessage.Sent = time.Now()
	broadMessage.From = chatevent.From
	broadMessage.Message = chatevent.Message
	data, err := json.Marshal(broadMessage)
	if err != nil {
		return fmt.Errorf("error marshaling broadcast message: %w", err)
	}
	var outgoingEvent = Event{
		Payload: data,
		Type: EventNewMessage,
	}

	for client := range c.manager.clients {
		if client.chatroom == "" {
			continue // skip clients who never joined a room
		}
		if client.chatroom == c.chatroom {
			client.egress <- outgoingEvent
		}
	}
	return nil
}

func (m *Manager) routeEvent(event Event, c *Client) error {
	if handler, exists := m.handlers[event.Type]; exists {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	}else {
		return errors.New("there is no such event type")
	}
}

func (m *Manager) serveWS(w http.ResponseWriter, r *http.Request) {
	otp := r.URL.Query().Get("otp")
	if otp == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if !m.otps.VerifyOTP(otp) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	log.Println("new connection")

	/// Upgrade the HTTP connection to a WebSocket connection
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error upgrading connection:", err)
		return
	}

	client := NewClient(conn, m)
	m.addClient(client)

	// Start client processes
	go client.readMessages()
	go client.writeMessages()
}

func (m *Manager) LoginHandler(w http.ResponseWriter, r *http.Request) {
	type userLoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var req userLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}
	if req.Username == "admin" || req.Password == "admin1234" {
		type response struct{
			OTP string `json:"otp"`
		}
		otp := m.otps.NewOTP()
		resp := response{OTP: otp.Key}
		data, err := json.Marshal(resp)
		if err != nil {
			log.Println("error marshaling response:", err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	m.clients[client] = true
	log.Println("client added, total clients:", len(m.clients))
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	if _, exists := m.clients[client]; exists {
		client.connection.Close()
		delete(m.clients, client)
		log.Println("client removed, total clients:", len(m.clients))
	}
}

func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")

	switch origin {
		case "https://localhost:8080":
			return true
		default:
			log.Println("invalid origin:", origin)
			return false
	}
}