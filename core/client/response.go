// Package client
// Create on 2025/4/16
// @author xuzhuoxi
package client

import (
	"github.com/xuzhuoxi/Rabbit-Home/core"
	"github.com/xuzhuoxi/Rabbit-Home/core/utils"
	"github.com/xuzhuoxi/infra-go/cryptox"
	"net/http"
)

// ParseLinkBackInfo 解析连接返回信息
// sucBackInfo: 连接成功返回信息，因为无重要信息，不需要加密数据
// failBackInfo: 连接失败返回信息
// err: 数据处理失败错误
func ParseLinkBackInfo(res *http.Response, body *[]byte, decryptCipher cryptox.IDecryptCipher) (sucBackInfo *core.LinkBackInfo, failBackInfo *core.HomeResponseInfo, err error) {
	if res.StatusCode != http.StatusOK {
		failBackInfo, err = utils.ParseHomeResponseInfo(*body)
		return
	}
	sucBackInfo = &core.LinkBackInfo{}
	err = utils.ParseDate(*body, sucBackInfo, decryptCipher)
	if nil != err {
		return nil, nil, err
	}
	return
}

// ParseUnlinkBackInfo 解析取消连接返回信息
// sucBackInfo: 连接成功返回信息，因为无重要信息，不需要加密数据
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

// ParseUpdateBackInfo 解析更新返回信息
// failBackInfo: 更新失败返回信息
// err: 数据处理失败错误
// 更新成功时不返回任何信息
func ParseUpdateBackInfo(res *http.Response, body *[]byte) (failBackInfo *core.HomeResponseInfo, err error) {
	if res.StatusCode != http.StatusOK {
		failBackInfo, err = utils.ParseHomeResponseInfo(*body)
		return
	}
	return nil, nil
}

// ParseQueryRouteBackInfo
// 解析查询路由返回信息
// sucBackInfo: 查询路由成功返回信息
// failBackInfo: 查询路由失败返回信息
// err: 数据处理失败错误
func ParseQueryRouteBackInfo(res *http.Response, body *[]byte) (sucBackInfo *core.QueryRouteBackInfo, failBackInfo *core.HomeResponseInfo, err error) {
	if res.StatusCode != http.StatusOK {
		failBackInfo, err = utils.ParseHomeResponseInfo(*body)
		return
	}
	sucBackInfo = &core.QueryRouteBackInfo{}
	err = utils.ParseDate(*body, sucBackInfo, nil)
	if nil != err {
		return nil, nil, err
	}
	return
}
