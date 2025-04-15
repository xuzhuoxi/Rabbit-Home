// Package home
// Create on 2025/4/12
// @author xuzhuoxi
package home

import (
	"github.com/xuzhuoxi/Rabbit-Home/core"
	"github.com/xuzhuoxi/Rabbit-Home/core/utils"
	"github.com/xuzhuoxi/infra-go/cryptox/asymmetric"
	"github.com/xuzhuoxi/infra-go/filex"
	"github.com/xuzhuoxi/infra-go/logx"
	"io/fs"
	"sync"
)

func NewRabbitHomeKeys() *RabbitHomeKeys {
	return &RabbitHomeKeys{
		logger:                 GlobalLogger,
		defaultId2PublicCipher: make(map[string]asymmetric.IRSAPublicCipher),
		defaultPublicCiphers:   make([]asymmetric.IRSAPublicCipher, 0),
	}
}

// RabbitHomeKeys 密钥管理器
type RabbitHomeKeys struct {
	logger logx.ILogger

	selfPrivateCipher      asymmetric.IRSAPrivateCipher           // 自身私钥的处理器
	defaultPublicCiphers   []asymmetric.IRSAPublicCipher          // 预存公钥的处理器
	defaultId2PublicCipher map[string]asymmetric.IRSAPublicCipher // 预存公钥处理器的映射

	verifyLock sync.RWMutex // 验签处理并发锁
}

// PrivateCipher 专属于当前Rabbit-Home程序的RSA私钥处理器
func (o *RabbitHomeKeys) PrivateCipher() asymmetric.IRSAPrivateCipher {
	return o.selfPrivateCipher
}

// LoadInternalDefaultPublicKeys
// 加载预设的密钥文件
func (o *RabbitHomeKeys) LoadInternalDefaultPublicKeys() {
	o.defaultId2PublicCipher = make(map[string]asymmetric.IRSAPublicCipher)
	kv := GlobalHomeConfig.InternalVerifier.KeyVerifier
	fixPath := filex.FixDirPath(kv.KeysPath)
	_ = filex.WalkAllFiles(fixPath, func(path string, info fs.FileInfo, err error) error {
		if nil != err {
			return nil
		}
		publicCipher, err1 := utils.LoadPublicCipher(path, kv.KeyType)
		if nil != err1 {
			o.logger.Warnln("[RabbitHomeKeys.LoadInternalDefaultPublicKeys]", err1, path)
			return nil
		}
		o.defaultPublicCiphers = append(o.defaultPublicCiphers, publicCipher)
		o.logger.Infoln("[RabbitHomeKeys.LoadInternalDefaultPublicKeys] suc at", path)
		return nil
	})
}

// LoadSelfPrivateKey
// 加载当前程序专用的RSA私钥
func (o *RabbitHomeKeys) LoadSelfPrivateKey() {
	kv := GlobalHomeConfig.ExternalVerifier.KeyVerifier
	fixPath := filex.FixFilePath(kv.KeyPath)
	privateCipher, err1 := utils.LoadPrivateCipher(fixPath, kv.KeyType)
	if nil != err1 {
		o.logger.Warnln("[RabbitHomeKeys.LoadInternalDefaultPublicKeys]", err1, kv.KeyPath)
		return
	}
	o.selfPrivateCipher = privateCipher
	o.logger.Infoln("[RabbitHomeKeys.LoadSelfPrivateKey] suc at", fixPath)
}

// VerifyLinkSign 对连接信息进行验证签名
func (o *RabbitHomeKeys) VerifyLinkSign(linkInfo *core.LinkInfo) (rsaPublicCipher asymmetric.IRSAPublicCipher, SK []byte, pass bool) {
	o.verifyLock.Lock()
	defer o.verifyLock.Unlock()
	original := linkInfo.OriginalSignData()
	if rsaPublicCipher, pass = o.verifyFromDefault(original, linkInfo.Id, linkInfo.Signature); pass {
		o.defaultId2PublicCipher[linkInfo.Id] = rsaPublicCipher
		goto Pass
	}
	if GlobalHomeConfig.InternalVerifier.KeyVerifier.HotKeysEnable {
		if rsaPublicCipher, pass = o.verifyFromHot(original, linkInfo.Signature); pass {
			goto Pass
		}
	}
	return nil, nil, false
Pass:
	SK = utils.DeriveTempKey(GlobalHomeConfig.InternalVerifier.KeyVerifier.Share, linkInfo)
	return rsaPublicCipher, SK, pass
}

// VerifyUnlinkSign 对断开连接信息进行验证签名
func (o *RabbitHomeKeys) VerifyUnlinkSign(unlinkInfo *core.UnlinkInfo) (rsaPublicCipher asymmetric.IRSAPublicCipher, pass bool) {
	o.verifyLock.Lock()
	defer o.verifyLock.Unlock()

	original := unlinkInfo.OriginalSignData()
	if rsaPublicCipher, pass = o.verifyFromDefault(original, unlinkInfo.Id, unlinkInfo.Signature); pass {
		delete(o.defaultId2PublicCipher, unlinkInfo.Id)
		goto Pass
	}
	if GlobalHomeConfig.InternalVerifier.KeyVerifier.HotKeysEnable {
		if rsaPublicCipher, pass = o.verifyFromHot(original, unlinkInfo.Id); pass {
			goto Pass
		}
	}
	return nil, false
Pass:
	return rsaPublicCipher, pass
}

func (o *RabbitHomeKeys) verifyFromDefault(original []byte, id string, signature string) (rsaPublicCipher asymmetric.IRSAPublicCipher, pass bool) {
	if defaultCipher, ok := o.defaultId2PublicCipher[id]; ok {
		pass, _ = defaultCipher.VerifySignBase64(original, signature, core.Base64Encoding)
		if pass {
			return defaultCipher, pass
		}
	}
	for _, cipher := range o.defaultPublicCiphers {
		pass, _ = cipher.VerifySignBase64(original, signature, core.Base64Encoding)
		if pass {
			return cipher, pass
		}
	}
	return nil, false
}
func (o *RabbitHomeKeys) verifyFromHot(original []byte, signature string) (rsaPublicCipher asymmetric.IRSAPublicCipher, pass bool) {
	kv := GlobalHomeConfig.InternalVerifier.KeyVerifier
	if !kv.HotKeysEnable {
		return nil, false
	}
	hotPath := filex.FixDirPath(kv.HotKeysPath)
	var paths []string
	_ = filex.WalkAllFiles(hotPath, func(path string, info fs.FileInfo, err error) error {
		paths = append(paths, path)
		return nil
	})
	for _, path := range paths {
		publicCipher, err1 := utils.LoadPublicCipher(path, kv.KeyType)
		if nil != err1 {
			continue
		}
		pass, _ = publicCipher.VerifySignBase64(original, signature, core.Base64Encoding)
		if pass {
			return publicCipher, true
		}
	}
	return nil, false
}
