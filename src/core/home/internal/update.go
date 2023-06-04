// Package internal
// Create on 2023/6/4
// @author xuzhuoxi
package internal

import (
	"fmt"
	"net/http"
)

func NewUpdateHandler() http.Handler {
	return &updateHandler{}
}

type updateHandler struct {
}

func (l *updateHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("updateHandler")
}
