package main

import (
	"fmt"
	"os"

	"github.com/koss-null/passy/internal/command"
)

func main() {
	rootCmd := command.NewCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
