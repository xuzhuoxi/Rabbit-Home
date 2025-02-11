// Package conf
// Create on 2023/6/4
// @author xuzhuoxi
package conf

import (
	"github.com/xuzhuoxi/infra-go/logx"
)

const (
	SepIp      = "."
	SepIp6     = ":"
	SepIpRange = "-"
)

type HttpConfig struct {
	Addr string `yaml:"addr"` // 服务器启动监听地址
}

type ServerConfig struct {
	Http     HttpConfig   `yaml:"http"`     // Http服务
	Internal *IpVerifier  `yaml:"internal"` // 内部Ip控制
	External *IpVerifier  `yaml:"external"` // 外部IP控制
	Timeout  int64        `yaml:"timeout"`  // 超时参数
	CfgLog   *logx.CfgLog `yaml:"log"`      // 日志记录参数
}

// PreProcess 预处理
func (o *ServerConfig) PreProcess() {
	if o.Internal != nil {
		o.Internal.PreProcess()
	}
	if o.External != nil {
		o.External.PreProcess()
	}
}

// VerifyInternalIP 检查是否为内部IP
func (o *ServerConfig) VerifyInternalIP(ipAddr string) bool {
	if len(ipAddr) == 0 {
		return false
	}
	if nil == o.Internal {
		return true
	}
	return o.Internal.CheckIpAddr(ipAddr)
}

// VerifyExternalIP 检查是否为外部IP
func (o *ServerConfig) VerifyExternalIP(ipAddr string) bool {
	if len(ipAddr) == 0 {
		return false
	}
	if len(ipAddr) == 0 {
		return false
	}
	if nil == o.External {
		return true
	}
	return o.External.CheckIpAddr(ipAddr)
}
