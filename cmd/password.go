package vault_go

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

func SetPassword(opts *cmdOpts) error {
	if password, err := GetPassword(opts); err != nil {
		return err
	} else {
		opts.Password = password
		return nil
	}
}

func GetPassword(opts *cmdOpts) (string, error) {
	if opts.PasswordFile != "" {
		return readPassword(opts.PasswordFile)
	} else {
		password, err := inputPassword("パスワード: ")
		if err != nil {
			return "", err
		}

		if opts.Sub == "encrypt" {
			confirm, err := inputPassword("パスワード(確認): ")
			if err != nil {
				return "", err
			}
			if confirm != password {
				return "", errors.New("パスワードが一致しません。")
			}
		}

		return password, nil
	}
}

func readPassword(path string) (string, error) {
	result, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(result), "\n"), err
}

func inputPassword(prompt string) (string, error) {
	// パスワード入力後、プロンプトを元に戻す
	fd := int(syscall.Stdin)
	state, err := terminal.GetState(fd)
	if err != nil {
		return "", err
	}
	defer terminal.Restore(fd, state)

	// パスワード入力中に強制終了された場合も。プロンプトを元に戻す
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		s := <-c
		terminal.Restore(fd, state)
		fmt.Println("signal:", s)
		os.Exit(1)
	}()
	defer signal.Stop(c)

	print(prompt)
	input, err := terminal.ReadPassword(fd)
	print("\n")

	return string(input), err
}
