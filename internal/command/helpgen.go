package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/term"
)

func buildHelpText(cmd *cobra.Command) string {
	var helpText strings.Builder

	// Get terminal width
	width := 80
	if fd := int(os.Stdout.Fd()); term.IsTerminal(fd) {
		if w, _, err := term.GetSize(fd); err == nil && w > 0 {
			width = w
		}
	}

	// Header
	helpText.WriteString(cmd.Long + "\n\n")

	// Usage
	helpText.WriteString("Usage:\n  passy [flags]\n\n")

	// Calculate column widths dynamically
	flagColWidth := 20                       // minimum space for flags
	descColWidth := width - flagColWidth - 4 // 4 = 2 spaces before + 2 spaces after flags

	// Password Management Section
	helpText.WriteString("Password Management:\n")
	printFlagSection(&helpText, cmd.Flags(), []string{
		"add",
		"password",
		"get",
		"delete",
		"list",
		"show-all",
	}, flagColWidth, descColWidth)

	// Password Generation Section
	helpText.WriteString("\nPassword Generation:\n")
	printFlagSection(&helpText, cmd.Flags(), []string{
		"generate",
		"readable",
		"safe",
		"strong",
	}, flagColWidth, descColWidth)

	// Configuration Section
	helpText.WriteString("\nConfiguration:\n")
	printFlagSection(&helpText, cmd.Flags(), []string{
		"config",
		"edit-config",
	}, flagColWidth, descColWidth)

	// Advanced Section
	helpText.WriteString("\nAdvanced:\n")
	printFlagSection(&helpText, cmd.Flags(), []string{
		"generate-key",
		"interactive",
	}, flagColWidth, descColWidth)

	// Help flag
	helpText.WriteString("\nHelp:\n")
	printFlagSection(&helpText, cmd.Flags(), []string{"help"}, flagColWidth, descColWidth)

	return helpText.String()
}

func printFlagSection(helpText *strings.Builder, flagSet *pflag.FlagSet, flagNames []string, flagColWidth, descColWidth int) {
	for _, name := range flagNames {
		flag := flagSet.Lookup(name)
		if flag == nil {
			continue
		}

		// Build flag representation
		flagRepr := buildFlagRepresentation(flag)

		// Wrap the description
		desc := flag.Usage
		wrappedDesc := wordWrap(desc, descColWidth, flagColWidth)

		// Print first line
		fmt.Fprintf(helpText, "  %-*s  %s\n", flagColWidth, flagRepr, wrappedDesc[0])
		// Print additional lines if description was wrapped
		for _, line := range wrappedDesc[1:] {
			fmt.Fprintf(helpText, "  %-*s  %s\n", flagColWidth, "", line)
		}
	}
}

func buildFlagRepresentation(flag *pflag.Flag) string {
	if flag.Shorthand != "" && flag.Shorthand != " " {
		return fmt.Sprintf("-%s, --%s", flag.Shorthand, flag.Name)
	}
	return fmt.Sprintf("    --%s", flag.Name)
}

func wordWrap(text string, lineWidth, indent int) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	currentLine := ""
	currentLength := 0

	for _, word := range words {
		if currentLength+len(word)+1 > lineWidth && currentLength > 0 {
			lines = append(lines, currentLine)
			currentLine = ""
			currentLength = 0
		}

		if currentLength == 0 {
			// First word in line
			currentLine = word
			currentLength = len(word)
			continue
		}
		currentLine += " " + word
		currentLength += len(word) + 1
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}
