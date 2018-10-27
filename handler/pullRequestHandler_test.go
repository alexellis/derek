// Copyright (c) Derek Author(s) 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handler

import (
	"testing"

	"github.com/google/go-github/github"
)

func Test_isSigned(t *testing.T) {

	var signOffOpts = []struct {
		title        string
		message      string
		expectedBool bool
	}{
		{
			title:        "Correctly signed off full string",
			message:      "Signed-off-by:",
			expectedBool: true,
		},
		{
			title:        "Correctly signed off within string",
			message:      "This PR was Signed-off-by: rgee0",
			expectedBool: true,
		},
		{
			title:        "Not hypenated signed off within string",
			message:      "This PR was Signed off by: rgee0",
			expectedBool: false,
		},
		{
			title:        "Not hypenated signed full string",
			message:      "Signed off by:",
			expectedBool: false,
		},
	}
	for _, test := range signOffOpts {
		t.Run(test.title, func(t *testing.T) {

			containsSignoff := isSigned(test.message)

			if containsSignoff != test.expectedBool {
				t.Errorf("Is signed off - Testing '%s'  - wanted: %t, found %t", test.message, test.expectedBool, containsSignoff)
			}
		})
	}
}

func Test_hasNoDcoLabel(t *testing.T) {

	var labelOpts = []struct {
		title        string
		labels       []string
		expectedBool bool
	}{
		{
			title:        "Has the no-dco label",
			labels:       []string{"no-dco"},
			expectedBool: true,
		},
		{
			title:        "Doesnt have the no-dco label",
			labels:       []string{"proposal", "bug", "question"},
			expectedBool: false,
		},
		{
			title:        "Has the no-dco label amongst others",
			labels:       []string{"proposal", "bug", "question", "no-dco"},
			expectedBool: true,
		},
		{
			title:        "Has no labels",
			labels:       []string{},
			expectedBool: false,
		},
	}
	for _, test := range labelOpts {
		t.Run(test.title, func(t *testing.T) {

			var ghLabels []github.Label
			for _, label := range test.labels {
				ghLabels = append(ghLabels, github.Label{Name: &label})
			}

			inputIssue := &github.Issue{Labels: ghLabels}

			hasLabel := hasNoDcoLabel(inputIssue)

			if hasLabel != test.expectedBool {
				t.Errorf("Has no-dco label - wanted: %t, found %t", test.expectedBool, hasLabel)
			}
		})
	}
}
