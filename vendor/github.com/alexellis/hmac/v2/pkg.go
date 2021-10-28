package hmac

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"strings"
)

// CheckMAC verifies hash checksum
func CheckMAC(message, messageMAC, key []byte, sha func() hash.Hash) bool {
	mac := hmac.New(sha, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)

	return hmac.Equal(messageMAC, expectedMAC)
}

// Sign a message with the key and return bytes.
// Note: for human readable output see encoding/hex and
// encode string functions.
func Sign(message, key []byte, sha func() hash.Hash) []byte {
	mac := hmac.New(sha, key)
	mac.Write(message)
	signed := mac.Sum(nil)
	return signed
}

// Validate validate an encodedHash taken
// from GitHub via X-Hub-Signature HTTP Header.
// Note: if using another source, just add a 5 letter prefix such as "sha1="
func Validate(bytesIn []byte, encodedHash string, secretKey string) error {
	var validated error

	var hashFn func() hash.Hash
	var payload string

	if strings.HasPrefix(encodedHash, "sha1=") {
		payload = strings.TrimPrefix(encodedHash, "sha1=")

		hashFn = sha1.New

	} else if strings.HasPrefix(encodedHash, "sha256=") {
		payload = strings.TrimPrefix(encodedHash, "sha256=")

		hashFn = sha256.New
	} else {
		return fmt.Errorf("valid hash prefixes: [sha1=, sha256=], got: %s", encodedHash)
	}

	messageMAC := payload
	messageMACBuf, _ := hex.DecodeString(messageMAC)

	res := CheckMAC(bytesIn, []byte(messageMACBuf), []byte(secretKey), hashFn)
	if !res {
		validated = fmt.Errorf("invalid message digest or secret")
	}

	return validated
}

func init() {

}
