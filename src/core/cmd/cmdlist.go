// Package cmd
// Create on 2023/6/3
// @author xuzhuoxi
package cmd

import (
	"github.com/xuzhuoxi/Rabbit-Home/src/core/home"
	"github.com/xuzhuoxi/infra-go/cmdx"
)

const (
	listName       = "name"
	listOn         = "on"
	listPlatformId = "pid"
)

// OnCmdList list
// list name=Name pid = PID
func OnCmdList(flagSet *cmdx.FlagSetExtend, args []string) {
	name := flagSet.String(listName, "", "-name=Name")
	on := flagSet.Bool(listOn, false, "-on=[true|false]")
	pId := flagSet.String(listPlatformId, "", "-pid=PlatformId")
	flagSet.Parse(args)
	bName := flagSet.CheckKey(listName)
	bOn := flagSet.CheckKey(listOn)
	bPId := flagSet.CheckKey(listPlatformId)

	list := home.Server.GetEntityList()
	home.Logger.Infoln("999", list)
	entities := list.GetEntities(func(each home.RegisteredEntity) bool {
		if bName && each.Name != *name {
			return false
		}
		if bOn && each.IsTimeout() != *on {
			return false
		}
		if bPId && each.PlatformId == *pId {
			return false
		}
		return true
	})
	if len(entities) == 0 {
		home.Logger.Infoln("0000000")
		return
	}
	for _, entity := range entities {
		home.Logger.Infoln(entity)
	}
}
