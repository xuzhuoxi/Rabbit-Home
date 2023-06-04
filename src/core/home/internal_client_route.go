// Package home
// Create on 2023/6/4
// @author xuzhuoxi
package home

import (
	"net/http"
)

func NewClientRouteHandler() http.Handler {
	return &clientRouteHandler{post: clientPost}
}

type clientRouteHandler struct {
	post bool
}

func (l *clientRouteHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	Logger.Infoln("clientRouteHandler")
}
