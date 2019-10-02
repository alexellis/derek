// Copyright (c) Derek Author(s) 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handler

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/alexellis/derek/auth"
	"github.com/alexellis/derek/config"
	"github.com/alexellis/derek/factory"
	"github.com/alexellis/derek/types"
	"github.com/google/go-github/github"
)

const (
	prDescriptionRequiredLabel = "invalid"
	openedPRAction             = "opened"
)

var anonymousSign = regexp.MustCompile("Signed-off-by:(.*)noreply.github.com")

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
			return
		}
	}
}

func HandlePullRequest(req types.PullRequestOuter, contributingURL string, config config.Config) {
	ctx := context.Background()
	token, tokenErr := getAccessToken(config, req.Installation.ID)

	if tokenErr != nil {
		fmt.Printf("Error getting installation token: %s\n", tokenErr.Error())
		return
	}

	client := factory.MakeClient(ctx, token, config)

	if req.Action == openedPRAction {
		if req.PullRequest.FirstTimeContributor() == true {
			_, res, assignLabelErr := client.Issues.AddLabelsToIssue(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number, []string{"new-contributor"})
			if assignLabelErr != nil {
				log.Fatalf("%s limit: %d, remaining: %d", assignLabelErr, res.Limit, res.Remaining)
			}
		}
	}

	commits, err := fetchPullRequestCommits(req, client)
	if err != nil {
		log.Fatalf("unable to fetch pull request commits for PR %d: %s", req.PullRequest.Number, err)
	}

	issue, res, labelErr := client.Issues.Get(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number)
	if labelErr != nil {
		log.Fatalf("unable to fetch labels for PR %d: %s", req.PullRequest.Number, err)
	}
	fmt.Println("Rate limiting", res.Rate)
	res.Body.Close()

	anonymousSign := hasAnonymousSign(commits)
	unsignedCommits := hasUnsigned(commits)
	noDcoLabelExists := hasNoDcoLabel(issue)

	if !anonymousSign && !unsignedCommits {
		fmt.Println("Things look OK right now.")
		if noDcoLabelExists {
			fmt.Println("Removing label")
			res, removeLabelErr := client.Issues.RemoveLabelForIssue(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number, "no-dco")
			if removeLabelErr != nil {
				log.Fatalf("unable to remove DCO label from PR %d: %s", req.PullRequest.Number, err)
			}
			fmt.Println("Rate limiting", res.Rate)
			res.Body.Close()
		}
		return
	}

	var body string
	if anonymousSign {
		body = anonymousCommitComment(contributingURL)
	} else {
		body = unsignedCommitComment(contributingURL)
	}

	if !noDcoLabelExists {
		fmt.Println("Applying DCO label")
		_, res, assignLabelErr := client.Issues.AddLabelsToIssue(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number, []string{"no-dco"})
		if assignLabelErr != nil {
			log.Fatalf("unable to add DCO label to PR %d: %s", req.PullRequest.Number, err)
		}
		fmt.Println("Rate limiting", res.Rate)
		res.Body.Close()

		if err = createPullRequestComment(ctx, body, req, client); err != nil {
			log.Fatalf("unable to add comment on PR %d: %s", req.PullRequest.Number, err)
		}
	} else {
		fmt.Println("DCO label was previously applied")
	}
}

