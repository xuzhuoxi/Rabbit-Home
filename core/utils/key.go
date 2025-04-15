// Package utils
// Create on 2025/4/13
// @author xuzhuoxi
package utils

import (
	"github.com/xuzhuoxi/Rabbit-Home/core"
	"github.com/xuzhuoxi/infra-go/cryptox/asymmetric"
	"github.com/xuzhuoxi/infra-go/cryptox/key"
	"github.com/xuzhuoxi/infra-go/timex"
	"strconv"
)

// LoadPublicCipher
// 加载公钥密钥
func LoadPublicCipher(path string, keyType string) (asymmetric.IRSAPublicCipher, error) {
	var publicCipher asymmetric.IRSAPublicCipher
	var err error
	if keyType == core.KeyTypePem {
		publicCipher, err = asymmetric.LoadPublicCipherPEM(path)
	} else if keyType == core.KeyTypeRSA {
		publicCipher, err = asymmetric.LoadPublicCipherRSA(path)
	} else if keyType == core.KeyTypeSSH {
		publicCipher, err = asymmetric.LoadPublicCipherSSH(path)
	}
	if nil != err {
		return nil, err
	}
	return publicCipher, nil
}

// LoadPrivateCipher
// 加载私钥密钥
func LoadPrivateCipher(path string, keyType string) (asymmetric.IRSAPrivateCipher, error) {
	var privateCipher asymmetric.IRSAPrivateCipher
	var err error
	if keyType == core.KeyTypePem {
		privateCipher, err = asymmetric.LoadPrivateCipherPEM(path)
	} else if keyType == core.KeyTypeRSA || keyType == core.KeyTypeSSH {
		privateCipher, err = asymmetric.LoadPrivateCipherRSA(path)
	}
	if nil != err {
		return nil, err
	}
	return privateCipher, nil
}

// DeriveTempKey
// 派生临时密钥
// 算法：
//  1. passphrase: shareKey + linkInfo.Id + 当前时间戳
//  2. 派生密钥: key.DeriveKeyPbkdf2(passphrase, salt, iterations, keyLen)
//		salt       = []byte("infra-go:cryptox.key")
//		iterations = 100000
//		keyLen     = 32
func DeriveTempKey(shareKey string, linkInfo *core.LinkInfo) []byte {
	now := int64(timex.NowDuration1970())
	buf := make([]byte, 0, 20) // 预留最多 20 字节（int64 最大长度）
	buf = strconv.AppendInt(buf, now, 10)
	passphrase := append(append([]byte(shareKey), []byte(linkInfo.Id)...), buf...)
	return key.DeriveKeyPbkdf2Default(passphrase)
}
