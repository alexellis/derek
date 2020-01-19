// Copyright (c) OpenFaaS Author(s) 2019. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package types

import (
	"reflect"
	"testing"
)

func Test_mergeDerekRepoConfigs_OnlyOneItem(t *testing.T) {

	config1 := DerekRepoConfig{
		ContributingURL: "http://example.com",
	}

	config2 := DerekRepoConfig{
		ContributingURL: "http://two.example.com",
	}

	configOut, err := MergeDerekRepoConfigs(config1, config2)

	if err != nil {
		t.Errorf("Got error for a single plan, expected no error: %s", err.Error())
		t.Fail()
	}

	if config2.ContributingURL != configOut.ContributingURL {
		t.Errorf("ContributingURL want: %s, but got: %s", config2.ContributingURL, configOut.ContributingURL)
	}
}

func Test_mergeDerekRepoConfigs_MergeEmptyItemsFromBoth(t *testing.T) {

	config1 := DerekRepoConfig{
		ContributingURL: "http://example.com",
	}

	config2 := DerekRepoConfig{
		Redirect: "http://two.example.com",
	}

	configOut, err := MergeDerekRepoConfigs(config1, config2)

	if err != nil {
		t.Errorf("Got error for a single plan, expected no error: %s", err.Error())
		t.Fail()
	}

	if config1.ContributingURL != configOut.ContributingURL {
		t.Errorf("Redirect want: %s, but got: %s", config1.Redirect, configOut.Redirect)
	}
	if config2.Redirect != configOut.Redirect {
		t.Errorf("ContributingURL want: %s, but got: %s", config2.ContributingURL, configOut.ContributingURL)
	}
}

func Test_mergeDerekRepoConfigs_ConfigValuesAppendedToList(t *testing.T) {

	config1 := DerekRepoConfig{
		Maintainers: []string{"Waterdrips"},
	}

	config2 := DerekRepoConfig{
		Maintainers: []string{"alexellis"},
	}

	configOut, err := MergeDerekRepoConfigs(config1, config2)

	if err != nil {
		t.Errorf("Got error for a single plan, expected no error: %s", err.Error())
		t.Fail()
	}

	want := []string{"alexellis", "Waterdrips"}
	if !reflect.DeepEqual(want, configOut.Maintainers) {
		t.Errorf("OpenFaaSCloudVersion want: %s, but got: %s", want, configOut.Maintainers)
	}

}
