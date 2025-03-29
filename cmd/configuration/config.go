package configuration

import (
	"time"

	"github.com/spf13/viper"
)


type Configuration struct {
	Server ServerConfiguration
	Database DatabaseConfiguration
}

type ServerConfiguration struct {
	SERVER_ADDRESS string	`mapstructure:"SERVER_ADDRESS"`
	SERVER_PORT string		`mapstructure:"SERVER_PORT"`
	ENVIRONMENT string		`mapstructure:"ENVIRONMENT"`
	VERSION string			`mapstructure:"VERSION"`
}

type DatabaseConfiguration struct {
	ENGINE					string			`mapstructure:"DB_ENGINE"`
	HOST 					string			`mapstructure:"DB_HOST"`
	PORT 					string			`mapstructure:"DB_PORT"`
	USER 					string			`mapstructure:"DB_USER"`
	PASSWORD 				string			`mapstructure:"DB_PASSWORD"`
	DB_NAME 				string			`mapstructure:"DB_NAME"`
	DB_MAX_OPEN_CONNS 		int				`mapstructure:"DB_MAX_OPEN_CONNS"`
	DB_MAX_IDLE_CONNS 		int				`mapstructure:"DB_MAX_IDLE_CONNS"`
	DB_MAX_IDLE_TIME 		time.Duration	`mapstructure:"DB_MAX_IDLE_TIME"`
}


func LoadConfig(path string) (cfg Configuration, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return Configuration{}, err
	}

	server_cfg := ServerConfiguration{
		SERVER_ADDRESS: viper.GetString("SERVER_ADDRESS"),
		SERVER_PORT:    viper.GetString("SERVER_PORT"),
		ENVIRONMENT:    viper.GetString("ENVIRONMENT"),
		VERSION:        viper.GetString("VERSION"),
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

	return Configuration{
		Server: server_cfg,
		Database: database_cfg,
	}, nil
}
