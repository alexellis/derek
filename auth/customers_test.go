// Copyright (c) Derek Author(s) 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

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

func Test_customerValidationEnabled(t *testing.T) {
	tests := []struct {
		title        string
		value        string
		expectedBool bool
	}{
		{
			title:        "`validate_customers` is unset, defaults to on",
			value:        "",
			expectedBool: true,
		},
		{
			title:        "`validate_customers` is set with random value, defaults to on",
			value:        "random",
			expectedBool: true,
		},
		{
			title:        "`validate_customers` is set with explicit `0`",
			value:        "0",
			expectedBool: false,
		},
		{
			title:        "`validate_customers` is set with explicit `false`",
			value:        "false",
			expectedBool: false,
		},
	}
	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			os.Setenv("validate_customers", test.value)
			value := customerValidationEnabled()
			if value != test.expectedBool {
				t.Errorf("Expected value: %v got: %v", test.expectedBool, value)
			}
		})
	}
}
