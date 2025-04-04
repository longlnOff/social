package configuration

import (
	"time"

	"github.com/spf13/viper"
)

type Configuration struct {
	Server   ServerConfiguration
	Database DatabaseConfiguration
	Mail     MailConfiguration
}

type MailConfiguration struct {
	EXP        time.Duration `mapstructure:"MAIL_EXP"`
	FROM_EMAIL string        `mapstructure:"FROM_EMAIL"`
	SendGrid   SendGridConfiguration
	MailTrap   MailTrapConfiguration
}

type SendGridConfiguration struct {
	API_KEY string `mapstructure:"SEND_GRID_API_KEY"`
}

type MailTrapConfiguration struct {
	API_KEY string `mapstructure:"MAIL_TRAP_API_KEY"`
}

type ServerConfiguration struct {
	SERVER_ADDRESS   string `mapstructure:"SERVER_ADDRESS"`
	SERVER_PORT      string `mapstructure:"SERVER_PORT"`
	ENVIRONMENT      string `mapstructure:"ENVIRONMENT"`
	VERSION          string `mapstructure:"VERSION"`
	EXTERNAL_ADDRESS string `mapstructure:"EXTERNAL_ADDRESS"`
	EXTERNAL_PORT    string `mapstructure:"EXTERNAL_PORT"`
	FRONTEND_URL     string `mapstructure:"FRONTEND_URL"`
}

type DatabaseConfiguration struct {
	ENGINE            string        `mapstructure:"DB_ENGINE"`
	HOST              string        `mapstructure:"DB_HOST"`
	PORT              string        `mapstructure:"DB_PORT"`
	USER              string        `mapstructure:"DB_USER"`
	PASSWORD          string        `mapstructure:"DB_PASSWORD"`
	DB_NAME           string        `mapstructure:"DB_NAME"`
	DB_MAX_OPEN_CONNS int           `mapstructure:"DB_MAX_OPEN_CONNS"`
	DB_MAX_IDLE_CONNS int           `mapstructure:"DB_MAX_IDLE_CONNS"`
	DB_MAX_IDLE_TIME  time.Duration `mapstructure:"DB_MAX_IDLE_TIME"`
}

func LoadConfig(path string) (cfg Configuration, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return Configuration{}, err
	}

	mail_cfg := MailConfiguration{
		EXP:        viper.GetDuration("MAIL_EXP"),
		FROM_EMAIL: viper.GetString("FROM_EMAIL"),
		SendGrid: SendGridConfiguration{
			API_KEY: viper.GetString("SENDGRID_API_KEY"),
		},
		MailTrap: MailTrapConfiguration{
			API_KEY: viper.GetString("MAILTRAP_API_KEY"),
		},
	}

	server_cfg := ServerConfiguration{
		SERVER_ADDRESS:   viper.GetString("SERVER_ADDRESS"),
		SERVER_PORT:      viper.GetString("SERVER_PORT"),
		ENVIRONMENT:      viper.GetString("ENVIRONMENT"),
		VERSION:          viper.GetString("VERSION"),
		EXTERNAL_ADDRESS: viper.GetString("EXTERNAL_ADDRESS"),
		EXTERNAL_PORT:    viper.GetString("EXTERNAL_PORT"),
		FRONTEND_URL:     viper.GetString("FRONTEND_URL"),
	}

	database_cfg := DatabaseConfiguration{
		ENGINE:            viper.GetString("DB_ENGINE"),
		HOST:              viper.GetString("DB_HOST"),
		PORT:              viper.GetString("DB_PORT"),
		USER:              viper.GetString("DB_USER"),
		PASSWORD:          viper.GetString("DB_PASSWORD"),
		DB_NAME:           viper.GetString("DB_NAME"),
		DB_MAX_OPEN_CONNS: viper.GetInt("DB_MAX_OPEN_CONNS"),
		DB_MAX_IDLE_CONNS: viper.GetInt("DB_MAX_IDLE_CONNS"),
		DB_MAX_IDLE_TIME:  viper.GetDuration("DB_MAX_IDLE_TIME"),
	}

	return Configuration{
		Server:   server_cfg,
		Database: database_cfg,
		Mail:     mail_cfg,
	}, nil
}
