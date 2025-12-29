package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Dinuka-Dilshan/go-web-dev/internal/auth"
	"github.com/Dinuka-Dilshan/go-web-dev/internal/db"
	"github.com/Dinuka-Dilshan/go-web-dev/internal/mailer"
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

	emailApiKey, ok := os.LookupEnv("EMAIL_API_KEY")
	if !ok {
		logger.Fatal("cannot find sendgrid api key")
	}

	fromEmail, ok := os.LookupEnv("FROM_EMAIL")
	if !ok {
		logger.Fatal("cannot find from email")
	}

	jwtSecret, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		logger.Fatal("cannot find jwt secret")
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
			sendGrid: sendGridConfig{
				apiKey:    emailApiKey,
				fromEmail: fromEmail,
			},
		},
		auth: authConfig{
			aud:    "GopherSocial",
			iss:    "GopherSocial",
			secret: jwtSecret,
			exp:    time.Hour * 24 * 3, //3days
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

	mailtrapMailer, err := mailer.NewMailtrapClient(
		config.mail.sendGrid.fromEmail,
		config.mail.sendGrid.apiKey,
	)
	if err != nil {
		panic(err)
	}

	jwtAuthenticator := auth.NewJwtAuthenticator(config.auth.secret, config.auth.aud, config.auth.iss)

	app := &application{
		config:        *config,
		store:         *store,
		logger:        logger,
		mailer:        mailtrapMailer,
		authenticator: &jwtAuthenticator,
	}

	mux := app.mount()

	if err := app.run(&mux); err != nil {
		logger.Fatal("Start failed")
	}
}
