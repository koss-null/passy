package storage

import (
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
)

const defaultCommitMessage = "nothing important"

// Encrypt encodes the folder paths and stores the data in a repository
func (s *Storage) Encrypt(key, pass, encryptionPass string, commitMessage *string) error {
	data, err := s.Decrypt()
	if err != nil {
		return errors.Wrap(err, "failed to get current passwords")
	}

	path := strings.Split(key, ".")
	for _, folder := range path {
		for _, sf := range data.SubFolder {
			if sf.Name == folder {
				data = sf
				break
			}
		}
		newFolder := &Folder{Name: folder, Key2Pass: make(map[string]string)}
		data.SubFolder = append(data.SubFolder, newFolder)
		data = newFolder
	}

	data.Key2Pass[path[len(path)-1]] = pass
	byteData, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal new password data")
	}
	// TODO: encode data
	s.Data = byteData

	message := defaultCommitMessage
	if commitMessage != nil {
		message = *commitMessage
	}
	return s.Store(message)
}
