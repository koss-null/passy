package storage

import (
	"errors"
	"strings"
)

type Folder struct {
	Name      string
	SubFolder []*Folder
	Pass      string
}

func (f *Folder) String() string {
	const line = "-------------------------------\n"

	sb := strings.Builder{}

	if f.Pass != "" {
		sb.WriteString(line)
		sb.WriteString(f.Name + "\n")
		sb.WriteString(line)
		sb.WriteString(f.Pass + "\n")
	}

	if f.SubFolder != nil {
		for _, sf := range f.SubFolder {
			sb.WriteString(sf.String())
		}
	}

	return sb.String()
}

func (f *Folder) Add(folderPath, pass string) error {
	if f.Name != "" {
		return errors.New("should add to the head only")
	}

	path := strings.Split(folderPath, "->")
	cf := f

	for j, folderName := range path {
		for i := range cf.SubFolder {
			if cf.SubFolder[i].Name == folderName {
				cf = cf.SubFolder[i]
				break
			}
		}
		for _, newFolderName := range path[j:] {
			newFolder := &Folder{
				Name:      newFolderName,
				SubFolder: make([]*Folder, 0),
				Pass:      "",
			}
			cf.SubFolder = append(cf.SubFolder, newFolder)
			cf = newFolder
		}
		break
	}

	cf.Pass = pass
	return nil
}
