// Package client
// Create on 2023/6/16
// @author xuzhuoxi
package client

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/xuzhuoxi/Rabbit-Home/core"
	"github.com/xuzhuoxi/Rabbit-Home/core/home"
	"github.com/xuzhuoxi/infra-go/cryptox/asymmetric"
	"github.com/xuzhuoxi/infra-go/cryptox/symmetric"
	"github.com/xuzhuoxi/infra-go/netx/httpx"
	"strconv"
)

var (
	isDetail = core.Base64Encoding.EncodeToString([]byte("1"))
)

// LinkWithGet 注册到 Rabbit-Home 服务器
// homeAddrUrl: Rabbit-Home 服务器地址，不需要包含Pattern, 实际Pattern是home.PatternLink，即"/link"
// info: 游戏服务器实例基本信息
//   注意：info中的Signature数据在启用签名认证时，必须设置
// weight: 实例压力系数重，值越大服务器压力越高
// cb: 回调，传入nil表示不处理
// 返回值: 如果调用出现错误，则返回错误信息
func LinkWithGet(homeAddrUrl string, info core.LinkInfo, weight float64, cb httpx.ReqCallBack) error {
	bs, err := jsoniter.Marshal(info)
	if nil != err {
		return err
	}
	data := core.Base64Encoding.EncodeToString(bs)
	w := core.Base64Encoding.EncodeToString([]byte(strconv.FormatFloat(weight, 'f', -1, 64)))
	httpUrl := homeAddrUrl + home.PatternLink + fmt.Sprintf("?%s=%s&%s=%v", core.HttpKeyData, data, core.HttpKeyWeight, w)
	//fmt.Println("LinkUrl:", httpUrl)
	return httpx.HttpGet(httpUrl, cb)
}

// UnlinkWithGet 从 Rabbit-Home 服务器上注销
// homeAddrUrl: Rabbit-Home 服务器地址，不需要包含Pattern, 实际Pattern是home.PatternUnlink，即"/unlink"
// info: 移除服务器必要信息
//   注意：info中的Signature数据在启用签名认证时，必须设置
// cb: 回调，传入nil表示不处理
// 返回值: 如果调用出现错误，则返回错误信息
func UnlinkWithGet(homeAddrUrl string, info core.UnlinkInfo, cb httpx.ReqCallBack) error {
	bs, err := jsoniter.Marshal(info)
	if nil != err {
		return err
	}
	data := core.Base64Encoding.EncodeToString(bs)
	httpUrl := homeAddrUrl + home.PatternUnlink + fmt.Sprintf("?%s=%s", core.HttpKeyData, data)
	return httpx.HttpGet(httpUrl, cb)
}

// UpdateWithGet 更新Rabbit-Home服务器上的简单信息
// homeAddrUrl: Rabbit-Home 服务器地址，不需要包含Pattern, 实际Pattern是home.PatternUpdate，即"/update"
// info: 实例基本状态
// aesCipher: AES加密处理器，传入nil表示不加密
// cb: 回调，传入nil表示不处理
// 返回值: 如果调用出现错误，则返回错误信息
func UpdateWithGet(homeAddrUrl string, info core.UpdateInfo, aesCipher symmetric.IAESCipher, cb httpx.ReqCallBack) error {
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
	//fmt.Println("UpdateWithGet: bs=", info, bs)
	data := core.Base64Encoding.EncodeToString(bs)
	httpUrl := homeAddrUrl + home.PatternUpdate + fmt.Sprintf("?%s=%s&%s=%s", core.HttpKeyId, id, core.HttpKeyData, data)
	//fmt.Println("UpdateWithGet: url=", info, nil != aesCipher, bs, httpUrl)
	return httpx.HttpGet(httpUrl, cb)
}

// UpdateDetailWithGet 更新Rabbit-Home服务器上的详细信息
// homeAddrUrl: Rabbit-Home 服务器地址，不需要包含Pattern, 实际Pattern是home.PatternUpdate，即"/update"
// detail: 实例详细状态
// cb: 回调，传入nil表示不处理
// 返回值: 如果调用出现错误，则返回错误信息
func UpdateDetailWithGet(homeAddrUrl string, detail core.UpdateDetailInfo, aesCipher symmetric.IAESCipher, cb httpx.ReqCallBack) error {
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
	httpUrl := homeAddrUrl + home.PatternUpdate + fmt.Sprintf("?%s=%s&%s=%s&%s=%s", core.HttpKeyId, id, core.HttpKeyData, data, core.HttpKeyDetail, isDetail)
	//fmt.Println("DetailUrl:", detail, nil != aesCipher, httpUrl)
	return httpx.HttpGet(httpUrl, cb)
}

// QueryRouteWithGet 路由请求，获得合适的服务器实例信息
// homeAddrUrl: Rabbit-Home服务器地址，不需要包含Pattern, 实际Pattern是home.PatternRoute，即"/route"
// queryRoute: 路由请求信息
// rsaPublicCipher: RSA加密处理器，传入nil表示不加密. 公钥属于Rabbit-Home
// cb: 回调，传入nil表示不处理
// 返回值: 如果调用出现错误，则返回错误信息
func QueryRouteWithGet(homeAddrUrl string, queryRoute core.QueryRouteInfo, publicRsaCipher asymmetric.IRSAPublicCipher, cb httpx.ReqCallBack) error {
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
	httpUrl := homeAddrUrl + home.PatternRoute + fmt.Sprintf("?%s=%s", core.HttpKeyQuery, date)
	return httpx.HttpGet(httpUrl, cb)
}
