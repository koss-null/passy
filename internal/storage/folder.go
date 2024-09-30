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
	const folderColor = "\033[1;34m"   // Blue color for folder names
	const passwordColor = "\033[1;33m" // Yellow color for passwords
	const resetColor = "\033[0m"       // Reset color

	sb := &strings.Builder{}

	if f.Name != "" {
		sb.WriteString(prefix + symbol + folderColor + f.Name + resetColor + "\n")
	}

	if f.Pass != "" {
		sb.WriteString(prefix + tab + passwordColor + "Password: " + resetColor + f.Pass + "\n")
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
	const folderColor = "\033[1;34m"   // Blue color for folder names
	const passwordColor = "\033[1;31m" // Red color for passwords
	const resetColor = "\033[0m"       // Reset color

	sb := &strings.Builder{}

	if f.Name != "" {
		sb.WriteString(prefix + symbol + folderColor + f.Name + resetColor + "\n")
	}

	if f.Pass != "" {
		sb.WriteString(prefix + tab + passwordColor + "Password: " + lockSymbol + resetColor + "\n")
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

const folderSeparator = "/"

func (f *Folder) Add(folderPath, pass string) error {
	if f.Name != "" {
		return errors.New("should add to the head only")
	}

	path := strings.Split(folderPath, folderSeparator)
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

func (f *Folder) Delete(folderPath string) error {
	if f.Name == "" {
		return errors.New("cannot delete from an empty folder")
	}

	path := strings.Split(folderPath, folderSeparator)
	cf := f
	var parent *Folder
	var folderToDelete *Folder

	for _, folderName := range path {
		found := false
		for i := range cf.SubFolder {
			if cf.SubFolder[i].Name == folderName {
				parent = cf
				cf = cf.SubFolder[i]
				found = true
				break
			}
		}
		if !found {
			return errors.New("folder not found")
		}
	}

	// Now cf is the folder to delete
	folderToDelete = cf

	// Remove folderToDelete from parent's SubFolder
	for i, subFolder := range parent.SubFolder {
		if subFolder == folderToDelete {
			parent.SubFolder = append(parent.SubFolder[:i], parent.SubFolder[i+1:]...)
			return nil
		}
	}

	return errors.New("folder not found in parent's subfolders")
}

func (f *Folder) GetSubFolder(key string) (*Folder, bool) {
	path := strings.Split(key, folderSeparator)
	cf := f

	for _, folderName := range path {
		found := false
		for i := range cf.SubFolder {
			if cf.SubFolder[i].Name == folderName {
				cf = cf.SubFolder[i]
				found = true
				break
			}
		}
		if !found {
			return nil, false
		}
	}

	return cf, true
}
