// Copyright (c) Derek Author(s) 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// JWTAuth token issued by Github in response to signed JWT Token
type JWTAuth struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// MakeAccessTokenForInstallation makes an access token for an installation / private key
func MakeAccessTokenForInstallation(appID string, installation int, privateKey string) (string, error) {
	signed, err := GetSignedJwtToken(appID, privateKey)

	if err != nil {
		msg := fmt.Sprintf("can't run GetSignedJwtToken for app_id: %s and installation_id: %d, error: %v", appID, installation, err)

		fmt.Printf("Error %s\n", msg)
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("https://api.github.com/app/installations/%d/access_tokens", installation), nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", signed))
	req.Header.Add("Accept", "application/vnd.github.machine-man-preview+json")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		msg := fmt.Sprintf("can't get access_token for app_id: %s and installation_id: %d error: %v", appID, installation, err)
		fmt.Printf("Error: %s\n", msg)
		return "", fmt.Errorf("%s", msg)
	}

	defer res.Body.Close()

	bytesOut, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return "", readErr
	}

	jwtAuth := JWTAuth{}
	jsonErr := json.Unmarshal(bytesOut, &jwtAuth)
	if jsonErr != nil {
		return "", jsonErr
	}
	return jwtAuth.Token, nil
}

// GetSignedJwtToken get a tokens signed with private key
func GetSignedJwtToken(appID string, privateKey string) (string, error) {

	keyBytes := []byte(privateKey)

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
