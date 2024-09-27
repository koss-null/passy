package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/koss-null/passy/internal/passgen"
)

type Command struct {
	Do func() string
}

func Parse() Command {
	help := flag.Bool("h", false, "print help page")
	helpLong := flag.Bool("help", false, "print help page")

	interactive := flag.Bool("i", false, "run Passy in interactive mode [not implemented yet]")

	showKeys := flag.Bool("k", false, "show keys for all existing passwords")
	getPass := flag.String("p", "", "show pass by key")

	composePass := flag.Bool("c", false, "compose password (safe level by default)")
	passLevelReadable := flag.Bool("readable", false, "[-c ] compose password that is readable, easy to remember and pretty safe")
	passLevelSafe := flag.Bool("safe", false, "compose password that is safe and have chances to be remembered")
	passLevelInsane := flag.Bool("insane", false, "compose password that is insanly complex")

	flag.Parse()

	if (help != nil && *help) || (helpLong != nil && *helpLong) {
		return Command{helpString}
	}

	if composePass != nil && *composePass {
		switch {
		case passLevelReadable != nil && *passLevelReadable:
			return Command{passgen.GenReadablePass}
		case passLevelSafe != nil && *passLevelSafe:
			return Command{passgen.GenSafePass}
		case passLevelInsane != nil && *passLevelInsane:
			return Command{passgen.GenInsanePass}
		default:
			return Command{passgen.GenSafePass}
		}
	}

	if interactive != nil && *interactive {
		return Command{func() string { return "not implemented" }}
	}

	if showKeys != nil && *showKeys {
		return Command{func() string { return "not implemented" }}
	}

	if getPass != nil && *getPass != "" {
		return Command{func() string { return "not implemented" }}
	}

	return Command{helpString}
}

func helpString() string {
	sb := strings.Builder{}
	sb.WriteString("Usage:\n")
	flag.VisitAll(func(f *flag.Flag) {
		sb.WriteString(fmt.Sprintf("  -%s, --%s  %s\n", f.Name, f.Name, f.Usage))
	})
	return sb.String()
}
