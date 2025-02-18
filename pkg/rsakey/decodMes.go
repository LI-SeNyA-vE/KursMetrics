// Package rsakey содержит функции для шифровки/расшифровки запроса.
// Шифрует/расшифровывает только необходимые данные для AES шифрования
package rsakey

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

// DecryptMessage Функция дешифрования сообщения
func DecryptMessage(pathPrivateKeyPath string, ciphertext []byte) ([]byte, error) {
	keyData, err := loadKey(pathPrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки ключа: %w", err)
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("не удалось декодировать PEM")
	}

	// Парсим PKCS#8 приватный ключ
	privInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга приватного ключа (PKCS#8): %w", err)
	}

	// Приводим к *rsa.PrivateKey
	priv, ok := privInterface.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("загруженный ключ не является RSA")
	}

	hash := sha256.New()
	decryptedBytes, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка расшифровки AES-ключа через RSA: %w", err)
	}

	fmt.Printf("RSA Дешифрованный ключ AES: %x\n", decryptedBytes) // Лог для проверки

	return decryptedBytes, nil
}
