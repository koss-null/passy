package storage

import (
	"crypto/rand"
)

// GenerateAESKey generates a random AES key of the specified size (16, 24, or 32 bytes)
func GenerateAESKey(size int) ([]byte, error) {
	key := make([]byte, size)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}
