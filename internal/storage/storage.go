package storage

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	git "github.com/go-git/go-git/v5"
)

// Storage struct to hold the data read from the file
type Storage struct {
	PrivKey []byte
	Data    string
	Cfg     *Config
	updated bool
}

// New initializes a new Storage instance
func New(cfg *Config) (*Storage, error) {
	// Read the private key
	privKey, err := readKey(cfg.PrivKeyPath)
	if err != nil {
		return nil, fmt.Errorf("error reading private key: %v", err)
	}

	storage := &Storage{
		PrivKey: privKey,
		Cfg:     cfg,
	}
	return storage, nil
}

// Update updates data inside of a storage from the git repo.
func (s *Storage) Update() error {
	if s.updated {
		return nil
	}
	s.updated = true

	// Clone the Git repository to a temporary directory
	tempDir, err := os.MkdirTemp("", "repo")
	if err != nil {
		return fmt.Errorf("error creating temporary directory: %v", err)
	}

	// Clone the repository
	if err := cloneRepo(s.Cfg.GitRepoPath, tempDir); err != nil {
		return err
	}

	// Read the data.dat file
	dataFilePath := filepath.Join(tempDir, "data.dat")
	data, err := os.ReadFile(dataFilePath)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			s.Data = ""
			return nil
		}
		return fmt.Errorf("error reading data.dat: %v", err)
	}
	s.Data = base64.StdEncoding.EncodeToString(data)
	return nil
}

// Store stores s.Data in the git repo.
func (s *Storage) Store(message *string) error {
	// Clone the Git repository to a temporary directory
	tempDir, err := os.MkdirTemp("", "repo")
	if err != nil {
		return fmt.Errorf("error creating temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Clone the repository
	if err := cloneRepo(s.Cfg.GitRepoPath, tempDir); err != nil {
		return err
	}

	// Write the data.dat file
	dataFilePath := filepath.Join(tempDir, "data.dat")
	if err := os.WriteFile(dataFilePath, []byte(s.Data), fs.ModePerm); err != nil {
		return fmt.Errorf("error reading data.dat: %v", err)
	}

	msg := defaultCommitMessage
	if message != nil {
		msg = *message
	}
	return commitRepo(tempDir, msg)
}

// readKey reads a key from a file or downloads it if it's a URL
func readKey(path string) ([]byte, error) {
	if isURL(path) {
		return downloadFile(path)
	}
	return os.ReadFile(path)
}

// isURL checks if the given string is a valid URL
func isURL(path string) bool {
	return len(path) > 5 && path[:5] == "https"
}

// downloadFile downloads a file from the given URL
func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download file: status code %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// cloneRepo clones the specified Git repository into the given directory
func cloneRepo(repoPath, destDir string) error {
	_, err := git.PlainClone(destDir, false, &git.CloneOptions{
		URL: repoPath,
	})
	if err != nil {
		return fmt.Errorf("failed to clone repository: %v", err)
	}
	return nil
}

// commitRepo commits changes to the repository with the specified commit message
func commitRepo(repoPath, commitMsg string) error {
	// Open the existing repository
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return fmt.Errorf("failed to open repository: %v", err)
	}

	// Stage the changes
	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %v", err)
	}

	// Add changes to the staging area
	_, err = w.Add("data.dat")
	if err != nil {
		return fmt.Errorf("failed to add changes to the repository: %v", err)
	}

	// Commit the changes
	_, err = w.Commit(commitMsg, &git.CommitOptions{
		All: true,
	})
	if err != nil {
		return fmt.Errorf("failed to commit changes to the repository: %v", err)
	}

	// Push the changes
	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
	})
	if err != nil {
		return fmt.Errorf("failed to push changes to the repository: %v", err)
	}

	return nil
}
