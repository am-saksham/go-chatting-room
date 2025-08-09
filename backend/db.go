package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func mustAsset(name string) []byte {
	p := filepath.Join(".", name)
	b, _ := os.ReadFile(p)
	return b
}

func initDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL env var required")
	}
	var err error
	DB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatalf("db ping: %v", err)
	}
	fmt.Println("✅ Connected to Postgres")
	// create tables if missing (simple)
	schema := string(mustAsset("schema.sql"))
	if _, err := DB.Exec(schema); err != nil {
		log.Printf("warning: creating tables failed: %v", err)
	}
}
