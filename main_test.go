package main

import "testing"

func Test_removeNewLine(t *testing.T) {
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
			stringNoLines := removeNewLine([]byte(test.secret))
			if test.expectedByte != string(stringNoLines) {
				t.Errorf("String after removal - wanted: \"%s\", got \"%s\"", test.expectedByte, test.secret)
			}
		})
	}
}
