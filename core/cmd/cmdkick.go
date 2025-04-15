// Package cmd
// Create on 2023/6/3
// @author xuzhuoxi
package cmd

import (
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Home/core/home"
	"github.com/xuzhuoxi/infra-go/cmdx"
)

const (
	kickId = "id"
)

// OnCmdKick -id=Name
func OnCmdKick(flagSet *cmdx.FlagSetExtend, args []string) {
	id := flagSet.String(kickId, "", "-n=Name")
	flagSet.Parse(args)
	if *id == "" && !flagSet.CheckKey(kickId) {
		fmt.Println("Command \"" + flagSet.Name() + "\" args error!")
		return
	}

	entity, ok := home.GlobalHomeServer.RemoveEntity(*id)
	if !ok || nil == entity {
		fmt.Println(fmt.Sprintf("Kick Entity(%s) fail! Entity is not exist!", *id))
		return
	}
	fmt.Println(fmt.Sprintf("Kick Entity(%s) Succ!", *id))
}
