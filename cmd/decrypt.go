package vault_go

import (
	"io/ioutil"

	vault "github.com/kino1134/ansible-vault-go/vault"
)

func decrypt(opts *cmdOpts) error {
	content, err := ioutil.ReadFile(opts.Path)
	if err != nil {
		return err
	}

	text, _, err := vault.Decrypt(string(content), opts.Password)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(opts.Path, []byte(text), 0644); err != nil {
		return err
	}

	return nil
}
