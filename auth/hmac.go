package auth

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
)

const derekSecretKey = "/run/secrets/derek-secret-key"

// CheckMAC verifies hash checksum
func CheckMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha1.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)

	return hmac.Equal(messageMAC, expectedMAC)
}

// ValidateHMAC validate a digest from Github via xHubSignature
func ValidateHMAC(bytesIn []byte, xHubSignature string) error {

	var validated error

	secretKey, err := ioutil.ReadFile(derekSecretKey)

	if err != nil {
		return fmt.Errorf("unable to read GitHub symmetrical secret: %s, error: %s", derekSecretKey, err)
	}

	if len(xHubSignature) > 5 {

		messageMAC := xHubSignature[5:] // first few chars are: sha1=
		messageMACBuf, _ := hex.DecodeString(messageMAC)
		secretKey = bytes.TrimRight(secretKey, "\n")

		res := CheckMAC(bytesIn, []byte(messageMACBuf), secretKey)
		if res == false {
			validated = fmt.Errorf("invalid message digest or secret")
		}
	}

	return validated
}
