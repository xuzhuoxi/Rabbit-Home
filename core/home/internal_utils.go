// Package home
// Create on 2023/6/4
// @author xuzhuoxi
package home

import (
	"encoding/base64"
	jsoniter "github.com/json-iterator/go"
	"net/http"
)

const (
	serverPost = false
	clientPost = false
)

func getValueWithPost(request *http.Request, key string, value interface{}) error {
	str, err := getStringWithPost(request, key)
	if nil != err {
		return err
	}
	return jsoniter.Unmarshal([]byte(str), value)
}

func getStringWithPost(request *http.Request, key string) (value string, err error) {
	if err := request.ParseForm(); err != nil {
		return "", err
	}
	val, err1 := fromBase64(request.PostFormValue(key))
	if nil != err1 {
		return "", err1
	}
	return val, nil
}

func getValueWithGet(request *http.Request, key string, value interface{}) error {
	str, err := getStringWithGet(request, key)
	if nil != err {
		return err
	}
	return jsoniter.Unmarshal([]byte(str), value)
}

func getStringWithGet(request *http.Request, key string) (value string, err error) {
	base64Str := request.FormValue(key)
	if len(base64Str) == 0 {
		return "", nil
	}
	val, err1 := fromBase64(base64Str)
	if nil != err1 {
		return "", err1
	}
	return val, nil
}

func fromBase64(base64Str string) (str string, err error) {
	s, err1 := base64.StdEncoding.DecodeString(base64Str)
	if nil != err {
		return "", err1
	}
	return string(s), nil
}

func toBase64(str string) (base64Str string, err error) {
	return base64.StdEncoding.EncodeToString([]byte(str)), nil
}
