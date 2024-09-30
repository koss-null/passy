package storage

import "github.com/pkg/errors"

func (s *Storage) AddPass(key, pass string) (err error) {
	if err = s.Update(); err != nil {
		return errors.Wrap(err, "failed to update data from the repo ")
	}

	data := &Folder{Name: "", SubFolder: []*Folder{}}
	if s.Data != "" {
		data, err = s.Decrypt()
		if err != nil {
			return errors.Wrap(err, "failed to get current passwords")
		}
	}

	if err := data.Add(key, pass); err != nil {
		return errors.Wrap(err, "failed to add key and pass to the data map")
	}

	if err = s.Encrypt(data); err != nil {
		return errors.Wrap(err, "failed to encrypt message with a new pass")
	}

	return nil
}
