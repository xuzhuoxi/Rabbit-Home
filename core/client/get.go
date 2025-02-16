// Package client
// Create on 2023/6/16
// @author xuzhuoxi
package client

import (
	"encoding/base64"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/xuzhuoxi/Rabbit-Home/core"
	"github.com/xuzhuoxi/Rabbit-Home/core/home"
	"github.com/xuzhuoxi/infra-go/netx/httpx"
)

// LinkWithGet 连接到 Rabbit-Home 服务器
// homeAddrUrl: Rabbit-Home 服务器地址，不需要包含Pattern, 实际Pattern是home.PatternLink，即"/link"
// info: 游戏服务器实例基本信息
// weight: 实例压力系数重，值越大服务器压力越高
// cb: 回调，传入nil表示不处理
// 返回值: 如果调用出现错误，则返回错误信息
func LinkWithGet(homeAddrUrl string, info core.LinkEntity, weight float64, cb httpx.ReqCallBack) error {
	bs, err := jsoniter.Marshal(info)
	if nil != err {
		return err
	}
	data := base64.StdEncoding.EncodeToString(bs)
	httpUrl := homeAddrUrl + home.PatternLink + fmt.Sprintf("?%s=%s&%s=%v", home.PatternDataKey, data, home.PatternEntityWeightKey, weight)
	return httpx.HttpGet(httpUrl, cb)
}

// UnlinkWithGet 断开与 Rabbit-Home 服务器的连接
// homeAddrUrl: Rabbit-Home 服务器地址，不需要包含Pattern, 实际Pattern是home.PatternUnlink，即"/unlink"
// id: 本地服务器id
// cb: 回调，传入nil表示不处理
// 返回值: 如果调用出现错误，则返回错误信息
func UnlinkWithGet(homeAddrUrl string, id string, cb httpx.ReqCallBack) error {
	data := base64.StdEncoding.EncodeToString([]byte(id))
	httpUrl := homeAddrUrl + home.PatternUnlink + fmt.Sprintf("?%s=%s", home.PatternDataKey, data)
	return httpx.HttpGet(httpUrl, cb)
}

// UpdateWithGet 更新服务器状态
// homeAddrUrl: Rabbit-Home 服务器地址，不需要包含Pattern, 实际Pattern是home.PatternUpdate，即"/update"
// info: 实例基本状态
// cb: 回调，传入nil表示不处理
// 返回值: 如果调用出现错误，则返回错误信息
func UpdateWithGet(homeAddrUrl string, info core.EntityStatus, cb httpx.ReqCallBack) error {
	bs, err := jsoniter.Marshal(info)
	if nil != err {
		return err
	}
	data := base64.StdEncoding.EncodeToString(bs)
	httpUrl := homeAddrUrl + home.PatternUpdate + fmt.Sprintf("?%s=%s", home.PatternDataKey, data)
	return httpx.HttpGet(httpUrl, cb)
}

// UpdateDetailWithGet 更新服务器详细状态
// homeAddrUrl: Rabbit-Home 服务器地址，不需要包含Pattern, 实际Pattern是home.PatternUpdate，即"/update"
// detail: 实例详细状态
// cb: 回调，传入nil表示不处理
// 返回值: 如果调用出现错误，则返回错误信息
func UpdateDetailWithGet(homeAddrUrl string, detail core.EntityDetailStatus, cb httpx.ReqCallBack) error {
	bs, err := jsoniter.Marshal(detail)
	if nil != err {
		return err
	}
	data := base64.StdEncoding.EncodeToString(bs)
	httpUrl := homeAddrUrl + home.PatternUpdate + fmt.Sprintf("?%s=%s&%s=true", home.PatternDataKey, data, home.PatternEntityDetailKey)
	return httpx.HttpGet(httpUrl, cb)
}

// RouteWithGet 路由请求，获得合适的服务器实例信息
// homeAddrUrl: Rabbit-Home 服务器地址，不需要包含Pattern, 实际Pattern是home.PatternRoute，即"/route"
// cb: 回调，传入nil表示不处理
// 返回值: 如果调用出现错误，则返回错误信息
func RouteWithGet(homeAddrUrl string, cb httpx.ReqCallBack) error {
	httpUrl := homeAddrUrl + home.PatternRoute
	return httpx.HttpGet(httpUrl, cb)
}
