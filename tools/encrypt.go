package tools

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
)

func EncryptBytes(data []byte, pass string) ([]byte, error) {
	passHash := sha256.Sum256([]byte(pass))

	block, err := aes.NewCipher(passHash[:])
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	cipher.NewCFBEncrypter(block, iv).XORKeyStream(ciphertext[aes.BlockSize:], data)

	return ciphertext, nil
}

func DecryptBytes(encrypted []byte, pass string) ([]byte, error) {
	passHash := sha256.Sum256([]byte(pass))

	block, err := aes.NewCipher(passHash[:])
	if err != nil {
		return nil, err
	}

	if len(encrypted) < aes.BlockSize {
		return nil, errors.New("encrypted data too short")
	}

	iv := encrypted[:aes.BlockSize]
	encrypted = encrypted[aes.BlockSize:]

	cipher.NewCFBDecrypter(block, iv).XORKeyStream(encrypted, encrypted)

	return encrypted, nil
}
