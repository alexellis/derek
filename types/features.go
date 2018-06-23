package types

const DCOCheckFeature = "dco_check"
const CommitLintingFeature = "commit_linting"
const CommentFeature = "comments"

// PullRequestFeatures
type PullRequestFeatures struct {
	DCOCheckFeature      bool
	CommitLintingFeature bool
}

func (p *PullRequestFeatures) Enabled() bool {
	return p.DCOCheckFeature || p.CommitLintingFeature
}
