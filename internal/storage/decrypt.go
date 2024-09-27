package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

// returns map key-value
func (s *Storage) Decrypt() (*Folder, error) {
	err := s.Update()
	if err != nil {
		return nil, err
	}

	decrypted, err := s.decrypt(s.Data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode pass data")
	}

	var head Folder
	if err := json.Unmarshal(decrypted, &head); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal decoded pass list")
	}
	return &head, nil
}

func (s *Storage) decrypt(dataStr string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(dataStr)
	if err != nil {
		return nil, err
	}

	// Create a new AES cipher
	block, err := aes.NewCipher(s.PrivKey)
	if err != nil {
		return nil, err
	}

	// GCM mode requires a nonce (number used once)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Split the nonce and the ciphertext
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce, cipherText := data[:nonceSize], data[nonceSize:]

	// Decrypt the ciphertext
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}
