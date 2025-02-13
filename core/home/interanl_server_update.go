// Package home
// Create on 2023/6/4
// @author xuzhuoxi
package home

import (
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Home/core"
	"net/http"
)

func newServerUpdateHandler() http.Handler {
	return &serverUpdateHandler{post: serverPost}
}

type serverUpdateHandler struct {
	post bool
}

func (l *serverUpdateHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if !ServerConfig.VerifyInternalIP(getClientIpAddr(request)) { // 验证是否内部IP
		return
	}
	state := &core.EntityStatus{}
	var err error
	if l.post {
		err = getValueWithPost(request, PatternDataKey, state)
	} else {
		err = getValueWithGet(request, PatternDataKey, state)
	}
	if nil != err {
		warnInfo := fmt.Sprintf("Update State Fail: %v", err)
		warnAndResponse(writer, http.StatusBadRequest, warnInfo, Logger)
		return
	}
	if state.IsNotValid() {
		warnInfo := fmt.Sprintf("Update State Fail: State is not valid. %v", state)
		warnAndResponse(writer, http.StatusBadRequest, warnInfo, Logger)
		return
	}
	ok := Server.UpdateState(*state)
	if !ok {
		warnInfo := fmt.Sprintf("Update State Fail: Id[%s] unregistered! ", state.Id)
		warnAndResponse(writer, http.StatusNotFound, warnInfo, Logger)
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write(empty)
	Logger.Infoln("[serverUpdateHandler.ServeHTTP]", fmt.Sprintf("Update State Succ: %v", state))
}
