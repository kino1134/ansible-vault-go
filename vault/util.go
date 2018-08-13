package vault_go

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/pbkdf2"
)

// Ansible Vaultヘッダ(1.1)
const FILE_HEADER_11 = "$ANSIBLE_VAULT;1.1;AES256"

// Ansible Vaultヘッダ(1.2)
const FILE_HEADER_12 = "$ANSIBLE_VAULT;1.2;AES256;"

// Ansible Vaultヘッダを表す正規表現
const FILE_HEADER_PATTERN = "\\A\\$ANSIBLE_VAULT;1\\.\\d;AES(?:256)?(?:;.+)?"

// 共通鍵・HMAC鍵の長さ
const KEY_LENGTH = 32

// Initialization Vectorの長さ
const IV_LENGTH = 16

// 共通鍵生成時の繰り返し回数
const KDF_ITERATIONS = 10000

// 共通鍵・HMACを作成する際のハッシュ関数名
var HASH_ALGORITHM = sha256.New

// ソルト値をランダム生成する
func CreateSalt() ([]byte, error) {
	salt := make([]byte, KEY_LENGTH)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	return salt, nil
}

// 16進ハッシュ文字列を生成する
func CalcHMAC(cipherText, hmacKey []byte) string {
	algorithm := hmac.New(HASH_ALGORITHM, hmacKey)
	algorithm.Write(cipherText)
	digest := algorithm.Sum(nil)
	return hex.EncodeToString(digest)
}

// 鍵を生成する
func Derivekeys(password string, salt []byte) (cipherKey, hmacKey, iv []byte) {
	length := (2*KEY_LENGTH + IV_LENGTH)
	key := pbkdf2.Key([]byte(password), salt, KDF_ITERATIONS, length, HASH_ALGORITHM)

	return key[0:KEY_LENGTH],
		key[KEY_LENGTH : KEY_LENGTH*2],
		key[KEY_LENGTH*2 : KEY_LENGTH*2+IV_LENGTH]
}

// 暗号化・復号化を交互に行う
func ThroughCipher(src, cipherKey, iv []byte) ([]byte, error) {
	stream, err := setupCipher(cipherKey, iv)
	if err != nil {
		return nil, err
	}

	dst := make([]byte, len(src))
	stream.XORKeyStream(dst, src)

	return dst, err
}

// 暗号化器を生成する
func setupCipher(cipherKey, iv []byte) (cipher.Stream, error) {
	block, err := aes.NewCipher(cipherKey)
	if err != nil {
		return nil, err
	}
	return cipher.NewCTR(block, iv), nil
}
