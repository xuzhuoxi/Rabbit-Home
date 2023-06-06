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
	infoId = "id"
)

// OnCmdInfo info -id=ID
func OnCmdInfo(flagSet *cmdx.FlagSetExtend, args []string) {
	id := flagSet.String(infoId, "", "-id=Id")
	flagSet.Parse(args)
	if *id == "" && !flagSet.CheckKey(infoId) {
		fmt.Println("Command \"" + flagSet.Name() + "\" args error!")
		return
	}
	entity, ok := home.Server.GetEntityById(*id)
	if !ok {
		return
	}
	fmt.Println(entity.DetailString())
}
