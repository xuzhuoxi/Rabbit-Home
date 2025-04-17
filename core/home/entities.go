// Package home
// Create on 2023/6/4
// @author xuzhuoxi
package home

import (
	"errors"
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Home/core"
	"sort"
	"sync"
)

type FuncEach = func(each RegisteredEntity) bool
type funcEach = func(each *RegisteredEntity) bool

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
	AddEntity(entity *RegisteredEntity) error
	// ReplaceEntity 替换实例
	ReplaceEntity(entity *RegisteredEntity) error
	// AddOrReplaceEntity 添加或替换实例
	AddOrReplaceEntity(entity *RegisteredEntity) (add bool, err error)
	// RemoveEntity 移除实例
	RemoveEntity(id string) (entity *RegisteredEntity, ok bool)
}

// IEntityStateUpdate 更新实例状态接口
type IEntityStateUpdate interface {
	// UpdateState 更新实例状态信息
	UpdateState(state core.UpdateInfo) bool
	// UpdateDetailState 更新实例状态详细信息
	UpdateDetailState(detail core.UpdateDetailInfo) bool
}

// IEntityQuery 查询实例接口
type IEntityQuery interface {
	// SetSortFunc 设置排序函数
	SetSortFunc(f FuncSortEntity)
	// QuerySmartEntity 查询最好的实例
	QuerySmartEntity(platformId string, typeName string) (entity RegisteredEntity, ok bool)
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
	return &EntityList{
		Entities: make([]*RegisteredEntity, 0, 128),
		funcSort: DefaultFuncSortEntity,
	}
}

type EntityList struct {
	Entities []*RegisteredEntity // 实例列表
	funcSort FuncSortEntity
	lock     sync.RWMutex // 读写锁
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
	index := o.findEntity(func(each *RegisteredEntity) bool {
		if each.Id == id {
			return true
		}
		return false
	})
	if index != -1 {
		return *o.Entities[index], true
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
		if o.Entities[index].TypeName == name {
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

func (o *EntityList) AddEntity(entity *RegisteredEntity) error {
	if entity.IsInvalid() {
		return errors.New("AddEntity Error: Entity not valid! ")
	}
	o.lock.Lock()
	defer o.lock.Unlock()
	index := o.findEntity(func(each *RegisteredEntity) bool {
		return each.Id == entity.Id
	})
	if index != -1 {
		return errors.New("AddEntity Error: Duplicate Entity exist! ")
	}
	o.Entities = append(o.Entities, entity)
	return nil
}

func (o *EntityList) ReplaceEntity(entity *RegisteredEntity) error {
	if entity.IsInvalid() {
		return errors.New("ReplaceEntity Error: Entity not valid! ")
	}
	o.lock.Lock()
	defer o.lock.Unlock()
	index := o.findEntity(func(each *RegisteredEntity) bool {
		return each.Id == entity.Id
	})
	if index == -1 {
		return errors.New("ReplaceEntity Error: No Entity exist! ")
	}
	o.Entities[index] = entity
	return nil
}

func (o *EntityList) AddOrReplaceEntity(entity *RegisteredEntity) (add bool, err error) {
	if entity.IsInvalid() {
		return false, errors.New("AddOrReplaceEntity Error: Entity not valid! ")
	}
	o.lock.Lock()
	defer o.lock.Unlock()
	index := o.findEntity(func(each *RegisteredEntity) bool {
		return each.Id == entity.Id
	})
	if index != -1 {
		o.Entities[index] = entity
	} else {
		o.Entities = append(o.Entities, entity)
		add = true
	}
	return
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

func (o *EntityList) UpdateState(state core.UpdateInfo) bool {
	if state.IsNotValid() {
		return false
	}
	o.lock.Lock()
	defer o.lock.Unlock()
	index := o.findEntity(func(each *RegisteredEntity) bool {
		return each.Id == state.Id
	})
	if index != -1 {
		o.Entities[index].UpdateState(state)
		return true
	}
	return false
}

func (o *EntityList) UpdateDetailState(detail core.UpdateDetailInfo) bool {
	if detail.IsNotValid() {
		return false
	}
	o.lock.Lock()
	defer o.lock.Unlock()
	index := o.findEntity(func(each *RegisteredEntity) bool {
		return each.Id == detail.Id
	})
	if index != -1 {
		o.Entities[index].UpdateDetailState(detail)
		return true
	}
	return false
}

func (o *EntityList) SetSortFunc(f FuncSortEntity) {
	if nil == f {
		return
	}
	o.lock.Lock()
	defer o.lock.Unlock()
	o.funcSort = f
}

func (o *EntityList) QuerySmartEntity(platformId string, typeName string) (entity RegisteredEntity, ok bool) {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.clearTimeoutEntities()
	entities := o.getEntities(func(each *RegisteredEntity) bool {
		return each.PlatformId == platformId && each.TypeName == typeName
	})
	return o.queryEntities(entities)
}

func (o *EntityList) findEntity(funcEach funcEach) (index int) {
	for index := range o.Entities {
		if funcEach(o.Entities[index]) {
			return index
		}
	}
	index = -1
	return
}

func (o *EntityList) getEntities(funcEach funcEach) (entities []*RegisteredEntity) {
	for index := range o.Entities {
		if funcEach(o.Entities[index]) {
			entities = append(entities, o.Entities[index])
		}
	}
	return
}

func (o *EntityList) queryEntities(entities []*RegisteredEntity) (entity RegisteredEntity, ok bool) {
	if len(entities) == 0 {
		return EntityEmpty, false
	}
	sort.Slice(entities, func(i, j int) bool {
		return o.funcSort(entities[i], entities[j])
	})
	fmt.Println("Sort:", entities)
	entities[0].AddHit()
	entity = *entities[0]
	return entity, true
}

func (o *EntityList) clearTimeoutEntities() {
	for index := len(o.Entities) - 1; index >= 0; index -= 1 {
		if o.Entities[index].IsTimeout() {
			o.Entities = append(o.Entities[:index], o.Entities[index+1:]...)
		}
	}
}
