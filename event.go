package main

import (
	"encoding/json"
	"time"
)

type Event struct {
	Type string `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, c *Client) error

const (
	EventSendMessage = "send message"
	EventNewMessage = "new_message"
	EventChangeRoom = "change_room"
)

type SendMessageEvent struct {
	Message string `json:"message"`
	From   string `json:"from"`
}

type NewMessageEvent struct {
	SendMessageEvent
	Sent time.Time `json:"sent"`
}

type ChangeRoomEvent struct {
	Name string `json:"name"`
}