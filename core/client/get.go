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
// httpUrl: http请求连接
// info: 本地服务器基本信息
// weight: 权重
func LinkWithGet(httpUrl string, info core.LinkEntity, weight float64) error {
	bs, err := jsoniter.Marshal(info)
	if nil != err {
		return err
	}
	data := base64.StdEncoding.EncodeToString(bs)
	httpUrl = fmt.Sprintf("%s?%s=%s&%s=%v", httpUrl, home.PatternDataKey, data, home.PatternEntityWeightKey, weight)
	return httpx.HttpGet(httpUrl, nil)
}

func UnlinkWithGet(httpUrl string, id string) error {
	data := base64.StdEncoding.EncodeToString([]byte(id))
	httpUrl = fmt.Sprintf("%s?%s=%s", httpUrl, home.PatternDataKey, data)
	return httpx.HttpGet(httpUrl, nil)
}

func UpdateWithGet(httpUrl string, info core.EntityStatus, cb httpx.ReqCallBack) error {
	bs, err := jsoniter.Marshal(info)
	if nil != err {
		return err
	}
	data := base64.StdEncoding.EncodeToString(bs)
	httpUrl = fmt.Sprintf("%s?%s=%s", httpUrl, home.PatternDataKey, data)
	return httpx.HttpGet(httpUrl, cb)
}
