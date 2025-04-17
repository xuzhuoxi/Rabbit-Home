// Package main
// Create on 2023/6/3
// @author xuzhuoxi
package main

import (
	"github.com/xuzhuoxi/Rabbit-Home/core/cmd"
	"github.com/xuzhuoxi/Rabbit-Home/core/home"
	"github.com/xuzhuoxi/Rabbit-Home/core/home/service"
)

func main() {
	setSortFunc()
	showHomeInfo()
	registerHomeServices()
	runRabbitHome()
}

func showHomeInfo() {
}

func registerHomeServices() {
	home.RegisterMapHandler(home.PatternLink, service.NewServiceLinkHandler)
	home.RegisterMapHandler(home.PatternUpdate, service.NewServiceUpdateHandler)
	home.RegisterMapHandler(home.PatternUnlink, service.NewServiceUnlinkHandler)
	home.RegisterMapHandler(home.PatternRoute, service.NewServiceRouteHandler)
}

func setSortFunc() {
	home.SetDefaultFuncSortEntity(home.WeightLess)
}

func runRabbitHome() {
	go home.StartHomeServer()
	cmd.StartCmdListener()
}
