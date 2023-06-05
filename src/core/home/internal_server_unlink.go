// Package home
// Create on 2023/6/4
// @author xuzhuoxi
package home

import (
	"net/http"
)

func newServerUnlinkHandler() http.Handler {
	return &serverUnlinkHandler{post: serverPost}
}

type serverUnlinkHandler struct {
	post bool
}

func (l *serverUnlinkHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	id := ""
	var err error
	if l.post {
		id, err = getStringWithPost(request, PatternDataKey)
	} else {
		id, err = getStringWithGet(request, PatternDataKey)
	}
	if nil != err {
		return
	}
	Server.GetEntityList().RemoveEntity(id)
}
