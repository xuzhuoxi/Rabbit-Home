// Package module
// Create on 2025/4/18
// @author xuzhuoxi
package module

import (
	"github.com/xuzhuoxi/Rabbit-Home/core/home"
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/serialx"
)

func NewModuleRunner() serialx.IStartupModule {
	return &moduleRunner{}
}

type moduleRunner struct {
	eventx.EventDispatcher
}

func (o *moduleRunner) Name() string {
	return "Rabbit-Home:Init"
}

func (o *moduleRunner) StartModule() {
	go home.StartHomeServer()
	o.DispatchEvent(serialx.EventOnStartupModuleStarted, o, nil)
}

func (o *moduleRunner) StopModule() {
	o.DispatchEvent(serialx.EventOnStartupModuleStopped, o, nil)
}

func (o *moduleRunner) SaveModule() {
	o.DispatchEvent(serialx.EventOnStartupModuleSaved, o, nil)
	return
}
