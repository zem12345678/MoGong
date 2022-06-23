package make_password

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/rand"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

var iterations = 260000
var algorithm = "pbkdf2_sha256"
var salt_length = 22
var salt_chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func genSalt() []byte {
	var bytes = make([]byte, salt_length)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = salt_chars[v%byte(len(salt_chars))]
	}
	return bytes
}

func PBKDF2PasswordHasher(password, salt string) string {
	digest := sha256.New
	var dk []byte
	if salt == "" {
		dk = pbkdf2.Key([]byte(password), genSalt(), iterations, 32, digest)
	} else {
		dk = pbkdf2.Key([]byte(password), []byte(salt), iterations, 32, digest)
	}
	return fmt.Sprintf("%s$%d$%s$%s", algorithm, iterations, salt, base64.StdEncoding.EncodeToString(dk))
}

func CheckPassword(raw_password, password string) bool {
	salt := strings.Split(password, "$")[2]
	raw_hash := PBKDF2PasswordHasher(raw_password, salt)
	return strings.EqualFold(raw_hash, password)
}
