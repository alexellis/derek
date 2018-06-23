// Copyright (c) Derek Author(s) 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/alexellis/derek/types"

	log "github.com/Sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

const configFile = ".DEREK.yml"

func enabledFeature(attemptedFeature string, config *types.DerekConfig) bool {

	featureEnabled := false

	for _, availableFeature := range config.Features {
		if strings.EqualFold(attemptedFeature, availableFeature) {
			featureEnabled = true
			break
		}
	}
	return featureEnabled
}

func permittedUserFeature(attemptedFeature string, config *types.DerekConfig, user string) bool {

	permitted := false

	if enabledFeature(attemptedFeature, config) {
		for _, maintainer := range config.Maintainers {
			if strings.EqualFold(user, maintainer) {
				permitted = true
				break
			}
		}
	}

	return permitted
}

func readConfigFromURL(client http.Client, url string) []byte {
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	res, resErr := client.Do(req)
	if resErr != nil {
		log.Fatalln(resErr)
	}

	if res.StatusCode != http.StatusOK {
		log.Fatalln(fmt.Sprintf("HTTP Status code: %d while fetching config (%s)", res.StatusCode, url))
	}

	defer res.Body.Close()
	bytesOut, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return bytesOut
}

func getValidRedirectDomains() []string {
	return []string{"github.com", "www.github.com", "raw.githubusercontent.com"}
}

func validateRedirectURL(url string) error {
	for _, d := range getValidRedirectDomains() {
		if strings.HasPrefix(url, d) || strings.HasPrefix(url, "http://"+d) || strings.HasPrefix(url, "https://"+d) {
			return nil
		}
	}
	return fmt.Errorf("the redirect URL doesn't seem to be GitHub based")
}

func getConfig(owner string, repository string) (*types.DerekConfig, error) {
	var config types.DerekConfig

	client := http.Client{
		Timeout: 30 * time.Second,
	}
	configFile := fmt.Sprintf("https://github.com/%s/%s/raw/master/%s", owner, repository, configFile)
	bytesConfig := readConfigFromURL(client, configFile)

	err := parseConfig(bytesConfig, &config)
	if err != nil {
		return nil, err
	}

	// The config contains a redirect URL. Load the config from there.
	if len(config.Redirect) > 0 {
		err = validateRedirectURL(config.Redirect)
		if err != nil {
			return nil, err
		}
		bytesConfig = readConfigFromURL(client, config.Redirect)
		err = parseConfig(bytesConfig, &config)
		if err != nil {
			return nil, err
		}
	}

	return &config, nil
}

func parseConfig(bytesOut []byte, config *types.DerekConfig) error {
	err := yaml.Unmarshal(bytesOut, &config)

	if len(config.Maintainers) == 0 && len(config.Curators) > 0 {
		config.Maintainers = config.Curators
	}

	return err
}
