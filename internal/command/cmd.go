package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/koss-null/passy/internal/command/impl"
	"github.com/koss-null/passy/internal/command/interactive"
)

const defaultConfigPath = "~/.config/passy/config.toml"

func NewCommand() *cobra.Command {
	var (
		isInteractive    bool
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
	advancedFlags.BoolVarP(&isInteractive, "interactive", "i", false,
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
			isInteractive,
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

func executeCommand(
	isInteractive bool,
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
		configPath = defaultConfigPath
	}

	if isInteractive {
		return interactive.Run(configPath)
	}

	if composePass {
		return impl.HandlePasswordComposition(passLevelReadable, passLevelSafe, passLevelInsane)
	}

	if showKeys {
		return impl.HandleShowKeys(configPath, showAll)
	}

	if editConfig {
		return impl.HandleEditConfig(configPath)
	}

	if getPass != "" {
		return impl.HandleGetPass(configPath, getPass)
	}

	if addPass != "" {
		return impl.HandleAddPassword(configPath, addPass, thePass, passLevelReadable, passLevelSafe, passLevelInsane)
	}

	if deletePass != "" {
		return impl.HandleDeletePassword(configPath, deletePass)
	}

	if keyGen != "" {
		return impl.HandleKeyGeneration(keyGen)
	}

	return fmt.Errorf("no valid command provided")
}
