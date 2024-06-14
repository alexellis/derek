package handler

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/alexellis/derek/config"
	"github.com/google/go-github/github"
)

type WorkingRelease struct {
	CurrentTag     string
	CurrentDate    time.Time
	CurrentRelease *github.RepositoryRelease

	PreviousTag  string
	PreviousDate time.Time
}

type ReleaseHandler interface {
	Handle(github.ReleaseEvent) error
}

type UpdatingReleaseHandler struct {
	Config         config.Config
	InstallationID int
}

func NewReleaseHandler(config config.Config, installationID int) ReleaseHandler {
	return &UpdatingReleaseHandler{
		InstallationID: installationID,
		Config:         config,
	}
}

func (h *UpdatingReleaseHandler) Handle(req github.ReleaseEvent) error {
	client, _ := makeClient(h.InstallationID, h.Config)

	err := updateReleaseNotes(client, req.Repo.Owner.GetLogin(), req.Repo.GetName(), req.Release.GetTagName())

	return err
}

func updateReleaseNotes(client *github.Client, owner, repo, latestTag string) error {

	listOptions := &github.ListOptions{}
	releases, _, err := client.Repositories.ListReleases(context.Background(), owner, repo, listOptions)
	if err != nil {
		return err
	}

	workingReleases := getWorkingReleases(releases, owner, repo, latestTag)

	if workingReleases.CurrentRelease == nil {
		return fmt.Errorf("unable to detect current release, retry webhook after a few seconds")

	}

	includedPRs, err := buildClosedPRs(client, workingReleases, owner, repo, latestTag)
	if err != nil {
		return err
	}

	includedCommits, err := buildCommits(client, workingReleases, owner, repo, latestTag)
	if err != nil {
		return err
	}

	output := fmt.Sprintf("Changelog for %s:\n", workingReleases.CurrentTag)

	for _, pr := range includedPRs {
		output = output + fmt.Sprintf("* PR #%d %s by @%s\n",
			pr.GetNumber(),
			pr.GetTitle(),
			pr.GetUser().GetLogin())
	}

	output += "\nCommits\n"

	for _, c := range includedCommits {
		id := c.GetSHA()

		author := c.GetAuthor().GetLogin()

		if len(author) > 0 {
			author = " by @" + author
		} else {
			author = " by " + c.GetCommit().GetAuthor().GetName()
		}

		title := c.GetCommit().GetMessage()
		if index := strings.Index(title, "\n"); index > -1 {
			title = title[:index]
		}

		output = output + fmt.Sprintf("%s %s%s\n", id, title, author)
	}

	var releaseDiff string
	if workingReleases.PreviousTag != "" {
		releaseDiff = fmt.Sprintf("Changes: https://github.com/%s/%s/compare/%s...%s", owner, repo, workingReleases.PreviousTag, workingReleases.CurrentTag)
	} else {
		releaseDiff = fmt.Sprintf("Changes: https://github.com/%s/%s/commits/%s", owner, repo, workingReleases.CurrentTag)
	}
	output = fmt.Sprintf("%s\n%s\n\nGenerated by [Derek](https://github.com/alexellis/derek/)\n", output, releaseDiff)

	log.Printf("Release notes: %q", output)

	err = updateRelease(client, workingReleases.CurrentRelease, owner, repo, workingReleases.CurrentTag, output)

	return err
}

func buildCommits(client *github.Client, workingReleases WorkingRelease, owner, repo, latestTag string) ([]github.RepositoryCommit, error) {
	var err error
	var commits []github.RepositoryCommit

	opts := github.CommitsListOptions{
		Since: workingReleases.PreviousDate,
		Until: workingReleases.CurrentDate,
	}

	res, _, err := client.Repositories.ListCommits(context.Background(), owner, repo, &opts)

	for _, c := range res {
		if includeCommit(*c, workingReleases.PreviousDate, workingReleases.CurrentDate) {
			commits = append(commits, *c)
		}
	}

	if err != nil {
		return nil, err
	}

	return commits, err
}

func buildClosedPRs(client *github.Client, workingReleases WorkingRelease, owner, repo, latestTag string) ([]github.PullRequest, error) {
	opts := &github.PullRequestListOptions{
		State:     "closed",
		Base:      "master",
		Sort:      "updated",
		Direction: "desc",
	}

	prs, _, err := client.PullRequests.List(context.Background(), owner, repo, opts)
	if err != nil {
		return nil, err
	}

	log.Printf("Release [%s/%s:%s] start: %s\tend: %q\n",
		owner, repo, latestTag, workingReleases.CurrentDate.String(), workingReleases.PreviousDate.String())

	included := []github.PullRequest{}
	for _, pr := range prs {
		if includePR(*pr, workingReleases.PreviousDate, workingReleases.CurrentDate) {
			included = append(included, *pr)
		}
	}

	sort.SliceStable(included, func(i, j int) bool {
		return included[i].GetClosedAt().After(included[j].GetClosedAt())
	})

	return included, nil
}

func getWorkingReleases(releases []*github.RepositoryRelease, owner, repo, tag string) WorkingRelease {
	rel := WorkingRelease{}

	for position, r := range releases {
		if r.GetTagName() == tag {
			rel.CurrentDate = r.CreatedAt.Time
			rel.CurrentTag = tag
			rel.CurrentRelease = r

			prevRel := position + 1

			for prevRel < len(releases) {

				prior := releases[prevRel]

				if !*prior.Prerelease {
					rel.PreviousDate = prior.CreatedAt.Time
					rel.PreviousTag = prior.GetTagName()
					break
				}

				prevRel++

			}

			break
		}
	}

	return rel
}

func includePR(pr github.PullRequest, previous, current time.Time) bool {
	window := time.Minute * 1 // GitHub is an async system and events can end up with times that are out by a few seconds.

	// pr.Merged returns false all the time, therefore use GetMergedAt()
	merged := pr.GetMergedAt().Equal(time.Time{}) == false

	return pr.GetClosedAt().After(previous.Add(window)) &&
		pr.GetClosedAt().Before(current.Add(window)) &&
		merged
}

func updateRelease(client *github.Client, release *github.RepositoryRelease, owner, repo, tag, body string) error {
	release.Body = &body

	_, _, err := client.Repositories.EditRelease(context.Background(), owner, repo, *release.ID, release)
	return err
}

func includeCommit(cm github.RepositoryCommit, previous, current time.Time) bool {
	window := time.Minute * 1 // GitHub is an async system and events can end up with times that are out by a few seconds.

	committed := cm.GetCommit().GetCommitter().GetDate().Equal(time.Time{}) == false

	return cm.GetCommit().GetCommitter().GetDate().After(previous.Add(window)) &&
		cm.GetCommit().GetCommitter().GetDate().Before(current.Add(window)) &&
		committed
}
