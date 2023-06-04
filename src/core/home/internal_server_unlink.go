// Package home
// Create on 2023/6/4
// @author xuzhuoxi
package home

import (
	"net/http"
)

const (
	unlinkKey = "unlink"
)

func newServerUnlinkHandler() http.Handler {
	return &serverUnlinkHandler{post: serverPost}
}

type serverUnlinkHandler struct {
	post bool
}

func (l *serverUnlinkHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	id := ""
	ok := false
	if l.post {
		id, ok = getStringWithPost(request, linkKey)
	} else {
		id, ok = getStringWithGet(request, linkKey)
	}
	if !ok {
		return
	}
	Server.GetEntityList().RemoveEntity(id)
}
