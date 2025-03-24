package main

import (
	"log"

	"github.com/spf13/viper"
)

func main() {
	cfg, err := LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	app := &application{
		configuration: cfg,
	}

	mux := app.routes()

	log.Fatal(app.run(mux))
}

func LoadConfig(path string) (cfg configuration,err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return configuration{}, err
	}

	if err = viper.Unmarshal(&cfg); err != nil {
		return configuration{}, err
	}

	return cfg, nil
}
