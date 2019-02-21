// Copyright (c) Derek Author(s) 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"

	log "github.com/Sirupsen/logrus"
	"github.com/alexellis/derek/auth"
	"github.com/alexellis/derek/config"
	"github.com/alexellis/derek/types"
)

const (
	configFile             = ".DEREK.yml"
	configURLFormat        = "https://github.com/%s/%s/raw/master/%s"
	privateConfigURLFormat = "https://api.github.com/repos/%s/%s/contents/%s"
)

func EnabledFeature(attemptedFeature string, config *types.DerekRepoConfig) bool {

	featureEnabled := false

	for _, availableFeature := range config.Features {
		if strings.EqualFold(attemptedFeature, availableFeature) {
			featureEnabled = true
			break
		}
	}
	return featureEnabled
}

func PermittedUserFeature(attemptedFeature string, config *types.DerekRepoConfig, user string) bool {

	permitted := false

	if EnabledFeature(attemptedFeature, config) {
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
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	return fetchConfig(client, req)
}

func readConfigFromURLWithToken(client http.Client, url string, accessToken string) []byte {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))
	req.Header.Set("Accept", "application/vnd.github.v3.raw'")

	return fetchConfig(client, req)
}

func fetchConfig(client http.Client, req *http.Request) []byte {
	res, resErr := client.Do(req)
	if resErr != nil {
		log.Fatalln(resErr)
	}

	if res.StatusCode != http.StatusOK {
		log.Fatalln(fmt.Sprintf("HTTP Status code: %d while fetching config (%s)", res.StatusCode, req.URL.String()))
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

// GetPrivateRepoConfig returns the configuration for derek
// for the specified repository. Since the repository is
// private we use the github API to fetch `.DEREK.yml`.
func GetPrivateRepoConfig(owner, repository string, installation int, config config.Config) (*types.DerekRepoConfig, error) {
	accessToken, err := auth.MakeAccessTokenForInstallation(config.ApplicationID, installation, config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("unable to get a signed JWT token: %s", err)
	}

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	configFile := fmt.Sprintf(privateConfigURLFormat, owner, repository, configFile)
	bytesConfig := readConfigFromURLWithToken(client, configFile, accessToken)

	return buildDerekConfig(client, bytesConfig)
}

// GetRepoConfig returns derek's configuration for the specified
// repository. The repository has to be public since this function
// will fetch the file from the CDN. If you are trying to fetch
// the config from a private repo use `GetPrivateRepoConfig` instead.
func GetRepoConfig(owner string, repository string) (*types.DerekRepoConfig, error) {
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	configFile := fmt.Sprintf(configURLFormat, owner, repository, configFile)
	bytesConfig := readConfigFromURL(client, configFile)

	return buildDerekConfig(client, bytesConfig)
}

func buildDerekConfig(client http.Client, bytesConfig []byte) (*types.DerekRepoConfig, error) {
	var config types.DerekRepoConfig

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

func parseConfig(bytesOut []byte, config *types.DerekRepoConfig) error {
	err := yaml.Unmarshal(bytesOut, &config)

	if len(config.Maintainers) == 0 && len(config.Curators) > 0 {
		config.Maintainers = config.Curators
	}

	return err
}
