// Package home
// Create on 2023/6/4
// @author xuzhuoxi
package home

import (
	"errors"
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Home/core"
	"github.com/xuzhuoxi/infra-go/netx/httpx"
	"net/http"
	"sync"
)

// IRabbitHomeServer Rabbit-Home 服务器接口
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

	IEntityGetter
	IEntitySetter
	IEntityStateUpdate
	IEntityQuery
	// GetHomeKeys 获取密钥管理器
	GetHomeKeys() *RabbitHomeKeys
}

// NewRabbitHomeServer
// 创建一个 Rabbit-Home 服务器实例
func NewRabbitHomeServer() IRabbitHomeServer {
	return &RabbitHomeServer{EntityList: NewEntityList(), Keys: NewRabbitHomeKeys()}
}

// RabbitHomeServer 服务器实例
type RabbitHomeServer struct {
	EntityList IEntityList       // 已经注册的实体列表
	Keys       *RabbitHomeKeys   // 密钥管理器
	HttpServer *httpx.HttpServer // Http服务器
	lock       sync.RWMutex
}

func (o *RabbitHomeServer) String() string {
	return fmt.Sprintf("{Running=%v, ListenAddr=%s, Size=%d}",
		o.HttpServer.Running(), o.HttpServer.Server.Addr, o.EntityList.Size())
}

func (o *RabbitHomeServer) GetHomeKeys() *RabbitHomeKeys {
	return o.Keys
}

func (o *RabbitHomeServer) Init() {
	o.lock.Lock()
	defer o.lock.Unlock()
	if nil != o.HttpServer {
		return
	}
	o.HttpServer = httpx.NewHttpServer().(*httpx.HttpServer)
	for _, handler := range MapHandlerList {
		o.HttpServer.MapHandle(handler.Pattern, handler.Handler())
	}
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

// -------------------------

func (o *RabbitHomeServer) GetEntityById(id string) (entity RegisteredEntity, ok bool) {
	return o.EntityList.GetEntityById(id)
}

func (o *RabbitHomeServer) GetEntities(funcEach FuncEach) (entities []RegisteredEntity) {
	return o.EntityList.GetEntities(funcEach)
}

func (o *RabbitHomeServer) GetEntityByName(name string) (entity []RegisteredEntity) {
	return o.EntityList.GetEntityByName(name)
}

func (o *RabbitHomeServer) GetEntitiesByPlatform(platformId string) (entities []RegisteredEntity) {
	return o.GetEntitiesByPlatform(platformId)
}

func (o *RabbitHomeServer) AddEntity(entity *RegisteredEntity) error {
	return o.EntityList.AddEntity(entity)
}

func (o *RabbitHomeServer) ReplaceEntity(entity *RegisteredEntity) error {
	return o.EntityList.ReplaceEntity(entity)
}

func (o *RabbitHomeServer) AddOrReplaceEntity(entity *RegisteredEntity) (add bool, err error) {
	return o.EntityList.AddOrReplaceEntity(entity)
}

func (o *RabbitHomeServer) RemoveEntity(id string) (entity *RegisteredEntity, ok bool) {
	return o.EntityList.RemoveEntity(id)
}

func (o *RabbitHomeServer) UpdateState(state core.UpdateInfo) bool {
	return o.EntityList.UpdateState(state)
}

func (o *RabbitHomeServer) UpdateDetailState(detail core.UpdateDetailInfo) bool {
	return o.EntityList.UpdateDetailState(detail)
}

func (o *RabbitHomeServer) QuerySmartEntity() (entity RegisteredEntity, ok bool) {
	return o.EntityList.QuerySmartEntity()
}

func (o *RabbitHomeServer) QueryEntity(name string, platformId string) (entity RegisteredEntity, ok bool) {
	return o.EntityList.QueryEntity(name, platformId)
}

// -------------------------

func (o *RabbitHomeServer) start(addr string) error {
	if nil == o.HttpServer {
		return errors.New("HttpServer is not exist! ")
	}
	if o.HttpServer.Running() {
		return errors.New("HttpServer is running! ")
	}
	o.loadKeys()
	return o.HttpServer.StartServer(addr)
}
func (o *RabbitHomeServer) loadKeys() {
	if nil != o.Keys {
		internalVerifier := GlobalHomeConfig.InternalVerifier.KeyVerifier
		if internalVerifier.Enable {
			o.Keys.LoadInternalDefaultPublicKeys()
		}
		externalVerifier := GlobalHomeConfig.ExternalVerifier.KeyVerifier
		if externalVerifier.Enable {
			o.Keys.LoadSelfPrivateKey()
		}
	}
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

// Private

type sortWeightList []*RegisteredEntity

func (o sortWeightList) Len() int {
	return len(o)
}

func (o sortWeightList) Less(i, j int) bool {
	bi := o[i].IsTimeout()
	bj := o[j].IsTimeout()
	if bi == bj {
		if o[i].State.Weight != o[j].State.Weight {
			return o[i].State.Weight < o[j].State.Weight
		}
		return o[i].hit < o[j].hit
	} else {
		return bj
	}
}

func (o sortWeightList) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

type sortLinkList []*RegisteredEntity

func (o sortLinkList) Len() int {
	return len(o)
}

func (o sortLinkList) Less(i, j int) bool {
	if o[i].Detail.Links != o[j].Detail.Links {
		return o[i].Detail.Links < o[j].Detail.Links
	}
	return o[i].hit < o[j].hit
}

func (o sortLinkList) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}
