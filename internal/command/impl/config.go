package impl

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/pkg/errors"

	"github.com/koss-null/passy/internal/storage"
)

// Open the file in the default editor
func HandleEditConfig(configPath string) error {
	// create new config to open it with default values, if not existed
	if err := storage.CheckConfigExistOrCreateNew(configPath); err != nil {
		return errors.Wrapf(err, "unable to achieve file")
	}

	totalFailMsg := fmt.Sprintf(
		"Unable to find any favorite editor, please edit the file by yourself.\n"+
			"Expected config to be in %q", configPath,
	)

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		fmt.Println("Trying Windows callouts")
		cmd = exec.Command("notepad.exe", configPath)
	case "darwin": // macOS
		fmt.Println("Trying macOS callouts")
		cmd = exec.Command("open", "-e", configPath)
	default: // *nix
		fmt.Println("Trying *nix callouts")

		editors := []string{"editor", "nano", "vim", "vi"}
		var lastErr error

		for _, editor := range editors {
			cmd = exec.Command(editor, configPath)

			// interactive term use
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			// check if the editor exists
			if path, err := exec.LookPath(editor); err == nil {
				fmt.Printf("Trying editor: %s (%s)\n", editor, path)
				err = cmd.Run()
				if err == nil {
					return nil
				}
				lastErr = err
				fmt.Printf("Editor %s failed: %v\n", editor, err)
			} else {
				fmt.Printf("Editor %s not found\n", editor)
			}
		}

		if lastErr != nil {
			return errors.Wrapf(lastErr, totalFailMsg)
		}
		return errors.New(totalFailMsg)
	}

	// For Windows and macOS
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return errors.Wrapf(cmd.Run(), totalFailMsg)
}
