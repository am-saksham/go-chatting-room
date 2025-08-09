package main

import (
    "os"
    "testing"
)

func TestInitDB(t *testing.T) {
    // Set a valid test database URL
    os.Setenv("DATABASE_URL", "postgres://user:password@localhost:5433/chatdb?sslmode=disable")

    defer func() {
        if r := recover(); r != nil {
            t.Errorf("initDB panicked: %v", r)
        }
    }()

    initDB()

    if db == nil {
        t.Errorf("Expected db to be initialized, got nil")
    }
}