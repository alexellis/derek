// Copyright (c) Derek Author(s) 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handler

import (
	"context"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/alexellis/derek/config"
	"github.com/alexellis/derek/factory"
	"github.com/alexellis/derek/types"
	"github.com/google/go-github/github"
)

const (
	prDescriptionRequiredLabel = "invalid"
	openedPRAction             = "opened"
)

// HandleHacktoberfestPR checks for opened PR, first time contributor. If only .MD files are changed, issue is closed and invalid label is added
// The goal of this function is to mark pull requests invalid and close them from people only making typo changes without signing their commit (flybys)
func HandleHacktoberfestPR(req types.PullRequestOuter, contributingURL string, config config.Config) {
	ctx := context.Background()
	token, tokenErr := getAccessToken(config, req.Installation.ID)

	if tokenErr != nil {
		fmt.Printf("Error getting installation token: %s\n", tokenErr.Error())
		return
	}

	client := factory.MakeClient(ctx, token, config)

	if req.Action == openedPRAction {

		if isHacktoberfestSpam(req, client) {
			_, res, assignLabelErr := client.Issues.AddLabelsToIssue(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number, []string{"invalid"})
			if assignLabelErr != nil {
				log.Fatalf("%s limit: %d, remaining: %d", assignLabelErr, res.Limit, res.Remaining)
			}

			closeState := "close"
			input := &github.IssueRequest{State: &closeState}

			_, _, err := client.Issues.Edit(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number, input)
			if err != nil {
				log.Fatalf("unable to close pull request %d: %s", req.PullRequest.Number, err)
				return
			}

			fmt.Println(fmt.Sprintf("Request to close issue #%d was successful.\n", req.PullRequest.Number))

			body := hacktoberfestSpamComment(contributingURL)

			if err = createPullRequestComment(ctx, body, req, client); err != nil {
				log.Fatalf("unable to add comment on PR %d: %s", req.PullRequest.Number, err)
			}

			return
		}
	}
}

func isHacktoberfestSpam(req types.PullRequestOuter, client *github.Client) bool {
	commits, err := fetchPullRequestCommits(req, client)
	if err != nil {
		log.Fatalf("unable to fetch pull request commits for PR %d: %s", req.PullRequest.Number, err)
		return false
	}

	anonymousSign := hasAnonymousSign(commits)
	unsignedCommits := hasUnsigned(commits)

	files, err := fetchPullRequestFiles(req, client)
	if err != nil {
		log.Fatalf("unable to fetch pull request files for PR %d: %s", req.PullRequest.Number, err)
		return false
	}

	onlyMD := onlyMarkdownFiles(files)

	return onlyMD && req.PullRequest.FirstTimeContributor() && (anonymousSign || unsignedCommits)
}

func onlyMarkdownFiles(files []*github.CommitFile) bool {
	if len(files) == 0 {
		return false
	}

	for _, f := range files {
		fileName := f.GetFilename()
		ext := fileName[strings.LastIndex(fileName, ".")+1:]
		if !strings.EqualFold(ext, "md") {
			return false
		}
	}
	return true
}
