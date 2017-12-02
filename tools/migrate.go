package main

import (
	"flag"
	"fmt"
	"log"
	_ "net/http/pprof"
	"os"

	"github.com/go-pg/migrations"
	"github.com/go-pg/pg"
	"github.com/spf13/viper"
)

const usageText = `This program runs command on the postgre db. Supported commands are:
  - init - creates pg_migrations table.
  - up - runs all available migrations.
  - down - reverts last migration.
  - reset - reverts all migrations.
  - version - prints current db version.
  - set_version [version] - sets db version without running migrations.

Usage:
go run migrations/foo.go <command> [args]
`

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

	flag.Usage = func() {
		fmt.Printf(usageText)
		flag.PrintDefaults()
		os.Exit(2)
	}

	flag.Parse()

	db := pg.Connect(&pg.Options{
		User:     viper.GetString("pg.user"),
		Database: viper.GetString("pg.database"),
	})

	oldVersion, newVersion, err := migrations.Run(db, flag.Args()...)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if newVersion != oldVersion {
		log.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		log.Printf("version is %d\n", oldVersion)
	}
}
