# Passy - Password Manager

**Passy** is a command-line password management tool designed to securely store and manage your passwords. It allows you to create, retrieve, and manage passwords efficiently, ensuring your sensitive information is protected. This README provides an overview of the available commands, their usage, and examples.

## Usage

```bash
passy [flags]
```

## Comprehensive Usage

Passy allows you to manage your passwords through various commands. Below are the flags you can use along with detailed examples to help you get started.

## Flags

### -a, --add
Add a new password associated with a specified key. The key separator is '>', allowing for hierarchical key structures (supports pass level key to generate the password automatically).

**Example:**
```bash
passy --add myKey>subKey
```
This command adds a new password under the key structure `myKey>subKey`.

### --pass
Specify the password to be added (requires the `-a` flag).

**Example:**
```bash
passy --add myKey --pass "mySecretPassword"
```
This command adds the password "mySecretPassword" to `myKey`.

### -p, --get-pass
Retrieve and display the password associated with the specified key.

**Example:**
```bash
passy --get-pass myKey
```
This command retrieves the password for `myKey`.

### -k, --show-keys
List all keys for existing passwords, allowing you to see available entries in the password manager.

**Example:**
```bash
passy --show-keys
```
This command displays all existing keys.

### --show-all
Display all existing keys and their associated passwords (requires the `-k` flag).

**Example:**
```bash
passy --show-keys --show-all
```
This command shows all keys along with their passwords.

### -c, --compose
Generate a new password based on specified criteria, defaulting to a safe level of complexity.

**Example:**
```bash
passy --compose
```
This command generates a new password using the default safe level.

### --readable
Create a password that is easy to read and remember, while still providing a moderate level of security (can be used with `-c` or `-a` for composition).

**Example:**
```bash
passy --compose --readable
```
This command generates a readable password.

### --safe
Generate a password that balances security and memorability, suitable for general use (can be used with `-c` or `-a`).

**Example:**
```bash
passy --compose --safe
```
This command generates a safe password.

### --insane
Compose a highly complex password that maximizes security but may be difficult to remember (can be used with `-c` or `-a`).

**Example:**
```bash
passy --compose --insane
```
This command generates an extremely complex password.

### -i, --interactive
Launch the Passy application in interactive mode for a guided password management experience [not implemented yet].

**Example:**
```bash
passy --interactive
```
This command starts the interactive mode.

### --keygen
Generate a private encryption key and save it to the specified file path for secure password storage.

**Example:**
```bash
passy --keygen /path/to/keyfile
```
This command generates a private encryption key and saves it to the specified path.

### -h, --help
Display this help message with available commands and their descriptions.

**Example:**
```bash
passy --help
```
This command shows the help message.

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

