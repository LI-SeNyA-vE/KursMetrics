package aesKey

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

func EncryptMessage(message []byte) (encryptMessage, nonce, aesKey []byte, err error) {
	// Генерируем ключ AES-256
	aesKey = make([]byte, 32)
	if _, err = rand.Read(aesKey); err != nil {
		return nil, nil, nil, fmt.Errorf("ошибка генерации AES-ключа: %w", err)
	}

	// Шифруем данные AES-256
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("ошибка создания AES-блока: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("ошибка создания AES-GCM: %w", err)
	}

	nonce = make([]byte, aesGCM.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, nil, nil, fmt.Errorf("ошибка генерации nonce: %w", err)
	}

	//Зашифрованные данные aes256
	ciphertext := aesGCM.Seal(nil, nonce, message, nil)
	return ciphertext, nonce, aesKey, nil
}

func DecryptMessage(encryptedMessage, keyAES, AesNode []byte) ([]byte, error) {
	block, err := aes.NewCipher(keyAES)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания AES-блока: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания AES-GCM: %w", err)
	}

	fmt.Printf("Длина nonce: %d (ожидалось %d)\n", len(AesNode), aesGCM.NonceSize()) // Проверяем размер nonce

	message, err := aesGCM.Open(nil, AesNode, encryptedMessage, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка расшифровки AES: %w", err)
	}

	return message, nil
}
