// Package home
// Create on 2023/6/4
// @author xuzhuoxi
package home

import (
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
	if !ServerConfig.VerifyExternalIP(getClientIpAddr(request)) { // 验证是否外部IP
		return
	}

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
		warnInfo := fmt.Sprintf("Qurey Fail: %v", err)
		warnAndResponse(writer, http.StatusBadRequest, warnInfo, Logger)
		return
	}
	entity, ok := Server.QueryEntity(query.Name, query.PlatformId)
	if !ok {
		warnInfo := fmt.Sprintf("Query Fail from %s!", request.RemoteAddr)
		warnAndResponse(writer, http.StatusNotFound, warnInfo, Logger)
		return
	}
	writer.WriteHeader(http.StatusOK)
	bs, _ := jsoniter.Marshal(entity)
	value := []byte(toBase64(bs))
	writer.Write(value)
	Logger.Infoln("[clientRouteHandler.ServeHTTP]", fmt.Sprintf("Query Succ from %s. Return %s", request.RemoteAddr, entity.Id))
}
