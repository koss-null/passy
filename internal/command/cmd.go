package command

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/koss-null/passy/internal/passgen"
	"github.com/koss-null/passy/internal/storage"
)

func helpString() string {
	return `Usage:
  passy [flags]

Flags:
  -a, --add                add password by key (key separator is '>') (supports pass level key to generate the pass automatically)
  --pass                   set password (use with -a)
  -p, --get-pass           show pass by key
  -k, --show-keys          show keys for all existing passwords
  --show-all               show all existing keys and passwords (use with -k)
  
  -c, --compose            compose password (safe level by default)
  --readable               compose password that is readable, easy to remember and pretty safe (use with -c or -a)
  --safe                   compose password that is safe and have chances to be remembered (use with -c or -a)
  --insane                 compose password that is insanly complex (use with -c or -a)

  -i, --interactive        run Passy in interactive mode [not implemented yet]
  --keygen                 generate the private encryption key on given path
  -h, --help               print help page
`
}

func NewCommand() *cobra.Command {
	var (
		interactive       bool
		showKeys          bool
		showAll           bool
		getPass           string
		addPass           string
		thePass           string
		keyGen            string
		composePass       bool
		passLevelReadable bool
		passLevelSafe     bool
		passLevelInsane   bool
	)

	cmd := &cobra.Command{
		Use:   "passy",
		Short: "A command-line password manager",
		Long:  `Passy is a password manager that allows you to generate, store, and retrieve passwords securely from your git repo.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Default action if no subcommand is specified
			_ = cmd.Help()
		},
	}

	cmd.SetHelpTemplate(helpString())

	cmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "run Passy in interactive mode [not implemented yet]")
	cmd.Flags().BoolVarP(&showKeys, "show-keys", "k", false, "show keys for all existing passwords")
	cmd.Flags().BoolVar(&showAll, "show-all", false, "[-k] show all existing keys and passwords")
	cmd.Flags().StringVarP(&getPass, "get-pass", "p", "", "show pass by key")
	cmd.Flags().StringVarP(&addPass, "add", "a", "", "add password by key, key separator is '>' (supports pass level key to generate the pass automatically)")
	cmd.Flags().StringVar(&thePass, "pass", "", "[-a] set password")
	cmd.Flags().StringVar(&keyGen, "keygen", "", "generate the private encryption key on given path")
	cmd.Flags().BoolVarP(&composePass, "compose", "c", false, "compose password (safe level by default)")
	cmd.Flags().BoolVar(&passLevelReadable, "readable", false, "[-c|-a] compose password that is readable, easy to remember and pretty safe")
	cmd.Flags().BoolVar(&passLevelSafe, "safe", false, "[-c|-a] compose password that is safe and have chances to be remembered")
	cmd.Flags().BoolVar(&passLevelInsane, "insane", false, "[-c|-a] compose password that is insanly complex")

	// Parse the command
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return executeCommand(interactive, showKeys, showAll, getPass, addPass, thePass, keyGen, composePass, passLevelReadable, passLevelSafe, passLevelInsane)
	}

	return cmd
}

func executeCommand(interactive, showKeys, showAll bool, getPass, addPass, thePass, keyGen string, _, passLevelReadable, passLevelSafe, passLevelInsane bool) error {
	if interactive {
		return fmt.Errorf("interactive mode is not implemented")
	}

	if showKeys {
		flds, err := folders()
		if err != nil {
			return err
		}
		if showAll {
			fmt.Println(flds.String("")())
		} else {
			fmt.Println(flds.SecureString("")())
		}
		return nil
	}

	if getPass != "" {
		return fmt.Errorf("not implemented")
	}

	if addPass != "" {
		var pass string
		switch {
		case passLevelReadable:
			pass = passgen.GenReadablePass()
		case passLevelSafe:
			pass = passgen.GenSafePass()
		case passLevelInsane:
			pass = passgen.GenInsanePass()
		default:
			return fmt.Errorf("please set the password strength option or [--pass] flag")
		}

		if thePass != "" {
			pass = thePass
		}
		return savePass(addPass, pass)
	}

	if keyGen != "" {
		key, err := passgen.GenerateAESKey(32)
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

	return fmt.Errorf("no valid command provided")
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

func savePass(key, pass string) error {
	cfg, err := storage.ParseConfig()
	if err != nil {
		return errors.Wrap(err, "failed to parse config")
	}

	st, err := storage.New(cfg)
	if err != nil {
		return errors.Wrap(err, "failed to init storage")
	}

	if err := st.Encrypt(key, pass, "", nil); err != nil {
		return errors.Wrap(err, "failed to encrypt")
	}

	fmt.Printf("the password %q was added successfully\n", pass)
	return nil
}
