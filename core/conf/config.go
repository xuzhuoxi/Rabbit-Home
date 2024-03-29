// Package conf
// Create on 2023/6/4
// @author xuzhuoxi
package conf

import "github.com/xuzhuoxi/infra-go/logx"

const (
	SepIp      = ","
	SepIpRange = "-"
)

type IpControl struct {
	AllowIPs []string `yaml:"allows"`    // IP白名单
	AllowOn  bool     `yaml:"allow_on"`  // 是否启用白名单
	BlockIPs []string `yaml:"blocks"`    // IP黑名单
	BlockOn  bool     `yaml:"blocks_on"` // 是否启用黑名单

	allowIps []*ipGroup
	blockIps []*ipGroup
}

func (o *IpControl) Check(ipAddr string) bool {
	if len(ipAddr) == 0 {
		return false
	}
	group := newIPGroupFromAddr(ipAddr)
	if o.BlockOn && o.contains(o.blockIps, group) {
		return false
	}
	if o.AllowOn && !o.contains(o.allowIps, group) {
		return false
	}
	return true
}

func (o *IpControl) contains(ipGroupArr []*ipGroup, ipGroup *ipGroup) bool {
	if len(ipGroupArr) == 0 || ipGroup == nil {
		return false
	}
	for index := range ipGroupArr {
		if !ipGroupArr[index].ContainsGroup(ipGroup) {
			return false
		}
	}
	return true
}

type HttpConfig struct {
	Addr string `yaml:"addr"` // 服务器启动监听地址
}

type ServerConfig struct {
	Http     HttpConfig   `yaml:"http"`     // Http服务
	Internal *IpControl   `yaml:"rabbit"`   // 内部Ip控制
	External *IpControl   `yaml:"external"` // 外部IP控制
	Timeout  int64        `yaml:"timeout"`  // 超时参数
	CfgLog   *logx.CfgLog `yaml:"log"`      // 日志记录参数
}

func (o *ServerConfig) PreProcess() {
	if o.Internal != nil {
		if o.Internal.AllowOn {
			o.Internal.allowIps = newMultiIPGroupFromAddr(o.Internal.AllowIPs)
		}
		if o.Internal.BlockOn {
			o.Internal.blockIps = newMultiIPGroupFromAddr(o.Internal.BlockIPs)
		}
	}
	if o.External != nil {
		if o.External.AllowOn {
			o.External.allowIps = newMultiIPGroupFromAddr(o.External.AllowIPs)
		}
		if o.External.BlockOn {
			o.External.blockIps = newMultiIPGroupFromAddr(o.External.BlockIPs)
		}
	}
}

func (o *ServerConfig) CheckInternalIP(ipAddr string) bool {
	if nil == o.Internal {
		return true
	}
	return o.Internal.Check(ipAddr)
}

func (o *ServerConfig) CheckExternalIP(ipAddr string) bool {
	if nil == o.External {
		return true
	}
	return o.External.Check(ipAddr)
}
