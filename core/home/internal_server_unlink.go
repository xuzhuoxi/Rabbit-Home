// Package home
// Create on 2023/6/4
// @author xuzhuoxi
package home

import (
	"fmt"
	"net/http"
)

func newServerUnlinkHandler() http.Handler {
	return &serverUnlinkHandler{post: serverPost}
}

type serverUnlinkHandler struct {
	post bool
}

func (l *serverUnlinkHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if !ServerConfig.VerifyInternalIP(getClientIpAddr(request)) { // 验证是否内部IP
		return
	}
	var bsId []byte
	var err error
	if l.post {
		bsId, err = getStringWithPost(request, PatternDataKey)
	} else {
		bsId, err = getStringWithGet(request, PatternDataKey)
	}
	if nil != err {
		warnInfo := fmt.Sprintf("Unlink Entity Fail: %v", err)
		warnAndResponse(writer, http.StatusBadRequest, warnInfo, Logger)
		return
	}
	id := string(bsId)
	entity, ok := Server.RemoveEntity(id)
	if !ok || nil == entity {
		warnInfo := fmt.Sprintf("Unlink Entity(%s) fail! Entity is not exist!", bsId)
		warnAndResponse(writer, http.StatusNotFound, warnInfo, Logger)
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write(empty)
	fmt.Println(fmt.Sprintf("Unlink Entity(%s) Succ!", bsId))
}
