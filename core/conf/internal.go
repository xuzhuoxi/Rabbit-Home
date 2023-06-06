// Package conf
// Create on 2023/6/4
// @author xuzhuoxi
package conf

import (
	"errors"
	"strconv"
	"strings"
)

func newIPGroup() *ipGroup {
	return &ipGroup{Values: [8]int{}}
}

func newIPGroupFromAddr(addr string) *ipGroup {
	rs := &ipGroup{Values: [8]int{}}
	rs.FromStringOverride(addr)
	return rs
}

func newMultiIPGroupFromAddr(addrArr []string) []*ipGroup {
	var rs []*ipGroup
	for _, addr := range addrArr {
		rs = append(rs, newIPGroupFromAddr(addr))
	}
	return rs
}

type ipGroup struct {
	Values [8]int
}

func (o *ipGroup) ContainsAddr(ipAddr string) bool {
	temp := newIPGroup()
	temp.FromStringOverride(ipAddr)
	return o.ContainsGroup(temp)
}

func (o *ipGroup) ContainsGroup(group *ipGroup) bool {
	return group.GetIPMinValue(0) >= o.GetIPMinValue(0) &&
		group.GetIPMaxValue(0) <= o.GetIPMaxValue(0) &&
		group.GetIPMinValue(1) >= o.GetIPMinValue(1) &&
		group.GetIPMaxValue(1) <= o.GetIPMaxValue(1) &&
		group.GetIPMinValue(2) >= o.GetIPMinValue(2) &&
		group.GetIPMaxValue(2) <= o.GetIPMaxValue(2) &&
		group.GetIPMinValue(3) >= o.GetIPMinValue(3) &&
		group.GetIPMaxValue(3) <= o.GetIPMaxValue(3)
}

func (o *ipGroup) FromStringOverride(ipAddr string) error {
	if len(ipAddr) == 0 {
		return errors.New("ipAddr is empty. ")
	}
	arr := strings.Split(ipAddr, SepIp)
	if len(arr) != 4 {
		return errors.New("ipAddr format error. ")
	}
	for index := range arr {
		idx := strings.Index(arr[index], SepIpRange)
		if -1 == idx {
			v, err := strconv.Atoi(arr[index])
			if nil != err {
				return err
			}
			o.Values[index], o.Values[index+1] = v, v
		} else {
			v0, err0 := strconv.Atoi(arr[index][:idx])
			if nil != err0 {
				return err0
			}
			v1, err1 := strconv.Atoi(arr[index][idx+1:])
			if nil != err1 {
				return err1
			}
			o.Values[index], o.Values[index+1] = v0, v1
		}
	}
	return nil
}

func (o *ipGroup) GetIPValue(index int, min bool) int {
	if min {
		return o.Values[index*2]
	} else {
		return o.Values[index*2+1]
	}
}

func (o *ipGroup) GetIPMinValue(index int) int {
	return o.Values[index*2]
}

func (o *ipGroup) GetIPMaxValue(index int) int {
	return o.Values[index*2+1]
}
