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

func Test_validatePullRequestBranch(t *testing.T) {

	pullBranchOptions := []struct {
		title             string
		headBranchName    string
		baseBranchName    string
		defaultBranchName string
		expectedResult    bool
		expectedMessage   string
		expectedLabels    string
	}{
		{
			title:             "Incorrectly named master head branch. Base branch equal to default.",
			headBranchName:    "master",
			baseBranchName:    "master",
			defaultBranchName: "master",
			expectedResult:    false,
			expectedMessage: "Thank you for your contribution. It appears that you are submitting changes directly from your master branch." +
				"Please raise a new pull request from a named branch i.e. `git checkout -b my_feature`.\n",
			expectedLabels: "review/source-branch",
		},
		{
			title:             "Correctly named non-master head branch. Base branch equal to default.",
			headBranchName:    "test_branch",
			baseBranchName:    "master",
			defaultBranchName: "master",
			expectedResult:    true,
			expectedMessage:   "",
			expectedLabels:    "",
		},
		{
			title:             "Incorrectly named master head branch. Base branch not equal to default.",
			headBranchName:    "master",
			baseBranchName:    "master",
			defaultBranchName: "development",
			expectedResult:    false,
			expectedMessage: "Thank you for your contribution. It appears that you are submitting changes directly from your master branch." +
				"Please raise a new pull request from a named branch i.e. `git checkout -b my_feature`.\n" +
				"Thank you for your contribution. It appears that you are submitting changes not against the default repository branch." +
				"Please raise a new pull request againgst the default branch: development\n",
			expectedLabels: "review/source-branchreview/target-branch",
		},
		{
			title:             "Correctly named non-master head branch. Base branch not equal to default.",
			headBranchName:    "test_branch",
			baseBranchName:    "master",
			defaultBranchName: "development",
			expectedResult:    false,
			expectedMessage: "Thank you for your contribution. It appears that you are submitting changes not against the default repository branch." +
				"Please raise a new pull request againgst the default branch: development\n",
			expectedLabels: "review/target-branch",
		},
		{
			title:             "Correctly named non-master head branch. Base branch not equal to default.",
			headBranchName:    "development",
			baseBranchName:    "master",
			defaultBranchName: "development",
			expectedResult:    false,
			expectedMessage: "Thank you for your contribution. It appears that you are submitting changes not against the default repository branch." +
				"Please raise a new pull request againgst the default branch: development\n",
			expectedLabels: "review/target-branch",
		},
		{
			title:             "Correctly named non-master head branch. Base branch equal to default.",
			headBranchName:    "development",
			baseBranchName:    "development",
			defaultBranchName: "development",
			expectedResult:    true,
			expectedMessage:   "",
			expectedLabels:    "",
		},
		{
			title:             "Incorrectly named master head branch. Base branch equal to default.",
			headBranchName:    "master",
			baseBranchName:    "development",
			defaultBranchName: "development",
			expectedResult:    false,
			expectedMessage: "Thank you for your contribution. It appears that you are submitting changes directly from your master branch." +
				"Please raise a new pull request from a named branch i.e. `git checkout -b my_feature`.\n",
			expectedLabels: "review/source-branch",
		},
		{
			title:             "Correctly named master head branch. Base branch equal to default.",
			headBranchName:    "test_branch",
			baseBranchName:    "development",
			defaultBranchName: "development",
			expectedResult:    true,
			expectedMessage:   "",
			expectedLabels:    "",
		},
	}
	for _, test := range pullBranchOptions {
		t.Run(test.title, func(t *testing.T) {
			repo := types.Repository{
				Name:          "test_repo",
				DefaultBranch: test.defaultBranchName,
			}
			headBranch := types.Branch{
				Repository: repo,
				Name:       test.headBranchName,
			}
			baseBranch := types.Branch{
				Repository: repo,
				Name:       test.baseBranchName,
			}
			req := types.PullRequestOuter{
				Repository: repo,
				BaseBranch: baseBranch,
				HeadBranch: headBranch,
			}

			validHeadAndBaseBranch, message, headBranchLabel, baseBranchLabel := validatePullRequestBranch(req)

			if validHeadAndBaseBranch != test.expectedResult {
				t.Errorf("Valid head and base branch check (head: %s, base: %s, default: %s) - want: %t, got %t",
					test.headBranchName, test.baseBranchName, test.defaultBranchName, test.expectedResult, validHeadAndBaseBranch)
			}
			if message != test.expectedMessage {
				t.Errorf("Valid head and base branch message check (head: %s, base: %s, default: %s) - want: %s, got %s",
					test.headBranchName, test.baseBranchName, test.defaultBranchName, test.expectedMessage, message)
			}
			if headBranchLabel+baseBranchLabel != test.expectedLabels {
				t.Errorf("Valid head and base branch message check (head: %s, base: %s, default: %s) - want: %s, got %s",
					test.headBranchName, test.baseBranchName, test.defaultBranchName, test.expectedLabels, headBranchLabel+baseBranchLabel)
			}
		})
	}
}
