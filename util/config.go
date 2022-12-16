package util

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DB_DRIVER              string
	DB_SOURCE              string
	DB_SOURCE_TEST         string
	DB_USER                string
	DB_NAME                string
	DB_NAME_TEST           string
	SERVER_PORT            string
	JWT_SECRET_KEY         string
	REFRESH_TOKEN_DURATION time.Duration
	ACCESS_TOKEN_DURATION  time.Duration
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		log.Printf("can't read environment variables: %v\n", err)
		return
	}
	err = viper.Unmarshal(&config)
	return
}
