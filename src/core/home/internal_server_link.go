// Package home
// Create on 2023/6/4
// @author xuzhuoxi
package home

import (
	"github.com/xuzhuoxi/Rabbit-Home/src/core"
	"net/http"
)

const (
	linkKey = "link"
)

func newServerLinkHandler() http.Handler {
	return &serverLinkHandler{post: serverPost}
}

type serverLinkHandler struct {
	post bool
}

func (l *serverLinkHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	data := &core.LinkEntity{}
	ok := false
	if l.post {
		ok = getValueWithPost(request, linkKey, data)
	} else {
		ok = getValueWithGet(request, linkKey, data)
	}
	if !ok || data.IsNotValid() {
		return
	}
	entity := NewRegisteredEntity(*data)
	Server.GetEntityList().AddEntity(*entity)
}
