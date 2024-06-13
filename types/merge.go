package types

import (
	"github.com/imdario/mergo"
)

func MergeDerekRepoConfigs(localConfig, remoteConfig DerekRepoConfig) (DerekRepoConfig, error) {

	mergeErr := mergo.Merge(&remoteConfig, &localConfig, mergo.WithAppendSlice)
	if mergeErr != nil {
		return remoteConfig, mergeErr
	}

	if len(localConfig.RequiredInIssues) > 0 {
		remoteConfig.RequiredInIssues = localConfig.RequiredInIssues
	}

	return remoteConfig, nil
}
