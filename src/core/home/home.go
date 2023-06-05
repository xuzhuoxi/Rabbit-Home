// Package home
// Create on 2023/6/3
// @author xuzhuoxi
package home

import (
	"github.com/xuzhuoxi/Rabbit-Home/src/core/conf"
	"github.com/xuzhuoxi/infra-go/logx"
)

const (
	DefaultPort = 80
	DefaultAddr = ":80"

	PatternLink   = "/link"
	PatternUnlink = "/unlink"
	PatternUpdate = "/update"
	PatternRoute  = "/route"

	PatternDataKey = "data"
)

var (
	Server       IRabbitHomeServer
	ServerConfig *conf.ServerConfig
	Logger       logx.ILogger
)

func init() {
	Logger = logx.NewLogger()
	Logger.SetConfig(logx.LogConfig{Type: logx.TypeConsole, Level: logx.LevelAll})
}

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
