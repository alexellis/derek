// Copyright (c) Derek Author(s) 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handler

import (
	"testing"

	"github.com/alexellis/derek/types"
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

func Test_hasDescription(t *testing.T) {
	var pr = []struct {
		title        string
		body         string
		expectedBool bool
	}{
		{
			title:        "PR with body",
			body:         "This PR has a body",
			expectedBool: true,
		},
		{
			title:        "This PR has no body",
			body:         "",
			expectedBool: false,
		},
	}

	for _, test := range pr {
		testPr := types.PullRequest{Body: test.body}
		hasDescription := hasDescription(testPr)

		if hasDescription != test.expectedBool {
			t.Errorf("PR missing body - wanted: %t, found: %t", test.expectedBool, hasDescription)
		}
	}
}

func Test_isAnonymous(t *testing.T) {
	var anonymousSignOpts = []struct {
		title        string
		message      string
		expectedBool bool
	}{
		{
			title:        "Correctly signed off full string",
			message:      "Signed-off-by: openfaas",
			expectedBool: false,
		},
		{
			title:        "Correctly signed off within string",
			message:      "This PR was Signed-off-by: rgee0",
			expectedBool: false,
		},
		{
			title:        "Signed off with an anonymous email",
			message:      "Signed-off-by: users@noreply.github.com",
			expectedBool: true,
		},
	}
	for _, test := range anonymousSignOpts {
		t.Run(test.title, func(t *testing.T) {

			anonymousSign := isAnonymousSign(test.message)

			if anonymousSign != test.expectedBool {
				t.Errorf("Is anonymous sign - Testing '%s'  - wanted: %t, found %t", test.message, test.expectedBool, anonymousSign)
			}
		})
	}
}

func Test_hasAnonymousSign(t *testing.T) {
	var anonymousSignOpts = []struct {
		title    string
		commits  []*github.RepositoryCommit
		expected bool
	}{
		{
			title: "Does not have an anonymous commit",
			commits: []*github.RepositoryCommit{
				&github.RepositoryCommit{
					Commit: &github.Commit{
						Message: stringPtr("Signed-off-by: test"),
					},
				},
			},
			expected: false,
		},
		{
			title: "Has an anonymous commit",
			commits: []*github.RepositoryCommit{
				&github.RepositoryCommit{
					Commit: &github.Commit{
						Message: stringPtr("Signed-off-by: User users@users.noreply.github.com"),
					},
				},
			},
			expected: true,
		},
		{
			title: "Has an unsigned commit",
			commits: []*github.RepositoryCommit{
				&github.RepositoryCommit{
					Commit: &github.Commit{
						Message: stringPtr("Commit message"),
					},
				},
			},
			expected: false,
		},
	}

	for _, test := range anonymousSignOpts {
		t.Run(test.title, func(t *testing.T) {
			hasAnonymous := hasAnonymousSign(test.commits)
			if hasAnonymous != test.expected {
				t.Errorf("Has anonymous sign - wanted: %t, found %t", test.expected, hasAnonymous)
			}
		})
	}
}

func Test_onlyMarkdownFiles(t *testing.T) {
	mdFileName1 := "readme.md"
	mdFileName2 := "README.MD"
	nonMDFileName := "main.go"

	var testCommits = []struct {
		files    []*github.CommitFile
		expected bool
	}{
		{
			files: []*github.CommitFile{
				&github.CommitFile{
					Filename: &mdFileName1,
				},
			},
			expected: true,
		},
		{
			files: []*github.CommitFile{
				&github.CommitFile{
					Filename: &mdFileName2,
				},
			},
			expected: true,
		},
		{
			files: []*github.CommitFile{
				&github.CommitFile{
					Filename: &mdFileName1,
				},
				&github.CommitFile{
					Filename: &mdFileName2,
				},
			},
			expected: true,
		},
		{
			files: []*github.CommitFile{
				&github.CommitFile{
					Filename: &nonMDFileName,
				},
			},
			expected: false,
		},
		{

			files: []*github.CommitFile{
				&github.CommitFile{
					Filename: &mdFileName1,
				},
				&github.CommitFile{
					Filename: &mdFileName2,
				},
				&github.CommitFile{
					Filename: &nonMDFileName,
				},
			},
			expected: false,
		},
	}

	for _, test := range testCommits {
		onlyMD := onlyMarkdownFiles(test.files)
		if onlyMD != test.expected {
			t.Errorf("Only markdown files - wanted %t, found %t", test.expected, onlyMD)
		}
	}
}

func stringPtr(s string) *string {
	return &s
}
