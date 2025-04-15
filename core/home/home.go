// Package home
// Create on 2023/6/3
// @author xuzhuoxi
package home

import (
	"github.com/xuzhuoxi/Rabbit-Home/core/conf"
	"github.com/xuzhuoxi/infra-go/logx"
	"net/http"
	"sync"
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
	GlobalHomeServer IRabbitHomeServer
	GlobalHomeConfig *conf.HomeConfig
	GlobalLogger     logx.ILogger

	GlobalLock sync.RWMutex
)

type mapHandleInfo struct {
	Pattern string
	Handler func() http.Handler
}

var (
	MapHandlerList []mapHandleInfo
)

func init() {
	GlobalLogger = logx.NewLogger()
	GlobalLogger.SetConfig(logx.LogConfig{Type: logx.TypeConsole, Level: logx.LevelAll})
}

// RegisterMapHandler 注册处理响应器
func RegisterMapHandler(pattern string, newHandler func() http.Handler) {
	MapHandlerList = append(MapHandlerList, mapHandleInfo{Pattern: pattern, Handler: newHandler})
}

// StartHomeServer 启动服务器
func StartHomeServer() {
	if nil == GlobalHomeServer {
		GlobalHomeServer = NewRabbitHomeServer()
		GlobalHomeServer.Init()
	}
	initConfig()
	updateLogger()
	err := GlobalHomeServer.StartByAddr(GlobalHomeConfig.Http.Addr)
	if nil != err {
		panic(err)
	}
}

// StopHomeServer 关闭服务器
func StopHomeServer() {
	if nil == GlobalHomeServer {
		return
	}
	err := GlobalHomeServer.Stop()
	if nil != err {
		panic(err)
	}
}
