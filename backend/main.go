package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
    "github.com/joho/godotenv"
)

func main() {
	// load env]
    err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Warning: could not load ../.env file")
	}

	// ensure DATABASE_URL and JWT_SECRET exist (or use defaults in getEnv)
	initDB()
	go hub.run()

	r := mux.NewRouter()

	// public
	r.HandleFunc("/signup", signupHandler).Methods("POST")
	r.HandleFunc("/login", loginHandler).Methods("POST")

	// protected
	api := r.PathPrefix("/api").Subrouter()
	api.Use(jwtMiddleware)
	api.HandleFunc("/rooms", createRoomHandler).Methods("POST")
	api.HandleFunc("/rooms", listRoomsHandler).Methods("GET")
	api.HandleFunc("/ws", wsHandler).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("starting server on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}