// Package home
// Create on 2023/6/4
// @author xuzhuoxi
package home

import (
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Home/core"
	"strings"
	"time"
)

var (
	EntityEmpty = RegisteredEntity{}
)

func NewRegisteredEntity(entity core.LinkEntity) *RegisteredEntity {
	rs := &RegisteredEntity{LinkEntity: entity, lastUpdateNano: time.Now().UnixNano()}
	rs.State.Id = rs.Id
	return rs
}

// RegisteredEntity 已注册实例
type RegisteredEntity struct {
	core.LinkEntity
	State  core.EntityStatus       // 实例简单状态
	Detail core.EntityDetailStatus // 实例详细状态

	lastUpdateNano int64 // 上一次更新时间戳
	hit            int   // 命中次数
}

func (o *RegisteredEntity) String() string {
	return fmt.Sprintf("{Base=%s, State=%s}", o.LinkEntity.String(), o.State.String())
}

func (o *RegisteredEntity) DetailString() string {
	return fmt.Sprintf("{Base=%s, State=%s, Detail=%s}", o.LinkEntity.String(), o.State.String(), o.Detail.String())
}

// IsTimeout 是否已经超时
func (o *RegisteredEntity) IsTimeout() bool {
	if core.LinkedTimeout <= 0 {
		return false
	}
	return (time.Now().UnixNano() - o.lastUpdateNano) >= core.LinkedTimeout
}

// UpdateState 更新实例状态信息
func (o *RegisteredEntity) UpdateState(state core.EntityStatus) {
	o.State.Weight = state.Weight
}

// UpdateDetailState 更新实例详细状态信息
func (o *RegisteredEntity) UpdateDetailState(detail core.EntityDetailStatus) {
	if o.Id != detail.Id || len(detail.Keys) == 0 {
		return
	}
	keys := strings.Split(detail.Keys, ",")
	for idx := range keys {
		switch keys[idx] {
		case "start":
			o.Detail.StartTimestamp = detail.StartTimestamp
		case "max_links":
			o.Detail.MaxLinks = detail.MaxLinks
		case "total_reg":
			o.Detail.TotalReqCount = detail.TotalReqCount
		case "total_resp":
			o.Detail.TotalRespTime = detail.TotalRespTime
		case "max_resp":
			o.Detail.MaxRespTime = detail.MaxRespTime
		case "links":
			o.Detail.Links = detail.Links

		case "sta_start":
			o.Detail.StatsTimestamp = detail.StatsTimestamp
		case "sta_req":
			o.Detail.StatsReqCount = detail.StatsReqCount
		case "sta_resp":
			o.Detail.StatsRespUnixNano = detail.StatsRespUnixNano
		case "sta_interval":
			o.Detail.StatsInterval = detail.StatsInterval
		}
	}
}

// AddHit 增加命中次数
func (o *RegisteredEntity) AddHit() {
	o.hit += 1
}
