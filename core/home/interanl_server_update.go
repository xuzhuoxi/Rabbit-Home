// Package home
// Create on 2023/6/4
// @author xuzhuoxi
package home

import (
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Home/core"
	"net/http"
	"strings"
)

func newServerUpdateHandler() http.Handler {
	return &serverUpdateHandler{}
}

type serverUpdateHandler struct{}

func (l *serverUpdateHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if !ServerConfig.Internal.VerifyPost(request) {
		return
	}
	if !ServerConfig.VerifyInternalIP(getClientIpAddr(request)) { // 验证是否内部IP
		return
	}
	isDetail, status, detail := l.getStatus(request)
	fmt.Println("serverUpdateHandler:", isDetail, status, detail)
	if isDetail {
		l.updateDetailStatus(writer, detail)
	} else {
		l.updateStatus(writer, status)
	}
}

func (l *serverUpdateHandler) getStatus(request *http.Request) (isDetail bool, status *core.EntityStatus, detail *core.EntityDetailStatus) {
	isDetail = l.checkDetail(request)
	if isDetail {
		detail = &core.EntityDetailStatus{}
		if request.Method == http.MethodPost {
			getValueWithPost(request, PatternDataKey, detail)
		} else {
			getValueWithGet(request, PatternDataKey, detail)
		}
	} else {
		status = &core.EntityStatus{}
		if request.Method == http.MethodPost {
			getValueWithPost(request, PatternDataKey, status)
		} else {
			getValueWithGet(request, PatternDataKey, status)
		}
	}
	return
}

func (l *serverUpdateHandler) checkDetail(request *http.Request) (isDetail bool) {
	var value []byte
	var err error
	if request.Method == http.MethodPost {
		value, err = getStringWithPost(request, PatternEntityDetailKey)
	} else {
		value, err = getStringWithGet(request, PatternEntityDetailKey)
	}
	if nil != err {
		dStr := strings.ToLower(string(value))
		return dStr == "true" || dStr == "1"
	} else {
		return false
	}
}

func (l *serverUpdateHandler) updateStatus(writer http.ResponseWriter, status *core.EntityStatus) {
	if status.IsNotValid() {
		warnInfo := fmt.Sprintf("Update State Fail: State is not valid. %v", status)
		warnAndResponse(writer, http.StatusBadRequest, warnInfo, Logger)
		return
	}
	ok := Server.UpdateState(*status)
	if !ok {
		warnInfo := fmt.Sprintf("Update State Fail: Id[%s] unregistered! ", status.Id)
		warnAndResponse(writer, http.StatusNotFound, warnInfo, Logger)
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write(empty)
	Logger.Infoln("[serverUpdateHandler.ServeHTTP]", fmt.Sprintf("Update State Succ: %v", status))
}

func (l *serverUpdateHandler) updateDetailStatus(writer http.ResponseWriter, detail *core.EntityDetailStatus) {
	if detail.IsNotValid() {
		warnInfo := fmt.Sprintf("Update Detail State Fail: State is not valid. %v", detail)
		warnAndResponse(writer, http.StatusBadRequest, warnInfo, Logger)
		return
	}
	ok := Server.UpdateDetailState(*detail)
	if !ok {
		warnInfo := fmt.Sprintf("Update Detail State Fail: Id[%s] unregistered! ", detail.Id)
		warnAndResponse(writer, http.StatusNotFound, warnInfo, Logger)
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write(empty)
	Logger.Infoln("[serverUpdateHandler.ServeHTTP]", fmt.Sprintf("Update Detail State Succ: %v", detail))
}
