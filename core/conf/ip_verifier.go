// Package conf
// Create on 2025/2/12
// @author xuzhuoxi
package conf

import (
	"net"
	"net/http"
)

type IpVerifier struct {
	Post     bool     `yaml:"post"`      // 接受的HTTP请求方式
	AllowOn  bool     `yaml:"allow_on"`  // 是否启用白名单
	AllowIPs []string `yaml:"allows"`    // IP白名单
	BlockOn  bool     `yaml:"blocks_on"` // 是否启用黑名单
	BlockIPs []string `yaml:"blocks"`    // IP黑名单

	allowIps []*ipRange
	blockIps []*ipRange
}

// PreProcess 对原始数据进行预处理
func (o *IpVerifier) PreProcess() {
	if len(o.AllowIPs) > 0 {
		o.allowIps = newMultiIPRangeFromAddr(o.AllowIPs)
	}
	if len(o.BlockIPs) > 0 {
		o.blockIps = newMultiIPRangeFromAddr(o.BlockIPs)
	}
}

// CheckIpAddr 检查IP地址是否合格
// 合格要求：
// 1. 不为空
// 2. 合法IP地址
// 3. 如果启用了黑名单，则必须不属于黑名单
// 4. 如果启用了白名单，则必须属于白名单
func (o *IpVerifier) CheckIpAddr(ipAddr string) bool {
	if len(ipAddr) == 0 {
		return false
	}
	ip := net.ParseIP(ipAddr)
	if nil == ip {
		return false
	}
	ip6 := ip.To16()
	if o.BlockOn && o.contains(o.blockIps, ip6) {
		return false
	}
	if o.AllowOn && !o.contains(o.allowIps, ip6) {
		return false
	}
	return true
}

// VerifyPost 检查请求是否为POST请求
// 合格要求：
// 1. 如果配置启用了POST请求，则必须为POST请求
// 2. 如果配置禁用了POST请求，则必须为GET请求
func (o *IpVerifier) VerifyPost(req *http.Request) bool {
	if o.Post {
		return req.Method == http.MethodPost
	}
	return req.Method == http.MethodGet
}

func (o *IpVerifier) contains(ipGroupArr []*ipRange, ip6 net.IP) bool {
	if len(ipGroupArr) == 0 {
		return false
	}
	for index := range ipGroupArr {
		if ipGroupArr[index].ContainsIP6(ip6) {
			return true
		}
	}
	return false
}
