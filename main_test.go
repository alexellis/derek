package main

import "testing"

func Test_removeNewLine(t *testing.T) {
	var exampleSecrets = []struct {
		secret       string
		expectedByte string
	}{
		{
			secret:       "nline \n",
			expectedByte: "nline ",
		},
		{
			secret: `newline
			`,
			expectedByte: "newline",
		},
		{
			secret:       `noline`,
			expectedByte: `noline`,
		},
		{
			secret:       "",
			expectedByte: "",
		},
	}
	for _, test := range exampleSecrets {

		t.Run(string(test.secret), func(t *testing.T) {
			stringNoLines := removeNewLine([]byte(test.secret))
			if test.expectedByte != string(stringNoLines) {
				t.Errorf("String after removal - wanted: \"%s\", got \"%s\"", test.expectedByte, test.secret)
			}
		})
	}
}
