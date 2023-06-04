// Package internal
// Create on 2023/6/4
// @author xuzhuoxi
package internal

import (
	"fmt"
	"net/http"
)

func NewServerUnlinkHandler() http.Handler {
	return &serverUnlinkHandler{}
}

type serverUnlinkHandler struct {
}

func (l *serverUnlinkHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("serverUnlinkHandler")
}
