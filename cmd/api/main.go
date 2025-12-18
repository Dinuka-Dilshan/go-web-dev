package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Dinuka-Dilshan/go-web-dev/internal/db"
	"github.com/Dinuka-Dilshan/go-web-dev/internal/store"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

//	@title			Gopher Social API
//	@description	social api for ghophers.

//	@contact.name	Gopher API Support

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:3000
// @BasePath	/v1/
func main() {
	zap, err := zap.NewProduction()
	logger := zap.Sugar()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	err = godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env file")
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		logger.Fatal("cannot find port")
	}

	databaseUrl, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		logger.Fatal("cannot find database url")
	}

	config := &config{
		address: port,
		dbConfig: dbConfig{
			address:            databaseUrl,
			maxOpenConnections: 2,
			maxIdleTime:        time.Second * 30,
		},
		apiUrl: "localhost:3000",
		mail: mailConfig{
			exp: time.Hour * 24 * 3,
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
		logger: logger,
	}

	mux := app.mount()

	if err := app.run(&mux); err != nil {
		logger.Fatal("Start failed")
	}
}
