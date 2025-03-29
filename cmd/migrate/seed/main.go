package main

import (
	"fmt"
	"log"

	"github.com/longlnOff/social/cmd/configuration"
	"github.com/longlnOff/social/internal/db"
	"github.com/longlnOff/social/internal/store"
)

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

	db.Seed(&store)
}
