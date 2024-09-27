package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"io"

	"github.com/pkg/errors"
)

const defaultCommitMessage = "nothing important"

// Encrypt encodes the folder paths and stores the data in a repository
func (s *Storage) Encrypt(key, pass, encryptionPass string, commitMessage *string) error {
	err := s.Update()
	if err != nil {
		return errors.Wrap(err, "failed to update data from the repo during encryption")
	}

	data := &Folder{Name: "", SubFolder: []*Folder{}}
	if s.Data != "" {
		var err error
		data, err = s.Decrypt()
		if err != nil {
			return errors.Wrap(err, "failed to get current passwords")
		}
	}

	if err := data.Add(key, pass); err != nil {
		return errors.Wrap(err, "failed to add key and pass to the data map")
	}
	byteData, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal new password data")
	}

	encryptedData, err := s.encrypt(byteData)
	if err != nil {
		return errors.Wrap(err, "failed to encrypt new password data")
	}
	s.Data = string(encryptedData)

	message := defaultCommitMessage
	if commitMessage != nil {
		message = *commitMessage
	}
	return s.Store(message)
}

func (s *Storage) encrypt(data []byte) ([]byte, error) {
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

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Encrypt the plaintext
	cipherText := gcm.Seal(nonce, nonce, data, nil)
	return []byte(cipherText), nil
}
