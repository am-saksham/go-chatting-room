package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte(getEnv("JWT_SECRET", "px1rUdWA4X/08sHi26IjjqYqqqQV+Pb1RAofsIga+Ww="))

func getEnv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Signup handler
func signupHandler(w http.ResponseWriter, r *http.Request) {
	var cred Credentials
	if err := json.NewDecoder(r.Body).Decode(&cred); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if cred.Username == "" || cred.Password == "" {
		http.Error(w, "username and password required", http.StatusBadRequest)
		return
	}

	// hash
	hash, err := bcrypt.GenerateFromPassword([]byte(cred.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	res, err := DB.Exec("INSERT INTO users (username, password_hash) VALUES ($1,$2)", cred.Username, string(hash))
	if err != nil {
		http.Error(w, "username exists or error", http.StatusBadRequest)
		return
	}
	id, _ := res.LastInsertId()
	_ = id
	w.WriteHeader(http.StatusCreated)
}

// Login handler
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var cred Credentials
	if err := json.NewDecoder(r.Body).Decode(&cred); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	var id int
	var hash string
	err := DB.QueryRow("SELECT id, password_hash FROM users WHERE username=$1", cred.Username).Scan(&id, &hash)
	if err == sql.ErrNoRows {
		http.Error(w, "invalid", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(cred.Password)) != nil {
		http.Error(w, "invalid", http.StatusUnauthorized)
		return
	}

	// create token
	claims := jwt.MapClaims{
		"sub": id,
		"usr": cred.Username,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := t.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": ss})
}
