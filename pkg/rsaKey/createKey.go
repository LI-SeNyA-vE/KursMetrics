package rsaKey

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
)

// Функция генерации и сохранения ключей
func GenerateAndSaveKeys(path string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	publicKey := &privateKey.PublicKey

	// Сохраняем приватный ключ в PKCS#8
	split := strings.Split(path, "/")
	privFile, err := os.Create(fmt.Sprint(split[0], "/privateKey.pem"))
	if err != nil {
		return err
	}
	defer privFile.Close()

	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return err
	}

	privBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privBytes,
	}
	pem.Encode(privFile, privBlock)

	fmt.Println("Приватный ключ сохранён (PKCS#8).")

	// Сохраняем публичный ключ в PKIX (X.509)
	pubFile, err := os.Create(fmt.Sprint(split[0], "/publicKey.pem"))
	if err != nil {
		return err
	}
	defer pubFile.Close()

	pubBytes, err := x509.MarshalPKIXPublicKey(publicKey) // Используем PKIX
	if err != nil {
		return err
	}

	pubBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	}
	pem.Encode(pubFile, pubBlock)

	fmt.Println("Публичный ключ сохранён (PKIX).")

	return nil
}
