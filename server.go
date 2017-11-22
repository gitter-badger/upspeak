package main

import (
	"log"
	_ "net/http/pprof"

	"github.com/spf13/viper"
)

func init() {
	// Set viper path and read configuration
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()

	// Handle errors reading the config file
	if err != nil {
		log.Fatalln("Fatal error config file", err)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
