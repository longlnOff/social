package main

import (
	"fmt"
	"log"

	"github.com/longlnOff/social/internal/store"
	"github.com/longlnOff/social/internal/db"
	"github.com/spf13/viper"
)

func main() {
	cfg, err := LoadConfig(".")
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

func LoadConfig(path string) (cfg configuration, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return configuration{}, err
	}

	server_cfg := ServerConfiguration{
		SERVER_ADDRESS: viper.GetString("SERVER_ADDRESS"),
		SERVER_PORT:    viper.GetString("SERVER_PORT"),
	}

	database_cfg := DatabaseConfiguration{
		ENGINE:				     viper.GetString("DB_ENGINE"),
		HOST:                    viper.GetString("DB_HOST"),
		PORT:                    viper.GetString("DB_PORT"),
		USER:                    viper.GetString("DB_USER"),
		PASSWORD:                viper.GetString("DB_PASSWORD"),
		DB_NAME:                 viper.GetString("DB_NAME"),
		DB_MAX_OPEN_CONNS:       viper.GetInt("DB_MAX_OPEN_CONNS"),
		DB_MAX_IDLE_CONNS:       viper.GetInt("DB_MAX_IDLE_CONNS"),
		DB_MAX_IDLE_TIME:        viper.GetDuration("DB_MAX_IDLE_TIME"),
	}

	return configuration{
		Server: server_cfg,
		Database: database_cfg,
	}, nil
}
