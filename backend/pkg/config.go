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
	MPESA_CONSUMER_KEY string `mapstructure:"MPESA_CONSUMER_KEY"`
	MPESA_CONSUMER_SECRET string `mapstructure:"MPESA_CONSUMER_SECRET"`
	MPESA_SHORT_CODE string `mapstructure:"MPESA_SHORT_CODE"`
	MPESA_PASSKEY string `mapstructure:"MPESA_PASSKEY"`
}

// Loads app configuration from .env file.
func LoadConfig(path ,name, configType string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType(configType)
	setDefaults()

	viper.AutomaticEnv()

	// if file not found use authomatic env
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, using environment variables")
		} else {
			return Config{}, Errorf(INTERNAL_ERROR, "failed to read config: %s", err.Error())
		}
	}

	var config Config

	return config, viper.Unmarshal(&config)
}

func setDefaults() {
	viper.SetDefault("CONFIG_PATH", "")
	viper.SetDefault("HTTP_PORT", "")
	viper.SetDefault("MYSQL_USER", "")
	viper.SetDefault("MYSQL_PASSWORD", "")
	viper.SetDefault("MYSQL_DB", "")
	viper.SetDefault("DB_DSN", "")
	viper.SetDefault("MIGRATION_PATH", "")
	viper.SetDefault("TOKEN_DURATION", 0)
	viper.SetDefault("PASSWORD_RESET_DURATION", 0)
	viper.SetDefault("REFRESH_TOKEN_DURATION", 0)
	viper.SetDefault("TOKEN_SYMMETRY_KEY", "")
	viper.SetDefault("PASSWORD_COST", 0)
	viper.SetDefault("RSA_PRIVATE_KEY", "")
	viper.SetDefault("RSA_PUBLIC_KEY", "")
	viper.SetDefault("EMAIL_SENDER_PASSWORD", "")
	viper.SetDefault("EMAIL_SENDER_NAME", "")
	viper.SetDefault("EMAIL_SENDER_ADDRESS", "")
	viper.SetDefault("REDIS_ADDRESS", "")
	viper.SetDefault("REDIS_PASSWORD", "")
}
