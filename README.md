# Passy - Password Manager

Passy is a command-line password management tool designed to securely store and manage your passwords. It allows you to create, retrieve, and manage passwords efficiently, ensuring your sensitive information is protected. This README provides an overview of the available commands, their usage, and examples.

## Getting Started/Installation

Currently avaliable only installation with go: 
`go install github.com/koss-null/passy@latest`

If you want to use Passy as a password generator, just go ahead:
`passy --create --insane # supports --readable and --safe option, both are pretty safe though`

If you want to store your passwords in your git repo, you may want to generate a new secret key:
`passy --keygen /path/to/the/key.aes`

To continue setup you need to open (or create) file:
`~/.config/passy/config.toml`

With the following content:
```toml
PrivKeyPath = "/path/to/the/key.aes" # can be https link
GitRepoPath = "git@github.com:your-gh-account/your-repo-name.git"
```

Now you can try to store new password in your keystorage:
```bash
passy -a google.com --pass ChangeMe123
# also you may generate new password and save it in a single line
passy -a google.com --insane
# for folders just use > separator
passy -a "socials>facebook.com" --readable
```

## Flags

Passy allows you to manage your passwords through various commands. Below are the flags you can use along with additional links for further details.

### -a, --add
Add a new password associated with a specified key. The key separator is '>', allowing for hierarchical key structures (supports pass level key to generate the password automatically).

### --pass
Specify the password to be added (requires the `-a` flag).

### -p, --get-pass
Retrieve and display the password associated with the specified key.

### -k, --show-keys
List all keys for existing passwords, allowing you to see available entries in the password manager.

### --show-all
Display all existing keys and their associated passwords (requires the `-k` flag).

### -c, --compose
Generate a new password based on specified criteria, defaulting to a safe level of complexity.

### --readable
Create a password that is easy to read and remember, while still providing a moderate level of security (can be used with `-c` or `-a` for composition).

### --safe
Generate a password that balances security and memorability, suitable for general use (can be used with `-c` or `-a`).

### --insane
Compose a highly complex password that maximizes security but may be difficult to remember (can be used with `-c` or `-a`).

### -i, --interactive
Launch the Passy application in interactive mode for a guided password management experience [not implemented yet].

### --keygen
Generate a private encryption key and save it to the specified file path for secure password storage.

### -h, --help
Display this help message with available commands and their descriptions.

## Examples

Here are some common examples of how to use Passy with links to each command for further details:

0. **Generate secret key**
   ```bash
   # this will save new key by the given path
   passy --keygen /path/to/the/key.aes
   ```

1. **Add a new password**: 
   ```bash
   # do not forget to put key in " since > is interpreted as an operator in bash
   passy --add "myKey>subKey" --pass "mySecretPassword"
   ```
   [Details on Add Command](#-a--add)

2. **Retrieve a password**: 
   ```bash
   passy --get-pass myKey
   ```
   [Details on Get Pass Command](#-p--get-pass)

3. **Show all keys**: 
   ```bash
   passy --show-keys
   ```
   [Details on Show Keys Command](#-k--show-keys)

4. **Generate a password**: 
   ```bash
   passy --compose # generates --safe password
   passy --compose --readable # eg.: dEFvOSY3M3dlLW9nalU=
   passy --compose --safe # eg.: aWJlZTg1dkktRWt+ZXV6LU9Ob1VEd3U=
   passy --compose --insane # eg.: ¢@E?P¥Æ+a.ÀleZ©º.7Ì0Â$+;Ö&XÎ?$¸±-
   ```
   [Details on Compose Command](#-c--compose)

## Creating and Saving Passwords

To create a password and save it securely in your repository, follow these steps:

1. **Generate a Password**: Use the `--compose` flag to generate a password. For example:
   ```bash
   passy --compose --safe
   ```
   This will create a password that balances security and memorability.

2. **Add the Password**: Once you have generated a password, you can save it by using the `--add` flag. For example:
   ```bash
   passy --add myKey --pass "GeneratedPassword"
   ```
   Replace `"GeneratedPassword"` with the password you created.

3. **Verify the Password**: To ensure that your password has been saved correctly, you can retrieve it using:
   ```bash
   passy --get-pass myKey
   ```
   This command will display the password associated with `myKey`.

By following these steps, you can effectively manage your passwords and ensure they are stored securely in your repository.
