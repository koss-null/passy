package command

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/term"

	"github.com/koss-null/passy/internal/passgen"
	"github.com/koss-null/passy/internal/storage"
)

const defaultConfigFile = "~/.config/passy/config.toml"

func NewCommand() *cobra.Command {
	var (
		interactive      bool
		listKeys         bool
		showAll          bool
		getPassword      string
		addPassword      string
		deletePassword   string
		passwordValue    string
		configPath       string
		editConfig       bool
		generateKey      string
		generatePassword bool
		passReadable     bool
		passSafe         bool
		passStrong       bool
	)

	cmd := &cobra.Command{
		Use:   "passy",
		Short: "Secure command-line password manager",
		Long: `Passy is an encrypted password manager that helps you:
- Generate strong, customizable passwords
- Securely store and organize credentials
- Quickly retrieve passwords when needed
- Manage hierarchical entries using '/' as separator

All data is encrypted and can be synced across devices.`,
	}

	// Disable default alphabetical sorting
	cmd.Flags().SortFlags = false

	passwordManagementFlags := cmd.Flags()
	passwordGenerationFlags := cmd.Flags()
	configurationFlags := cmd.Flags()
	advancedFlags := cmd.Flags()

	passwordManagementFlags.StringVarP(&addPassword, "add", "a", "",
		"Add new password entry with specified key path (e.g., 'email/gmail')")
	passwordManagementFlags.StringVar(&passwordValue, "password", "",
		"Specify password value directly (must be used with --add)")
	passwordManagementFlags.StringVarP(&getPassword, "get", "g", "",
		"Retrieve password by its key path (e.g., 'email/gmail')")
	passwordManagementFlags.StringVarP(&deletePassword, "delete", "d", "",
		"Permanently remove password entry or folder by key path")
	passwordManagementFlags.BoolVarP(&listKeys, "list", "l", false,
		"Display all stored password keys (hides passwords by default)")
	passwordManagementFlags.BoolVar(&showAll, "show-all", false,
		"Reveal passwords when listing (must be used with --list)")

	passwordGenerationFlags.BoolVarP(&generatePassword, "generate", "n", false,
		"Create new random password (uses 'safe' level by default)")
	passwordGenerationFlags.BoolVar(&passReadable, "readable", false,
		"Generate memorable password (combine with --generate or --add)")
	passwordGenerationFlags.BoolVar(&passSafe, "safe", false,
		"Generate balanced security/memorability password (default)")
	passwordGenerationFlags.BoolVar(&passStrong, "strong", false,
		"Generate maximum security password (harder to remember)")

	configurationFlags.StringVar(&configPath, "config", "",
		"Specify alternative configuration file path")
	configurationFlags.BoolVar(&editConfig, "edit-config", false,
		"Open configuration file in your default editor")

	advancedFlags.StringVar(&generateKey, "generate-key", "",
		"Generate new encryption key file at specified location")
	advancedFlags.BoolVarP(&interactive, "interactive", "i", false,
		"Launch interactive mode (menu-driven interface)")

	// Help Configuration
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Print(buildHelpText(cmd))
	})

	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Println("Basic Usage:")
		fmt.Println("  passy [command] [flags]")
		fmt.Println("\nCommon Commands:")
		fmt.Println("  passy --add email/gmail --password 'mypass'  # Add new password")
		fmt.Println("  passy --get email/gmail                     # Retrieve password")
		fmt.Println("  passy --generate --strong                   # Create strong password")
		return nil
	})

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return executeCommand(
			interactive,
			listKeys,
			showAll,
			getPassword,
			addPassword,
			deletePassword,
			passwordValue,
			configPath,
			editConfig,
			generateKey,
			generatePassword,
			passReadable,
			passSafe,
			passStrong,
		)
	}

	return cmd
}

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
