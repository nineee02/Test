package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AESUtil interface {
	AES256Encrypt(plaintext []byte, key []byte) (string, error)
	AES256Decrypt(ciphertext string, key []byte) (string, error)
	Difference(x, y []string) (diff []string)
}

type AESUtilImpl struct{}

func (a *AESUtilImpl) AES256Encrypt(plaintext []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func (a *AESUtilImpl) AES256Decrypt(ciphertextBase64 string, key []byte) (string, error) {
	if len(ciphertextBase64) < 24 {
		log.Printf("Warning: Password may not be encrypted")
		return ciphertextBase64, nil
	}
	ciphertext, err := base64.URLEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		log.Println("Warning: Password is not encoded in Base64, returning as plaintext")
		return ciphertextBase64, nil
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		log.Println("Warning: Ciphertext too short, returning as plaintext")
		return ciphertextBase64, nil
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

func (a *AESUtilImpl) Difference(x, y []string) (diff []string) {
	m := make(map[string]bool)

	for _, item := range y {
		m[item] = true
	}

	for _, item := range x {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}

func GenerateJWT(userID string, secret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
