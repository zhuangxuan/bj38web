package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// 加密密钥
const secret = "ouzhsh"

// EncryptPassword 密码MD5加密
func EncryptPassword(oPWD string) string {
	hash := md5.New()
	// 对应原密码进行hash
	hash.Write([]byte(oPWD))
	// 再用系统自定义的秘钥再加密
	return hex.EncodeToString(hash.Sum([]byte(secret)))
}
