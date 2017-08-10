package main

import (
	"context"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/alexellis/derek/types"
	"github.com/google/go-github/github"
)

func handle(req types.PullRequestOuter) {
	client := github.NewClient(nil)

	hasUnsignedCommits, err := hasUnsigned(req, client)

	if err != nil {
		fmt.Println("Something went wrong: ", err)
	} else if hasUnsignedCommits {
		fmt.Println("Need to apply labels on item.")
	} else {
		fmt.Println("Things look OK right now.")
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
