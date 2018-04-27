package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

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

// MakeClient make a HTTP client with a signed access token
func MakeClient(ctx context.Context, accessToken string) *github.Client {
	if len(accessToken) == 0 {
		return github.NewClient(nil)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return client
}
