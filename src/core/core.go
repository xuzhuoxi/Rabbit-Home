// Package core
// Create on 2023/6/4
// @author xuzhuoxi
package core

import "fmt"

const (
	// CmcKeyState 查看服务器状态
	CmcKeyState = "state"
	// CmcKeyList 查看列表
	CmcKeyList = "list"
	// CmdKeyInfo 查看信息
	CmdKeyInfo = "info"
	// CmdKeyKick 踢除实例
	CmdKeyKick = "kick"
)

var (
	LinkedTimeout int64 // 超时时间
)

// LinkEntity 连接的实例信息
type LinkEntity struct {
	Id         string `json:"id"`      // 实例Id
	PlatformId string `json:"pid"`     // 平台Id
	Name       string `json:"name"`    // 实例名称
	Network    string `json:"network"` // 连接类型
	Addr       string `json:"addr"`    // 连接地址
}

func (o *LinkEntity) String() string {
	return fmt.Sprintf("{Id=%s,PId=%s,Name=%s,Network=%s,Addr=%s}",
		o.Id, o.PlatformId, o.Name, o.Network, o.Addr)
}

// IsNotValid 是否为未验证
func (o *LinkEntity) IsNotValid() bool {
	return len(o.Id) == 0 || len(o.Name) == 0
}

// EntityState 实例状态
type EntityState struct {
	Id     string  `json:"id"`     // 实例Id
	Weight float64 `json:"weight"` // 压力系数
}

func (o *EntityState) String() string {
	return fmt.Sprintf("{Id=%s,Weight=%v}", o.Id, o.Weight)
}

// IsNotValid 是否为未验证
func (o *EntityState) IsNotValid() bool {
	return len(o.Id) == 0
}

// EntityDetailState 实例详细状态
type EntityDetailState struct {
	Id             string `json:"id"`         // 实例Id
	StartTimestamp int64  `json:"start"`      // 启动时间戳(纳秒)
	MaxLinks       uint64 `json:"max_links"`  // 最大连接数
	TotalReqCount  int64  `json:"total_reg"`  // 总请求数
	TotalRespTime  int64  `json:"total_resp"` // 总响应时间
	MaxRespTime    int64  `json:"max_resp"`   // 最大响应时间(纳秒)
	Links          uint64 `json:"links"`      // 连接数

	StatsTimestamp    int64 `json:"sta_start"`    // 统计开始时间戳(纳秒)
	StatsReqCount     int64 `json:"sta_req"`      // 统计请求数
	StatsRespUnixNano int64 `json:"sta_resp"`     // 统计响应时间(纳称)
	StatsInterval     int64 `json:"sta_interval"` // 统计间隔

	Keys string `json:"sta_interval"` // 属性启用标记
}

func (o *EntityDetailState) String() string {
	return fmt.Sprintf("{Id=%s,Start=%v,MaxLink=%v,Link=%v",
		o.Id, o.StartTimestamp, o.MaxLinks, o.Links)
}

// IsNotValid 是否为未验证
func (o *EntityDetailState) IsNotValid() bool {
	return len(o.Id) == 0
}
