package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

// Decrypt get's data from the repository, decrypts in and unmarshal.
func (s *Storage) Decrypt() (*Folder, error) {
	err := s.Update()
	if err != nil {
		return nil, err
	}

	decrypted, err := s.decrypt(s.Data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode pass data")
	}

	if len(decrypted) < 2 {
		return nil, errors.Wrap(err, "pass data file is too short")
	}
	startGarbageLen, endGarbageLen := binary.BigEndian.Uint32(decrypted[:4]), binary.BigEndian.Uint32(decrypted[4:8])
	if len(decrypted) < int(startGarbageLen)+int(endGarbageLen)+8 {
		return nil, errors.Wrap(err, "pass data file encoded incorrectly")
	}

	decrypted = decrypted[8+startGarbageLen : len(decrypted)-int(endGarbageLen)]

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
