// Package core
// Create on 2025/4/12
// @author xuzhuoxi
package core

import (
	"errors"
	"fmt"
	"github.com/xuzhuoxi/infra-go/cryptox/symmetric"
	"strconv"
	"time"
)

var (
	LinkedTimeout = int64(5 * time.Minute) // 超时时间
)

type ISignatureInfo interface {
	OriginalSignData() []byte
	SignatureData() string
}

// Link ---------- ---------- ---------- ---------- ----------

// LinkInfo
// 连接到Rabbit-Home要求的信息
type LinkInfo struct {
	Id          string `json:"id"`           // 实例Id(唯一)
	PlatformId  string `json:"pid"`          // 平台Id
	TypeName    string `json:"type-name"`    // 实例类型名称(不唯一)
	OpenNetwork string `json:"open-network"` // 开放连接通信协议
	OpenAddr    string `json:"open-addr"`    // 开放连接地址
	OpenKeyOn   bool   `json:"open-key-on"`  // 针对客户端是否启用密钥验证
	Signature   string `json:"signature"`    // 签名
}

func (o *LinkInfo) String() string {
	return fmt.Sprintf("{Id=%s,PId=%s,Name=%s,Network=%s,Addr=%s,Signature=%s}",
		o.Id, o.PlatformId, o.TypeName, o.OpenNetwork, o.OpenAddr, o.Signature)
}

// IsInvalid 是否为未验证
func (o *LinkInfo) IsInvalid() bool {
	return len(o.Id) == 0 || len(o.TypeName) == 0
}

func (o *LinkInfo) OriginalSignData() []byte {
	original := []byte(Base64Encoding.EncodeToString([]byte(
		fmt.Sprintf("I=%s,P=%s,N=%s,ON=%s,OA=%s", o.Id, o.PlatformId, o.TypeName, o.OpenNetwork, o.OpenAddr))))
	return original
}

func (o *LinkInfo) SignatureData() string {
	return o.Signature
}

// LinkBackInfo
// 连接结果信息，从Rabbit-Home返回, 经过RSA加密
type LinkBackInfo struct {
	Id         string `json:"id"`          // 实例Id(唯一), 明文
	InternalSK []byte `json:"internal-sk"` // 临时RSA密钥, 用于内部加密, []byte
	OpenSK     []byte `json:"open-sk"`     // 临时RSA密钥, 用于外部加密, []byte
	Extend     string `json:"extend"`      // 扩展信息, 明文
}

// Unlink ---------- ---------- ---------- ---------- ----------

// UnlinkInfo
// 通知Rabbit-Home断开连接的信息
type UnlinkInfo struct {
	Id        string `json:"id"`        // 实例Id(唯一),明文
	Signature string `json:"signature"` // 签名
}

func (o *UnlinkInfo) OriginalSignData() []byte {
	return []byte(Base64Encoding.EncodeToString([]byte(o.Id)))
}

func (o *UnlinkInfo) SignatureData() string {
	return o.Signature
}

// UnlinkBackInfo
// 通知Rabbit-Home断开连接的结果信息，从Rabbit-Home返回
type UnlinkBackInfo struct {
	Id     string `json:"id"`     // 实例Id(唯一)，明文
	Extend string `json:"extend"` // 扩展信息，明文
}

// Update ---------- ---------- ---------- ---------- ----------

// UpdateInfo 实例状态
type UpdateInfo struct {
	Id     string  `json:"id"`     // 实例Id
	Weight float64 `json:"weight"` // 压力系数
}

func (o *UpdateInfo) String() string {
	return fmt.Sprintf("{Id=%s,Weight=%v}", o.Id, o.Weight)
}

// IsNotValid 是否为未验证
func (o *UpdateInfo) IsNotValid() bool {
	return len(o.Id) == 0
}

// UpdateDetailInfo 实例详细状态
type UpdateDetailInfo struct {
	Id             string `json:"id"`           // 实例Id
	StartTimestamp int64  `json:"start"`        // 启动时间戳(纳秒)
	StatsInterval  int64  `json:"sta-interval"` // 统计间隔

	MaxLinks      uint64 `json:"max-links"`  // 最大连接数
	TotalReqCount int64  `json:"total-reg"`  // 总请求数
	TotalRespTime int64  `json:"total-resp"` // 总响应时间
	MaxRespTime   int64  `json:"max-resp"`   // 最大响应时间(纳秒)
	Links         uint64 `json:"links"`      // 连接数

	StatsTimestamp    int64 `json:"sta-start"` // 统计开始时间戳(纳秒)
	StatsReqCount     int64 `json:"sta-req"`   // 统计请求数
	StatsRespUnixNano int64 `json:"sta-resp"`  // 统计响应时间(纳称)

	EnableKeys string `json:"enable-keys"` // 属性启用标记
}

func (o *UpdateDetailInfo) String() string {
	return fmt.Sprintf("{Id=%s,Start=%v,MaxLink=%v,Link=%v}",
		o.Id, o.StartTimestamp, o.MaxLinks, o.Links)
}

// IsNotValid 是否为未验证
func (o *UpdateDetailInfo) IsNotValid() bool {
	return len(o.Id) == 0
}

// Route ---------- ---------- ---------- ---------- ----------

type QueryRouteInfo struct {
	PlatformId string `json:"pid"`       // 服务平台Id
	TypeName   string `json:"type-name"` // 类型名称
	TempAesKey []byte `json:"temp-key"`  // 临时AES密钥，用于Rabbit-Home返回数据时加密, 如果不提供，返回的密钥数据将以Base64字符串返回
}

type QueryRouteBackInfo struct {
	Id           string `json:"id"`           // 实例Id(唯一)
	PlatformId   string `json:"pid"`          // 服务平台Id
	TypeName     string `json:"type-name"`    // 实例类型名称(不唯一)
	OpenNetwork  string `json:"open-network"` // 开放连接通信协议
	OpenAddr     string `json:"open-addr"`    // 开放连接地址
	OpenKeyOn    bool   `json:"open-key-on"`  // 针对客户端是否启用密钥验证
	OpenBase64SK string `json:"open-sk"`      // 临时密钥的Base64字符串表示，用于对称加密数据。如果请求时有设置临时密钥，则经过加密
	OpenSK       []byte // 临时密钥，执行ComputeOpenSK后更新，用于对称加密数据
}

func (q *QueryRouteBackInfo) ComputeOpenSK(tempAesKey []byte) error {
	if len(q.OpenBase64SK) == 0 || !q.OpenKeyOn {
		return nil
	}
	openSK, err := Base64Encoding.DecodeString(q.OpenBase64SK)
	if nil != err {
		return err
	}
	if len(tempAesKey) == 0 {
		q.OpenSK = openSK
		return nil
	}
	if len(tempAesKey) == 32 {
		openSK, err = symmetric.NewAESCipher(tempAesKey).Decrypt(openSK)
		if nil != err {
			return err
		}
		q.OpenSK = openSK
		return nil
	}
	return errors.New("TempAesKey size error:" + strconv.Itoa(len(tempAesKey)))
}

func (q *QueryRouteBackInfo) String() string {
	return fmt.Sprintf("{Id=%s，PId=%s, T=%s, Network=%s, Addr=%s, OpenBase64SK=%s}",
		q.Id, q.PlatformId, q.TypeName, q.OpenNetwork, q.OpenAddr, q.OpenBase64SK)
}
