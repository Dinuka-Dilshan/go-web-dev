package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Dinuka-Dilshan/go-web-dev/internal/db"
	"github.com/Dinuka-Dilshan/go-web-dev/internal/store"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Fatal("cannot find port")
	}

	databaseUrl, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		log.Fatal("cannot find database url")
	}

	config := &config{
		address: port,
		dbConfig: dbConfig{
			address:            databaseUrl,
			maxOpenConnections: 2,
			maxIdleTime:        time.Second * 30,
		},
	}

	db, err := db.New(context.Background(), db.DBConfig{
		Address:         config.dbConfig.address,
		MaxConns:        config.dbConfig.maxOpenConnections,
		MaxConnIdleTime: config.dbConfig.maxIdleTime,
	})

	if err != nil {
		panic(err)
	}
	defer db.Close()

	store := store.NewStorage(db)

	app := &application{
		config: *config,
		store:  *store,
	}

	mux := app.mount()

	log.Fatal(app.run(&mux))
}
