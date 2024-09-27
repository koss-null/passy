package storage

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// Storage struct to hold the data read from the file
type Storage struct {
	PrivKey []byte
	PubKey  []byte
	Data    []byte
	Cfg     *Config
}

// NewStorage initializes a new Storage instance
func NewStorage(cfg *Config) (*Storage, error) {
	// Read the private key
	privKey, err := readKey(cfg.PrivKeyPath)
	if err != nil {
		return nil, fmt.Errorf("error reading private key: %v", err)
	}

	// Read the public key
	pubKey, err := readKey(cfg.PubKeyPath)
	if err != nil {
		return nil, fmt.Errorf("error reading public key: %v", err)
	}

	// Clone the Git repository to a temporary directory
	tempDir, err := os.MkdirTemp("", "repo")
	if err != nil {
		return nil, fmt.Errorf("error creating temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up the temporary directory

	// Clone the repository
	if err := cloneRepo(cfg.GitRepoPath, tempDir); err != nil {
		return nil, err
	}

	// Read the data.dat file
	dataFilePath := filepath.Join(tempDir, "data.dat")
	data, err := os.ReadFile(dataFilePath)
	if err != nil {
		return nil, fmt.Errorf("error reading data.dat: %v", err)
	}

	// Create and return the Storage instance
	return &Storage{
		PrivKey: privKey,
		PubKey:  pubKey,
		Data:    data,
		Cfg:     cfg,
	}, nil
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

	return os.ReadAll(resp.Body)
}

// cloneRepo clones the specified Git repository into the given directory
func cloneRepo(repoPath, destDir string) error {
	cmd := exec.Command("git", "clone", repoPath, destDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %v", err)
	}
	return nil
}
