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

func Test_hmacValidationDefault(t *testing.T) {
	os.Setenv("hmac_validation", "")
	enabled := hmacValidation()
	want := false
	if enabled != want {
		t.Errorf("want %t got %t", want, enabled)
		t.Fail()
	}
}
