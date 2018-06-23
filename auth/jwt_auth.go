package auth

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// GetSignedJwtToken get a tokens signed with private key
func GetSignedJwtToken(appID string, privateKey string) (string, error) {

	keyBytes, err := getSecret(privateKey)
	if err != nil {
		return "", fmt.Errorf("unable to read private key path: %s, error: %s", privateKey, err)
	}

	key, keyErr := jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
	if keyErr != nil {
		return "", keyErr
	}

	now := time.Now()
	claims := jwt.StandardClaims{
		Issuer:    appID,
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
