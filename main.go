package main

import (
	"github.com/upspeak/upspeak/cmd"
	"github.com/upspeak/upspeak/core"
)

func main() {
	c := cmd.InitConfig()
	core.ConnectDB(c.GetString("PostgresURL"))
}
