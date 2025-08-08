package main

import (
    "database/sql"
    _ "github.com/lib/pq"
    "fmt"
    "os"
)

var db *sql.DB

func initDB() {
    connStr := os.Getenv("DATABASE_URL")
    var err error

    db, err = sql.Open("postgres", connStr)
    if err != nil {
        panic(err)
    }

    // Optional: check connection
    err = db.Ping()
    if err != nil {
        panic(err)
    }

    fmt.Println("✅ Connected to PostgreSQL!")
}