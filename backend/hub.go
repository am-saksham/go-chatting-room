package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn   *websocket.Conn
	send   chan []byte
	userID int
	room   string
}

type Hub struct {
	rooms      map[string]map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan broadcastMsg
}

type broadcastMsg struct {
	Room string
	Data []byte
}

var hub = &Hub{
	rooms:      make(map[string]map[*Client]bool),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	broadcast:  make(chan broadcastMsg),
}

func (h *Hub) run() {
	for {
		select {
		case c := <-h.register:
			if h.rooms[c.room] == nil {
				h.rooms[c.room] = make(map[*Client]bool)
			}
			h.rooms[c.room][c] = true
		case c := <-h.unregister:
			if _, ok := h.rooms[c.room][c]; ok {
				delete(h.rooms[c.room], c)
				close(c.send)
			}
		case m := <-h.broadcast:
			for cl := range h.rooms[m.Room] {
				select {
				case cl.send <- m.Data:
				default:
					close(cl.send)
					delete(h.rooms[m.Room], cl)
				}
			}
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(60 * time.Second)); return nil })
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("ws read err: %v", err)
			}
			break
		}
		// attach user id and timestamp on server side - broadcast JSON
		payload := map[string]interface{}{
			"user_id": c.userID,
			"room":    c.room,
			"content": string(msg),
			"ts":      time.Now().UTC(),
		}
		data, _ := json.Marshal(payload)
		// persist message optionally, e.g., DB.Exec(...)
		hub.broadcast <- broadcastMsg{Room: c.room, Data: data}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.WriteMessage(websocket.TextMessage, msg)
		case <-ticker.C:
			c.conn.WriteMessage(websocket.PingMessage, []byte{})
		}
	}
}