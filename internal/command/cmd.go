package command

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/koss-null/passy/internal/passgen"
	"github.com/koss-null/passy/internal/storage"
)

func helpString() string {
	return `Usage:
	passy [flag] [*value] [*flag] [*value]
* - optional
Flags:
	-a, --add [key_name]	Add a new password associated with a specified key. The key separator is '/', allowing for hierarchical key structures (supports pass level key to generate the password automatically).
	    --pass [password]	Specify the password to be added (requires -a flag).
	-p, --get-pass [key_name]	Retrieve and display the password associated with the specified key.
	-d, --delete [key_name]	Remove key or folder.
	-k, --show-keys		List all keys for existing passwords, allowing you to see available entries in the password manager.
	    --show-all	Display all existing keys and their associated passwords (requires -k flag).
	-c, --compose	Generate a new password based on specified criteria, defaulting to a safe level of complexity.
		--readable	Create a password that is easy to read and remember, while still providing a moderate level of security (can be used with -c or -a for composition).
		--safe	Generate a password that balances security and memorability, suitable for general use (can be used with -c or -a).
		--insane	Compose a highly complex password that maximizes security but may be difficult to remember (can be used with -c or -a).
	-i, --interactive	Launch the Passy application in interactive mode for a guided password management experience [not implemented yet].
	    --config [path]	Specify custom config file.
	    --config-edit	Start yor favorive editor to edit config.
	    --keygen	Generate a private encryption key and save it to the specified file path for secure password storage.
	-h, --help	Display this help message with available commands and their descriptions.
`
}

const defaultConfigFile = "~/.config/passy/config.toml"

func NewCommand() *cobra.Command {
	var (
		interactive       bool
		showKeys          bool
		showAll           bool
		getPass           string
		addPass           string
		deletePass        string
		thePass           string
		configPath        string
		editConfig        bool
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

	cmd.SetHelpTemplate(fmt.Sprint(helpString()))

	cmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "run Passy in interactive mode [not implemented yet]")
	cmd.Flags().BoolVarP(&showKeys, "show-keys", "k", false, "show keys for all existing passwords")
	cmd.Flags().BoolVar(&showAll, "show-all", false, "[-k] show all existing keys and passwords")
	cmd.Flags().StringVarP(&getPass, "get-pass", "p", "", "show pass by key")
	cmd.Flags().StringVarP(&addPass, "add", "a", "", "add password by key, key separator is '/' (supports pass level key to generate the pass automatically)")
	cmd.Flags().StringVarP(&deletePass, "delete", "d", "", "delete key or key folder, key separator is '/'")
	cmd.Flags().StringVar(&thePass, "pass", "", "[-a] set password")
	cmd.Flags().StringVar(&configPath, "config", "", "Specify custom config file")
	cmd.Flags().BoolVar(&editConfig, "config-edit", false, "Start your favorite editor to edit config")
	cmd.Flags().StringVar(&keyGen, "keygen", "", "generate the private encryption key on given path")
	cmd.Flags().BoolVarP(&composePass, "compose", "c", false, "compose password (safe level by default)")
	cmd.Flags().BoolVar(&passLevelReadable, "readable", false, "[-c|-a] compose password that is readable, easy to remember and pretty safe")
	cmd.Flags().BoolVar(&passLevelSafe, "safe", false, "[-c|-a] compose password that is safe and have chances to be remembered")
	cmd.Flags().BoolVar(&passLevelInsane, "insane", false, "[-c|-a] compose password that is insanly complex")

	// Parse the command
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return executeCommand(
			interactive,
			showKeys,
			showAll,
			getPass,
			addPass,
			deletePass,
			thePass,
			configPath,
			editConfig,
			keyGen,
			composePass,
			passLevelReadable,
			passLevelSafe,
			passLevelInsane,
		)
	}

	return cmd
}

func executeCommand(
	interactive bool,
	showKeys bool,
	showAll bool,
	getPass string,
	addPass string,
	deletePass string,
	thePass string,
	configPath string,
	editConfig bool,
	keyGen string,
	composePass bool,
	passLevelReadable bool,
	passLevelSafe bool,
	passLevelInsane bool,
) error {
	if configPath == "" {
		configPath = defaultConfigFile
	}

	if interactive {
		return fmt.Errorf("interactive mode is not implemented")
	}

	if composePass {
		return handlePasswordComposition(passLevelReadable, passLevelSafe, passLevelInsane)
	}

	if showKeys {
		return handleShowKeys(configPath, showAll)
	}

	if editConfig {
		return handleEditConfig(configPath)
	}

	if getPass != "" {
		return handleGetPass(configPath, getPass)
	}

	if addPass != "" {
		return handleAddPassword(configPath, addPass, thePass, passLevelReadable, passLevelSafe, passLevelInsane)
	}

	if deletePass != "" {
		return handleDeletePassword(configPath, deletePass)
	}

	if keyGen != "" {
		return handleKeyGeneration(keyGen)
	}

	return fmt.Errorf("no valid command provided")
}

