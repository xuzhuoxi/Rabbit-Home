// Package home
// Create on 2023/6/4
// @author xuzhuoxi
package home

import (
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Home/src/core"
	"net/http"
)

func newServerLinkHandler() http.Handler {
	return &serverLinkHandler{post: serverPost}
}

type serverLinkHandler struct {
	post bool
}

func (l *serverLinkHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	linkEntity := &core.LinkEntity{}
	var err0 error
	if l.post {
		err0 = getValueWithPost(request, PatternDataKey, linkEntity)
	} else {
		err0 = getValueWithGet(request, PatternDataKey, linkEntity)
	}
	if nil != err0 {
		Logger.Warnln(fmt.Sprintf("LinkEntity Fail: %v", err0))
		return
	}
	if linkEntity.IsNotValid() {
		Logger.Warnln(fmt.Sprintf("LinkEntity Fail: Entity is not valid. %v", linkEntity))
		return
	}
	entity := NewRegisteredEntity(*linkEntity)
	err := Server.GetEntityList().AddEntity(*entity)
	if nil != err {
		Logger.Warnln(fmt.Sprintf("LinkEntity Fail: %v", err))
		return
	}
	Logger.Infoln(fmt.Sprintf("LinkEntity Succ: %v", linkEntity))
}
