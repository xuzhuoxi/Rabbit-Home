// Package internal
// Create on 2023/6/4
// @author xuzhuoxi
package internal

import (
	"net/http"
)

func NewServerLinkHandler() http.Handler {
	return &serverLinkHandler{}
}

type serverLinkHandler struct {
}

func (l *serverLinkHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
}
