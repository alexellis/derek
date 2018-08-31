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

// GetSignedJwtToken get a tokens signed with private key
func GetSignedJwtToken(appID string, privateKeyPath string) (string, error) {

	keyBytes, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return "", fmt.Errorf("unable to read private key path: %s, error: %s", privateKeyPath, err)
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

// JwtAuth token issued by Github in response to signed JWT Token
type JwtAuth struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// MakeAccessTokenForInstallation makes an access token for an installation / private key
func MakeAccessTokenForInstallation(appID string, installation int, privateKeyPath string) (string, error) {
	signed, err := GetSignedJwtToken(appID, privateKeyPath)

	if err == nil {
		c := http.Client{}
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("https://api.github.com/installations/%d/access_tokens", installation), nil)

		req.Header.Add("Authorization", "Bearer "+signed)
		req.Header.Add("Accept", "application/vnd.github.machine-man-preview+json")

		res, err := c.Do(req)

		if err == nil {
			defer res.Body.Close()
			bytesOut, readErr := ioutil.ReadAll(res.Body)
			if readErr != nil {
				return "", readErr
			}
			jwtAuth := JwtAuth{}
			jsonErr := json.Unmarshal(bytesOut, &jwtAuth)
			if jsonErr != nil {
				return "", jsonErr
			}

			return string(jwtAuth.Token), nil
		}
	}

	return "", err
}
