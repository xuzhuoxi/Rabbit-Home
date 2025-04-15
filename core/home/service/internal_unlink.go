// Package service
// Create on 2023/6/4
// @author xuzhuoxi
package service

import (
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Home/core"
	"github.com/xuzhuoxi/Rabbit-Home/core/home"
	"github.com/xuzhuoxi/infra-go/logx"
	"net/http"
)

func NewServiceUnlinkHandler() http.Handler {
	return &serverUnlinkHandler{
		logPrefix: "[serverLinkHandler.ServeHTTP]",
		logger:    home.GlobalLogger}
}

type serverUnlinkHandler struct {
	logPrefix string
	logger    logx.ILogger
}

func (l *serverUnlinkHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if !home.GlobalHomeConfig.InternalVerifier.VerifyPost(request) {
		return
	}
	if !home.GlobalHomeConfig.VerifyInternalIP(getClientIpAddr(request)) { // 验证是否内部IP
		return
	}
	unlinkInfo := &core.UnlinkInfo{}
	var err0 error
	if request.Method == http.MethodPost {
		err0 = loadValueFromPost(request, core.HttpKeyData, unlinkInfo)
	} else {
		err0 = loadValueFromGet(request, core.HttpKeyData, unlinkInfo)
	}
	if nil != err0 {
		warnInfo := fmt.Sprintf("%s Unlink Entity Fail: %v", l.logPrefix, err0)
		warnResponse(writer, http.StatusBadRequest, core.CodeParamError, warnInfo, l.logger)
		return
	}

	home.GlobalLock.Lock()
	defer home.GlobalLock.Unlock()

	if _, has := home.GlobalHomeServer.GetEntityById(unlinkInfo.Id); !has {
		warnInfo := fmt.Sprintf("%s Unlink Entity Fail: No such id='%s'", l.logPrefix, unlinkInfo.Id)
		warnResponse(writer, http.StatusBadRequest, core.CodeParamInvalid, warnInfo, l.logger)
		return
	}
	// 验证签名
	kv := home.GlobalHomeConfig.InternalVerifier.KeyVerifier
	var pass bool
	if kv.Enable {
		_, pass = home.GlobalHomeServer.GetHomeKeys().VerifyUnlinkSign(unlinkInfo)
		if !pass {
			warnInfo := fmt.Sprintf("%s Unlink Entity(%s) Fail: Key verify fail.", l.logPrefix, unlinkInfo.Id)
			warnResponse(writer, http.StatusBadRequest, core.CodeVerifyUnlinkSignFail, warnInfo, l.logger)
			return
		}
	}

	// 删除实例
	entity, ok := home.GlobalHomeServer.RemoveEntity(unlinkInfo.Id)
	if !ok || nil == entity {
		warnInfo := fmt.Sprintf("%s Unlink Entity(%s) fail! Entity is not exist!", l.logPrefix, unlinkInfo.Id)
		warnResponse(writer, http.StatusNotFound, core.CodeUnlinkEntityFail, warnInfo, l.logger)
		return
	}

	backInfo := &core.UnlinkBackInfo{Id: unlinkInfo.Id}
	sucResponse(writer, backInfo, l.logger)
	l.logger.Infoln(l.logPrefix, "Unlink Entity Suc:", unlinkInfo.Id)
}
