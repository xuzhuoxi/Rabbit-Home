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
	var weightStr string
	if l.post {
		err0 = getValueWithPost(request, PatternDataKey, linkEntity)
		weightStr = request.PostFormValue(PatternEntityWeightKey)
	} else {
		err0 = getValueWithGet(request, PatternDataKey, linkEntity)
		weightStr = request.FormValue(PatternEntityWeightKey)
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
	err1 := Server.AddEntity(*entity)
	if nil != err1 {
		err2 := Server.ReplaceEntity(*entity)
		if nil != err2 {
			Logger.Warnln(fmt.Sprintf("Link Entity Fail: %v", err2))
			return
		}
		Logger.Infoln(fmt.Sprintf("Relink Entity Succ: %v", linkEntity))
	} else {
		Logger.Infoln(fmt.Sprintf("Link Entity Succ: %v", linkEntity))
	}

	if weightStr != "" {
		weight, err := strconv.ParseFloat(weightStr, 64)
		if nil != err {
			Logger.Warnln(fmt.Sprintf("Update State After Link Fail: %v", err))
			return
		}
		state := core.EntityStatus{Id: linkEntity.Id, Weight: weight}
		Server.UpdateState(state)
		Logger.Infoln(fmt.Sprintf("Update State After Link Succ: %s", state.String()))
	}
}
