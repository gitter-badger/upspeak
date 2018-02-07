package rpc

import (
	"github.com/applait/upspeak/config"
)

type RpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var Conf config.Config

func InitRpc(conf config.Config) {
	Conf = conf
}