// VerifyPullRequestDescription checks that the PR has anything in the body.
// If there is no body, a label is added and comment posted to the PR with a link to the contributing guide.
func VerifyPullRequestDescription(req types.PullRequestOuter, contributingURL string, config config.Config) {
	ctx := context.Background()
	token, tokenErr := getAccessToken(config, req.Installation.ID)

	if tokenErr != nil {
		fmt.Printf("Error getting installation token: %s\n", tokenErr.Error())
		return
	}

	client := factory.MakeClient(ctx, token, config)

	if req.Action == openedPRAction {
		if !hasDescription(req.PullRequest) {
			fmt.Printf("Applying label: %s", prDescriptionRequiredLabel)
			_, res, assignLabelErr := client.Issues.AddLabelsToIssue(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number, []string{prDescriptionRequiredLabel})
			if assignLabelErr != nil {
				log.Fatalf("%s limit: %d, remaining: %d", assignLabelErr, res.Limit, res.Remaining)
			}

			body := emptyDescriptionComment(contributingURL)

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
	}
}

func anonymousCommitComment(contributingURL string) string {
	return `Thank you for your contribution. It seems that one or more of your commits have an anonymous email address. Please consider signing your commits with a valid email address. Please see our [contributing guide](` + contributingURL + `).`
}

func unsignedCommitComment(contributingURL string) string {
	return `Thank you for your contribution. I've just checked and your commit doesn't appear to be signed-off. That's something we need before your Pull Request can be merged. Please see our [contributing guide](` + contributingURL + `).
Tip: if you only have one commit so far then run: ` + "`" + `git commit --amend --signoff` + "`" + ` and then ` + "`" + `git push --force` + "`."
}

func emptyDescriptionComment(contributingURL string) string {
	return `Thank you for your contribution. I've just checked and your Pull Request doesn't appear to have any description.
That's something we need before your Pull Request can be merged. Please see our [contributing guide](` + contributingURL + `).`
}

func getAccessToken(config config.Config, installationID int) (string, error) {
	token := os.Getenv("personal_access_token")
	if len(token) == 0 {

		installationToken, tokenErr := auth.MakeAccessTokenForInstallation(
			config.ApplicationID,
			installationID,
			config.PrivateKey)

		if tokenErr != nil {
			return "", tokenErr
		}

		token = installationToken
	}

	return token, nil
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

func createPullRequestComment(ctx context.Context, body string, req types.PullRequestOuter, client *github.Client) error {
	comment := &github.IssueComment{
		Body: &body,
	}
	comment, resp, err := client.Issues.CreateComment(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number, comment)
	if err != nil {
		return err
	}

	fmt.Println("Rate limiting", resp.Rate)
	resp.Body.Close()
	return nil
}

func fetchPullRequestCommits(req types.PullRequestOuter, client *github.Client) ([]*github.RepositoryCommit, error) {
	ctx := context.Background()
	listOpts := &github.ListOptions{
		Page: 0,
	}
	commits, resp, err := client.PullRequests.ListCommits(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number, listOpts)
	if err != nil {
		log.Fatalf("Error getting PR %d\n%s", req.PullRequest.Number, err.Error())
		return nil, err
	}

	fmt.Println("Rate limiting", resp.Rate)
	resp.Body.Close()
	return commits, nil
}

func hasUnsigned(commits []*github.RepositoryCommit) bool {
	for _, commit := range commits {
		if commit.Commit != nil && commit.Commit.Message != nil {
			if isSigned(*commit.Commit.Message) == false {
				return true
			}
		}
	}
	return false
}

func hasAnonymousSign(commits []*github.RepositoryCommit) bool {
	for _, commit := range commits {
		if commit.Commit != nil && commit.Commit.Message != nil {
			if isAnonymousSign(*commit.Commit.Message) {
				return true
			}
		}
	}
	return false
}

func isAnonymousSign(msg string) bool {
	return anonymousSign.Match([]byte(msg))
}

func isSigned(msg string) bool {
	return strings.Contains(msg, "Signed-off-by:")
}

func hasDescription(pr types.PullRequest) bool {
	return len(strings.TrimSpace(pr.Body)) > 0
}

func isHacktoberfestSpam(req types.PullRequestOuter, client *github.Client) bool {
	commits, err := fetchPullRequestCommits(req, client)
	if err != nil {
		log.Fatalf("unable to fetch pull request commits for PR %d: %s", req.PullRequest.Number, err)
		return false
	}

	anonymousSign := hasAnonymousSign(commits)
	unsignedCommits := hasUnsigned(commits)

	onlyMD := onlyMarkdownFiles(commits)

	return onlyMD && req.PullRequest.FirstTimeContributor() && (anonymousSign || unsignedCommits)
}

func onlyMarkdownFiles(commits []*github.RepositoryCommit) bool {
	for _, c := range commits {
		if len(c.Files) == 0 {
			return false
		}

		for _, f := range c.Files {
			fileName := f.GetFilename()
			ext := fileName[strings.LastIndex(fileName, ".")+1:]
			if !strings.EqualFold(ext, "md") {
				return false
			}
		}
	}
	return true
}
