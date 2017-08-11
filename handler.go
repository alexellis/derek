package main

import (
	"context"
	"fmt"
	"strings"

	"os"

	"golang.org/x/oauth2"

	log "github.com/Sirupsen/logrus"
	"github.com/alexellis/derek/types"
	"github.com/google/go-github/github"
)

func makeClient(ctx context.Context, accessToken string) *github.Client {
	if len(accessToken) == 0 {
		return github.NewClient(nil)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return client
}

func handle(req types.PullRequestOuter) {
	ctx := context.Background()

	client := makeClient(ctx, os.Getenv("access_token"))

	hasUnsignedCommits, err := hasUnsigned(req, client)
	hasNoDcoLabel := false

	if err != nil {
		fmt.Println("Something went wrong: ", err)
	} else if hasUnsignedCommits {
		fmt.Println("May need to apply labels on item.")

		issue, _, labelErr := client.Issues.Get(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number)

		if labelErr != nil {
			log.Fatalln(labelErr)
		}
		fmt.Println("Current labels ", issue.Labels)

		for _, label := range issue.Labels {
			if label.GetName() == "no-dco" {
				hasNoDcoLabel = true
			}
		}

		if !hasNoDcoLabel {
			fmt.Println("Applying label")
			assignResult, _, assignLabelErr := client.Issues.AddLabelsToIssue(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number, []string{"no-dco"})
			fmt.Println(assignResult, assignLabelErr)

			link := fmt.Sprintf("https://github.com/%s/%s/blob/master/CONTRIBUTING.md", req.Repository.Owner.Login, req.Repository.Name)
			body := `Thank you for your contribution. I've just checked and your commit doesn't appear to be signed-off.
That's something we need before your Pull Request can be merged. Please see our [contributing guide](` + link + `).`

			comment := &github.IssueComment{
				Body: &body,
			}

			comment, resp, err := client.Issues.CreateComment(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number, comment)
			fmt.Println(comment, resp.Rate, err)
		}
	} else {
		fmt.Println("Things look OK right now.")
		issue, _, labelErr := client.Issues.Get(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number)

		if labelErr != nil {
			log.Fatalln(labelErr)
		}

		fmt.Println("Current labels ", issue.Labels)

		for _, label := range issue.Labels {
			if label.GetName() == "no-dco" {
				hasNoDcoLabel = true
			}
		}

		if hasNoDcoLabel {
			fmt.Println("Removing label")
			_, removeLabelErr := client.Issues.RemoveLabelForIssue(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number, "no-dco")
			fmt.Println(removeLabelErr)
		}

	}
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

			fmt.Printf("Commit - %s - signed-text: %t\n", commit.GetSHA(), isSigned(*commit.Commit.Message))
			fmt.Println(commit.Commit.Verification)

			if commit.Commit.Verification != nil {
				fmt.Println("Verification element")

				fmt.Printf("IsVerified? %t\n", commit.Commit.Verification.GetVerified())
				if commit.Commit.Verification.Signature != nil {
					fmt.Printf("Signature value: %s\n", *commit.Commit.Verification.Signature)
				}
			} else {
				fmt.Println("No verification")
			}
			fmt.Printf("Commit msg:\n'%s'\n", *commit.Commit.Message)
		}
	}

	return hasUnsigned, err
}

func isSigned(msg string) bool {
	return strings.Contains(msg, "Signed-off-by:")
}
