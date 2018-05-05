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

func getConfig(owner string, repository string) (*types.DerekConfig, error) {

	maintainersFile := fmt.Sprintf("https://github.com/%s/%s/raw/master/%s", owner, repository, configFile)

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	req, _ := http.NewRequest(http.MethodGet, maintainersFile, nil)

	res, resErr := client.Do(req)
	if resErr != nil {
		log.Fatalln(resErr)
	}

	if res.StatusCode != http.StatusOK {
		log.Fatalln(fmt.Sprintf("HTTP Status code: %d while fetching maintainers list (%s)", res.StatusCode, maintainersFile))
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	bytesOut, _ := ioutil.ReadAll(res.Body)
	var config types.DerekConfig

	err := parseConfig(bytesOut, &config)
	if err != nil {
		return nil, err
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
