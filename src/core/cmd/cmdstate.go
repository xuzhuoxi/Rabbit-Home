// Package cmd
// Create on 2023/6/3
// @author xuzhuoxi
package cmd

import (
	"github.com/xuzhuoxi/Rabbit-Home/src/core/home"
	"github.com/xuzhuoxi/infra-go/cmdx"
)

// OnCmdState state
func OnCmdState(flagSet *cmdx.FlagSetExtend, args []string) {
	home.Logger.Infoln("Home Server:", home.Server)
}
