package handler

import (
	"log"
	"strings"

	"github.com/alexellis/derek/config"
	"github.com/alexellis/derek/types"
	"github.com/google/go-github/github"
)

func CheckIssueTemplateHeadings(req types.IssuesOuter, derekConfig *types.DerekRepoConfig, config config.Config) error {

	maintainer := false
	for _, u := range derekConfig.Maintainers {
		if u == req.Sender.Login {
			maintainer = true
			break
		}
	}
	if maintainer {
		return nil
	}

	body := req.Issue.Body

	found := 0
	for _, heading := range derekConfig.RequiredInIssues {
		if strings.Contains(body, heading) {
			found++
		}
	}

	client, ctx := makeClient(req.Installation.ID, config)

	if found != len(derekConfig.RequiredInIssues) {
		log.Printf("Issue headings found: %d, wanted: %d", found, len(derekConfig.RequiredInIssues))
		if _, _, err := client.Issues.AddLabelsToIssue(ctx, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number, []string{"invalid"}); err != nil {
			return err
		}

		messageValue, err := createIssueComment(derekConfig.Messages, "template")
		if err != nil {
			msg := "Please complete the whole issue template, without deleting any headings."
			messageValue = &github.IssueComment{
				Body: &msg,
			}
		}
		if _, _, err = client.Issues.CreateComment(ctx, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number, messageValue); err != nil {
			return err
		}
	}

	return nil
}
