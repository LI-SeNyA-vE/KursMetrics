// Package rsakey содержит функции для шифровки/расшифровки запроса.
// Шифрует/расшифровывает только необходимые данные для AES шифрования
package rsakey

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/LI-SeNyA-vE/KursMetrics/pkg/aeskey"
	"os"
)

// EncryptMessage Функция шифрования сообщения
func EncryptMessage(publicKeyPath string, message []byte) ([]byte, error) {
	keyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("не удалось декодировать PEM публичного ключа")
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("не удалось декодировать PEM публичного ключа")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга публичного ключа (PKIX): %w", err)
	}

	// Приводим к типу *rsa.PublicKey
	pub, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("ошибка: загруженный ключ не является RSA")
	}

	// Шифруем данные AES-256
	ciphertext, nonce, keyAES, err := aeskey.EncryptMessage(message)
	if err != nil {
		return nil, fmt.Errorf("ошибка шифрования AES: %w", err)
	}

	// Шифруем AES-256 ключ через публичный ключ RSA
	hash := sha256.New()
	encryptedAESKey, err := rsa.EncryptOAEP(hash, rand.Reader, pub, keyAES, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка шифрования AES-ключа через RSA: %w", err)
	}

	// Создаём JSON для передачи
	result, err := json.Marshal(struct {
		AesKey           string `json:"AES-KEY_Encode-RSA"`
		AesNode          string `json:"Rand_valu_AES-GCM"`
		EncryptedMessage string `json:"Encrypted_Message"`
	}{
		AesKey:           base64.StdEncoding.EncodeToString(encryptedAESKey),
		AesNode:          base64.StdEncoding.EncodeToString(nonce),
		EncryptedMessage: base64.StdEncoding.EncodeToString(ciphertext),
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка при маршлинге зашифрованного запроса: %v", err)
	}
	return result, nil
}
