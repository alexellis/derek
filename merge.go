package main

import (
	"context"
	"fmt"
	"log"

	"github.com/alexellis/derek/auth"
	"github.com/alexellis/derek/types"
	"github.com/google/go-github/github"
)

type merge struct {
}

func (m *merge) Merge(req types.IssueCommentOuter, cmdType string, cmdValue string) (string, error) {
	result := ""

	if req.Issue.PullRequest != nil && len(req.Issue.PullRequest.URL) > 0 {
		log.Println("Wants to merge a PR")

		token := getAccessToken(req.Installation.ID)
		client := auth.MakeClient(context.Background(), token)
		pr, _, err := client.PullRequests.Get(context.Background(), req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number)
		if err != nil {

			if pr.GetMerged() == false {
				if pr.GetMergeable() == true {

					pullRequestOptions := github.PullRequestOptions{
						MergeMethod: "rebase",
						CommitTitle: fmt.Sprintf("Merge PR #%d", req.Issue.Number),
					}
					mergeRes, _, err := client.PullRequests.Merge(context.Background(),
						req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number,
						fmt.Sprintf(`Merging PR #%d by Derek
This is an automated merge by the bot Derek, find more
https://github.com/alexellis/derek/

Signed-off-by: derek@openfaas.com`, req.Issue.Number), &pullRequestOptions)

					if err != nil {

						body := fmt.Sprintf(`I have been unable to merge the requested PR: %s`, err.Error())
						comment := &github.IssueComment{
							Body: &body,
						}

						client.Issues.CreateComment(context.Background(),
							req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number, comment)

						return fmt.Sprintf("Merge issue: %s, %t", mergeRes.GetMessage(), mergeRes.GetMerged()), err
					}

					body := `I have merged the pull request using the rebase strategy.`
					comment := &github.IssueComment{
						Body: &body,
					}

					_, _, commentErr := client.Issues.CreateComment(context.Background(),
						req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number, comment)

					if commentErr != nil {
						return "Unable to create successful merge comment", commentErr
					}

				}
			}

			return result, err
		}

	} else {
		log.Println("Can't merge a non-PR issue")
	}

	return result, nil
}
