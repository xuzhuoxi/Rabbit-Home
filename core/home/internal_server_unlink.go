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
	entity, ok := Server.RemoveEntity(id)
	if !ok || nil == entity {
		Logger.Warnln(fmt.Sprintf("Unlink Entity(%s) fail! Entity is not exist!", id))
		return
	}
	fmt.Println(fmt.Sprintf("Unlink Entity(%s) Succ!", id))
}
