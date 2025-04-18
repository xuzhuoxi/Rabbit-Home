// Package module
// Create on 2025/4/18
// @author xuzhuoxi
package module

import (
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/serialx"
)

func NewModuleCMD() serialx.IStartupModule {
	return &moduleCMD{}
}

type moduleCMD struct {
	eventx.EventDispatcher
}

func (o *moduleCMD) Name() string {
	return "Rabbit-Home:Init"
}

func (o *moduleCMD) StartModule() {
	o.DispatchEvent(serialx.EventOnStartupModuleStarted, o, nil)
}

func (o *moduleCMD) StopModule() {
	o.DispatchEvent(serialx.EventOnStartupModuleStopped, o, nil)
}

func (o *moduleCMD) SaveModule() {
	o.DispatchEvent(serialx.EventOnStartupModuleSaved, o, nil)
}
