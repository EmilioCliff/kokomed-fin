package pkg

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	CONFIG_PATH				string		  `mapstructure:"CONFIG_PATH"`
	HTTP_PORT               string        `mapstructure:"HTTP_PORT"`
	MYSQL_USER              string        `mapstructure:"MYSQL_USER"`
	MYSQL_PASSWORD          string        `mapstructure:"MYSQL_PASSWORD"`
	MYSQL_DB                string        `mapstructure:"MYSQL_DB"`
	DB_DSN                  string        `mapstructure:"DB_DSN"`
	MIGRATION_PATH          string        `mapstructure:"MIGRATION_PATH"`
	TOKEN_DURATION          time.Duration `mapstructure:"TOKEN_DURATION"`
	PASSWORD_RESET_DURATION time.Duration `mapstructure:"PASSWORD_RESET_DURATION"`
	REFRESH_TOKEN_DURATION  time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	TOKEN_SYMMETRY_KEY      string        `mapstructure:"TOKEN_SYMMETRY_KEY"`
	PASSWORD_COST           int           `mapstructure:"PASSWORD_COST"`
	RSA_PRIVATE_KEY         string        `mapstructure:"RSA_PRIVATE_KEY"`
	RSA_PUBLIC_KEY          string        `mapstructure:"RSA_PUBLIC_KEY"`
	EMAIL_SENDER_PASSWORD	string	`mapstructure:"EMAIL_SENDER_PASSWORD"`
	EMAIL_SENDER_NAME	string	`mapstructure:"EMAIL_SENDER_NAME"`
	EMAIL_SENDER_ADDRESS	string	`mapstructure:"EMAIL_SENDER_ADDRESS"`
	REDIS_ADDRESS	string	`mapstructure:"REDIS_ADDRESS"`
	REDIS_PASSWORD	string	`mapstructure:"REDIS_PASSWORD"`
}

// Loads app configuration from .env file.
func LoadConfig(path string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, using environment variables")
		} 

		return Config{}, Errorf(INTERNAL_ERROR, "failed to read config: %s", err.Error())
	}

	var config Config
	err := viper.Unmarshal(&config)

	return config, err
}
