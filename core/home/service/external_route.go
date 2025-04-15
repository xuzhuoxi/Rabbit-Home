// Package service
// Create on 2023/6/4
// @author xuzhuoxi
package service

import (
	"fmt"
	"github.com/json-iterator/go"
	"github.com/xuzhuoxi/Rabbit-Home/core"
	"github.com/xuzhuoxi/Rabbit-Home/core/home"
	"github.com/xuzhuoxi/infra-go/logx"
	"net/http"
)

func NewServiceRouteHandler() http.Handler {
	return &clientRouteHandler{
		logPrefix: "[serverLinkHandler]",
		logger:    home.GlobalLogger}
}

type clientRouteHandler struct {
	logPrefix string
	logger    logx.ILogger
}

func (l *clientRouteHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// 验证数据模式与IP
	if !home.GlobalHomeConfig.ExternalVerifier.VerifyPost(request) {
		return
	}
	if !home.GlobalHomeConfig.VerifyExternalIP(getClientIpAddr(request)) { // 验证是否外部IP
		return
	}

	var base64Data []byte
	var err error
	if request.Method == http.MethodPost {
		base64Data, err = base64FromPost(request, core.HttpKeyQuery)
	} else {
		base64Data, err = base64FromGet(request, core.HttpKeyQuery)
	}
	if nil != err {
		warnInfo := fmt.Sprintf("Param base64 error! %v", err)
		warnResponse(writer, http.StatusBadRequest, core.CodeParamBase64, warnInfo, l.logger)
		return
	}
	data := base64Data
	kv := home.GlobalHomeConfig.ExternalVerifier.KeyVerifier
	if kv.Enable {
		privateCipher := home.GlobalHomeServer.GetHomeKeys().PrivateCipher()
		if nil == privateCipher {
			warnInfo := "Private Key Not Exist!"
			warnResponse(writer, http.StatusBadRequest, core.CodePrivateKeyLack, warnInfo, l.logger)
			return
		}
		data, err = privateCipher.Decrypt(base64Data)
		if nil != err {
			warnInfo := fmt.Sprintf("Decrypt Fail: %v", err)
			warnResponse(writer, http.StatusBadRequest, core.CodeParamDecrypt, warnInfo, l.logger)
			return
		}
	}

	query := &core.QueryInfo{}
	err = jsoniter.Unmarshal(data, query)
	if nil != err {
		warnInfo := fmt.Sprintf("Data Unmarshal Fail: %v", err)
		warnResponse(writer, http.StatusBadRequest, core.CodeParamJson, warnInfo, l.logger)
		return
	}
	entity, ok := home.GlobalHomeServer.QueryEntity(query.Name, query.PlatformId)
	if !ok {
		warnInfo := fmt.Sprintf("Query Fail from %s!", request.RemoteAddr)
		warnResponse(writer, http.StatusNotFound, core.CodeEntityQueryFail, warnInfo, l.logger)
		return
	}
	sucResponse(writer, query, home.GlobalLogger)
	home.GlobalLogger.Infoln("[clientRouteHandler.ServeHTTP]", fmt.Sprintf("Query Succ from %s. Return %s", request.RemoteAddr, entity.Id))
}
