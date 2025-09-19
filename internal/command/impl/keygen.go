package impl

import (
	"fmt"
	"os"

	"github.com/koss-null/passy/internal/storage"
)

func HandleKeyGeneration(keyGen string) error {
	key, err := storage.GenerateAESKey(32)
	if err != nil {
		return err
	}

	err = os.WriteFile(keyGen, key, 0o644)
	if err != nil {
		return err
	}
	fmt.Printf("the file was successfully created: %s\n", keyGen)
	return nil
}
