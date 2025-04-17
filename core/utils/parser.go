// Package utils
// Create on 2025/4/14
// @author xuzhuoxi
package utils

import (
	"fmt"
	"github.com/json-iterator/go"
	"github.com/xuzhuoxi/Rabbit-Home/core"
	"github.com/xuzhuoxi/infra-go/cryptox"
)

// ParseHomeResponseInfo
// 解析字节数据为 HomeResponseInfo
func ParseHomeResponseInfo(bs []byte) (info *core.HomeResponseInfo, err error) {
	info = &core.HomeResponseInfo{}
	err = ParseDate(bs, info, nil)
	if nil != err {
		return nil, err
	}
	return info, nil
}

// SerializeHomeResponseInfo
// 序列化 HomeResponseInfo 为字节数据
func SerializeHomeResponseInfo(info core.HomeResponseInfo) []byte {
	rs, _ := SerializeData(info, nil)
	return rs
}

// ParseDate
// 解析字节数据为数据
// bs: 字节数据
// data: 结果数据结构体，要求是引用
// decryptCipher: 解密处理器, 如果存在则在转为base后进行解密
func ParseDate(bs []byte, data interface{}, decryptCipher cryptox.IDecryptCipher) error {
	json, err := core.Base64Encoding.DecodeString(string(bs))
	if nil != err {
		return err
	}
	if nil != decryptCipher {
		json, err = decryptCipher.Decrypt(json)
		if nil != err {
			fmt.Println("2")
			return err
		}
	}
	err = jsoniter.Unmarshal(json, data)
	if nil != err {
		return err
	}
	return nil
}

// SerializeData
// 序列化数据为字节数据
// data: 结构体数据
// encryptCipher: 加密处理器, 如果存在则在转为base前进行加密
func SerializeData(data interface{}, encryptCipher cryptox.IEncryptCipher) ([]byte, error) {
	json, err := jsoniter.Marshal(data)
	if nil != err {
		return nil, err
	}
	if nil != encryptCipher {
		json, err = encryptCipher.Encrypt(json)
		if nil != err {
			return nil, err
		}
	}
	base64Data := core.Base64Encoding.EncodeToString(json)
	return []byte(base64Data), nil
}
