package jsapi

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func generateNonceStr() string {
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		return ""
	}
	return hex.EncodeToString(nonce)
}

// 计算签名
func calculatePaySign(appId, timeStamp, nonceStr, pkg string, privateKey *rsa.PrivateKey) (string, error) {
	// 构造签名串
	signString := fmt.Sprintf("%s\n%s\n%s\n%s\n", appId, timeStamp, nonceStr, pkg)

	// 计算SHA256散列值
	hash := sha256.New()
	hash.Write([]byte(signString))
	hashed := hash.Sum(nil)

	// 进行SHA256 with RSA签名
	signature, err := rsa.SignPKCS1v15(nil, privateKey, crypto.SHA256, hashed)
	if err != nil {
		return "", err
	}

	// 对签名结果进行Base64编码
	return base64.StdEncoding.EncodeToString(signature), nil
}
