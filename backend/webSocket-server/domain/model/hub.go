package model

import "github.com/gorilla/websocket"

// Hub handles the low-level client state
type Hub struct {
	Clients   map[*websocket.Conn]bool
	Broadcast chan interface{}
}

func NewHub() *Hub {
	return &Hub{
		Clients:   make(map[*websocket.Conn]bool),
		Broadcast: make(chan interface{}),
	}
}
