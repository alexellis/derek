// Copyright (c) Derek Author(s) 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handler

import (
	"testing"

	"github.com/google/go-github/github"
)

func Test_onlyMarkdownFiles(t *testing.T) {
	mdFileName1 := "readme.md"
	mdFileName2 := "README.MD"
	nonMDFileName := "main.go"

	var testCommits = []struct {
		files    []*github.CommitFile
		expected bool
	}{
		{
			files: []*github.CommitFile{
				&github.CommitFile{
					Filename: &mdFileName1,
				},
			},
			expected: true,
		},
		{
			files: []*github.CommitFile{
				&github.CommitFile{
					Filename: &mdFileName2,
				},
			},
			expected: true,
		},
		{
			files: []*github.CommitFile{
				&github.CommitFile{
					Filename: &mdFileName1,
				},
				&github.CommitFile{
					Filename: &mdFileName2,
				},
			},
			expected: true,
		},
		{
			files: []*github.CommitFile{
				&github.CommitFile{
					Filename: &nonMDFileName,
				},
			},
			expected: false,
		},
		{

			files: []*github.CommitFile{
				&github.CommitFile{
					Filename: &mdFileName1,
				},
				&github.CommitFile{
					Filename: &mdFileName2,
				},
				&github.CommitFile{
					Filename: &nonMDFileName,
				},
			},
			expected: false,
		},
	}

	for _, test := range testCommits {
		onlyMD := onlyMarkdownFiles(test.files)
		if onlyMD != test.expected {
			t.Errorf("Only markdown files - wanted %t, found %t", test.expected, onlyMD)
		}
	}
}
