package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type createRoomReq struct {
	Name      string `json:"name"`
	IsPrivate bool   `json:"is_private"`
}

func createRoomHandler(w http.ResponseWriter, r *http.Request) {
	var req createRoomReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad req", http.StatusBadRequest)
		return
	}
	uid := r.Context().Value(ctxUserID).(int)
	slug := uuid.New().String()
	var invite sql.NullString
	if req.IsPrivate {
		invite = sql.NullString{String: uuid.New().String(), Valid: true}
	}
	_, err := DB.Exec("INSERT INTO rooms (name,slug,is_private,invite_token,created_by) VALUES ($1,$2,$3,$4,$5)",
		req.Name, slug, req.IsPrivate, invite, uid)
	if err != nil {
		http.Error(w, "err", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"slug": slug, "invite_token": invite.String})
}

func listRoomsHandler(w http.ResponseWriter, r *http.Request) {
	rooms := []Room{}
	if err := DB.Select(&rooms, "SELECT id,name,slug,is_private,invite_token,created_by FROM rooms WHERE is_private = false ORDER BY created_at DESC"); err != nil {
		http.Error(w, "err", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(rooms)
}