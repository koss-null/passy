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

func (f *Folder) String(prefix string) func() string {
	const tab = "    "
	const symbol = "└── "
	const verticalLine = "│   "

	sb := &strings.Builder{}

	if f.Name != "" {
		sb.WriteString(prefix + symbol + f.Name + "\n")
	}

	if f.Pass != "" {
		sb.WriteString(prefix + tab + "Password: " + f.Pass + "\n")
	}

	if f.SubFolder != nil {
		for _, sf := range f.SubFolder {
			newPrefix := prefix
			if f.Name != "" {
				newPrefix += verticalLine
			}
			sb.WriteString(sf.String(newPrefix)())
		}
	}

	return sb.String
}

func (f *Folder) SecureString(prefix string) func() string {
	const tab = "    "
	const symbol = "└── "
	const verticalLine = "│   "
	const lockSymbol = "🔒"

	sb := &strings.Builder{}

	if f.Name != "" {
		sb.WriteString(prefix + symbol + f.Name + "\n")
	}

	if f.Pass != "" {
		sb.WriteString(prefix + tab + "Password: " + lockSymbol + "\n")
	}

	if f.SubFolder != nil {
		for _, sf := range f.SubFolder {
			newPrefix := prefix
			if f.Name != "" {
				newPrefix += verticalLine
			}
			sb.WriteString(sf.SecureString(newPrefix)())
		}
	}

	return sb.String
}

func (f *Folder) Add(folderPath, pass string) error {
	if f.Name != "" {
		return errors.New("should add to the head only")
	}

	path := strings.Split(folderPath, ">")
	cf := f

	for j, folderName := range path {
		found := false
		for i := range cf.SubFolder {
			if cf.SubFolder[i].Name == folderName {
				cf = cf.SubFolder[i]
				found = true
				break
			}
		}
		if !found {
			// Create new folders for the remaining path
			for _, newFolderName := range path[j:] {
				newFolder := &Folder{
					Name:      newFolderName,
					SubFolder: make([]*Folder, 0),
					Pass:      "",
				}
				cf.SubFolder = append(cf.SubFolder, newFolder)
				cf = newFolder
			}
			// Exit the loop after adding new folders
			break
		}
	}

	cf.Pass = pass
	return nil
}
