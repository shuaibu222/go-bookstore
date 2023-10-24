package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	MongoUsername string `mapstructure:"MONGO_INITDB_ROOT_USERNAME"`
	MongoPassword string `mapstructure:"MONGO_INITDB_ROOT_PASSWORD"`

	WebPort   string `mapstructure:"WEB_PORT"`
	JWTSecret string `mapstructure:"JWT_SECRET"`
}

func LoadConfig() (*Config, error) {
	env := Config{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	return &env, nil
}
