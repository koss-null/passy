package impl

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/koss-null/passy/internal/passgen"
	"github.com/koss-null/passy/internal/storage"
)

func HandlePasswordComposition(passLevelReadable, passLevelSafe, passLevelInsane bool) error {
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

func HandleGetPass(configPath, key string) error {
	flds, err := Folders(configPath)
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

func HandleAddPassword(
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

func HandleDeletePassword(configPath, key string) error {
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

func HandleShowKeys(configPath string, showAll bool) error {
	flds, err := Folders(configPath)
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

func Folders(configPath string) (*storage.Folder, error) {
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
