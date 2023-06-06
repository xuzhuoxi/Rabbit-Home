// Package home
// Create on 2023/6/4
// @author xuzhuoxi
package home

import (
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Home/core"
	"net/http"
	"strconv"
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
	var weigthStr string
	if l.post {
		err0 = getValueWithPost(request, PatternDataKey, linkEntity)
		weigthStr = request.PostFormValue(PatternEntityWeightKey)
	} else {
		err0 = getValueWithGet(request, PatternDataKey, linkEntity)
		weigthStr = request.FormValue(PatternEntityWeightKey)
	}
	if nil != err0 {
		Logger.Warnln(fmt.Sprintf("Link Entity Fail: %v", err0))
		return
	}
	if linkEntity.IsNotValid() {
		Logger.Warnln(fmt.Sprintf("Link Entity Fail: Entity is not valid. %v", linkEntity))
		return
	}
	entity := NewRegisteredEntity(*linkEntity)
	err := Server.AddEntity(*entity)
	if nil != err {
		Logger.Warnln(fmt.Sprintf("Link Entity Fail: %v", err))
		return
	}
	Logger.Infoln(fmt.Sprintf("Link Entity Succ: %v", linkEntity))

	if weigthStr != "" {
		weight, err := strconv.ParseFloat(weigthStr, 64)
		if nil != err {
			Logger.Warnln(fmt.Sprintf("Update State After Link Fail: %v", err))
			return
		}
		state := core.EntityState{Id: linkEntity.Id, Weight: weight}
		Server.UpdateState(state)
		Logger.Warnln(fmt.Sprintf("Update State After Link Succ: %s", state.String()))
	}
}
