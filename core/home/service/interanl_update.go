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
	"strings"
)

func NewServiceUpdateHandler() http.Handler {
	return &serverUpdateHandler{
		logPrefix: "[serverLinkHandler.ServeHTTP]",
		logger:    home.GlobalLogger}
}

type serverUpdateHandler struct {
	logPrefix string
	logger    logx.ILogger
}

func (l *serverUpdateHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// 验证数据模式与IP
	if !home.GlobalHomeConfig.InternalVerifier.VerifyPost(request) {
		return
	}
	if !home.GlobalHomeConfig.VerifyInternalIP(getClientIpAddr(request)) { // 验证是否内部IP
		return
	}

	// 提取Id数据
	id, ok := l.getId(request)
	if !ok {
		warnInfo := "Update Fail: Lack id params."
		warnResponse(writer, http.StatusBadRequest, core.CodeParamLack, warnInfo, l.logger)
		return
	}

	home.GlobalLock.Lock()
	defer home.GlobalLock.Unlock()

	// 查询已注册信息
	var entity home.RegisteredEntity
	if entity, ok = home.GlobalHomeServer.GetEntityById(id); !ok {
		warnInfo := fmt.Sprintf("%s Update Fail: Id[%s] unregistered! ", l.logPrefix, id)
		warnResponse(writer, http.StatusNotFound, core.CodeParamInvalid, warnInfo, l.logger)
		return
	}

	// 提取更新数据
	decryptedData, isDetail, err := l.getData(request, &entity)
	if nil != err {
		warnInfo := fmt.Sprintf("%s Update Fail: get data error! %v", l.logPrefix, err)
		warnResponse(writer, http.StatusBadRequest, core.CodeParamDecrypt, warnInfo, l.logger)
		return

	}
	info, detail := l.toInfo(decryptedData, isDetail)
	if isDetail {
		l.updateDetailStatus(writer, detail)
	} else {
		l.updateInfoStatus(writer, info)
	}
}

func (l *serverUpdateHandler) getId(request *http.Request) (id string, ok bool) {
	var bs []byte
	var err error
	if request.Method == http.MethodPost {
		bs, err = base64FromPost(request, core.HttpKeyId)
	} else {
		bs, err = base64FromGet(request, core.HttpKeyId)
	}
	if nil != err {
		return "", false
	}
	return string(bs), true
}

func (l *serverUpdateHandler) getData(request *http.Request, entity *home.RegisteredEntity) (decryptedData []byte, isDetail bool, err error) {
	var bs []byte
	if request.Method == http.MethodPost {
		bs, err = base64FromPost(request, core.HttpKeyData)
	} else {
		bs, err = base64FromGet(request, core.HttpKeyData)
	}
	if nil != err {
		return nil, false, err
	}
	cipher, aesOn := entity.GetAesCipher()
	if aesOn && nil != cipher {
		decryptedData, err = cipher.DecryptGCM(bs)
		if nil != err {
			return nil, false, err
		}
	} else {
		decryptedData = bs
	}
	isDetail = l.checkDetail(request)
	return decryptedData, isDetail, nil
}

func (l *serverUpdateHandler) toInfo(decryptedData []byte, isDetail bool) (info *core.UpdateInfo, detail *core.UpdateDetailInfo) {
	if isDetail {
		detail = &core.UpdateDetailInfo{}
		jsoniter.Unmarshal(decryptedData, detail)
	} else {
		info = &core.UpdateInfo{}
		jsoniter.Unmarshal(decryptedData, info)
	}
	return
}

func (l *serverUpdateHandler) checkDetail(request *http.Request) (isDetail bool) {
	var value []byte
	var err error
	if request.Method == http.MethodPost {
		value, err = base64FromPost(request, core.HttpKeyDetail)
	} else {
		value, err = base64FromGet(request, core.HttpKeyDetail)
	}
	if nil != err || len(value) == 0 { // 解释出错 | 无值
		return false
	}
	dStr := strings.ToLower(string(value))
	return dStr == "true" || dStr == "1"
}

func (l *serverUpdateHandler) updateInfoStatus(writer http.ResponseWriter, info *core.UpdateInfo) {
	if info.IsNotValid() {
		warnInfo := fmt.Sprintf("%s Update InfoStatus Fail: info is not valid. %v", l.logPrefix, info)
		warnResponse(writer, http.StatusBadRequest, core.CodeParamInvalid, warnInfo, l.logger)
		return
	}
	ok := home.GlobalHomeServer.UpdateState(*info)
	if !ok {
		warnInfo := fmt.Sprintf("%s Update InfoStatus Fail: Id[%s] unregistered! ", l.logPrefix, info.Id)
		warnResponse(writer, http.StatusNotFound, core.CodeParamInvalid, warnInfo, l.logger)
		return
	}
	sucResponseEmpty(writer)
	l.logger.Infoln(l.logPrefix, fmt.Sprintf("Update InfoStatus Succ: %v", info))
}

func (l *serverUpdateHandler) updateDetailStatus(writer http.ResponseWriter, detail *core.UpdateDetailInfo) {
	if detail.IsNotValid() {
		warnInfo := fmt.Sprintf("%s Update DetailStatus Fail: State is not valid. %v", l.logPrefix, detail)
		warnResponse(writer, http.StatusBadRequest, core.CodeParamInvalid, warnInfo, l.logger)
		return
	}
	ok := home.GlobalHomeServer.UpdateDetailState(*detail)
	if !ok {
		warnInfo := fmt.Sprintf("%s Update DetailStatus Fail: Id[%s] unregistered! ", l.logPrefix, detail.Id)
		warnResponse(writer, http.StatusNotFound, core.CodeParamInvalid, warnInfo, l.logger)
		return
	}
	sucResponseEmpty(writer)
	l.logger.Infoln(l.logPrefix, fmt.Sprintf("Update DetailStatus Succ: %v", detail))
}
