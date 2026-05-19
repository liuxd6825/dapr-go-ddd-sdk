package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

var encryptAESKey = []byte("*#06#1@abc$%^xyz)(&-=+_MmuU8~|[]")

// EncryptAES AES-GCM 加密
func EncryptAES(text string) (string, error) {
	// 1. 创建 Cipher Block
	block, err := aes.NewCipher(encryptAESKey)
	if err != nil {
		return "", err
	}

	// 2. 创建 GCM 模式
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 3. 生成随机 Nonce (GCM 标准 Nonce 长度为 12 字节)
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 4. 加密并拼接 Nonce (Seal: 将 nonce 和密文打包在一起)
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(text), nil)

	// 5. 转为 Base64 方便存储展示
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptAES AES-GCM 解密
func DecryptAES(cryptoText string) (string, error) {
	// 1. Base64 解码
	ciphertext, err := base64.StdEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}

	// 2. 创建 Cipher Block
	block, err := aes.NewCipher(encryptAESKey)
	if err != nil {
		return "", err
	}

	// 3. 创建 GCM 模式
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 4. 分离 Nonce 和 密文
	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// 5. 解密
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
