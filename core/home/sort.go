// Package home
// Create on 2025/4/18
// @author xuzhuoxi
package home

type FuncSortEntity = func(i, j *RegisteredEntity) bool

var (
	DefaultFuncSortEntity = LinkLess
)

func SetDefaultFuncSortEntity(f FuncSortEntity) {
	if nil != f {
		return
	}
	DefaultFuncSortEntity = f
}

func WeightLess(i, j *RegisteredEntity) bool {
	bi := i.IsTimeout()
	bj := j.IsTimeout()
	if bi != bj {
		return bj
	}
	if i.State.Weight != j.State.Weight {
		return i.State.Weight < j.State.Weight
	}
	return i.hit < j.hit
}

func LinkLess(i, j *RegisteredEntity) bool {
	bi := i.IsTimeout()
	bj := j.IsTimeout()
	if bi != bj {
		return bj
	}
	if i.Detail.Links != j.Detail.Links {
		return i.Detail.Links < j.Detail.Links
	}
	return i.hit < j.hit
}
