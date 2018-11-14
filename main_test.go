package main

import (
	"os"
	"testing"
)

func Test_getContributingURL(t *testing.T) {
	var TestCases = []struct {
		Name            string
		ContributingURL string
		Owner           string
		RepositoryName  string
		ExpectedOuptput string
	}{
		{
			Name:            "Empty contributing URL",
			ContributingURL: "",
			Owner:           "openfaas",
			RepositoryName:  "faas",
			ExpectedOuptput: "https://github.com/openfaas/faas/blob/master/CONTRIBUTING.md",
		},
		{
			Name:            "Non empty contributing URL",
			ContributingURL: "https://github.com/openfaas/faas/blob/master/CONTRIBUTING.md",
			Owner:           "openfaas",
			RepositoryName:  "faas",
			ExpectedOuptput: "https://github.com/openfaas/faas/blob/master/CONTRIBUTING.md",
		},
	}

	for _, test := range TestCases {
		actualContrinbutingURL := getContributingURL(test.ContributingURL, test.Owner, test.RepositoryName)
		if actualContrinbutingURL != test.ExpectedOuptput {
			t.Errorf("Testcase %s failed. want - %s, got - %s", test.Name, test.ExpectedOuptput, actualContrinbutingURL)
		}
	}
}

func Test_customerValidation(t *testing.T) {
	tests := []struct {
		title        string
		value        string
		expectedBool bool
	}{
		{
			title:        "`validate_hmac` is unset, defaults to on",
			value:        "",
			expectedBool: true,
		},
		{
			title:        "`validate_hmac` is set with random value, defaults to on",
			value:        "random",
			expectedBool: true,
		},
		{
			title:        "`validate_hmac` is set with explicit `0`",
			value:        "0",
			expectedBool: false,
		},
		{
			title:        "`validate_hmac` is set with explicit `false`",
			value:        "false",
			expectedBool: false,
		},
	}
	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			os.Setenv("validate_hmac", test.value)
			value := hmacValidation()
			if value != test.expectedBool {
				t.Errorf("Expected value: %v got: %v", test.expectedBool, value)
			}
		})
	}
}
