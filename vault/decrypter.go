package vault_go

import (
	"bytes"
	"encoding/hex"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// 復号化を行う
func Decrypt(text, password string) (string, []byte, error) {
	if err := isDecrypted(text); err != nil {
		return "", nil, err
	}

	_, salt, hmac, body, err := decodeText(text)
	if err != nil {
		return "", nil, err
	}

	cipherKey, hmacKey, iv := Derivekeys(password, salt)
	if hmac != CalcHMAC(body, hmacKey) {
		return "", nil, errors.New("HMACの値が一致しません。パスワードを確認してください。")
	}

	dst, err := ThroughCipher(body, cipherKey, iv)
	if err != nil {
		return "", nil, err
	}

	return string(unpadding(dst)), salt, nil
}

func isDecrypted(text string) error {
	matched, err := regexp.MatchString(FILE_HEADER_PATTERN, text)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("すでに復号化されています。")
	}
	return nil
}

// 各要素に分割する
func decodeText(text string) (header string, salt []byte, hmac string, body []byte, err error) {
	lines := strings.Split(text, "\n")
	header, rest := lines[0], strings.Join(lines[1:], "")

	restBytes, err := hex.DecodeString(rest)
	if err != nil {
		return
	}

	restList := bytes.SplitN(restBytes, []byte("\n"), 3)
	hmac = string(restList[1])
	salt, err = hex.DecodeString(string(restList[0]))
	if err != nil {
		return
	}
	body, err = hex.DecodeString(string(restList[2]))
	if err != nil {
		return
	}

	return
}

// 末尾からパディング文字を除く
func unpadding(src []byte) []byte {
	length := len(src)
	paddingLength := int(src[length-1])
	return src[0 : length-paddingLength]
}
