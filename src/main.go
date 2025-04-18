// Package main
// Create on 2023/6/3
// @author xuzhuoxi
package main

import (
	"github.com/xuzhuoxi/Rabbit-Home/core/cmd"
	"github.com/xuzhuoxi/Rabbit-Home/src/module"
	"github.com/xuzhuoxi/infra-go/serialx"
)

var (
	StartupManager serialx.IStartupManager
)

func init() {
	StartupManager = serialx.NewStartupManager()
	StartupManager.AppendModule(module.NewModuleLogo())
	StartupManager.AppendModule(module.NewModuleServer())
	StartupManager.AppendModule(module.NewModuleRunner())
}

func main() {
	StartupManager.StartManager()
	cmd.StartCmdListener()
}
