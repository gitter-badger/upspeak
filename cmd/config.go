package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// InitConfig reads config and returns a Viper instance
func InitConfig() *viper.Viper {
	v := viper.New()

	// Set config file name depending on the `ENV` env variable
	if strings.ToLower(os.Getenv("ENV")) == "production" {
		v.SetConfigName("upspeak")
	} else {
		v.SetConfigName("upspeak.dev")
	}

	// We will try only JSON for now
	v.SetConfigType("json")

	// Set viper path and read configuration
	v.AddConfigPath("$HOME/.upspeak")
	v.AddConfigPath(".")

	// Defaults
	v.SetDefault("Port", "8080")
	v.SetDefault("PostgresURL", "postgres://root@localhost/upspeak")
	v.SetDefault("Env", "DEV")

	// Bind environment values
	v.BindEnv("Port", "PORT")
	v.BindEnv("PostgresURL", "POSTGRES_URL")
	v.BindEnv("Env", "ENV")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Print("Config file not found.")
		} else {
			// Config file was found but another error was produced
			log.Print("Error in config", err)
		}
	}

	return v
}
