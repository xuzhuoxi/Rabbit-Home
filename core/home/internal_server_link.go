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
	if !ServerConfig.VerifyInternalIP(getClientIpAddr(request)) { // 验证是否内部IP
		return
	}
	funcName := "[serverLinkHandler.ServeHTTP]"
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
		warnInfo := fmt.Sprintf("Link Entity Fail: %v", err0)
		warnAndResponse(writer, http.StatusBadRequest, warnInfo, Logger)
		return
	}
	if linkEntity.IsNotValid() {
		warnInfo := fmt.Sprintf("Link Entity Fail: Entity is not valid. %v", linkEntity)
		warnAndResponse(writer, http.StatusBadRequest, warnInfo, Logger)
		return
	}
	entity := NewRegisteredEntity(*linkEntity)
	err1 := Server.AddEntity(*entity)
	if nil != err1 {
		err2 := Server.ReplaceEntity(*entity)
		if nil != err2 {
			warnInfo := fmt.Sprintf("Link Entity Fail: %v", err2)
			warnAndResponse(writer, http.StatusBadRequest, warnInfo, Logger)
			return
		}
		Logger.Infoln(funcName, fmt.Sprintf("Relink Entity Succ: %v", linkEntity))
	} else {
		Logger.Infoln(funcName, fmt.Sprintf("Link Entity Succ: %v", linkEntity))
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write(empty)

	if weightStr != "" {
		weight, err := strconv.ParseFloat(weightStr, 64)
		if nil != err {
			Logger.Warnln(funcName, fmt.Sprintf("Update State After Link Fail: %v", err))
			return
		}
		state := core.EntityStatus{Id: linkEntity.Id, Weight: weight}
		Server.UpdateState(state)
		Logger.Infoln(funcName, fmt.Sprintf("Update State After Link Succ: %s", state.String()))
	}
}
