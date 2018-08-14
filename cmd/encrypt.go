package vault_go

import (
	"errors"
	"io/ioutil"

	vault "github.com/kino1134/ansible-vault-go/vault"
)

func encrypt(opts *cmdOpts) error {
	content, err := ioutil.ReadFile(opts.Path)
	if err != nil {
		return err
	}

	if opts.PasswordFile == "" {
		confirm, err := inputPassword("パスワード(確認): ")
		if err != nil {
			return err
		}
		if confirm != opts.Password {
			return errors.New("パスワードが一致しません。")
		}
	}

	text, err := vault.Encrypt(string(content), opts.Password, opts.Label, nil)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(opts.Path, []byte(text), 0644); err != nil {
		return err
	}

	return nil
}
