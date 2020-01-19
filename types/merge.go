package types

import (
	"github.com/imdario/mergo"
)

func MergeDerekRepoConfigs(localConfig, remoteConfig DerekRepoConfig) (DerekRepoConfig, error) {

	mergeErr := mergo.Merge(&remoteConfig, &localConfig, mergo.WithAppendSlice)
	if mergeErr != nil {
		return remoteConfig, mergeErr
	}

	return remoteConfig, nil
}
