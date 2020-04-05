package main

import (
	"github.com/upspeak/upspeak/cmd"
	"github.com/upspeak/upspeak/upspeak"
)

func main() {
	c := cmd.InitConfig()
	upspeak.ConnectDB(c.GetString("PostgresURL"))
}
