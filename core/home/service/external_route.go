// Package service
// Create on 2023/6/4
// @author xuzhuoxi
package service

import (
	"fmt"
	"github.com/json-iterator/go"
	"github.com/xuzhuoxi/Rabbit-Home/core"
	"github.com/xuzhuoxi/Rabbit-Home/core/home"
	"github.com/xuzhuoxi/infra-go/cryptox/symmetric"
	"github.com/xuzhuoxi/infra-go/logx"
	"net/http"
)

func NewServiceRouteHandler() http.Handler {
	return &clientRouteHandler{
		logPrefix: "[serverLinkHandler.ServeHTTP]",
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

	// 从请求中获取查询数据并完成Base64解码
	var data []byte
	var err error
	if request.Method == http.MethodPost {
		data, err = base64FromPost(request, core.HttpKeyQuery)
	} else {
		data, err = base64FromGet(request, core.HttpKeyQuery)
	}
	if nil != err {
		warnInfo := "Param base64 error! " + err.Error()
		warnResponse(writer, http.StatusBadRequest, core.CodeParamBase64, warnInfo, l.logger)
		return
	}

	// 对数据进行解密
	kv := home.GlobalHomeConfig.ExternalVerifier.KeyVerifier
	if kv.Enable {
		privateCipher := home.GlobalHomeServer.GetHomeKeys().PrivateCipher()
		if nil == privateCipher {
			warnInfo := "Private Key Not Exist!"
			warnResponse(writer, http.StatusBadRequest, core.CodePrivateKeyLack, warnInfo, l.logger)
			return
		}
		data, err = privateCipher.Decrypt(data)
		if nil != err {
			warnInfo := "Private Key DecryptMode Fail: " + err.Error()
			warnResponse(writer, http.StatusBadRequest, core.CodeParamDecrypt, warnInfo, l.logger)
			return
		}
	}

	// 解析数据
	query := &core.QueryRouteInfo{}
	err = jsoniter.Unmarshal(data, query)
	if nil != err {
		warnInfo := "Data Unmarshal Fail: " + err.Error()
		warnResponse(writer, http.StatusBadRequest, core.CodeParamJson, warnInfo, l.logger)
		return
	}

	// 执行查询
	entity, ok := home.GlobalHomeServer.QuerySmartEntity(query.PlatformId, query.TypeName)
	fmt.Println("结果：", &entity)
	if !ok {
		warnInfo := "Query Fail. Cannot query a smart entity." + request.RemoteAddr
		warnResponse(writer, http.StatusNotFound, core.CodeEntityQueryFail, warnInfo, l.logger)
		return
	}

	openBase64SK := ""
	// 启用外部密钥处理
	if entity.OpenKeyOn {
		keyLen := len(query.TempAesKey)
		// 密钥长度不合法
		if keyLen != 0 && keyLen != 32 {
			warnInfo := fmt.Sprintf("Query Fail. TempAesKey len(%d) is invalid. %v", keyLen, request.RemoteAddr)
			warnResponse(writer, http.StatusBadRequest, core.CodeEntityQueryFailKey, warnInfo, l.logger)
			return
		}
		// 密钥合法，加密OpenSK
		if 32 == keyLen {
			cipher := symmetric.NewAESCipher(query.TempAesKey)
			encryptOpenSK, err := cipher.Encrypt(entity.GetOpenSK())
			if nil != err {
				warnInfo := "Query Fail. TempAesKey EncryptMode fail. " + request.RemoteAddr
				warnResponse(writer, http.StatusBadRequest, core.CodeEntityQueryFailKey, warnInfo, l.logger)
				return
			}
			openBase64SK = core.Base64Encoding.EncodeToString(encryptOpenSK)
		} else {
			openBase64SK = entity.GetOpenBase64SK()
		}
	}

	// 返回数据
	backInfo := &core.QueryRouteBackInfo{
		Id:           entity.Id,
		PlatformId:   entity.PlatformId,
		TypeName:     entity.TypeName,
		OpenNetwork:  entity.OpenNetwork,
		OpenAddr:     entity.OpenAddr,
		OpenKeyOn:    entity.OpenKeyOn,
		OpenBase64SK: openBase64SK,
	}
	sucResponse(writer, backInfo, nil, l.logger)
	home.GlobalLogger.Infoln(l.logPrefix, fmt.Sprintf("Query Succ from %s. Return %s", request.RemoteAddr, entity.Id))
}
