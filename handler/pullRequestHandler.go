// Copyright (c) Derek Author(s) 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handler

import (
	"context"
	"fmt"
	"strings"

	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/alexellis/derek/auth"
	"github.com/alexellis/derek/config"
	"github.com/alexellis/derek/factory"
	"github.com/alexellis/derek/types"
	"github.com/google/go-github/github"
)

func HandlePullRequest(req types.PullRequestOuter, contributingURL string, config config.Config) {
	ctx := context.Background()

	token := os.Getenv("personal_access_token")
	if len(token) == 0 {

		newToken, tokenErr := auth.MakeAccessTokenForInstallation(
			config.ApplicationID,
			req.Installation.ID,
			config.PrivateKey)

		if tokenErr != nil {
			log.Fatalln(tokenErr.Error())
		}

		token = newToken
	}

	client := factory.MakeClient(ctx, token, config)

	hasUnsignedCommits, err := hasUnsigned(req, client)

	if req.Action == "opened" {
		if req.PullRequest.FirstTimeContributor() == true {
			_, res, assignLabelErr := client.Issues.AddLabelsToIssue(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number, []string{"new-contributor"})
			if assignLabelErr != nil {
				log.Fatalf("%s limit: %d, remaining: %d", assignLabelErr, res.Limit, res.Remaining)
			}
		}
	}

	if err != nil {
		log.Fatal(err)
	} else if hasUnsignedCommits {
		fmt.Println("May need to apply labels on item.")

		issue, _, labelErr := client.Issues.Get(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number)

		if labelErr != nil {
			log.Fatalln(labelErr)
		}
		fmt.Println("Current labels ", issue.Labels)

		if hasNoDcoLabel(issue) == false {
			fmt.Println("Applying label")
			_, res, assignLabelErr := client.Issues.AddLabelsToIssue(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number, []string{"no-dco"})
			if assignLabelErr != nil {
				log.Fatalf("%s limit: %d, remaining: %d", assignLabelErr, res.Limit, res.Remaining)
			}

			body := `Thank you for your contribution. I've just checked and your commit doesn't appear to be signed-off.
That's something we need before your Pull Request can be merged. Please see our [contributing guide](` + contributingURL + `).`

			comment := &github.IssueComment{
				Body: &body,
			}

			comment, resp, err := client.Issues.CreateComment(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number, comment)
			if err != nil {
				log.Fatalf("%s limit: %d, remaining: %d", assignLabelErr, resp.Limit, resp.Remaining)
				log.Fatal(err)
			}
			fmt.Println(comment, resp.Rate)
		}
	} else {
		fmt.Println("Things look OK right now.")
		issue, res, labelErr := client.Issues.Get(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number)

		if labelErr != nil {
			log.Fatalf("%s limit: %d, remaining: %d", labelErr, res.Limit, res.Remaining)
			log.Fatalln()
		}

		if hasNoDcoLabel(issue) {
			fmt.Println("Removing label")
			_, removeLabelErr := client.Issues.RemoveLabelForIssue(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number, "no-dco")
			if removeLabelErr != nil {
				log.Fatal(removeLabelErr)
			}
		}
	}
}

func hasNoDcoLabel(issue *github.Issue) bool {
	if issue != nil {
		for _, label := range issue.Labels {
			if label.GetName() == "no-dco" {
				return true
			}
		}
	}
	return false
}

func hasUnsigned(req types.PullRequestOuter, client *github.Client) (bool, error) {
	hasUnsigned := false
	ctx := context.Background()

	var err error
	listOpts := &github.ListOptions{
		Page: 0,
	}

	commits, resp, err := client.PullRequests.ListCommits(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number, listOpts)
	if err != nil {
		log.Fatalf("Error getting PR %d\n%s", req.PullRequest.Number, err.Error())
		return hasUnsigned, err
	}

	fmt.Println("Rate limiting", resp.Rate)

	for _, commit := range commits {
		if commit.Commit != nil && commit.Commit.Message != nil {
			if isSigned(*commit.Commit.Message) == false {
				hasUnsigned = true
			}
		}
	}

	return hasUnsigned, err
}

func isSigned(msg string) bool {
	return strings.Contains(msg, "Signed-off-by:")
}
