package handler

import (
	"context"
	"fmt"

	"github.com/alexellis/derek/config"
	"github.com/alexellis/derek/types"
	"github.com/google/go-github/github"
)

type merge struct {
	Config *config.Config
}

func (m *merge) Merge(req types.IssueCommentOuter, cmdType string, cmdValue string) (string, error) {
	result := ""

	client, ctx := makeClient(req.Installation.ID, *m.Config)

	if req.Issue.PullRequest == nil {
		return "can't merge a non-PR issue", nil
	}

	pr, _, err := client.PullRequests.Get(ctx, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number)
	if err != nil {
		return "unable to get pull request", err
	}

	if pr.GetMerged() == false {

		if pr.GetMergeable() == true {

			if validMergePolicy(req) == false {
				sendComment(client, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number,
					"I am unable to merge this PR due to merge-policy exception(s)")

				return "invalid merge policy", nil
			}

			pullRequestOptions := github.PullRequestOptions{
				MergeMethod: "rebase",
				CommitTitle: fmt.Sprintf("Merge PR #%d", req.Issue.Number),
			}
			mergeRes, _, err := client.PullRequests.Merge(ctx,
				req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number,
				fmt.Sprintf(`Merging PR #%d by Derek
This is an automated merge by the bot Derek, find more
https://github.com/alexellis/derek/
Signed-off-by: derek@openfaas.com`, req.Issue.Number), &pullRequestOptions)

			if err != nil {

				body := fmt.Sprintf(`I have been unable to merge the requested PR: %s`, err.Error())

				sendComment(client, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number,
					body)

				return fmt.Sprintf("Merge issue: %s, %t", mergeRes.GetMessage(), mergeRes.GetMerged()), err
			}

			sendComment(client, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number,
				`I have merged the pull request using the rebase strategy.`)
		} else {
			sendComment(client, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number,
				"This pull request cannot be merged. Rebase your work and try again.")
		}
	}

	return result, err
}

func sendComment(client *github.Client, login string, repo string, issue int, comment string) {

	issueComment := &github.IssueComment{
		Body: &comment,
	}
	client.Issues.CreateComment(context.Background(),
		login, repo, issue, issueComment)
}

func validMergePolicy(req types.IssueCommentOuter) bool {
	validDCO := true
	for _, label := range req.Issue.Labels {
		if label.Name == "no-dco" {
			validDCO = false
			break
		}
	}

	return validDCO
}
