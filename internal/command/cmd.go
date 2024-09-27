package command

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/koss-null/passy/internal/passgen"
	"github.com/koss-null/passy/internal/storage"
)

type Command struct {
	Do func() string
}

func Parse() Command {
	help := flag.Bool("h", false, "print help page")
	helpLong := flag.Bool("help", false, "print help page")

	interactive := flag.Bool("i", false, "run Passy in interactive mode [not implemented yet]")

	showKeys := flag.Bool("k", false, "show keys for all existing passwords")
	showAll := flag.Bool("show-all", false, "[-k ] show all existing keys and passwords")
	getPass := flag.String("p", "", "show pass by key")
	addPass := flag.String("a", "", "add password by key, key separator is '>' (supports pass level key to generate the pass automatically)")
	thePass := flag.String("pass", "", "[-a ] set password")
	keyGen := flag.String("keygen", "", "generate the private encryption key on given path")

	composePass := flag.Bool("c", false, "compose password (safe level by default)")
	passLevelReadable := flag.Bool("readable", false, "[-c|-a ] compose password that is readable, easy to remember and pretty safe")
	passLevelSafe := flag.Bool("safe", false, "[-c|-a ] compose password that is safe and have chances to be remembered")
	passLevelInsane := flag.Bool("insane", false, "[-c|-a ] compose password that is insanly complex")

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
		flds, err := folders()
		if err != nil {
			return Command{err.Error}
		}
		if showAll != nil && *showAll {
			return Command{flds.String("")}
		}
		return Command{flds.SecureString("")}
	}

	if getPass != nil && *getPass != "" {
		return Command{func() string { return "not implemented" }}
	}

	if addPass != nil && *addPass != "" {
		var pass string
		switch {
		case passLevelReadable != nil && *passLevelReadable:
			pass = passgen.GenReadablePass()
		case passLevelSafe != nil && *passLevelSafe:
			pass = passgen.GenReadablePass()
		case passLevelInsane != nil && *passLevelInsane:
			pass = passgen.GenInsanePass()
		default:
			return Command{func() string { return "please set the password strength option or [--pass] flag" }}
		}

		if thePass != nil && *thePass != "" {
			pass = *thePass
		}
		return Command{savePass(*addPass, pass)}
	}

	if keyGen != nil && *keyGen != "" {
		key, err := passgen.GenerateAESKey(32)
		if err != nil {
			return Command{err.Error}
		}

		err = os.WriteFile(*keyGen, key, fs.ModePerm)
		if err != nil {
			return Command{err.Error}
		}
		return Command{func() string { return "the file was successfully created: " + *keyGen }}
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

func folders() (*storage.Folder, error) {
	cfg, err := storage.ParseConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse config")
	}

	st, err := storage.New(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init storage")
	}

	folders, err := st.Decrypt()
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt")
	}
	return folders, nil
}

func savePass(key, pass string) func() string {
	cfg, err := storage.ParseConfig()
	if err != nil {
		return errors.Wrap(err, "failed to parse config").Error
	}

	st, err := storage.New(cfg)
	if err != nil {
		return errors.Wrap(err, "failed to init storage").Error
	}

	if err := st.Encrypt(key, pass, "", nil); err != nil {
		return errors.Wrap(err, "failed to encrypt").Error
	}

	return func() string { return fmt.Sprintf("the password %q was added successfully", pass) }
}
