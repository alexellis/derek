package handler

import (
	"testing"
	"time"

	"github.com/google/go-github/github"
)

func Test_includePR_AfterCurrentPeriod(t *testing.T) {
	now := time.Now()

	prDate := now.Add(time.Hour * 48)
	pr := github.PullRequest{
		ClosedAt: &prDate,
		MergedAt: &prDate,
	}

	previous := now.Add(time.Hour * -24)
	current := now.Add(time.Hour * 24)

	got := includePR(pr, previous, current)
	want := false

	if got != want {
		t.Errorf("Included value for PR %s incorrect for range: [%s-%s] got: %v, want %v",
			pr.ClosedAt.String(), previous.String(), current.String(), got, want)
		t.Fail()
	}
}

func Test_includePR_WithinCurrentRange(t *testing.T) {
	now := time.Now()

	prDate := now.Add(time.Hour)
	pr := github.PullRequest{
		ClosedAt: &prDate,
		MergedAt: &prDate,
	}

	previous := now.Add(time.Hour * -24)
	current := now.Add(time.Hour * 24)

	got := includePR(pr, previous, current)
	want := true

	if got != want {
		t.Errorf("Included value for PR %s incorrect for range: [%s-%s] got: %v, want %v",
			pr.ClosedAt.String(), previous.String(), current.String(), got, want)
		t.Fail()
	}
}

func Test_includePR_WithinCurrentRange_ButNotMerged(t *testing.T) {
	now := time.Now()

	prDate := now.Add(time.Hour)
	pr := github.PullRequest{
		ClosedAt: &prDate,
	}

	previous := now.Add(time.Hour * -24)
	current := now.Add(time.Hour * 24)

	got := includePR(pr, previous, current)
	want := false

	if got != want {
		t.Errorf("Included value for PR %s incorrect for range: [%s-%s] got: %v, want %v",
			pr.ClosedAt.String(), previous.String(), current.String(), got, want)
		t.Fail()
	}
}

func Test_getWorkingReleases_TwoReleases(t *testing.T) {
	owner := "alexellis"
	repo := "derek"
	tag := "0.2.0"
	lastTag := "0.1.0"

	releases := []*github.RepositoryRelease{
		&github.RepositoryRelease{
			TagName: &tag,
			CreatedAt: &github.Timestamp{
				Time: time.Now(),
			},
		},
		&github.RepositoryRelease{
			TagName: &lastTag,
			CreatedAt: &github.Timestamp{
				Time: time.Now().Add(time.Hour * -1),
			},
		},
	}
	workingReleases := getWorkingReleases(releases, owner, repo, tag)

	gotCurrentDate := workingReleases.CurrentDate
	wantCurrentDate := releases[0].GetCreatedAt().Time

	if gotCurrentDate != wantCurrentDate {
		t.Errorf("current date, got: %s, want: %s", gotCurrentDate, wantCurrentDate)
	}

	gotPreviousDate := workingReleases.PreviousDate
	wantPreviousDate := releases[1].GetCreatedAt().Time

	if gotPreviousDate != wantPreviousDate {
		t.Errorf("previous date, got: %s, want: %s", gotPreviousDate, wantPreviousDate)
	}

	gotPreviousTag := workingReleases.PreviousTag
	wantPreviousTag := *releases[1].TagName

	if gotPreviousTag != wantPreviousTag {
		t.Errorf("previous tag, got: %s, want: %s", gotPreviousTag, wantPreviousTag)
	}
}

func Test_getWorkingReleases_OneRelease(t *testing.T) {
	owner := "alexellis"
	repo := "derek"
	tag := "0.2.0"

	releases := []*github.RepositoryRelease{
		&github.RepositoryRelease{
			TagName: &tag,
			CreatedAt: &github.Timestamp{
				Time: time.Now(),
			},
		},
	}
	workingReleases := getWorkingReleases(releases, owner, repo, tag)

	gotCurrentDate := workingReleases.CurrentDate
	wantCurrentDate := releases[0].GetCreatedAt().Time

	if gotCurrentDate != wantCurrentDate {
		t.Errorf("current date, got: %s, want: %s", gotCurrentDate, wantCurrentDate)
	}

	gotPreviousDate := workingReleases.PreviousDate
	wantPreviousDate := time.Time{}

	if gotPreviousDate != wantPreviousDate {
		t.Errorf("previous date, got: %s, want: %s", gotPreviousDate, wantPreviousDate)
	}

	gotPreviousTag := workingReleases.PreviousTag
	wantPreviousTag := ""

	if gotPreviousTag != wantPreviousTag {
		t.Errorf("previous tag, got: %s, want: %s", gotPreviousTag, wantPreviousTag)
	}

}

func Test_includeCommit_AfterCurrentPeriod(t *testing.T) {
	now := time.Now()

	cmDate := now.Add(time.Hour * 48)
	cm := github.RepositoryCommit{
		Commit: &github.Commit{
			Committer: &github.CommitAuthor{
				Date: &cmDate,
			},
		},
	}

	previous := now.Add(time.Hour * -24)
	current := now.Add(time.Hour * 24)

	got := includeCommits(cm, previous, current)
	want := false

	if got != want {
		t.Errorf("Included value for Commot %s incorrect for range: [%s-%s] got: %v, want %v",
			cm.GetCommit().GetCommitter().GetDate().String(), previous.String(), current.String(), got, want)
		t.Fail()
	}
}

func Test_includeCommit_WithinCurrentRange(t *testing.T) {
	now := time.Now()

	cmDate := now.Add(time.Hour)
	cm := github.RepositoryCommit{
		Commit: &github.Commit{
			Committer: &github.CommitAuthor{
				Date: &cmDate,
			},
		},
	}

	previous := now.Add(time.Hour * -24)
	current := now.Add(time.Hour * 24)

	got := includeCommits(cm, previous, current)
	want := true
	t.Logf("Running Test Case")
	if got != want {
		t.Errorf("Included value for Commit %s incorrect for range: [%s-%s] got: %v, want %v",
			cm.GetCommit().GetCommitter().GetDate().String(), previous.String(), current.String(), got, want)
		t.Fail()
	}
}
