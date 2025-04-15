// Package service
// Create on 2023/6/4
// @author xuzhuoxi
package service

import (
	"github.com/json-iterator/go"
	"github.com/xuzhuoxi/Rabbit-Home/core"
	"github.com/xuzhuoxi/Rabbit-Home/core/utils"
	"github.com/xuzhuoxi/infra-go/logx"
	"net"
	"net/http"
)

var (
	empty = []byte("")
)

func warnResponse(writer http.ResponseWriter, httpStatusCode int, extCode int, warnInfo string, logger logx.ILogger) {
	if nil != logger {
		logger.Warnln("[warnResponse]", warnInfo)
	}
	writer.WriteHeader(httpStatusCode)
	base64Data := utils.SerializeHomeResponseInfo(core.HomeResponseInfo{ExtCode: extCode, Info: warnInfo})
	writer.Write(base64Data)
}

func sucResponse(writer http.ResponseWriter, respData interface{}, logger logx.ILogger) {
	json, err := jsoniter.Marshal(respData)
	if nil != err {
		if nil != logger {
			logger.Warnln("[sucResponse]", err)
		}
		json = empty
	}
	base64Data := core.Base64Encoding.EncodeToString(json)
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(base64Data))
}
func sucResponseEmpty(writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusOK)
	writer.Write(empty)
}

func loadValueFromPost(request *http.Request, key string, value interface{}) error {
	bs, err := base64FromPost(request, key)
	if nil != err {
		return err
	}
	return jsoniter.Unmarshal(bs, value)
}

func loadValueFromGet(request *http.Request, key string, value interface{}) error {
	bs, err := base64FromGet(request, key)
	if nil != err {
		return err
	}
	return jsoniter.Unmarshal(bs, value)
}

func base64FromPost(request *http.Request, key string) (value []byte, err error) {
	if err := request.ParseForm(); err != nil {
		return empty, err
	}
	val, err1 := fromBase64(request.PostFormValue(key))
	if nil != err1 {
		return empty, err1
	}
	return val, nil
}

func base64FromGet(request *http.Request, key string) (value []byte, err error) {
	base64Str := request.FormValue(key)
	if len(base64Str) == 0 {
		return empty, nil
	}
	val, err1 := fromBase64(base64Str)
	if nil != err1 {
		return empty, err1
	}
	//fmt.Println("base64FromGet:", key, base64Str, val)
	return val, nil
}

func fromBase64(base64Str string) (value []byte, err error) {
	s, err1 := core.Base64Encoding.DecodeString(base64Str)
	if nil != err {
		return empty, err1
	}
	return s, nil
}

func toBase64(bs []byte) (base64Str string) {
	return core.Base64Encoding.EncodeToString(bs)
}

func getClientIpAddr(req *http.Request) string {
	// 从X-Forwarded-For请求头获取
	ip := req.Header.Get("X-Forwarded-For")
	if ip == "" {
		// 从X-Real-IP请求头获取
		ip = req.Header.Get("X-Real-IP")
	}
	if ip == "" {
		// 如果请求头中没有，直接从RemoteAddr获取
		ip, _, _ = net.SplitHostPort(req.RemoteAddr)
	}
	return ip
}
