package main

import (
	"fmt"
	"log"

	"github.com/longlnOff/social/cmd/configuration"
	"github.com/longlnOff/social/internal/db"
	"github.com/longlnOff/social/internal/store"
)

//	@title			GopherSocial API
//	@description	API for GopherSocial, a social network for gophers.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	LongLN
//	@contact.url	http://www.swagger.io/support
//	@contact.email	longlnofficial@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath	/v1

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				Type "Bearer" followed by a space and then your token
func main() {
	cfg, err := configuration.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	database, err := db.New(
		fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.Database.ENGINE,
			cfg.Database.USER,
			cfg.Database.PASSWORD,
			cfg.Database.HOST,
			cfg.Database.PORT,
			cfg.Database.DB_NAME),
		cfg.Database.DB_MAX_OPEN_CONNS,
		cfg.Database.DB_MAX_IDLE_CONNS,
		cfg.Database.DB_MAX_IDLE_TIME,
	)
	if err != nil {
		log.Panic(err)
	}
	defer database.Close()
	log.Printf("Connected to database: postgres://%s:%s@%s:%s/%s\n", cfg.Database.USER, cfg.Database.PASSWORD, cfg.Database.HOST, cfg.Database.PORT, cfg.Database.DB_NAME)

	store := store.NewStorage(database)

	app := &application{
		configuration: cfg,
		store:         store,
	}

	mux := app.routes()

	log.Fatal(app.run(mux))
}
