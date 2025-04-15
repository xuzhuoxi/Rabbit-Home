// Package conf
// Create on 2025/4/5
// @author xuzhuoxi
package conf

import (
	"github.com/xuzhuoxi/Rabbit-Home/core"
	"github.com/xuzhuoxi/Rabbit-Home/core/conf/verifier"
	"net/http"
)

// InternalVerifier 内部Rabbit-Server实例 服务配置
type InternalVerifier struct {
	Post        bool                          `yaml:"post"`       // 接受的HTTP请求方式
	Timeout     int64                         `yaml:"timeout"`    // 内部Rabbit-Server实例超时设置
	IpVerifier  *verifier.IpVerifier          `yaml:"ip-verify"`  // IP访问控制
	KeyVerifier *verifier.InternalKeyVerifier `yaml:"key-verify"` // 内部Rabbit-Server实例 密钥校验配置
}

// PreProcess 对原始数据进行预处理
func (o *InternalVerifier) PreProcess() {
	if 0 == o.Timeout {
		o.Timeout = core.LinkedTimeout
	}
	if nil != o.IpVerifier {
		o.IpVerifier.PreProcess()
	}
	if nil != o.KeyVerifier {
		o.KeyVerifier.PreProcess()
	}
}

// VerifyPost 检查请求是否为POST请求
// 合格要求：
// 1. 如果配置启用了POST请求，则必须为POST请求
// 2. 如果配置禁用了POST请求，则必须为GET请求
func (o *InternalVerifier) VerifyPost(req *http.Request) bool {
	if o.Post {
		return req.Method == http.MethodPost
	}
	return req.Method == http.MethodGet
}

// CheckIpAddr 检查请求的IP地址是否合法
// 合格要求：
// 1. 如果没有配置相关IP控制信息，则直接通过
// 2. 如果配置了IP控制信息，则根据IpVerifier.CheckIpAddr函数进行判断
func (o *InternalVerifier) CheckIpAddr(ipAddr string) bool {
	if nil == o.IpVerifier {
		return true
	}
	return o.IpVerifier.CheckIpAddr(ipAddr)
}
