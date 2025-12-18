package main

import (
	"context"
	"log"
	"os"

	"github.com/Dinuka-Dilshan/go-web-dev/internal/db"
	"github.com/Dinuka-Dilshan/go-web-dev/internal/store"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	databaseUrl, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		log.Fatal("cannot find database url")
	}
	conn, err := db.New(context.Background(), db.DBConfig{Address: databaseUrl, MaxConns: 3})
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	store := store.NewStorage(conn)
	db.Seed(*store, conn)
}
