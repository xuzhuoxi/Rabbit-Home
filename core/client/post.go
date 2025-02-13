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
	"net/url"
)

func LinkWithPost(httpUrl string, info core.LinkEntity, weight float64) error {
	bs, err := jsoniter.Marshal(info)
	if nil != err {
		return err
	}
	data := base64.StdEncoding.EncodeToString(bs)
	value := make(url.Values)
	value.Set(home.PatternDataKey, data)
	value.Set(home.PatternEntityWeightKey, fmt.Sprintf("%v", weight))
	return httpx.HttpPostForm(httpUrl, value, nil)
}

func UnlinkWithPost(httpUrl string, id string) error {
	data := base64.StdEncoding.EncodeToString([]byte(id))
	value := make(url.Values)
	value.Set(home.PatternDataKey, data)
	return httpx.HttpPostForm(httpUrl, value, nil)
}

func UpdateWithPost(httpUrl string, info core.EntityStatus) error {
	bs, err := jsoniter.Marshal(info)
	if nil != err {
		return err
	}
	data := base64.StdEncoding.EncodeToString(bs)
	value := make(url.Values)
	value.Set(home.PatternDataKey, data)
	return httpx.HttpPostForm(httpUrl, value, nil)
}
