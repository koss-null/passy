package main

import (
	"fmt"

	"github.com/koss-null/passy/internal/command"
)

func main() {
	cmd := command.Parse()
	output := cmd.Do()
	fmt.Println(output)
}
