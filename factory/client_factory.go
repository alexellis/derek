// Copyright (c) Derek Author(s) 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

// Package factory is used to generate http.Client instances
package factory

import (
	"context"

	"github.com/alexellis/derek/config"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// MakeClient makes a HTTP client with a signed access token
func MakeClient(ctx context.Context, accessToken string, config config.Config) *github.Client {
	if len(accessToken) == 0 {
		return github.NewClient(nil)
	}

	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tokenClient := oauth2.NewClient(ctx, tokenSource)

	client := github.NewClient(tokenClient)
	return client
}
