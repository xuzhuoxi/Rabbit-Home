// Package home
// Create on 2023/6/4
// @author xuzhuoxi
package home

import (
	"encoding/base64"
	"fmt"
	"github.com/json-iterator/go"
	"github.com/xuzhuoxi/Rabbit-Home/core"
	"net/http"
)

func newClientRouteHandler() http.Handler {
	return &clientRouteHandler{post: clientPost}
}

type clientRouteHandler struct {
	post bool
}

func (l *clientRouteHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	query := &core.HttpRequestQueryEntity{}
	var err error
	if l.post {
		if l.post {
			err = getValueWithPost(request, PatternDataKey, query)
		} else {
			err = getValueWithGet(request, PatternDataKey, query)
		}
	}
	if nil != err {
		Logger.Warnln(fmt.Sprintf("Qurey Fail: %v", err))
		return
	}
	entity, ok := Server.QueryEntity(query.Name, query.PlatformId)
	resp := &core.HttpResponse{}
	if !ok {
		resp.Status = StatusNotFound
		Logger.Warnln(fmt.Sprintf("Query Fail from %s!", request.RemoteAddr))
	} else {
		resp.Status = StatusFound
		bs, _ := jsoniter.Marshal(entity)
		value := base64.StdEncoding.EncodeToString(bs)
		resp.Value = value
		Logger.Warnln(fmt.Sprintf("Query Succ from %s. Return %s", request.RemoteAddr, entity.Id))
	}
	rs, _ := jsoniter.Marshal(resp)
	writer.Write(rs)
}
