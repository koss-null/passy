package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"io"
	mrand "math/rand"

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

	// each 4 bytes
	randStartLen, randEndLen := mrand.Intn(262144), mrand.Intn(262144)
	randomStartBytes, randomEndBytes := make([]byte, randStartLen), make([]byte, randEndLen)
	if _, err := rand.Read(randomStartBytes); err != nil {
		return errors.Wrap(err, "failed to read from random stream")
	}
	if _, err := rand.Read(randomEndBytes); err != nil {
		return errors.Wrap(err, "failed to read from random stream")
	}

	var startLenBytes, endLenBytes [4]byte
	binary.BigEndian.PutUint32(startLenBytes[:], uint32(randStartLen))
	binary.BigEndian.PutUint32(endLenBytes[:], uint32(randEndLen))
	byteData = append(
		make([]byte, 0, 8+randStartLen+len(byteData)+randEndLen),
		append(
			// first 4 bytes represent random data lengths
			append(startLenBytes[:], endLenBytes[:]...),
			// adding random data
			append(randomStartBytes, append(byteData, randomEndBytes...)...)...,
		)...,
	)

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
