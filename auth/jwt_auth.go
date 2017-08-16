package auth

import (
	"io/ioutil"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// GetSignedJwtToken get a tokens signed with private key
func GetSignedJwtToken(keyPath string) (string, error) {
	keyBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return "", err
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
