// Package module
// Create on 2025/4/18
// @author xuzhuoxi
package module

import (
	"github.com/xuzhuoxi/Rabbit-Home/core/home"
	"github.com/xuzhuoxi/Rabbit-Home/core/home/service"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/serialx"
)

func NewModuleServer() serialx.IStartupModule {
	return &moduleServer{}
}

type moduleServer struct {
	eventx.EventDispatcher
}

func (o *moduleServer) Name() string {
	return "Rabbit-Home:Init"
}

func (o *moduleServer) StartModule() {
	o.setDefault()
	o.registerServices()
	o.DispatchEvent(serialx.EventOnStartupModuleStarted, o, nil)
}

func (o *moduleServer) setDefault() {
	home.SetDefaultFuncSortEntity(home.WeightLess)
}
func (o *moduleServer) registerServices() {
	home.RegisterMapHandler(home.PatternLink, service.NewServiceLinkHandler)
	home.RegisterMapHandler(home.PatternUpdate, service.NewServiceUpdateHandler)
	home.RegisterMapHandler(home.PatternUnlink, service.NewServiceUnlinkHandler)
	home.RegisterMapHandler(home.PatternRoute, service.NewServiceRouteHandler)
}

func (o *moduleServer) StopModule() {
	o.DispatchEvent(serialx.EventOnStartupModuleStopped, o, nil)
}

func (o *moduleServer) SaveModule() {
	o.DispatchEvent(serialx.EventOnStartupModuleSaved, o, nil)
}
