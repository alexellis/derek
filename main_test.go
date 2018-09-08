package main

import "testing"

func Test_getFirstLine(t *testing.T) {
	var exampleSecrets = []struct {
		secret       string
		expectedByte string
	}{
		{
			secret:       "New-line \n",
			expectedByte: "New-line ",
		},
		{
			secret: `Newline and text 
			`,
			expectedByte: "Newline and text ",
		},
		{
			secret:       `Example secret2 `,
			expectedByte: `Example secret2 `,
		},
		{
			secret:       "\n",
			expectedByte: "",
		},
		{
			secret:       "",
			expectedByte: "",
		},
	}
	for _, test := range exampleSecrets {

		t.Run(string(test.secret), func(t *testing.T) {
			stringNoLines := getFirstLine([]byte(test.secret))
			if test.expectedByte != string(stringNoLines) {
				t.Errorf("String after removal - wanted: \"%s\", got \"%s\"", test.expectedByte, test.secret)
			}
		})
	}
}

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
