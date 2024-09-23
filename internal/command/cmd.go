package command

import (
	"flag"

	"github.com/koss-null/passy/internal/passgen"
)

type Command struct {
	Do func() string
}

func Parse() Command {
	interactive := flag.Bool("i", false, "runs Passy in interactive mode [not implemented yet]")

	showKeys := flag.Bool("k", false, "show keys for all existing passwords")
	getPass := flag.String("p", "", "show pass by key")

	composePass := flag.Bool("c", false, "composes password")
	passLevelReadable := flag.Bool("readable", false, "composes password that is readable, easy to remember and pretty safe")
	passLevelSafe := flag.Bool("safe", false, "composes password that is safe and have chances to be remembered")
	passLevelInsane := flag.Bool("insane", false, "composes password that is insanly complex")

	flag.Parse()

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

	return Command{func() string { return "no command" }}
}
