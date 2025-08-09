package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // in prod: validate origin
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// auth context
	uid := r.Context().Value(ctxUserID)
	if uid == nil {
		http.Error(w, "unauth", http.StatusUnauthorized)
		return
	}
	userID := uid.(int)
	room := r.URL.Query().Get("room")
	if room == "" {
		room = "world"
	}

	// TODO: check room private + invite token logic here

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	client := &Client{conn: conn, send: make(chan []byte, 256), userID: userID, room: room}
	hub.register <- client

	go client.writePump()
	client.readPump()
}