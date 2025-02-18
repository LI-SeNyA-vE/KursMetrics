// Package rsakey содержит функции для шифровки/расшифровки запроса.
// Шифрует/расшифровывает только необходимые данные для AES шифрования
package rsakey

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func CheckKey(pathKeyPath string) error {
	keyData, err := loadKey(pathKeyPath)
	if err != nil {
		return fmt.Errorf("не удалось прочитать файл: %v", err)
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return fmt.Errorf("не удалось декодировать PEM: %s", err)
	}

	_, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return err
	}
	return nil
}
