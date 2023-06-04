// Package main
// Create on 2023/6/3
// @author xuzhuoxi
package main

import (
	"github.com/xuzhuoxi/Rabbit-Home/src/core/cmd"
	"github.com/xuzhuoxi/Rabbit-Home/src/core/home"
)

func main() {
	showHomeInfo()
	runRabbitHome()
}

func showHomeInfo() {

}

func runRabbitHome() {
	go home.StartHomeServer()
	cmd.StartCmdListener()
}
