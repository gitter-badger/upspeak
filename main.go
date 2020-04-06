package main

import (
	"github.com/upspeak/upspeak/cmd"
	"github.com/upspeak/upspeak/models"
)

func main() {
	c := cmd.InitConfig()
	models.ConnectDB(c.GetString("PostgresURL"))
}
