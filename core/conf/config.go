// Package conf
// Create on 2023/6/4
// @author xuzhuoxi
package conf

import (
	"github.com/xuzhuoxi/infra-go/logx"
)

// HttpConfig http服务配置
// 用于启动http服务
// 服务于Internal与External的功能
type HttpConfig struct {
	Addr string `yaml:"addr"` // 服务器启动监听地址
}

// HomeConfig 配置
// Rabbit-Home的根配置
// 一般对应运行目标下名为config.yaml的文件
type HomeConfig struct {
	Http             HttpConfig        `yaml:"http"`     // http服务配置
	InternalVerifier *InternalVerifier `yaml:"internal"` // 内部Rabbit-Server访问控制配置
	ExternalVerifier *ExternalVerifier `yaml:"external"` // 外网查询访问控制配置
	CfgLog           *logx.CfgLog      `yaml:"log"`      // 日志记录配置
}

// PreProcess 预处理
func (o *HomeConfig) PreProcess() {
	if o.InternalVerifier != nil {
		o.InternalVerifier.PreProcess()
	}
	if o.ExternalVerifier != nil {
		o.ExternalVerifier.PreProcess()
	}
}

// VerifyInternalIP 检查是内部IP的合法性
func (o *HomeConfig) VerifyInternalIP(ipAddr string) bool {
	if len(ipAddr) == 0 {
		return false
	}
	if nil == o.InternalVerifier {
		return true
	}
	return o.InternalVerifier.CheckIpAddr(ipAddr)
}

// VerifyExternalIP 检查外部IP的合法性
func (o *HomeConfig) VerifyExternalIP(ipAddr string) bool {
	if len(ipAddr) == 0 {
		return false
	}
	if len(ipAddr) == 0 {
		return false
	}
	if nil == o.ExternalVerifier {
		return true
	}
	return o.ExternalVerifier.CheckIpAddr(ipAddr)
}
