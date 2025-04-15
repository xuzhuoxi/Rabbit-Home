// Package verifier
// Create on 2025/2/12
// @author xuzhuoxi
package verifier

import (
	"errors"
	"fmt"
	"github.com/xuzhuoxi/infra-go/mathx"
	"net"
	"strings"
)

const (
	SepIp      = "."
	SepIp6     = ":"
	SepIpRange = "-"
)

func newIPRangeFromAddr(addr string) *ipRange {
	rs := &ipRange{}
	rs.FromStringOverride(addr)
	return rs
}

func newMultiIPRangeFromAddr(addrArr []string) []*ipRange {
	var rs []*ipRange
	for _, addr := range addrArr {
		rs = append(rs, newIPRangeFromAddr(addr))
	}
	return rs
}

// ipRange ---------- ---------- ---------- ---------- ---------- ----------

type ipRange struct {
	MinIP net.IP
	MaxIP net.IP
}

// ContainsAddr 检测是否包含ip地址
// 如果 ipAddr 为不合法,则返回false
func (o *ipRange) ContainsAddr(ipAddr string) bool {
	ip := net.ParseIP(ipAddr)
	if nil == ip {
		return false
	}
	return o.ContainsIP6(ip.To16())
}

// ContainsIP6 检测是否包含ip地址
// 如果 ip6 为 nil,则返回false
func (o *ipRange) ContainsIP6(ip6 net.IP) bool {
	if nil == ip6 {
		return false
	}
	return o.containsIP(ip6)
}

// FromStringOverride 把字符串表示的ip地址组解释为ipGroup结构
// 支持格式:
//  1. 单个ipv4或ipv6地址: 192.168.0.1 or 2001:0db8:85a3:0000:0000:8a2e:0370:7334 or 2001:db8:85a3::8a2e:370:7334
//  2. 多个ipv4或ipv6地址(最后个一组使用'-'表示范围): 192.168.0.1-255 or 2001:0db8:85a3:0000:0000:8a2e:0370:0-7334
func (o *ipRange) FromStringOverride(ipAddr string) error {
	if len(ipAddr) == 0 {
		return errors.New("ipAddr is empty. ")
	}

	idx := strings.LastIndex(ipAddr, SepIpRange)
	if -1 != idx { // 范围表示
		return o.fromIpRange(ipAddr, idx)
	} else { // 单个表示
		return o.fromIp(ipAddr)
	}
}

func (o *ipRange) fromIpRange(ipAddr string, idxRange int) error {
	ipMin := net.ParseIP(ipAddr[:idxRange])
	if nil == ipMin {
		return errors.New(fmt.Sprintf("ipAddr format error. %s", ipAddr))
	}
	dot4 := strings.LastIndex(ipAddr, SepIp)
	dot6 := strings.LastIndex(ipAddr, SepIp6)
	if -1 == dot4 && -1 == dot6 {
		return errors.New(fmt.Sprintf("ipAddr format error. %s", ipAddr))
	}
	dot := mathx.MaxInt(dot4, dot6)
	ipMax := net.ParseIP(ipAddr[:dot+1] + ipAddr[idxRange+1:])
	if nil == ipMax {
		return errors.New(fmt.Sprintf("ipAddr format error. %s", ipAddr))
	}
	o.MinIP, o.MaxIP = ipMin.To16(), ipMax.To16()
	return nil
}

func (o *ipRange) fromIp(ipAddr string) error {
	ip := net.ParseIP(ipAddr)
	if nil == ip {
		return errors.New(fmt.Sprintf("ipAddr format error. %s", ipAddr))
	}
	o.MinIP = ip.To16()
	o.MaxIP = o.MinIP
	return nil
}

func (o *ipRange) containsIP(ip net.IP) bool {
	if nil == ip {
		return false
	}
	for idx := 16 - 1; idx >= 0; idx-- {
		if ip[idx] < o.MinIP[idx] || ip[idx] > o.MaxIP[idx] {
			return false
		}
	}
	return true
}

// IpVerifier ---------- ---------- ---------- ---------- ---------- ----------

type IpVerifier struct {
	AllowOn  bool     `yaml:"allow-on"`  // 是否启用白名单
	AllowIPs []string `yaml:"allows"`    // IP白名单
	BlockOn  bool     `yaml:"blocks-on"` // 是否启用黑名单
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
