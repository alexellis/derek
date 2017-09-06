package auth

import (
	"fmt"
	"io/ioutil"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// GetSignedJwtToken get a tokens signed with private key
func GetSignedJwtToken(keyPath string) (string, error) {
	if len(keyPath) == 0 {
		return "", fmt.Errorf("unable to read from empty keypath, try setting env: \"private_key\" to a filename and path")
	}

	keyBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return "", fmt.Errorf("unable to read keypath: %s, error: %s", keyPath, err)
	}

	key, keyErr := jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
	if keyErr != nil {
		return "", keyErr
	}

	now := time.Now()
	claims := jwt.StandardClaims{
		Issuer:    "4385",
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(time.Minute * 9).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signedVal, signErr := token.SignedString(key)
	if signErr != nil {
		return "", signErr
	}

	return string(signedVal), nil
}
