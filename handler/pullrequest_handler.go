// Copyright (c) Derek Author(s) 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handler

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

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
	actionRequiredConclusion   = "action_required"
	successConclusion          = "success"
)

// DCO is the check name
var DCO string = "DCO"

var anonymousSign = regexp.MustCompile("Signed-off-by:(.*)noreply.github.com")

func HandlePullRequest(req types.PullRequestOuter, contributingURL string, config config.Config) {
	ctx := context.Background()
	token, tokenErr := getAccessToken(config, req.Installation.ID)

	if tokenErr != nil {
		fmt.Printf("Error getting installation token: %s\n", tokenErr.Error())
		return
	}

	client := factory.MakeClient(ctx, token, config)

	checkErr := createSuccessfulCheck(req, client, ctx)
	if checkErr != nil {
		log.Fatalf("Error while creating successful DCO check: %s", checkErr.Error())
	}

	if req.Action == "review_requested" {
		log.Printf("[%s/%s] review_requested on PR %d, unable to process this request\n",
			req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number)
		return
	}

	if req.Action == openedPRAction {
		if req.PullRequest.FirstTimeContributor() == true {
			_, res, assignLabelErr := client.Issues.AddLabelsToIssue(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number, []string{"new-contributor"})
			if assignLabelErr != nil {
				log.Fatalf("[%s/%s] %s limit: %d, remaining: %d",
					req.Repository.Owner.Login, req.Repository.Name, assignLabelErr, res.Limit, res.Remaining)
			}
		}
	}

	commits, err := fetchPullRequestCommits(req, client)
	if err != nil {
		log.Fatalf("[%s/%s] unable to fetch pull request commits for PR %d: %s",
			req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number, err)
	}

	issue, res, labelErr := client.Issues.Get(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number)
	if labelErr != nil {
		log.Fatalf("[%s/%s] unable to fetch labels for PR %d: %s",
			req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number, err)
	}

	fmt.Printf("[%s/%s] %s rate limit\n",
		req.Repository.Owner.Login, req.Repository.Name, res.Rate)

	defer res.Body.Close()

	anonymousSign := hasAnonymousSign(commits)
	unsignedCommits := hasUnsigned(commits)
	noDcoLabelExists := hasNoDcoLabel(issue)

	if unsignedCommits {
		checkErr := updateExistingDCOCheck(req, client, ctx, actionRequiredConclusion)
		if checkErr != nil {
			log.Fatalf("Error while updating existing DCO check: %s", checkErr)
		}
	} else {
		checkErr := updateExistingDCOCheck(req, client, ctx, successConclusion)
		if checkErr != nil {
			log.Fatalf("Error while updating check: %s", checkErr.Error())
		}
	}

	if !anonymousSign && !unsignedCommits {
		if noDcoLabelExists {
			fmt.Printf("[%s/%s] Removing no-dco label: PR: %d\n",
				req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number)
			res, removeLabelErr := client.Issues.RemoveLabelForIssue(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number, "no-dco")
			if removeLabelErr != nil {
				log.Fatalf("unable to remove DCO label from PR %d: %s", req.PullRequest.Number, err)
			}
			fmt.Printf("[%s/%s] Removing no-dco label: PR: %d - %v\n",
				req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number, res.Rate)
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
		fmt.Printf("[%s/%s] Adding no-dco label: PR: %d\n",
			req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number)

		_, res, assignLabelErr := client.Issues.AddLabelsToIssue(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number, []string{"no-dco"})
		if assignLabelErr != nil {
			log.Fatalf("unable to add DCO label to PR %d: %v", req.PullRequest.Number, assignLabelErr)
		}
		fmt.Println("Rate limiting", res.Rate)

		if err = createPullRequestComment(ctx, body, req, client); err != nil {
			log.Fatalf("unable to add comment on PR %d: %v", req.PullRequest.Number, err)
		}
	} else {
		fmt.Printf("[%s/%s] DCO label already applied: PR: %d\n",
			req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number)
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
	_, resp, err := client.Issues.CreateComment(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number, comment)
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
		log.Fatalf("Error getting commits for PR %d\n%s", req.PullRequest.Number, err.Error())
		return nil, err
	}

	fmt.Println("Rate limiting", resp.Rate)
	resp.Body.Close()
	return commits, nil
}

func fetchPullRequestFileList(req types.PullRequestOuter, client *github.Client) ([]*github.CommitFile, error) {
	ctx := context.Background()
	listOpts := &github.ListOptions{
		Page: 0,
	}
	commitFiles, resp, err := client.PullRequests.ListFiles(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number, listOpts)
	if err != nil {
		log.Fatalf("Error getting files for PR %d\n%s", req.PullRequest.Number, err.Error())
		return nil, err
	}

	fmt.Println("Rate limiting", resp.Rate)
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	return commitFiles, nil
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

func createSuccessfulCheck(req types.PullRequestOuter, client *github.Client, ctx context.Context) error {
	checks, checksErr := determineExistingDCOCheck(req, client, ctx)
	if checksErr != nil {
		return fmt.Errorf("Error while creating successful DCO check: %s", checksErr.Error())
	}
	if *checks.Total > 1 {
		return fmt.Errorf("Error unexpected count of existing DCO checks: %d", *checks.Total)
	}
	if *checks.Total == 1 {
		return nil
	}

	check := createDCOCheck(req)
	_, apiResponse, apiErr := client.Checks.CreateCheckRun(ctx, req.Repository.Owner.Login, req.Repository.Name, check)
	if apiErr != nil {
		return fmt.Errorf("Error while creating successful DCO check: %s", apiErr.Error())
	}
	if apiResponse.StatusCode != 201 {
		return fmt.Errorf("Error while creating successful DCO check unexpected status code: %d", apiResponse.StatusCode)
	}
	return nil
}

func determineExistingDCOCheck(req types.PullRequestOuter, client *github.Client, ctx context.Context) (*github.ListCheckRunsResults, error) {
	checks, checkRes, checkErr := client.Checks.ListCheckRunsForRef(ctx,
		req.Repository.Owner.Login,
		req.Repository.Name,
		req.PullRequest.Head.SHA,
		&github.ListCheckRunsOptions{CheckName: &DCO})

	if checkRes.StatusCode != 200 {
		return nil, fmt.Errorf("Error unexpected status code while retreiving existing checks %d", checkRes.StatusCode)
	}
	if checkErr != nil {
		return nil, fmt.Errorf("Error while retreiving existing checks: %s", checkErr.Error())
	}
	return checks, nil
}

func createDCOCheck(req types.PullRequestOuter) github.CreateCheckRunOptions {
	now := github.Timestamp{time.Now()}
	status := "in_progress"
	text := "Thank you for the contribution, everything looks fine."
	title := "Signed commits"
	summary := "All of your commits are signed"
	check := github.CreateCheckRunOptions{
		StartedAt: &now,
		Name:      DCO,
		HeadSHA:   req.PullRequest.Head.SHA,
		Status:    &status,
		Output: &github.CheckRunOutput{
			Text:    &text,
			Title:   &title,
			Summary: &summary,
		},
	}
	conclusion := successConclusion
	check.Conclusion = &conclusion
	check.CompletedAt = &now
	return check
}

func updateExistingDCOCheck(req types.PullRequestOuter, client *github.Client, ctx context.Context, conclusion string) error {
	var check github.UpdateCheckRunOptions
	checks, checksErr := determineExistingDCOCheck(req, client, ctx)
	if checksErr != nil {
		return fmt.Errorf("Error while creating successful DCO check: %s", checksErr.Error())
	}
	if *checks.Total == 0 {
		return fmt.Errorf("Error unexpected count of existing DCO checks: %d", *checks.Total)
	}

	if conclusion == successConclusion {
		check = updateSuccessfulDCOCheck(checks)
	} else if conclusion == actionRequiredConclusion {
		check = updateUnsuccessfulDCOCheck(checks)
	}

	_, apiResponse, apiErr := client.Checks.UpdateCheckRun(ctx, req.Repository.Owner.Login, req.Repository.Name, *checks.CheckRuns[0].ID, check)
	if apiErr != nil {
		return fmt.Errorf("Error while updating the DCO check: %s", apiErr.Error())
	}
	if apiResponse.StatusCode != 200 {
		return fmt.Errorf("Error while updating the DCO check unexpected status code: %d", apiResponse.StatusCode)
	}
	return nil
}

func updateSuccessfulDCOCheck(checks *github.ListCheckRunsResults) github.UpdateCheckRunOptions {
	now := github.Timestamp{time.Now()}
	text := "Thank you for the contribution, everything looks fine."
	title := "Signed commits"
	summary := "All of your commits are signed"

	check := github.UpdateCheckRunOptions{
		Name: *checks.CheckRuns[0].Name,
		Output: &github.CheckRunOutput{
			Text:    &text,
			Title:   &title,
			Summary: &summary,
		},
	}
	conclusion := successConclusion
	check.Conclusion = &conclusion
	check.CompletedAt = &now
	return check
}

func updateUnsuccessfulDCOCheck(checks *github.ListCheckRunsResults) github.UpdateCheckRunOptions {
	now := github.Timestamp{time.Now()}
	text := `Thank you for your contribution. I've just checked and your commit doesn't appear to be signed-off.
	That's something we need before your Pull Request can be merged.`
	title := "Unsigned commits"
	summary := "One or more of the commits in this Pull Request are not signed-off."

	check := github.UpdateCheckRunOptions{
		Name: *checks.CheckRuns[0].Name,
		Output: &github.CheckRunOutput{
			Text:    &text,
			Title:   &title,
			Summary: &summary,
		},
	}
	conclusion := actionRequiredConclusion
	check.Conclusion = &conclusion
	check.CompletedAt = &now
	return check
}
