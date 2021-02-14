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

	"github.com/alexellis/derek/config"
	"github.com/alexellis/derek/types"
	github "github.com/google/go-github/github"
)

const (
	configFile      = ".DEREK.yml"
	configURLFormat = "https://github.com/%s/%s/raw/%s/%s"
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

func readConfigFromURL(client http.Client, url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to make request to %q, %e", url, err)
	}

	res, resErr := client.Do(req)
	if resErr != nil {
		return nil, fmt.Errorf("could not action request url: %q, err: %s", url, resErr.Error())
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP Status code: %d while fetching config (%s)", res.StatusCode, req.URL.String())
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	bytesOut, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return bytesOut, nil
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
func GetPrivateRepoConfig(owner, repository, branch string, installation int, config config.Config) (*types.DerekRepoConfig, error) {
	client, ctx := makeClient(installation, config)
	response, err := client.Repositories.DownloadContents(ctx, owner, repository, configFile, &github.RepositoryContentGetOptions{
		Ref: branch,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to download config file: %s", err)
	}
	defer response.Close()

	bytesConfig, err := ioutil.ReadAll(response)
	if err != nil {
		return nil, fmt.Errorf("unable to read github's response: %s", err)
	}

	httpClient := http.Client{
		Timeout: 30 * time.Second,
	}
	return buildDerekConfig(httpClient, bytesConfig)
}

// GetRepoConfig returns derek's configuration for the specified
// repository. The repository has to be public since this function
// will fetch the file from the CDN. If you are trying to fetch
// the config from a private repo use `GetPrivateRepoConfig` instead.
func GetRepoConfig(owner, repository, branch string) (*types.DerekRepoConfig, error) {
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	configFile := fmt.Sprintf(configURLFormat, owner, repository, branch, configFile)
	bytesConfig, err := readConfigFromURL(client, configFile)
	if err != nil {
		return nil, err
	}

	return buildDerekConfig(client, bytesConfig)
}

func buildDerekConfig(client http.Client, bytesConfig []byte) (*types.DerekRepoConfig, error) {
	var localConfig types.DerekRepoConfig
	var remoteConfig types.DerekRepoConfig

	err := parseConfig(bytesConfig, &localConfig)
	if err != nil {
		return nil, err
	}

	// The config contains a redirect URL. Load the config from there.
	if len(localConfig.Redirect) > 0 {
		err = validateRedirectURL(localConfig.Redirect)
		if err != nil {
			return nil, err
		}

		bytesConfig, err = readConfigFromURL(client, localConfig.Redirect)
		if err != nil {
			return nil, err
		}

		err = parseConfig(bytesConfig, &remoteConfig)
		if err != nil {
			return nil, err
		}
	}

	mergedConfig, err := types.MergeDerekRepoConfigs(localConfig, remoteConfig)

	if err != nil {
		return &mergedConfig, err
	}

	return &mergedConfig, nil
}

func parseConfig(bytesOut []byte, config *types.DerekRepoConfig) error {
	err := yaml.Unmarshal(bytesOut, &config)

	if len(config.Maintainers) == 0 && len(config.Curators) > 0 {
		config.Maintainers = config.Curators
	}

	return err
}
