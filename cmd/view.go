package vault_go

import (
	"fmt"
	"io/ioutil"

	vault "github.com/kino1134/ansible-vault-go/vault"
)

func view(opts *cmdOpts) error {
	content, err := ioutil.ReadFile(opts.Path)
	if err != nil {
		return err
	}

	text, _, err := vault.Decrypt(string(content), opts.Password)
	if err != nil {
		return err
	}

	fmt.Println(text)
	return err
}
