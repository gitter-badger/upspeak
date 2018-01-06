package main

import (
	"log"
	_ "net/http/pprof"
	"os"

	"net/http"

	"github.com/applait/upspeak/models"
	rpcService "github.com/applait/upspeak/rpc"

	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json2"
	"github.com/spf13/viper"
)

func init() {
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
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Initiate database connection
	models.ConnectDB(viper.GetString("pg.connStr"))

	server := rpc.NewServer()
	server.RegisterCodec(json2.NewCodec(), "application/json")

	// RPC Service: Node
	server.RegisterService(new(rpcService.NodeService), "node")

	// RPC Service: User
	server.RegisterService(new(rpcService.UserService), "user")

	http.Handle("/api/v1", server)

	log.Printf("Upspeak rig | Port %s", viper.GetString("server.port"))
	log.Fatal(http.ListenAndServe(":"+viper.GetString("server.port"), nil))
}
