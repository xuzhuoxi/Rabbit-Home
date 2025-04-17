// Package service
// Create on 2023/6/4
// @author xuzhuoxi
package service

import (
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Home/core"
	"github.com/xuzhuoxi/Rabbit-Home/core/home"
	"github.com/xuzhuoxi/Rabbit-Home/core/utils"
	"github.com/xuzhuoxi/infra-go/cryptox/asymmetric"
	"github.com/xuzhuoxi/infra-go/logx"
	"net/http"
	"strconv"
)

func NewServiceLinkHandler() http.Handler {
	return &serverLinkHandler{
		logPrefix: "[serverLinkHandler.ServeHTTP]",
		logger:    home.GlobalLogger}
}

type serverLinkHandler struct {
	logPrefix string
	logger    logx.ILogger
}

func (l *serverLinkHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// 验证数据模式与IP
	if !home.GlobalHomeConfig.InternalVerifier.VerifyPost(request) {
		return
	}
	if !home.GlobalHomeConfig.VerifyInternalIP(getClientIpAddr(request)) { // 验证是否内部IP
		return
	}

	// 从请求中获取注册数据并完成Base64解码
	linkInfo := &core.LinkInfo{}
	var err0 error
	var err1 error
	var weightBs []byte
	if request.Method == http.MethodPost {
		err0 = loadValueFromPost(request, core.HttpKeyData, linkInfo)
		weightBs, err1 = base64FromPost(request, core.HttpKeyWeight)
	} else {
		err0 = loadValueFromGet(request, core.HttpKeyData, linkInfo)
		weightBs, err1 = base64FromGet(request, core.HttpKeyWeight)
	}
	if nil != err0 || nil != err1 {
		warnInfo := fmt.Sprintf("%s Link Entity Fail: %v,%v", l.logPrefix, err0, err1)
		warnResponse(writer, http.StatusBadRequest, core.CodeParamError, warnInfo, l.logger)
		return
	}
	if linkInfo.IsInvalid() {
		warnInfo := fmt.Sprintf("%s Link Entity Fail: Entity is invalid. %v", l.logPrefix, linkInfo)
		warnResponse(writer, http.StatusBadRequest, core.CodeParamInvalid, warnInfo, l.logger)
		return
	}

	home.GlobalLock.Lock()
	defer home.GlobalLock.Unlock()

	// 验证签名
	kv := home.GlobalHomeConfig.InternalVerifier.KeyVerifier
	var rsa asymmetric.IRSAPublicCipher
	var SK []byte
	var pass bool
	if kv.Enable {
		rsa, SK, pass = home.GlobalHomeServer.GetHomeKeys().VerifyLinkSign(linkInfo)
		if !pass {
			warnInfo := fmt.Sprintf("%s Link Entity Fail: Key verify fail. %v", l.logPrefix, linkInfo)
			warnResponse(writer, http.StatusBadRequest, core.CodeVerifyLinkSignFail, warnInfo, l.logger)
			return
		}
		l.logger.Infoln(l.logPrefix, fmt.Sprintf("Id(%s)TempKey(%d):%v", linkInfo.Id, len(SK), SK))
	}

	// 更新实例列表
	entity := home.NewRegisteredEntity(*linkInfo)
	// 如果有w字段，更新权重
	if len(weightBs) > 1 {
		weight, err := strconv.ParseFloat(string(weightBs), 64)
		if nil == err {
			entity.UpdateState(core.UpdateInfo{Id: linkInfo.Id, Weight: weight})
		}
	}
	add, err1 := home.GlobalHomeServer.AddOrReplaceEntity(entity)
	if nil != err1 {
		warnInfo := fmt.Sprintf("%s Link Entity Fail: %v", l.logPrefix, err1)
		warnResponse(writer, http.StatusBadRequest, core.CodeLinkEntityFail, warnInfo, l.logger)
		return
	}
	backInfo := &core.LinkBackInfo{Id: linkInfo.Id}

	// 生成内部(Rabbit-Home与Rabbit-Server间通信)使用的临时32位密钥
	if kv.Enable {
		entity.SaveInternalSK(SK)
		backInfo.InternalSK = SK
	}

	// 生成外部(Rabbit-Server与客户端间通信)临时密钥
	if linkInfo.OpenKeyOn {
		openSK := utils.DeriveRandomKey32("", linkInfo.Id)
		entity.SaveOpenSK(openSK)
		backInfo.OpenSK = openSK
	}
	if add {
		l.logger.Infoln(l.logPrefix, fmt.Sprintf("Link Entity Succ: %v", linkInfo))
	} else {
		l.logger.Infoln(l.logPrefix, fmt.Sprintf("Relink Entity Succ: %v", linkInfo))
	}
	sucResponse(writer, backInfo, rsa, l.logger)
}
