// Package home
// Create on 2023/6/3
// @author xuzhuoxi
package home

import (
	"github.com/xuzhuoxi/Rabbit-Home/src/core/conf"
)

const (
	DefaultPort = 80
	DefaultAddr = ":80"

	PatternLink   = "/link"
	PatternUnlink = "/unlink"
	PatternUpdate = "/update"
	PatternRoute  = "/route"
)

var (
	Server       IRabbitHomeServer
	ServerConfig *conf.ServerConfig
)

// StartHomeServer 启动服务器
func StartHomeServer() {
	if nil == Server {
		Server = NewRabbitHomeServer()
		Server.Init()
	}
	initConfig()
	err := Server.StartByAddr(ServerConfig.Http.Addr)
	if nil != err {
		panic(err)
	}
}

// StopHomeServer 关闭服务器
func StopHomeServer() {
	if nil == Server {
		return
	}
	err := Server.Stop()
	if nil != err {
		panic(err)
	}
}
