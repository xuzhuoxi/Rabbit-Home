// Package home
// Create on 2023/6/4
// @author xuzhuoxi
package home

import (
	"errors"
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Home/src/core/home/internal"
	"github.com/xuzhuoxi/infra-go/netx"
	"net/http"
	"sync"
)

// IRabbitHomeServer 服务器接口
type IRabbitHomeServer interface {
	// Init 初始化
	Init()
	// MapHandle 绑定Pattern处理器
	MapHandle(pattern string, handler http.Handler)
	// MapFunc 绑定Pattern处理函数
	MapFunc(pattern string, f func(w http.ResponseWriter, r *http.Request))
	// Start 启用服务器
	Start() error
	// StartByPort 启用服务器
	StartByPort(port int) error
	// StartByAddr 启用服务器
	StartByAddr(addr string) error
	// Stop 停止服务器
	Stop() error

	// GetEntityList 取实例列表
	GetEntityList() IEntityList
}

func NewRabbitHomeServer() IRabbitHomeServer {
	return &RabbitHomeServer{EntityList: NewEntityList()}
}

// RabbitHomeServer 服务器实例
type RabbitHomeServer struct {
	EntityList IEntityList
	HttpServer *netx.HttpServer
	lock       sync.RWMutex
}

func (o *RabbitHomeServer) String() string {
	return fmt.Sprintf("{Running=%v, ListenAddr=%s, Size=%d}",
		o.HttpServer.Running(), o.HttpServer.Server.Addr, o.EntityList.Size())
}

func (o *RabbitHomeServer) GetEntityList() IEntityList {
	return o.EntityList
}

func (o *RabbitHomeServer) Init() {
	o.lock.Lock()
	defer o.lock.Unlock()
	if nil != o.HttpServer {
		return
	}
	o.HttpServer = netx.NewHttpServer().(*netx.HttpServer)
	o.HttpServer.MapHandle(PatternLogin, internal.NewLoginHandler())
	o.HttpServer.MapHandle(PatternLogout, internal.NewLogoutHandler())
	o.HttpServer.MapHandle(PatternUpdate, internal.NewUpdateHandler())
	o.HttpServer.MapHandle(PatternRoute, internal.NewRouteHandler())
}

func (o *RabbitHomeServer) MapHandle(pattern string, handler http.Handler) {
	o.lock.RLock()
	defer o.lock.RUnlock()
	o.HttpServer.MapHandle(pattern, handler)
}

func (o *RabbitHomeServer) MapFunc(pattern string, f func(w http.ResponseWriter, r *http.Request)) {
	o.lock.RLock()
	defer o.lock.RUnlock()
	o.HttpServer.MapFunc(pattern, f)
}

func (o *RabbitHomeServer) Start() error {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.start(DefaultAddr)
}

func (o *RabbitHomeServer) StartByPort(port int) error {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.start(fmt.Sprintf(":%d", port))
}

func (o *RabbitHomeServer) StartByAddr(addr string) error {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.start(addr)
}

func (o *RabbitHomeServer) Stop() error {
	o.lock.Lock()
	defer o.lock.Unlock()
	return o.stop()
}

func (o *RabbitHomeServer) start(addr string) error {
	if nil == o.HttpServer {
		return errors.New("HttpServer is not exist! ")
	}
	if o.HttpServer.Running() {
		return errors.New("HttpServer is running! ")
	}
	return o.HttpServer.StartServer(addr)
}

func (o *RabbitHomeServer) stop() error {
	if nil == o.HttpServer {
		return errors.New("HttpServer is not exist! ")
	}
	if !o.HttpServer.Running() {
		return errors.New("HttpServer is not running! ")
	}
	return o.HttpServer.StopServer()
}
