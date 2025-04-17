// Package client
// Create on 2023/6/16
// @author xuzhuoxi
package client

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/xuzhuoxi/Rabbit-Home/core"
	"github.com/xuzhuoxi/Rabbit-Home/core/home"
	"github.com/xuzhuoxi/infra-go/cryptox/asymmetric"
	"github.com/xuzhuoxi/infra-go/cryptox/symmetric"
	"github.com/xuzhuoxi/infra-go/netx/httpx"
	"net/url"
	"strconv"
)

// LinkWithPost 连接到 Rabbit-Home 服务器
// homeAddrUrl: Rabbit-Home 服务器地址，不需要包含Pattern, 实际Pattern是home.PatternLink，即"/link"
// info: 游戏服务器实例基本信息,
//   注意：info中的Signature数据在启用签名认证时，必须设置
// weight: 实例压力系数重，值越大服务器压力越高
// cb: 回调，传入nil表示不处理
// 返回值：如果调用出现错误，则返回错误信息
func LinkWithPost(homeAddrUrl string, info core.LinkInfo, weight float64, cb httpx.ReqCallBack) error {
	bs, err := jsoniter.Marshal(info)
	if nil != err {
		return err
	}
	w := core.Base64Encoding.EncodeToString([]byte(strconv.FormatFloat(weight, 'f', -1, 64)))
	data := core.Base64Encoding.EncodeToString(bs)
	value := make(url.Values)
	value.Set(core.HttpKeyData, data)
	value.Set(core.HttpKeyWeight, w)
	httpUrl := homeAddrUrl + home.PatternLink
	return httpx.HttpPostForm(httpUrl, value, cb)
}

// UnlinkWithPost 断开与 Rabbit-Home 服务器的连接
// homeAddrUrl: Rabbit-Home 服务器地址，不需要包含Pattern, 实际Pattern是home.PatternUnlink，即"/unlink"
// info: 移除服务器必要信息
//   注意：info中的Signature数据在启用签名认证时，必须设置
// cb: 回调，传入nil表示不处理
// 返回值：如果调用出现错误，则返回错误信息
func UnlinkWithPost(homeAddrUrl string, info core.UnlinkInfo, cb httpx.ReqCallBack) error {
	bs, err := jsoniter.Marshal(info)
	if nil != err {
		return err
	}
	data := core.Base64Encoding.EncodeToString(bs)
	value := make(url.Values)
	value.Set(core.HttpKeyData, data)
	httpUrl := homeAddrUrl + home.PatternUnlink
	return httpx.HttpPostForm(httpUrl, value, cb)
}

// UpdateWithPost 更新服务器状态
// homeAddrUrl: Rabbit-Home 服务器地址，不需要包含Pattern, 实际Pattern是home.PatternRoute，即"/route"
// info: 服务器状态信息
// cb: 回调，传入nil表示不处理
// 返回值：如果调用出现错误，则返回错误信息
func UpdateWithPost(homeAddrUrl string, info core.UpdateInfo, aesCipher symmetric.IAESCipher, cb httpx.ReqCallBack) error {
	id := core.Base64Encoding.EncodeToString([]byte(info.Id))
	bs, err := jsoniter.Marshal(info)
	if nil != err {
		return err
	}
	if nil != aesCipher {
		bs, err = aesCipher.Encrypt(bs)
		if nil != err {
			return err
		}
	}
	data := core.Base64Encoding.EncodeToString(bs)
	value := make(url.Values)
	value.Set(core.HttpKeyId, id)
	value.Set(core.HttpKeyData, data)
	httpUrl := homeAddrUrl + home.PatternUpdate
	return httpx.HttpPostForm(httpUrl, value, cb)
}

// UpdateDetailWithPost 更新服务器状态
// homeAddrUrl: Rabbit-Home 服务器地址，不需要包含Pattern, 实际Pattern是home.PatternRoute，即"/route"
// detail: 服务器详细状态信息
// cb: 回调，传入nil表示不处理
// 返回值：如果调用出现错误，则返回错误信息
func UpdateDetailWithPost(homeAddrUrl string, detail core.UpdateDetailInfo, aesCipher symmetric.IAESCipher, cb httpx.ReqCallBack) error {
	id := core.Base64Encoding.EncodeToString([]byte(detail.Id))
	bs, err := jsoniter.Marshal(detail)
	if nil != err {
		return err
	}
	if nil != aesCipher {
		bs, err = aesCipher.Encrypt(bs)
		if nil != err {
			return err
		}
	}
	data := core.Base64Encoding.EncodeToString(bs)
	value := make(url.Values)
	value.Set(core.HttpKeyId, id)
	value.Set(core.HttpKeyData, data)
	value.Set(core.HttpKeyDetail, "true")
	httpUrl := homeAddrUrl + home.PatternRoute
	return httpx.HttpPostForm(httpUrl, value, cb)
}

// QueryRouteWithPost 路由请求，获得合适的服务器实例信息
// homeAddrUrl: Rabbit-Home 服务器地址，不需要包含Pattern, 实际Pattern是home.PatternRoute，即"/route"
// queryRoute: 路由请求信息
// rsaPublicCipher: 签名公钥，如果启用签名认证，则必须设置
// cb: 回调，传入nil表示不处理
// 返回值：如果调用出现错误，则返回错误信息
func QueryRouteWithPost(homeAddrUrl string, queryRoute core.QueryRouteInfo, publicRsaCipher asymmetric.IRSAPublicCipher, cb httpx.ReqCallBack) error {
	bs, err := jsoniter.Marshal(queryRoute)
	if nil != err {
		return err
	}
	if nil != publicRsaCipher {
		bs, err = publicRsaCipher.Encrypt(bs)
		if nil != err {
			return err
		}
	}
	date := core.Base64Encoding.EncodeToString(bs)
	value := make(url.Values)
	value.Set(core.HttpKeyQuery, date)
	httpUrl := homeAddrUrl + home.PatternRoute
	return httpx.HttpPostForm(httpUrl, value, cb)
}
