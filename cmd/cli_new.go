package vault_go

import (
	"fmt"

	"github.com/urfave/cli"
)

// TODO
func RunNew(args []string) error {
	app := cli.NewApp()
	app.Name = "vault_go"
	app.Usage = "encrypt/decrypt Ansible data files by golang"
	app.Version = "0.0.1"

	app.Action = func(c *cli.Context) error {
		password, err := GetPassword(&cmdOpts{})
		fmt.Println(password)
		return err
	}

	return app.Run(args)
}
