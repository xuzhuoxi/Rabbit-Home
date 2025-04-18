// Package module
// Create on 2025/4/18
// @author xuzhuoxi
package module

import (
	"github.com/xuzhuoxi/infra-go/eventx"
	"github.com/xuzhuoxi/infra-go/serialx"
)

func NewModuleLogo() serialx.IStartupModule {
	return &moduleLogo{}
}

type moduleLogo struct {
	eventx.EventDispatcher
}

func (o *moduleLogo) Name() string {
	return "Rabbit-Home:Init"
}

func (o *moduleLogo) StartModule() {
	o.DispatchEvent(serialx.EventOnStartupModuleStarted, o, nil)
}

func (o *moduleLogo) StopModule() {
	o.DispatchEvent(serialx.EventOnStartupModuleStopped, o, nil)
}

func (o *moduleLogo) SaveModule() {
	o.DispatchEvent(serialx.EventOnStartupModuleSaved, o, nil)
}