func handlePasswordComposition(passLevelReadable, passLevelSafe, passLevelInsane bool) error {
	gen, err := passgen.New()
	if err != nil {
		return fmt.Errorf("unable to create generator: %v", err)
	}

	switch {
	case passLevelReadable:
		fmt.Println(gen.GenReadablePass())
	case passLevelSafe:
		fmt.Println(gen.GenSafePass())
	case passLevelInsane:
		fmt.Println(gen.GenInsanePass())
	default:
		fmt.Println(gen.GenSafePass())
	}
	return nil
}

func handleShowKeys(configPath string, showAll bool) error {
	flds, err := folders(configPath)
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

// Open the file in the default editor
func handleEditConfig(configPath string) error {
	// create new config to open it with default values, if not existed
	if err := storage.CheckConfigExistOrCreateNew(configPath); err != nil {
		return errors.Wrapf(err, "unable to achieve file")
	}

	totalFailMsg := fmt.Sprintf(
		"Unable to find any favorite editor, please edit the file by yourself.\n"+
			"Default config is in %q", defaultConfigFile,
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

func handleGetPass(configPath, key string) error {
	flds, err := folders(configPath)
	if err != nil {
		return err
	}

	sf, found := flds.GetSubFolder(key)
	if !found {
		fmt.Println("no such key")
	}
	fmt.Println(sf.String("")())
	return nil
}

func handleAddPassword(
	configPath, addPass, thePass string,
	passLevelReadable, passLevelSafe, passLevelInsane bool,
) error {
	gen, err := passgen.New()
	if err != nil {
		return fmt.Errorf("unable to create generator: %v", err)
	}

	var pass string
	switch {
	case passLevelReadable:
		pass = gen.GenReadablePass()
	case passLevelSafe:
		pass = gen.GenSafePass()
	case passLevelInsane:
		pass = gen.GenInsanePass()
	default:
		return fmt.Errorf("please set the password strength option or [--pass] flag")
	}

	if thePass != "" {
		pass = thePass
	}
	return savePass(configPath, addPass, pass)
}

func handleDeletePassword(configPath, key string) error {
	fmt.Printf("do you really want to delete %q [y/N]\n", key)
	var ans string
	_, err := fmt.Scanln(&ans)
	if err != nil {
		fmt.Printf("wow, we have some troubles reading your input")
		ans = "no"
	}

	if ans == "y" || ans == "Y" || ans == "yes" {
		cfg, err := storage.ParseConfig(configPath)
		if err != nil {
			return errors.Wrap(err, "failed to parse config")
		}

		st, err := storage.New(cfg)
		if err != nil {
			return errors.Wrap(err, "failed to init storage")
		}

		flds, err := st.Decrypt()
		if err != nil {
			return err
		}
		if err := flds.Delete(key); err != nil {
			return err
		}

		if err := st.Encrypt(flds); err != nil {
			return errors.Wrap(err, "failed to encrypt new password")
		}

		if err = st.Store(nil); err != nil {
			return errors.Wrap(err, "failed to store new password")
		}
	}

	fmt.Println("ok, leave everything as is")
	return nil
}

func handleKeyGeneration(keyGen string) error {
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

func folders(configPath string) (*storage.Folder, error) {
	cfg, err := storage.ParseConfig(configPath)
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

func savePass(configPath, key, pass string) error {
	cfg, err := storage.ParseConfig(configPath)
	if err != nil {
		return errors.Wrap(err, "failed to parse config")
	}

	st, err := storage.New(cfg)
	if err != nil {
		return errors.Wrap(err, "failed to init storage")
	}

	flds, err := st.Decrypt()
	if err != nil {
		return errors.Wrap(err, "failed to decrypt")
	}

	if err = flds.Add(key, pass); err != nil {
		return errors.Wrap(err, "failed to add a new key")
	}

	if err = st.Encrypt(flds); err != nil {
		return errors.Wrap(err, "failed to encrypt")
	}

	if err = st.Store(nil); err != nil {
		return errors.Wrap(err, "failed to store new password")
	}

	fmt.Printf("the password %q was added successfully\n", pass)
	return nil
}
