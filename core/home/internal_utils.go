// Package home
// Create on 2023/6/4
// @author xuzhuoxi
package home

import (
	"encoding/base64"
	"github.com/json-iterator/go"
	"github.com/xuzhuoxi/infra-go/logx"
	"net"
	"net/http"
)

var empty = []byte("")

func warnAndResponse(writer http.ResponseWriter, statusCode int, warnInfo string, logger logx.ILogger) {
	logger.Warnln("[warnAndResponse]", warnInfo)
	writer.WriteHeader(statusCode)
	writer.Write([]byte(warnInfo))
}

func getValueWithPost(request *http.Request, key string, value interface{}) error {
	bs, err := getStringWithPost(request, key)
	if nil != err {
		return err
	}
	return jsoniter.Unmarshal(bs, value)
}

func getStringWithPost(request *http.Request, key string) (value []byte, err error) {
	if err := request.ParseForm(); err != nil {
		return empty, err
	}
	val, err1 := fromBase64(request.PostFormValue(key))
	if nil != err1 {
		return empty, err1
	}
	return val, nil
}

func getValueWithGet(request *http.Request, key string, value interface{}) error {
	bs, err := getStringWithGet(request, key)
	if nil != err {
		return err
	}
	return jsoniter.Unmarshal(bs, value)
}

func getStringWithGet(request *http.Request, key string) (value []byte, err error) {
	base64Str := request.FormValue(key)
	if len(base64Str) == 0 {
		return empty, nil
	}
	val, err1 := fromBase64(base64Str)
	if nil != err1 {
		return empty, err1
	}
	return val, nil
}

func fromBase64(base64Str string) (value []byte, err error) {
	s, err1 := base64.StdEncoding.DecodeString(base64Str)
	if nil != err {
		return empty, err1
	}
	return s, nil
}

func toBase64(bs []byte) (base64Str string) {
	return base64.StdEncoding.EncodeToString(bs)
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
