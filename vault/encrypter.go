package vault_go

import (
	"encoding/hex"
	"errors"
	"regexp"
	"strings"
)

// AESにおける1ブロックあたりのバイト数
const BLOCK_SIZE = 16

// 一行あたりの文字数
const LINE_PER_CHARS = 80

// 暗号化を行う
func Encrypt(text, password, label string, srcSalt []byte) (string, error) {
	if err := isEncrypted(text); err != nil {
		return "", err
	}

	salt, err := setupSalt(srcSalt)
	if err != nil {
		return "", err
	}

	header := createHeader(label)
	src := padding(text)
	cipherKey, hmacKey, iv := Derivekeys(password, salt)

	dst, err := ThroughCipher(src, cipherKey, iv)
	if err != nil {
		return "", err
	}

	return encodeText(header, salt, CalcHMAC(dst, hmacKey), dst), nil
}

func isEncrypted(text string) error {
	matched, err := regexp.MatchString(FILE_HEADER_PATTERN, text)
	if err != nil {
		return err
	}
	if matched {
		return errors.New("すでに暗号化されています。")
	}
	return nil
}

func createHeader(label string) (header string) {
	header = FILE_HEADER_11
	if label != "" {
		header = FILE_HEADER_12 + label
	}
	return
}

func setupSalt(salt []byte) ([]byte, error) {
	if salt != nil {
		return salt, nil
	}
	return CreateSalt()
}

// ブロックサイズの倍数まで文字を埋める
func padding(text string) []byte {
	len := BLOCK_SIZE - len([]byte(text))%BLOCK_SIZE
	return []byte(text + strings.Repeat(string(len), len))
}

// 各要素を足し合わせる
func encodeText(header string, salt []byte, hmac string, dst []byte) string {
	rawBody := strings.Join([]string{
		hex.EncodeToString(salt),
		hmac,
		hex.EncodeToString(dst),
	}, "\n")

	body := split(hex.EncodeToString([]byte(rawBody)))

	result := strings.Join(append([]string{header}, body...), "\n")
	// 最終行にも改行を入れる
	return result + "\n"
}

// 80文字ごとにテキストを折り返す
func split(src string) []string {
	length := len(src)
	result := make([]string, 0)
	for i := 0; i*LINE_PER_CHARS < length; i++ {
		end := LINE_PER_CHARS * (i + 1)
		if length < end {
			end = length
		}
		result = append(result, src[LINE_PER_CHARS*i:end])
	}
	return result
}
