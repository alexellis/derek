// Copyright (c) Derek Author(s) 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handler

import (
	"testing"

	"github.com/alexellis/derek/types"
)

func Test_maintainersparsed(t *testing.T) {
	config := types.DerekRepoConfig{}
	parseConfig([]byte(`maintainers:
- alexellis
- rgee0
`), &config)
	actual := len(config.Maintainers)
	if actual != 2 {
		t.Errorf("want: %d maintainers, got: %d", 2, actual)
	}
}

func Test_validateRedirectURL(t *testing.T) {
	type redirectURLTest struct {
		URL         string
		expectedErr bool
	}
	// append invalid domain tests
	tests := []redirectURLTest{
		redirectURLTest{
			URL:         "http://somedomain.com",
			expectedErr: true,
		},
		redirectURLTest{
			URL:         "www.somedomain.com",
			expectedErr: true,
		},
		redirectURLTest{
			URL:         "www.somedomain.com/github.com",
			expectedErr: true,
		},
	}
	// append valid domain tests
	for _, d := range getValidRedirectDomains() {
		tests = append(tests,
			redirectURLTest{URL: d, expectedErr: false},
			redirectURLTest{URL: "http://" + d, expectedErr: false},
			redirectURLTest{URL: "https://" + d, expectedErr: false},
		)
	}
	for _, test := range tests {
		err := validateRedirectURL(test.URL)
		if (err != nil) != test.expectedErr {
			t.Fatalf("URL: %q, expected error: %v, got: %v", test.URL, err != nil, test.expectedErr)
		}
	}
}

func Test_redirectparsed(t *testing.T) {
	url := "some-url"
	config := types.DerekRepoConfig{}
	parseConfig([]byte(`redirect: `+url), &config)
	actual := len(config.Redirect)
	lenURL := len(url)
	if actual != lenURL {
		t.Errorf("want: redirect URL of size %d, got: %d", lenURL, actual)
	}
}

func Test_curatorequalsmaintainer(t *testing.T) {
	config := types.DerekRepoConfig{}
	parseConfig([]byte(`curators:
- alexellis
- rgee0
`), &config)
	actual := len(config.Maintainers)
	if actual != 2 {
		t.Errorf("want: %d maintainers, got: %d", 2, actual)
	}
}

func Test_EnabledFeature(t *testing.T) {

	var enableFeatureOpts = []struct {
		title            string
		attemptedFeature string
		configFeatures   []string
		expectedVal      bool
	}{
		{
			title:            "dco enabled try dco case sensitive",
			attemptedFeature: "dco_check",
			configFeatures:   []string{"dco_check"},
			expectedVal:      true,
		},
		{
			title:            "dco enabled try dco case insensitive",
			attemptedFeature: "DCO_check",
			configFeatures:   []string{"dco_check"},
			expectedVal:      true,
		},
		{
			title:            "dco enabled try comments",
			attemptedFeature: "comments",
			configFeatures:   []string{"dco_check"},
			expectedVal:      false,
		},
		{
			title:            "Comments enabled try comments case insensitive",
			attemptedFeature: "Comments",
			configFeatures:   []string{"comments"},
			expectedVal:      true,
		},
		{
			title:            "Comments enabled try comments case sensitive",
			attemptedFeature: "comments",
			configFeatures:   []string{"comments"},
			expectedVal:      true,
		},
		{
			title:            "Comments enabled try dco",
			attemptedFeature: "dco",
			configFeatures:   []string{"comments"},
			expectedVal:      false,
		},
		{
			title:            "Non-existent feature",
			attemptedFeature: "gibberish",
			configFeatures:   []string{"dco_check", "comments"},
			expectedVal:      false,
		},
	}
	for _, test := range enableFeatureOpts {
		t.Run(test.title, func(t *testing.T) {

			inputConfig := &types.DerekRepoConfig{Features: test.configFeatures}

			featureEnabled := EnabledFeature(test.attemptedFeature, inputConfig)

			if featureEnabled != test.expectedVal {
				t.Errorf("Enabled feature - wanted: %t, found %t", test.expectedVal, featureEnabled)
			}
		})
	}
}

func Test_PermittedUserFeature(t *testing.T) {

	var permittedUserFeatureOpts = []struct {
		title            string
		attemptedFeature string
		user             string
		config           types.DerekRepoConfig
		expectedVal      bool
	}{
		{
			title:            "Valid feature with valid maintainer",
			attemptedFeature: "comment",
			user:             "Burt",
			config: types.DerekRepoConfig{
				Features:    []string{"comment"},
				Maintainers: []string{"Burt", "Tarquin", "Blanche"},
			},
			expectedVal: true,
		},
		{
			title:            "Valid feature with valid maintainer case insensitive",
			attemptedFeature: "comment",
			user:             "burt",
			config: types.DerekRepoConfig{
				Features:    []string{"comment"},
				Maintainers: []string{"Burt", "Tarquin", "Blanche"},
			},
			expectedVal: true,
		},
		{
			title:            "Valid feature with invalid maintainer",
			attemptedFeature: "comment",
			user:             "ernie",
			config: types.DerekRepoConfig{
				Features:    []string{"comment"},
				Maintainers: []string{"Burt", "Tarquin", "Blanche"},
			},
			expectedVal: false,
		},
		{
			title:            "Valid feature with invalid maintainer case insensitive",
			attemptedFeature: "Comment",
			user:             "ernie",
			config: types.DerekRepoConfig{
				Features:    []string{"comment"},
				Maintainers: []string{"Burt", "Tarquin", "Blanche"},
			},
			expectedVal: false,
		},
		{
			title:            "Invalid feature with valid maintainer",
			attemptedFeature: "invalid",
			user:             "Burt",
			config: types.DerekRepoConfig{
				Features:    []string{"comment"},
				Maintainers: []string{"Burt", "Tarquin", "Blanche"},
			},
			expectedVal: false,
		},
		{
			title:            "Invalid feature with valid maintainer case insensitive",
			attemptedFeature: "invalid",
			user:             "burt",
			config: types.DerekRepoConfig{
				Features:    []string{"comment"},
				Maintainers: []string{"Burt", "Tarquin", "Blanche"},
			},
			expectedVal: false,
		},
	}

	for _, test := range permittedUserFeatureOpts {
		t.Run(test.title, func(t *testing.T) {

			inputConfig := &test.config

			permittedFeature := PermittedUserFeature(test.attemptedFeature, inputConfig, test.user)

			if permittedFeature != test.expectedVal {
				t.Errorf("Permitted user feature - wanted: %t, found %t", test.expectedVal, permittedFeature)
			}
		})
	}
}
