// Copyright (c) Derek Author(s) 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package factory

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

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
