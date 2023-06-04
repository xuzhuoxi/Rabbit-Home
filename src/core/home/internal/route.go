// Package internal
// Create on 2023/6/4
// @author xuzhuoxi
package internal

import (
	"fmt"
	"net/http"
)

func NewRouteHandler() http.Handler {
	return &routeHandler{}
}

type routeHandler struct {
}

func (l *routeHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("routeHandler")
}
