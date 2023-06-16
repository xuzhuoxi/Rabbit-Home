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
	state := &core.EntityStatus{}
	var err error
	if l.post {
		err = getValueWithPost(request, PatternDataKey, state)
	} else {
		err = getValueWithGet(request, PatternDataKey, state)
	}
	if nil != err {
		Logger.Warnln(fmt.Sprintf("Update State Fail: %v", err))
		return
	}
	if state.IsNotValid() {
		Logger.Warnln(fmt.Sprintf("Update State Fail: State is not valid. %v", state))
		return
	}
	ok := Server.UpdateState(*state)
	if !ok {
		Logger.Warnln("Update State Fail: not ok! ", *state)
		return
	}
	Logger.Infoln(fmt.Sprintf("Update State Succ: %v", state))
}
