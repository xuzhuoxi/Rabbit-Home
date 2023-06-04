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

type sortWeightList []*RegisteredEntity

func (o sortWeightList) Len() int {
	return len(o)
}

func (o sortWeightList) Less(i, j int) bool {
	bi := o[i].IsTimeout()
	bj := o[j].IsTimeout()
	if bi == bj {
		return o[i].State.Weight < o[j].State.Weight
	} else {
		return bj
	}
}

func (o sortWeightList) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

type sortLinkList []*RegisteredEntity

func (o sortLinkList) Len() int {
	return len(o)
}

func (o sortLinkList) Less(i, j int) bool {
	return o[i].Detail.Links < o[j].Detail.Links
}

func (o sortLinkList) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

func initConfig() {
	confPath := flag.String("conf", "", "ServerConfig file for running")
	addr := flag.String("addr", DefaultAddr, "ServerConfig file for running")
	flag.Parse()
	if *confPath != "" {
		err := initConfigWithFile(*confPath)
		if nil != err {
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
