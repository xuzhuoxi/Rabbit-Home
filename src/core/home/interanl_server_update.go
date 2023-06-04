// Package home
// Create on 2023/6/4
// @author xuzhuoxi
package home

import (
	"github.com/xuzhuoxi/Rabbit-Home/src/core"
	"net/http"
)

func NewServerUpdateHandler() http.Handler {
	return &serverUpdateHandler{post: serverPost}
}

type serverUpdateHandler struct {
	post bool
}

func (l *serverUpdateHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	state := &core.EntityState{}
	ok := false
	if l.post {
		ok = getValueWithPost(request, linkKey, state)
	} else {
		ok = getValueWithGet(request, linkKey, state)
	}
	if !ok || state.IsNotValid() {
		return
	}
	Server.GetEntityList().UpdateState(*state)
}
