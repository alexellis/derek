package auth

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/alexellis/derek/types"
)

func IsCustomer(repo types.Repository) (bool, error) {
	validate := os.Getenv("validate_customers")
	if len(validate) == 0 || (validate == "false" || validate == "0") {
		return true, nil
	}

	var err error
	var found bool
	c := http.Client{}
	request, _ := http.NewRequest(http.MethodGet, "https://raw.githubusercontent.com/alexellis/derek/master/.CUSTOMERS", nil)
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
