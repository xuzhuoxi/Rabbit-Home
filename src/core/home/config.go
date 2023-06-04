// Package home
// Create on 2023/6/4
// @author xuzhuoxi
package home

import (
	"flag"
	"github.com/xuzhuoxi/Rabbit-Home/src/core/conf"
	"github.com/xuzhuoxi/infra-go/filex"
	"github.com/xuzhuoxi/infra-go/osxu"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func initConfig() {
	confPath := flag.String("conf", "", "ServerConfig file for running")
	addr := flag.String("addr", DefaultAddr, "ServerConfig file for running")
	flag.Parse()
	if *confPath != "" {
		err := initConfigWithFile(*confPath)
		if nil != err {
			Logger.Errorln(err)
			panic(err)
		}
		return
	}
	if *addr != "" {
		initConfigWithAddr(*addr)
		return
	}
	initConfigDefault()
}

func initConfigWithFile(filePath string) error {
	if !filex.IsFile(filePath) {
		filePath = filex.Combine(osxu.GetRunningDir(), filePath)
	}
	bs, err := ioutil.ReadFile(filePath)
	if nil != err {
		return err
	}
	cfg := &conf.ServerConfig{}
	err = yaml.Unmarshal(bs, cfg)
	if nil != err {
		return err
	}
	ServerConfig = cfg
	return nil
}

func initConfigDefault() {
	initConfigWithAddr(DefaultAddr)
}

func initConfigWithAddr(addr string) {
	ServerConfig = &conf.ServerConfig{Http: conf.HttpConfig{Addr: addr}}
}
