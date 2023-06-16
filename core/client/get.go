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
	"github.com/xuzhuoxi/infra-go/netx"
)

func LinkWithGet(httpUrl string, info core.LinkEntity, weight float64) error {
	bs, err := jsoniter.Marshal(info)
	if nil != err {
		return err
	}
	data := base64.StdEncoding.EncodeToString(bs)
	httpUrl = fmt.Sprintf("%s?%s=%s&%s=%v", httpUrl, home.PatternDataKey, data, home.PatternEntityWeightKey, weight)
	return netx.HttpGet(httpUrl, nil)
}

func UnlinkWithGet(httpUrl string, id string) error {
	data := base64.StdEncoding.EncodeToString([]byte(id))
	httpUrl = fmt.Sprintf("%s?%s=%s", httpUrl, home.PatternDataKey, data)
	return netx.HttpGet(httpUrl, nil)
}

func UpdateWithGet(httpUrl string, info core.EntityStatus) error {
	bs, err := jsoniter.Marshal(info)
	if nil != err {
		return err
	}
	data := base64.StdEncoding.EncodeToString(bs)
	httpUrl = fmt.Sprintf("%s?%s=%s", httpUrl, home.PatternDataKey, data)
	return netx.HttpGet(httpUrl, nil)
}
