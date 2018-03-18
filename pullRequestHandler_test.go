package main

import (
	"testing"

	"github.com/google/go-github/github"
)

func Test_LintingResultGiven_WithLongSubject(t *testing.T) {

	var testCases = []struct {
		scenario   string
		message    string
		lintResult bool
	}{
		{
			scenario:   "Commit subject over 50 chars",
			message:    "This commit is necessary to make sure that all future commits conform to a certain set pattern\nSigned-off-by: Alex Ellis <alex@openfaas.com>",
			lintResult: false,
		},
		{
			scenario:   "Commit subject exactly 50 chars",
			message:    "This commit subject falls well within the boundar\nSigned-off-by: Alex Ellis <alex@openfaas.com>",
			lintResult: true,
		},
		{
			scenario:   "Commit subject starts with lowercase",
			message:    "has lowercase subject\nSigned-off-by: Alex Ellis <alex@openfaas.com>",
			lintResult: false,
		},
	}

	for _, testCase := range testCases {

		testMsg := testCase.message
		result := lintCommit(&testMsg)

		if result != testCase.lintResult {
			t.Logf("scenario: %s - want linting: %v, but got: %v\n  message: %s", testCase.scenario, testCase.lintResult, result, testCase.message)
			t.Fail()
		}
	}

}

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

			hasLabel := hasLabelAssigned("no-dco", inputIssue)

			if hasLabel != test.expectedBool {
				t.Errorf("Has no-dco label - wanted: %t, found %t", test.expectedBool, hasLabel)
			}
		})
	}
}
