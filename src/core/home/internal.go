// Package home
// Create on 2023/6/4
// @author xuzhuoxi
package home

type sortWeightList []*RegisteredEntity

func (o sortWeightList) Len() int {
	return len(o)
}

func (o sortWeightList) Less(i, j int) bool {
	bi := o[i].IsTimeout()
	bj := o[j].IsTimeout()
	if bi == bj {
		return o[i].State.Weight < o[j].State.Weight
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
	return o[i].Detail.Links < o[j].Detail.Links
}

func (o sortLinkList) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}
