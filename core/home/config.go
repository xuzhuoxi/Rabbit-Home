// Package home
// Create on 2023/6/4
// @author xuzhuoxi
package home

import (
	"flag"
	"github.com/xuzhuoxi/Rabbit-Home/core"
	"github.com/xuzhuoxi/Rabbit-Home/core/conf"
	"github.com/xuzhuoxi/infra-go/filex"
	"github.com/xuzhuoxi/infra-go/logx"
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
	if 0 == cfg.Timeout {
		cfg.Timeout = core.LinkedTimeout
	}
	ServerConfig = cfg
	return nil
}

func initConfigDefault() {
	initConfigWithAddr(DefaultAddr)
}

func initConfigWithAddr(addr string) {
	ServerConfig = &conf.ServerConfig{Http: conf.HttpConfig{Addr: addr}, Timeout: core.LinkedTimeout}
}

func updateLogger() {
	logCfg := ServerConfig.Logger
	if nil == logCfg {
		return
	}
	Logger.RemoveConfig(logx.TypeConsole)
	Logger.SetConfig(logx.LogConfig{Type: logCfg.Type, Level: logCfg.Level,
		FilePath: filex.Combine(osxu.GetRunningDir(), logCfg.FilePath), MaxSize: logCfg.Max})
}
