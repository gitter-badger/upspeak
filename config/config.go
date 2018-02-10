package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string
	}
	PG struct {
		ConnStr string
	}
	JWT struct {
		Issuer   string
		Audience string
		Secret   string
	}
}

// InitConfig unmarshals config into config struct
func InitConfig() Config {
	// Set viper path and read configuration
	viper.AddConfigPath(".")
	if os.Getenv("ENV") == "PRODUCTION" {
		viper.SetConfigName("config")
	} else {
		viper.SetConfigName("devconfig")
	}
	err := viper.ReadInConfig()

	// Handle errors reading the config file
	if err != nil {
		log.Fatalln("Fatal error config file", err)
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Unable to decode config into struct, %v", err)
	}

	return config
}