package types

import (
	"testing"
)

func Test_FirstTimeContributor(t *testing.T) {
	type pullRequestTest struct {
		pullRequest  PullRequest
		expectedBool bool
	}
	authorLabel := []pullRequestTest{
		pullRequestTest{
			pullRequest:  PullRequest{AuthorAssociation: "NONE"},
			expectedBool: true,
		},
		pullRequestTest{
			pullRequest:  PullRequest{AuthorAssociation: "CONTRIBUTOR"},
			expectedBool: false,
		},
		pullRequestTest{
			pullRequest:  PullRequest{AuthorAssociation: "OWNER"},
			expectedBool: false,
		},
	}
	for _, test := range authorLabel {
		t.Run(test.pullRequest.AuthorAssociation, func(t *testing.T) {
			isFirstTime := test.pullRequest.FirstTimeContributor()
			if isFirstTime != test.expectedBool {
				t.Errorf("First time contributor - %s - want %t, got %t", test.pullRequest.AuthorAssociation, test.expectedBool, isFirstTime)
			}
		})
	}
}
