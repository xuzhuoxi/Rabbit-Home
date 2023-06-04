// Package internal
// Create on 2023/6/4
// @author xuzhuoxi
package internal

import (
	"fmt"
	"net/http"
)

func NewServerUpdateHandler() http.Handler {
	return &serverUpdateHandler{}
}

type serverUpdateHandler struct {
}

func (l *serverUpdateHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("serverUpdateHandler")
}
