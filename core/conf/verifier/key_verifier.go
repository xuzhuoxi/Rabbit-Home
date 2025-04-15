// Package verifier
// Create on 2025/3/1
// @author xuzhuoxi
package verifier

type InternalKeyVerifier struct {
	Enable        bool   `yaml:"enable"`          // 是否开启密钥校验
	Share         string `yaml:"share"`           // 共享密钥
	KeyType       string `yaml:"key-type"`        // 密钥类型：pem,rsa,ssh
	KeysPath      string `yaml:"keys-path"`       // 公钥文件目录
	HotKeysEnable bool   `yaml:"hot-keys-enable"` // 是否开启密钥的热部署
	HotKeysPath   string `yaml:"hot-keys-path"`   // 热部署的密钥文件路径
}

// PreProcess 对原始数据进行预处理
func (o *InternalKeyVerifier) PreProcess() {
	if !o.Enable {
		return
	}
}

type ExternalKeyVerifier struct {
	Enable  bool   `yaml:"enable"`   // 是否开启密钥校验
	KeyType string `yaml:"key-type"` // 密钥类型：pem,rsa,ssh
	KeyPath string `yaml:"key-path"` // 私钥文件路径
}

// PreProcess 对原始数据进行预处理
func (o *ExternalKeyVerifier) PreProcess() {
	if !o.Enable {
		return
	}
}
