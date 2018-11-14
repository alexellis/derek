// Copyright (c) Derek Author(s) 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package auth

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	defaultCustomersURL string = "https://raw.githubusercontent.com/alexellis/derek/master/.CUSTOMERS"
	customersURLEnv     string = "customers_url"
)

func buildCustomerURL() string {

	if customURL, exists := os.LookupEnv(customersURLEnv); exists && (len(customURL) > 0) {

		if !strings.HasPrefix(strings.ToLower(customURL), "http") {
			customURL = fmt.Sprintf("https://%s", customURL)
		}

		return customURL
	}
	return defaultCustomersURL
}

// IsCustomer returns true if a customer is listed in the customers file.
// The validation is controlled by the 'validate_customers' env-var
func IsCustomer(ownerLogin string, c *http.Client) (bool, error) {
	validate := customerValidationEnabled()
	if validate == false {
		return true, nil
	}

	var err error
	var found bool

	customersURL := buildCustomerURL()

	request, _ := http.NewRequest(http.MethodGet, customersURL, nil)

	res, doErr := c.Do(request)
	if doErr != nil {
		err = doErr
		// Not sure how I feel about goto, but seems OK here (Alex Ellis)
		goto DO_RETURN
	}

	if res.Body != nil {
		defer res.Body.Close()
		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			err = readErr
			goto DO_RETURN
		}

		trimmedBody := strings.TrimSpace(string(body))
		lines := strings.Split(trimmedBody, "\n")

		for _, line := range lines {
			if line == ownerLogin {
				found = true
				break
			}
		}
	}

DO_RETURN:

	return found, err
}

func customerValidationEnabled() bool {
	validate := os.Getenv("validate_customers")
	return validate != "false" && validate != "0"
}
