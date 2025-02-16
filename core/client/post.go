// Package client
// Create on 2023/6/16
// @author xuzhuoxi
package client

import (
	"encoding/base64"
	jsoniter "github.com/json-iterator/go"
	"github.com/xuzhuoxi/Rabbit-Home/core"
	"github.com/xuzhuoxi/Rabbit-Home/core/home"
	"github.com/xuzhuoxi/infra-go/netx/httpx"
	"net/url"
	"strconv"
)

// LinkWithPost 连接到 Rabbit-Home 服务器
// homeAddrUrl: Rabbit-Home 服务器地址，不需要包含Pattern, 实际Pattern是home.PatternLink，即"/link"
// info: 游戏服务器实例基本信息
// weight: 实例压力系数重，值越大服务器压力越高
// cb: 回调，传入nil表示不处理
// 返回值：如果调用出现错误，则返回错误信息
func LinkWithPost(homeAddrUrl string, info core.LinkEntity, weight float64, cb httpx.ReqCallBack) error {
	bs, err := jsoniter.Marshal(info)
	if nil != err {
		return err
	}
	data := base64.StdEncoding.EncodeToString(bs)
	value := make(url.Values)
	value.Set(home.PatternDataKey, data)
	value.Set(home.PatternEntityWeightKey, strconv.FormatFloat(weight, 'f', -1, 64))
	httpUrl := homeAddrUrl + home.PatternLink
	return httpx.HttpPostForm(httpUrl, value, cb)
}

// UnlinkWithPost 断开与 Rabbit-Home 服务器的连接
// homeAddrUrl: Rabbit-Home 服务器地址，不需要包含Pattern, 实际Pattern是home.PatternUnlink，即"/unlink"
// id: 服务器实例ID
// cb: 回调，传入nil表示不处理
// 返回值：如果调用出现错误，则返回错误信息
func UnlinkWithPost(homeAddrUrl string, id string, cb httpx.ReqCallBack) error {
	data := base64.StdEncoding.EncodeToString([]byte(id))
	value := make(url.Values)
	value.Set(home.PatternDataKey, data)
	httpUrl := homeAddrUrl + home.PatternUnlink
	return httpx.HttpPostForm(httpUrl, value, cb)
}

// UpdateWithPost 更新服务器状态
// homeAddrUrl: Rabbit-Home 服务器地址，不需要包含Pattern, 实际Pattern是home.PatternRoute，即"/route"
// info: 服务器状态信息
// cb: 回调，传入nil表示不处理
// 返回值：如果调用出现错误，则返回错误信息
func UpdateWithPost(homeAddrUrl string, info core.EntityStatus, cb httpx.ReqCallBack) error {
	bs, err := jsoniter.Marshal(info)
	if nil != err {
		return err
	}
	data := base64.StdEncoding.EncodeToString(bs)
	value := make(url.Values)
	value.Set(home.PatternDataKey, data)
	httpUrl := homeAddrUrl + home.PatternUpdate
	return httpx.HttpPostForm(httpUrl, value, cb)
}

// UpdateDetailWithPost 更新服务器状态
// homeAddrUrl: Rabbit-Home 服务器地址，不需要包含Pattern, 实际Pattern是home.PatternRoute，即"/route"
// detail: 服务器详细状态信息
// cb: 回调，传入nil表示不处理
// 返回值：如果调用出现错误，则返回错误信息
func UpdateDetailWithPost(homeAddrUrl string, detail core.EntityDetailStatus, cb httpx.ReqCallBack) error {
	bs, err := jsoniter.Marshal(detail)
	if nil != err {
		return err
	}
	data := base64.StdEncoding.EncodeToString(bs)
	value := make(url.Values)
	value.Set(home.PatternDataKey, data)
	value.Set(home.PatternEntityDetailKey, "true")
	httpUrl := homeAddrUrl + home.PatternRoute
	return httpx.HttpPostForm(httpUrl, value, cb)
}

// RouteWithPost 路由请求，获得合适的服务器实例信息
// homeAddrUrl: Rabbit-Home 服务器地址，不需要包含Pattern, 实际Pattern是home.PatternRoute，即"/route"
// cb: 回调，传入nil表示不处理
// 返回值：如果调用出现错误，则返回错误信息
func RouteWithPost(homeAddrUrl string, cb httpx.ReqCallBack) error {
	value := make(url.Values)
	httpUrl := homeAddrUrl + home.PatternRoute
	return httpx.HttpPostForm(httpUrl, value, cb)
}
