// Package client
// Create on 2023/6/16
// @author xuzhuoxi
package client

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/xuzhuoxi/Rabbit-Home/core"
	"github.com/xuzhuoxi/Rabbit-Home/core/home"
	"github.com/xuzhuoxi/Rabbit-Home/core/utils"
	"github.com/xuzhuoxi/infra-go/cryptox/symmetric"
	"github.com/xuzhuoxi/infra-go/netx/httpx"
	"net/http"
	"strconv"
)

var (
	isDetail = core.Base64Encoding.EncodeToString([]byte("1"))
)

// LinkWithGet 连接到 Rabbit-Home 服务器
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

// ParseLinkBackInfo 解析连接返回信息
// sucBackInfo: 连接成功返回信息
// failBackInfo: 连接失败返回信息
// err: 数据处理失败错误
func ParseLinkBackInfo(res *http.Response, body *[]byte) (sucBackInfo *core.LinkBackInfo, failBackInfo *core.HomeResponseInfo, err error) {
	if res.StatusCode != http.StatusOK {
		failBackInfo, err = utils.ParseHomeResponseInfo(*body)
		return
	}
	sucBackInfo = &core.LinkBackInfo{}
	err = utils.ParseDate(*body, sucBackInfo, nil)
	if nil != err {
		return nil, nil, err
	}
	return
}

// UnlinkWithGet 断开与 Rabbit-Home 服务器的连接
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

// ParseUnlinkBackInfo 解析取消连接返回信息
// sucBackInfo: 连接成功返回信息
// failBackInfo: 连接失败返回信息
// err: 数据处理失败错误
func ParseUnlinkBackInfo(res *http.Response, body *[]byte) (sucBackInfo *core.UnlinkBackInfo, failBackInfo *core.HomeResponseInfo, err error) {
	if res.StatusCode != http.StatusOK {
		failBackInfo, err = utils.ParseHomeResponseInfo(*body)
		return
	}
	sucBackInfo = &core.UnlinkBackInfo{}
	err = utils.ParseDate(*body, sucBackInfo, nil)
	if nil != err {
		return nil, nil, err
	}
	return
}

// UpdateWithGet 更新服务器状态
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
		bs, err = aesCipher.EncryptGCM(bs)
		if nil != err {
			return err
		}
	}
	data := core.Base64Encoding.EncodeToString(bs)
	httpUrl := homeAddrUrl + home.PatternUpdate + fmt.Sprintf("?%s=%s&%s=%s", core.HttpKeyId, id, core.HttpKeyData, data)
	//fmt.Println("UpdateUrl:", info, nil != aesCipher, httpUrl)
	return httpx.HttpGet(httpUrl, cb)
}

// UpdateDetailWithGet 更新服务器详细状态
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
		bs, err = aesCipher.EncryptGCM(bs)
		if nil != err {
			return err
		}
	}
	data := core.Base64Encoding.EncodeToString(bs)
	httpUrl := homeAddrUrl + home.PatternUpdate + fmt.Sprintf("?%s=%s&%s=%s&%s=%s", core.HttpKeyId, id, core.HttpKeyData, data, core.HttpKeyDetail, isDetail)
	//fmt.Println("DetailUrl:", detail, nil != aesCipher, httpUrl)
	return httpx.HttpGet(httpUrl, cb)
}

// ParseUpdateBackInfo 解析更新返回信息
// failBackInfo: 更新失败返回信息
// err: 数据处理失败错误
// 更新成功时不返回任何信息
func ParseUpdateBackInfo(res *http.Response, body *[]byte, aesCipher symmetric.IAESCipher) (failBackInfo *core.HomeResponseInfo, err error) {
	if res.StatusCode != http.StatusOK {
		failBackInfo, err = utils.ParseHomeResponseInfo(*body)
		return
	}
	return nil, nil
}

// RouteWithGet 路由请求，获得合适的服务器实例信息
// homeAddrUrl: Rabbit-Home 服务器地址，不需要包含Pattern, 实际Pattern是home.PatternRoute，即"/route"
// cb: 回调，传入nil表示不处理
// 返回值: 如果调用出现错误，则返回错误信息
func RouteWithGet(homeAddrUrl string, cb httpx.ReqCallBack) error {
	httpUrl := homeAddrUrl + home.PatternRoute
	return httpx.HttpGet(httpUrl, cb)
}
