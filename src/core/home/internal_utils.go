// Package home
// Create on 2023/6/4
// @author xuzhuoxi
package home

import (
	"encoding/base64"
	"github.com/json-iterator/go"
	"net/http"
)

const (
	serverPost = false
	clientPost = false
)

func getValueWithPost(request *http.Request, key string, value interface{}) error {
	if err := request.ParseForm(); err != nil {
		return err
	}
	val := []byte(request.PostFormValue(key))
	return jsoniter.Unmarshal(val, value)
}

func getStringWithPost(request *http.Request, key string) (value string, err error) {
	if err := request.ParseForm(); err != nil {
		return "", err
	}
	return request.PostFormValue(key), nil
}

func getValueWithGet(request *http.Request, key string, value interface{}) error {
	val, err := fromBase64(request.FormValue(key))
	if nil != err {
		return err
	}
	return jsoniter.Unmarshal([]byte(val), value)
}

func getStringWithGet(request *http.Request, key string) (value string, err error) {
	val := request.FormValue(key)
	return fromBase64(val)
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
