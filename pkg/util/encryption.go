package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
)

const saltLength = 8

func Decrypt(payload []byte, secret string) ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	salt := payload[:saltLength]
	key := encryptionKeyToBytes(secret, string(salt))
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(payload) < aes.BlockSize {
		return nil, errors.New("payload too short")
	}
	iv := payload[saltLength : saltLength+aes.BlockSize]
	payload = payload[saltLength+aes.BlockSize:]
	payloadDst := make([]byte, len(payload))
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(payloadDst, payload)
	return payloadDst, nil
}
func Encrypt(payload []byte, secret string) ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	salt := GetRandomString(saltLength)
	key := encryptionKeyToBytes(secret, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, saltLength+aes.BlockSize+len(payload))
	copy(ciphertext[:saltLength], []byte(salt))
	iv := ciphertext[saltLength : saltLength+aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[saltLength+aes.BlockSize:], payload)
	return ciphertext, nil
}
func encryptionKeyToBytes(secret, salt string) []byte {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return PBKDF2([]byte(secret), []byte(salt), 10000, 32, sha256.New)
}
