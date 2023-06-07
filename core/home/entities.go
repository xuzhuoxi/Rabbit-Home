// Package home
// Create on 2023/6/4
// @author xuzhuoxi
package home

import (
	"errors"
	"github.com/xuzhuoxi/Rabbit-Home/core"
	"sort"
	"sync"
)

type FuncEach func(each RegisteredEntity) bool
type funcEach func(each *RegisteredEntity) bool

type IEntitySize interface {
	// Size 实例数量
	Size() int
}

// IEntityGetter 读取实例接口
type IEntityGetter interface {
	// GetEntityById 通过实例Id获取实例
	GetEntityById(id string) (entity RegisteredEntity, ok bool)
	// GetEntities 获取实例列表
	GetEntities(funcEach FuncEach) (entities []RegisteredEntity)
	// GetEntityByName 通过实例名称获取实例列表
	GetEntityByName(name string) (entity []RegisteredEntity)
	// GetEntitiesByPlatform 通过平台Id获取实例列表
	GetEntitiesByPlatform(platformId string) (entities []RegisteredEntity)
}

// IEntitySetter 设置实例接口
type IEntitySetter interface {
	// AddEntity 添加实例
	AddEntity(entity RegisteredEntity) error
	// RemoveEntity 移除实例
	RemoveEntity(id string) (entity *RegisteredEntity, ok bool)
}

// IEntityStateUpdate 更新实例状态接口
type IEntityStateUpdate interface {
	// UpdateState 更新实例状态信息
	UpdateState(state core.EntityState) bool
	// UpdateDetailState 更新实例状态详细信息
	UpdateDetailState(detail core.EntityDetailState) bool
}

// IEntityQuery 查询实例接口
type IEntityQuery interface {
	// QuerySmartEntity 查询最好的实例
	QuerySmartEntity() (entity RegisteredEntity, ok bool)
	// QueryEntity 根据name与platformId查询最好的实例
	QueryEntity(name string, platformId string) (entity RegisteredEntity, ok bool)
}

// IEntityList 实例列表
type IEntityList interface {
	IEntitySize
	IEntityGetter
	IEntitySetter
	IEntityStateUpdate
	IEntityQuery
}

func NewEntityList() IEntityList {
	return &EntityList{}
}

type EntityList struct {
	Entities []*RegisteredEntity // 实例列表
	lock     sync.RWMutex        // 读写锁
}

func (o *EntityList) Size() int {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return len(o.Entities)
}

func (o *EntityList) GetEntityById(id string) (entity RegisteredEntity, ok bool) {
	if len(id) == 0 {
		return EntityEmpty, false
	}
	o.lock.RLock()
	defer o.lock.RUnlock()
	found, ok := o.findEntity(func(each *RegisteredEntity) bool {
		if each.Id == id {
			return true
		}
		return false
	})
	if ok {
		return *found, true
	}
	return EntityEmpty, false
}

func (o *EntityList) GetEntities(funcEach FuncEach) (entities []RegisteredEntity) {
	if nil == funcEach {
		return nil
	}
	o.lock.RLock()
	defer o.lock.RUnlock()
	for index := range o.Entities {
		if funcEach(*o.Entities[index]) {
			entities = append(entities, *o.Entities[index])
		}
	}
	return
}

func (o *EntityList) GetEntityByName(name string) (entity []RegisteredEntity) {
	if len(name) == 0 {
		return nil
	}
	o.lock.RLock()
	defer o.lock.RUnlock()
	var rs []RegisteredEntity
	for index := range o.Entities {
		if o.Entities[index].Name == name {
			rs = append(rs, *o.Entities[index])
		}
	}
	return rs
}

func (o *EntityList) GetEntitiesByPlatform(platformId string) (entities []RegisteredEntity) {
	if len(platformId) == 0 {
		return nil
	}
	o.lock.RLock()
	defer o.lock.RUnlock()
	var rs []RegisteredEntity
	for index := range o.Entities {
		if o.Entities[index].PlatformId == platformId {
			rs = append(rs, *o.Entities[index])
		}
	}
	return rs
}

func (o *EntityList) AddEntity(entity RegisteredEntity) error {
	if entity.IsNotValid() {
		return errors.New("AddEntity Error: Entity not valid! ")
	}
	o.lock.Lock()
	defer o.lock.Unlock()
	_, ok := o.findEntity(func(each *RegisteredEntity) bool {
		return each.Id == entity.Id || each.Name == entity.Name
	})
	if ok {
		return errors.New("AddEntity Error: Duplicate Entity exist! ")
	}
	newEntity := entity
	o.Entities = append(o.Entities, &newEntity)
	return nil
}

func (o *EntityList) RemoveEntity(id string) (entity *RegisteredEntity, ok bool) {
	if len(id) == 0 {
		return
	}
	o.lock.Lock()
	defer o.lock.Unlock()
	for index := range o.Entities {
		if o.Entities[index].Id == id {
			entity = o.Entities[index]
			o.Entities = append(o.Entities[:index], o.Entities[index+1:]...)
			return entity, true
		}
	}
	return
}

func (o *EntityList) UpdateState(state core.EntityState) bool {
	if state.IsNotValid() {
		return false
	}
	o.lock.Lock()
	defer o.lock.Unlock()
	entity, ok := o.findEntity(func(each *RegisteredEntity) bool {
		return each.Id == state.Id
	})
	if ok {
		entity.UpdateState(state)
		return true
	}
	return false
}

func (o *EntityList) UpdateDetailState(detail core.EntityDetailState) bool {
	if detail.IsNotValid() {
		return false
	}
	o.lock.Lock()
	defer o.lock.Unlock()
	entity, ok := o.findEntity(func(each *RegisteredEntity) bool {
		return each.Id == detail.Id
	})
	if ok {
		entity.UpdateDetailState(detail)
		return true
	}
	return false
}

func (o *EntityList) QuerySmartEntity() (entity RegisteredEntity, ok bool) {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.clearTimeoutEntities()
	return o.queryEntities(o.Entities)
}

func (o *EntityList) QueryEntity(name string, platformId string) (entity RegisteredEntity, ok bool) {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.clearTimeoutEntities()
	entities := o.getEntities(func(each RegisteredEntity) bool {
		return each.Name == name && each.PlatformId == platformId
	})
	return o.queryEntities(entities)
}

func (o *EntityList) findEntity(funcEach funcEach) (entity *RegisteredEntity, ok bool) {
	for index := range o.Entities {
		if funcEach(o.Entities[index]) {
			return o.Entities[index], true
		}
	}
	ok = false
	return
}

func (o *EntityList) getEntities(funcEach FuncEach) (entities []*RegisteredEntity) {
	for index := range o.Entities {
		if funcEach(*o.Entities[index]) {
			entities = append(entities, o.Entities[index])
		}
	}
	return
}

func (o *EntityList) queryEntities(entities []*RegisteredEntity) (entity RegisteredEntity, ok bool) {
	if len(entities) == 0 {
		return EntityEmpty, false
	}
	sort.Sort(sortLinkList(entities))
	o.Entities[0].AddHit()
	entity = *o.Entities[0]
	return entity, true
}

func (o *EntityList) clearTimeoutEntities() {
	for index := len(o.Entities) - 1; index >= 0; index -= 1 {
		if o.Entities[index].IsTimeout() {
			o.Entities = append(o.Entities[:index], o.Entities[index+1:]...)
		}
	}
}
