package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"

	git "gopkg.in/src-d/go-git.v4"
)

// Encrypt encodes the folder paths and stores the data in a repository
func (s *Storage) Encrypt(folderPath, key, pass, encryptionPass string, commitMessage *string) (*Folder, error) {
	shaSum := sha256.Sum256([]byte(encryptionPass))
	aesBlock, err := aes.NewCipher(shaSum[:])
	if err != nil {
		return nil, err
	}

	gcmInstance, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcmInstance.NonceSize())
	_, _ = io.ReadFull(rand.Reader, nonce)

	cipheredText := gcmInstance.Seal(nonce, nonce, value, nil)
}

// gitCommitAndPush commits and pushes changes to the git repository
func gitCommitAndPush(repoPath string, commitMessage *string) error {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return fmt.Errorf("failed to open repository: %v", err)
	}

	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %v", err)
	}

	// Add changes to staging
	_, err = w.Add("encoded_paths.json")
	if err != nil {
		return fmt.Errorf("failed to add file to staging: %v", err)
	}

	// Create commit message
	message := "Update encoded paths"
	if commitMessage != nil {
		message = *commitMessage
	}

	// Commit
	_, err = w.Commit(message, &git.CommitOptions{})
	if err != nil {
		return fmt.Errorf("failed to commit changes: %v", err)
	}

	// Push changes
	if err := repo.Push(&git.PushOptions{}); err != nil {
		return fmt.Errorf("failed to push changes: %v", err)
	}

	return nil
}
