package storage

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const defaultConfigFile = "~/.config/passy/config.toml"

const mockConfig = `
PrivKeyPath = "/path/to/private/key.pem" # can be https link
GitRepoPath = "https://github.com/user/repository.git"
`

type Config struct {
	PrivKeyPath string
	GitRepoPath string
}

// ParseConfig reads the config file, fills config fields, and validates them.
func ParseConfig() (*Config, error) {
	var config Config

	// Expand the default config file path
	configFilePath, err := expandPath(defaultConfigFile)
	if err != nil {
		return nil, fmt.Errorf("error expanding config file path: %v", err)
	}

	// create mock config if not exist
	if err = checkConfigExistOrCreateNew(configFilePath); err != nil {
		return nil, err
	}

	// Read the config file
	if _, err := toml.DecodeFile(configFilePath, &config); err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	// Validate the config fields
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func checkConfigExistOrCreateNew(configFilePath string) (err error) {
	var configFile *os.File
	if configFile, err = os.OpenFile(configFilePath, os.O_RDONLY, fs.ModeType); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("error opening config file: %v", err)
		}

		if err := os.MkdirAll(filepath.Dir(configFilePath), os.ModePerm); err != nil {
			return fmt.Errorf("error creating config directory: %v", err)
		}
		var configFile *os.File
		if configFile, err = os.Create(configFilePath); err != nil {
			return fmt.Errorf("error creating new config: %v", err)
		}

		if _, err = configFile.WriteString(mockConfig); err != nil {
			configFile.Close()
			return fmt.Errorf("error writing mock config: %v", err)
		}
	}
	configFile.Close()
	return nil
}

// expandPath expands the user's home directory in the given path.
func expandPath(path string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, path[2:]), nil
}

// validateConfig checks if the paths are valid and if the Git repository is valid.
func validateConfig(config *Config) error {
	// Validate PrivKeyPath
	if err := validateFileExists(config.PrivKeyPath); err != nil {
		return fmt.Errorf("invalid private key path: %v", err)
	}

	// Validate GitRepoPath
	if err := validateGitRepo(config.GitRepoPath); err != nil {
		return fmt.Errorf("invalid Git repository path: %v", err)
	}

	return nil
}

// validateFileExists checks if a file exists at the given path.
func validateFileExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", path)
	}
	return nil
}

// validateGitRepo checks if the given path is a valid remote Git repository.
func validateGitRepo(repoPath string) error {
	cmd := exec.Command("git", "ls-remote", repoPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("not a valid remote Git repository: %s", repoPath)
	}
	return nil
}
