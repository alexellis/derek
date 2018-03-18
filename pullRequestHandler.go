package main

import (
	"context"
	"fmt"
	"strings"

	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/alexellis/derek/auth"
	"github.com/alexellis/derek/types"
	"github.com/google/go-github/github"
)

func handlePullRequest(req types.PullRequestOuter, prFeatures types.PullRequestFeatures) {
	ctx := context.Background()

	token := os.Getenv("access_token")
	if len(token) == 0 {
		newToken, tokenErr := auth.MakeAccessTokenForInstallation(
			os.Getenv("application"),
			req.Installation.ID,
			os.Getenv("private_key"))

		if tokenErr != nil {
			log.Fatalln(tokenErr.Error())
		}

		token = newToken
	}

	client := auth.MakeClient(ctx, token)

	commits, commitFetchErr := getCommits(req, client)

	if commitFetchErr != nil {
		log.Fatal(commitFetchErr)
	}

	issue, _, labelErr := client.Issues.Get(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number)
	if labelErr != nil {
		log.Fatalln("Unable to fetch labels for PR %s/%s/#%d", req.Repository.Owner, req.Repository.Name, req.PullRequest.Number)
	}

	if prFeatures.CommitLintingFeature {
		lintResult := lintCommits(commits)
		applyErr := applyLintingLabel(req, client, issue, lintResult)
		if applyErr != nil {
			log.Printf("Error applying linting rule: %s", applyErr)
		}
	}

	if prFeatures.DCOCheckFeature {
		if hasUnsigned(commits) == true {
			fmt.Println("May need to apply labels on item.")

			fmt.Println("Current labels ", issue.Labels)

			if hasLabelAssigned("no-dco", issue) == false {
				fmt.Println("Applying label")
				_, res, assignLabelErr := client.Issues.AddLabelsToIssue(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number, []string{"no-dco"})
				if assignLabelErr != nil {
					log.Fatalf("%s limit: %d, remaining: %d", assignLabelErr, res.Limit, res.Remaining)
				}

				link := fmt.Sprintf("https://github.com/%s/%s/blob/master/CONTRIBUTING.md", req.Repository.Owner.Login, req.Repository.Name)
				body := `Thank you for your contribution. I've just checked and your commit doesn't appear to be signed-off.
That's something we need before your Pull Request can be merged. Please see our [contributing guide](` + link + `).`

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

			if hasLabelAssigned("no-dco", issue) {
				fmt.Println("Removing label")
				_, removeLabelErr := client.Issues.RemoveLabelForIssue(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number, "no-dco")
				if removeLabelErr != nil {
					log.Fatal(removeLabelErr)
				}
			}
		}
	}
}

func hasLabelAssigned(labelName string, issue *github.Issue) bool {
	if issue != nil {
		for _, label := range issue.Labels {
			if label.GetName() == labelName {
				return true
			}
		}
	}

	return false
}

func getCommits(req types.PullRequestOuter, client *github.Client) ([]*github.RepositoryCommit, error) {
	ctx := context.Background()

	var responseErr error
	listOpts := &github.ListOptions{
		Page: 0,
	}

	commits, resp, err := client.PullRequests.ListCommits(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number, listOpts)
	if err != nil {
		responseErr = fmt.Errorf("unable to fetch commits for PR: %d, Error: %s", req.PullRequest.Number, err.Error())
	}

	log.Printf("Rate limiting remaining: %d", resp.Rate.Remaining)
	return commits, responseErr
}

func hasUnsigned(commits []*github.RepositoryCommit) bool {
	hasUnsigned := false

	for _, commit := range commits {
		if commit.Commit != nil && commit.Commit.Message != nil {
			if isSigned(*commit.Commit.Message) == false {
				hasUnsigned = true
			}
		}
	}

	return hasUnsigned
}

func isSigned(msg string) bool {
	return strings.Contains(msg, "Signed-off-by:")
}

func lintCommits(commits []*github.RepositoryCommit) bool {
	for _, commit := range commits {
		if commit.Commit != nil && commit.Commit.Message != nil {
			if lintCommit(commit.Commit.Message) == false {
				return false
			}
		}
	}
	return true
}

func lintCommit(message *string) bool {
	var valid bool

	if message == nil {
		return false
	}

	parts := strings.Split(*message, "\n")

	if len(parts) > 0 {
		lengthValid := len(parts[0]) <= 50
		var startsWithUpper bool

		firstCharacter := getFirstCharacter(parts[0])
		if firstCharacter != nil {
			startsWithUpper = len(*firstCharacter) > 0 && strings.ToUpper(*firstCharacter) == *firstCharacter
		}

		valid = lengthValid && startsWithUpper
	}

	return valid
}

func getFirstCharacter(msg string) *string {
	var ret *string

	for _, runeVal := range msg {
		asStr := string(runeVal)
		ret = &asStr
		break
	}

	return ret
}

func applyLintingLabel(req types.PullRequestOuter, client *github.Client, issue *github.Issue, lintResult bool) error {
	labelCaption := "review/commit-message"
	var actionErr error
	hasLabel := hasLabelAssigned(labelCaption, issue)

	ctx := context.Background()

	if hasLabel {
		if lintResult == true {
			res, assignLabelErr := client.Issues.RemoveLabelForIssue(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number, labelCaption)
			if assignLabelErr != nil {
				actionErr = fmt.Errorf("removeLabel: %s limit: %d, remaining: %d", assignLabelErr, res.Limit, res.Remaining)
			}
		}
	} else {
		if lintResult == false {
			_, res, assignLabelErr := client.Issues.AddLabelsToIssue(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number, []string{labelCaption})
			if assignLabelErr != nil {
				actionErr = fmt.Errorf("addLabel: %s limit: %d, remaining: %d", assignLabelErr, res.Limit, res.Remaining)
			}
			link := fmt.Sprintf("https://github.com/%s/%s/blob/master/CONTRIBUTING.md", req.Repository.Owner.Login, req.Repository.Name)

			body := `Please check that your commit messages fit within [these guidelines](` + link + `):
- Commit subject should not exceed 50 characters
- Commit subject should start with an uppercase letter
`

			comment := &github.IssueComment{
				Body: &body,
			}

			_, _, err := client.Issues.CreateComment(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number, comment)
			if err != nil {
				actionErr = fmt.Errorf("unable to create comment due to linting check: %s", err)
			}
		}
	}

	return actionErr
}
