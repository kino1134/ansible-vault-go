package vault_go

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestDecrypter(t *testing.T) {
	text, err := ioutil.ReadFile("../test/data/secret.txt")
	if err != nil {
		t.Error("Test File could not be read.", err)
	}
	password, err := ioutil.ReadFile("../test/data/vault_password")
	if err != nil {
		t.Error("Password File could not be read.", err)
	}

	password = bytes.TrimRight(password, "\n")
	result, salt, err := Decrypt(string(text), string(password))
	if err != nil {
		t.Error("Failed Decrypt.", err)
	}
	if salt == nil {
		t.Error("Not Exists salt.")
	}

	expected, err := ioutil.ReadFile("../test/data/plain.txt")
	if err != nil {
		t.Error("Expected File could not be read.", err)
	}
	if string(expected) != result {
		t.Error("UnExpected Decrypt Text")
	}
}
