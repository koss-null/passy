package storage

import (
	"errors"
	"strings"
)

type Folder struct {
	Name      string
	SubFolder []*Folder
	Key2Pass  map[string]string
}

func (f *Folder) String() string {
	const line = "-------------------------------\n"

	sb := strings.Builder{}

	sb.WriteString(line)
	sb.WriteString(f.Name + "\n")
	sb.WriteString(line)
	for k, v := range f.Key2Pass {
		sb.WriteString(k + ":\n" + v + "\n")
	}

	if f.SubFolder != nil {
		for _, sf := range f.SubFolder {
			sb.WriteString(sf.String())
		}
	}

	return sb.String()
}

func (f *Folder) Add(folderPath, key, pass string) error {
	path := strings.Split(folderPath, "/")
	if len(path) > 0 {
		if path[0] == f.Name {
			if len(path) == 1 {
				f.Key2Pass[key] = pass
				return nil
			}
			if f.SubFolder == nil {
				f.SubFolder = []*Folder{{Name: path[1]}}
			}
			return f.SubFolder[0].Add(strings.Join(path[1:], ""), key, pass)
		}
		return errors.New("")
	}
	return errors.New("")
}
