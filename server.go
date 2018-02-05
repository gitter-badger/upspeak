package main

import (
	"log"
	_ "net/http/pprof"

	"net/http"

	"github.com/applait/upspeak/config"
	"github.com/applait/upspeak/models"
	rpcService "github.com/applait/upspeak/rpc"

	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json2"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	Conf := config.InitConfig()

	// Init rpc package and pass in configuration for usage from within the package
	rpcService.InitRpc(Conf)

	// Initiate database connection
	models.ConnectDB(Conf.PG.ConnStr)

	server := rpc.NewServer()
	server.RegisterCodec(json2.NewCodec(), "application/json")

	// RPC Service: Node
	server.RegisterService(new(rpcService.NodeService), "node")

	// RPC Service: User
	server.RegisterService(new(rpcService.UserService), "user")

	http.Handle("/api/v1", server)

	log.Printf("Upspeak rig | Port %s", Conf.Server.Port)
	log.Fatal(http.ListenAndServe(":"+Conf.Server.Port, nil))
}
