// Package core
// Create on 2023/6/4
// @author xuzhuoxi
package core

import (
	"encoding/base64"
	"fmt"
)

type HomeResponseInfo struct {
	ExtCode int    `json:"code"`  // 扩展码
	Info    string `json:"value"` // 通常为通过base64转化的josn字符器
	Other   string `json:"other"` // 通常为通过base64转化的josn字符器
}

func (h HomeResponseInfo) String() string {
	return fmt.Sprintf("HomeResponseInfo{ExtCode=%d,Info='%s',Other='%s'}", h.ExtCode, h.Info, h.Other)
}

var (
	Base64Encoding = base64.RawURLEncoding
)

// Cmd ---------- ---------- ---------- ---------- ---------- ----------

const (
	// CmcState 查看服务器状态
	CmcState = "state"
	// CmcList 查看列表
	CmcList = "list"
	// CmdInfo 查看信息
	CmdInfo = "info"
	// CmdKick 踢除实例
	CmdKick = "kick"
)

// 与 Rabbit-Server 通信 ---------- ---------- ---------- ---------- ---------- ----------

const (
	HttpKeyId     = "id"
	HttpKeyData   = "d"
	HttpKeyDetail = "dt"
	HttpKeyWeight = "w"
)

const (
	// CodeParamError 参数错误
	CodeParamError = 1000 + iota
	// CodeParamLack 参数缺失
	CodeParamLack
	// CodeParamInvalid 参数无效
	CodeParamInvalid
	// CodeParamDecrypt 参数解密失败
	CodeParamDecrypt
	// CodeParamBase64 参数base64解析错误
	CodeParamBase64
	// CodeParamJson 参数json解析错误
	CodeParamJson

	// CodeVerifyLinkSignFail 注册时验证签名失败
	CodeVerifyLinkSignFail
	// CodeLinkEntityFail 注册失败
	CodeLinkEntityFail
	// CodeVerifyUnlinkSignFail  取消注册时验证签名失败
	CodeVerifyUnlinkSignFail
	// CodeUnlinkEntityFail   取消注册失败
	CodeUnlinkEntityFail

	CodePrivateKeyLack
	CodeEntityQueryFail
	CodeEntityQueryFailKey
)

// 与 Rabbit-Client 通信 ---------- ---------- ---------- ---------- ---------- ----------

const (
	HttpKeyQuery = "q"
)
