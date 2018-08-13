package vault_go

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"testing"
)

const saltString = "4a6b67ff79f7c495feede7d48cf3831694302eccf3e51c849626429d5473de8b"
const label = "label_test"

func TestEncrypt(t *testing.T) {
	text, err := ioutil.ReadFile("../test/data/plain.txt")
	if err != nil {
		t.Error("Test File could not be read.", err)
	}
	password, err := ioutil.ReadFile("../test/data/vault_password")
	if err != nil {
		t.Error("Password File could not be read.", err)
	}

	salt, err := hex.DecodeString(saltString)
	if err != nil {
		t.Error("Failed Decode salt.")
	}

	password = bytes.TrimRight(password, "\n")
	result, err := Encrypt(string(text), string(password), "", salt)
	if err != nil {
		t.Error("Failed Encrypt.", err)
	}

	expected, err := ioutil.ReadFile("../test/data/secret.txt")
	if err != nil {
		t.Error("Expected File could not be read.", err)
	}
	if string(expected) != result {
		t.Error("UnExpected Encrypt Text")
	}
}

func TestEncryptWithLabel(t *testing.T) {
	text, err := ioutil.ReadFile("../test/data/plain.txt")
	if err != nil {
		t.Error("Test File could not be read.", err)
	}
	password, err := ioutil.ReadFile("../test/data/vault_password")
	if err != nil {
		t.Error("Password File could not be read.", err)
	}

	salt, err := hex.DecodeString(saltString)
	if err != nil {
		t.Error("Failed Decode salt.")
	}

	password = bytes.TrimRight(password, "\n")
	result, err := Encrypt(string(text), string(password), label, salt)
	if err != nil {
		t.Error("Failed Encrypt.", err)
	}

	expected, err := ioutil.ReadFile("../test/data/secret_label.txt")
	if err != nil {
		t.Error("Expected File could not be read.", err)
	}
	if string(expected) != result {
		t.Error("UnExpected Encrypt Text")
	}
}
