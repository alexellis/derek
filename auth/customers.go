package auth

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/alexellis/derek/types"
)

const defaultCustomersURL string = "https://raw.githubusercontent.com/alexellis/derek/master/.CUSTOMERS"
const customersURLEnv string = "customers_url"

func findCustomersURL() string {

	if customURL, exists := os.LookupEnv(customersURLEnv); exists && (len(customURL) > 0) {

		if !strings.HasPrefix(strings.ToLower(customURL), "http") {
			customURL = fmt.Sprintf("https://%s", customURL)
		}

		return customURL
	}
	return defaultCustomersURL
}

func IsCustomer(repo types.Repository) (bool, error) {
	validate := os.Getenv("validate_customers")
	if len(validate) == 0 || (validate == "false" || validate == "0") {
		return true, nil
	}

	var err error
	var found bool
	c := http.Client{}
	customersURL := findCustomersURL()
	request, _ := http.NewRequest(http.MethodGet, customersURL, nil)
	res, doErr := c.Do(request)
	if err != nil {
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

		lines := strings.Split(strings.TrimSpace(string(body)), "\n")

		for _, line := range lines {
			if line == repo.Owner.Login {
				found = true
				break
			}
		}
	}

DO_RETURN:
	return found, err
}
