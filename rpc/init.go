package rpc

import (
	"github.com/applait/upspeak/config"
)

var Conf config.Config

func InitRpc(conf config.Config) {
	Conf = conf
}
