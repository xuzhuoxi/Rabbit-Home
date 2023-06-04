// Package internal
// Create on 2023/6/4
// @author xuzhuoxi
package internal

import (
	"fmt"
	"net/http"
)

func NewLogoutHandler() http.Handler {
	return &logoutHandler{}
}

type logoutHandler struct {
}

func (l *logoutHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("logoutHandler")
}
