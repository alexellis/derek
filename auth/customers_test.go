package auth

import (
	"os"
	"testing"
)

func Test_findCustomerURL(t *testing.T) {

	var URIOpts = []struct {
		title       string
		envVar      string
		expectedURL string
	}{
		{
			title:       "No ENV var",
			envVar:      "",
			expectedURL: defaultCustomersURL,
		},
		{
			title:       "ENV var exists",
			envVar:      "https://raw.githubusercontent.com/rgee0/derek/master/.CUSTOMERS",
			expectedURL: "https://raw.githubusercontent.com/rgee0/derek/master/.CUSTOMERS",
		},
		{
			title:       "Invalid ENV var exists",
			envVar:      "raw.githubusercontent.com/rgee0/derek/master/.CUSTOMERS",
			expectedURL: "https://raw.githubusercontent.com/rgee0/derek/master/.CUSTOMERS",
		},
	}

	for _, test := range URIOpts {
		t.Run(test.title, func(t *testing.T) {

			os.Setenv(customersURLEnv, test.envVar)

			customersURL := buildCustomerURL()

			os.Unsetenv(customersURLEnv)

			if customersURL != test.expectedURL {
				t.Errorf("customer URL - wrong URL found - wanted: %s, found %s", test.expectedURL, customersURL)
			}
		})
	}
}
