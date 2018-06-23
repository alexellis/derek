// Copyright (c) Derek Author(s) 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package auth

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

const derekSecretKey = "derek-secret-key"

func getSecret(secretName string) (secretBytes []byte, err error) {

	secretPaths := []string{"/var/openfaas/secrets/", "/run/secrets/"}

	secretDir := filepath.Dir(secretName)
	if len(secretDir) > 0 {
		secretPaths = append([]string{secretDir}, secretPaths...)
	}

	secretName = filepath.Base(secretName)

	for _, path := range secretPaths {
		secretFile := filepath.Join(path, secretName)
		if secret, err := ioutil.ReadFile(secretFile); err == nil {
			return secret, nil
		}
	}
	return nil, fmt.Errorf("unable to read secret: %s", secretName)
}

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

	secretKey, err := getSecret(derekSecretKey)

	if err != nil {
		return fmt.Errorf("unable to read GitHub symmetrical secret, error: %s", err)
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
