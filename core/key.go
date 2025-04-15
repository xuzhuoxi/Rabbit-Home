// Package core
// Create on 2025/4/14
// @author xuzhuoxi
package core

const (
	// KeyTypePem 通用密钥
	// 公钥使用x.509标准，私钥使用PKCS#8标准
	KeyTypePem = "pem"
	// KeyTypeRSA RSA专用密钥
	// 公钥使用PKCS#1标准，私钥使用PKCS#1标准
	KeyTypeRSA = "rsa"
	// KeyTypeSSH SSH专用密钥
	// 公钥使用OpenSSH标准，私钥使用PKCS#1标准
	KeyTypeSSH = "ssh"
)
