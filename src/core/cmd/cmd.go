// Package cmd
// Create on 2023/6/3
// @author xuzhuoxi
package cmd

import (
	"github.com/xuzhuoxi/Rabbit-Home/src/core"
	"github.com/xuzhuoxi/infra-go/cmdx"
)

func StartCmdListener() {
	cmdLine := cmdx.CreateCommandLineListener("请输入命令：", 0)
	cmdLine.MapCommand(core.CmcKeyState, OnCmdState)
	cmdLine.MapCommand(core.CmcKeyList, OnCmdList)
	cmdLine.MapCommand(core.CmdKeyInfo, OnCmdInfo)
	cmdLine.MapCommand(core.CmdKeyKick, OnCmdKick)

	cmdLine.StartListen() //这里会发生阻塞，保证程序不会结束
}
