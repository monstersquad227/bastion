package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"os"
)

func DecryptPassword(encrypt string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encrypt)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(os.Getenv("AesSecret")))
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 解析 nonce
	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("invalid ciphertext")
	}
	nonce, cipherText := data[:nonceSize], data[nonceSize:]

	// 解密数据
	plainText, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
