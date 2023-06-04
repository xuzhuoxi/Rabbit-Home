// Package internal
// Create on 2023/6/4
// @author xuzhuoxi
package internal

import (
	"fmt"
	"net/http"
)

func NewLoginHandler() http.Handler {
	return &loginHandler{}
}

type loginHandler struct {
}

func (l *loginHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("loginHandler")
}
