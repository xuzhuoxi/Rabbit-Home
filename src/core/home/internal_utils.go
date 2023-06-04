// Package home
// Create on 2023/6/4
// @author xuzhuoxi
package home

import (
	jsoniter "github.com/json-iterator/go"
	"net/http"
)

const (
	serverPost = false
	clientPost = false
)

func getValueWithGet(request *http.Request, key string, value interface{}) bool {
	val := []byte(request.FormValue(key))
	jsoniter.Unmarshal(val, value)
	return true
}

func getValueWithPost(request *http.Request, key string, value interface{}) bool {
	if err := request.ParseForm(); err != nil {
		Logger.Infoln(err)
		return false
	}
	val := []byte(request.PostFormValue(key))
	jsoniter.Unmarshal(val, value)
	return true
}

func getStringWithGet(request *http.Request, key string) (value string, err error) {
	return base62ToString(request.FormValue(key))
}

func getStringWithPost(request *http.Request, key string) (value string, err error) {
	if err := request.ParseForm(); err != nil {
		return "", err
	}
	return base62ToString(request.PostFormValue(key))
}

func base62ToString(base64 string) (str string, err error) {
	return base64, nil
}
