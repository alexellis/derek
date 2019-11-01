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
