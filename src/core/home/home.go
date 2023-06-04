// Package home
// Create on 2023/6/3
// @author xuzhuoxi
package home

const (
	DefaultPort = 80
	DefaultAddr = ":80"

	PatternRoute  = "/route"
	PatternUpdate = "/update"
	PatternLogin  = "/login"
	PatternLogout = "/logout"
)

var (
	Server IRabbitHomeServer
)

// StartHomeServer 启动服务器
func StartHomeServer() {
	if nil == Server {
		Server = NewRabbitHomeServer()
		Server.Init()
	}
	err := Server.Start()
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
