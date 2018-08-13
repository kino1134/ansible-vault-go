package vault_go

import (
	"io/ioutil"
	"os"
	"os/exec"

	vault "github.com/kino1134/ansible-vault-go/vault"
)

func edit(opts *cmdOpts) error {
	content, err := ioutil.ReadFile(opts.Path)
	if err != nil {
		return err
	}

	plainText, salt, err := vault.Decrypt(string(content), opts.Password)
	if err != nil {
		return err
	}

	ret, err := editText(opts, plainText)
	if err != nil {
		return err
	}

	if !opts.RestoreSalt {
		salt = nil
	}

	cipherText, err := vault.Encrypt(ret, opts.Password, opts.Label, salt)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(opts.Path, []byte(cipherText), 0644); err != nil {
		return err
	}

	return err
}

func editText(opts *cmdOpts, plainText string) (string, error) {
	tempFile, err := openTempFile(plainText)
	if err != nil {
		return "", err
	}
	defer os.Remove(tempFile.Name())

	if err = editor(opts, tempFile).Run(); err != nil {
		return "", err
	}

	ret, err := readTempFile(tempFile)
	if err != nil {
		return "", err
	}

	return string(ret), nil
}

func openTempFile(plainText string) (*os.File, error) {
	tempFile, err := ioutil.TempFile("", "vault")
	if err != nil {
		return nil, err
	}
	if err := writeTempFile(plainText, tempFile); err != nil {
		return nil, err
	}
	return tempFile, nil
}

func writeTempFile(plainText string, tempFile *os.File) error {
	if _, err := tempFile.Write([]byte(plainText)); err != nil {
		return err
	}
	if err := tempFile.Sync(); err != nil {
		return err
	}
	return nil
}

func readTempFile(tempFile *os.File) ([]byte, error) {
	ret, err := ioutil.ReadFile(tempFile.Name())
	if err != nil {
		return nil, err
	}
	if err = tempFile.Close(); err != nil {
		return nil, err
	}
	return ret, nil
}

func editor(opts *cmdOpts, tempFile *os.File) *exec.Cmd {
	cmd := exec.Command(opts.Editor, tempFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}
