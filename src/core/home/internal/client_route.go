// Package internal
// Create on 2023/6/4
// @author xuzhuoxi
package internal

import (
	"fmt"
	"net/http"
)

func NewClientRouteHandler() http.Handler {
	return &clientRouteHandler{}
}

type clientRouteHandler struct {
}

func (l *clientRouteHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("clientRouteHandler")
}
