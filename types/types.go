// Copyright (c) Derek Author(s) 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package types

type Repository struct {
	Owner         Owner  `json:"owner"`
	Name          string `json:"name"`
	DefaultBranch string `json:"default_branch"`
}

type Branch struct {
	Repository Repository `json:"repository"`
	Name       string     `json:"ref"`
}

type Owner struct {
	Login string `json:"login"`
	Type  string `json:"type"`
}

type PullRequest struct {
	Number            int    `json:"number"`
	AuthorAssociation string `json:"author_association"`
	Body              string `json:"body"`
	State             string `json:"state"`
}

type InstallationRequest struct {
	Installation ID `json:"installation"`
}

type ID struct {
	ID int `json:"id"`
}

type PullRequestOuter struct {
	Repository  Repository  `json:"repository"`
	PullRequest PullRequest `json:"pull_request"`
	BaseBranch  Branch      `json:"base"`
	HeadBranch  Branch      `json:"head"`
	Action      string      `json:"action"`
	InstallationRequest
}

/*
 * The default branch of a repository.
 * Usually set to `master`
 */
func (req *PullRequestOuter) GetDefaultBranch() string {
	return req.Repository.DefaultBranch
}

/*
 * The branch a pull request is open against.
 * It should be the default branch.
 */
func (req *PullRequestOuter) GetBaseBranch() string {
	return req.BaseBranch.Name
}

/*
 * The branch a pull request is open from
 */
func (req *PullRequestOuter) GetHeadBranch() string {
	return req.HeadBranch.Name
}

type IssueCommentOuter struct {
	Repository Repository `json:"repository"`
	Comment    Comment    `json:"comment"`
	Action     string     `json:"action"`
	Issue      Issue      `json:"issue"`
	InstallationRequest
}

type IssueLabel struct {
	Name string `json:"name"`
}

type Issue struct {
	Labels    []IssueLabel `json:"labels"`
	Number    int          `json:"number"`
	Title     string       `json:"title"`
	Locked    bool         `json:"locked"`
	State     string       `json:"state"`
	Milestone Milestone    `json:"milestone"`
	URL       string       `json:"url"`
}

type Milestone struct {
	Title string `json:"title"`
}

type Comment struct {
	Body     string `json:"body"`
	IssueURL string `json:"issue_url"`
	User     struct {
		Login string `json:"login"`
	}
}

type CommentAction struct {
	Type  string
	Value string
}

// DerekRepoConfig is a config for a Derek-enabled repository
type DerekRepoConfig struct {

	// A redirect URL to load the config from another location.
	Redirect string

	// Features can be turned on/off if needed.
	Features []string

	// Users who are enrolled to make use of Derek
	Maintainers []string

	// Curators is an alias for Maintainers and is only used if the Maintainers list is empty.
	Curators []string

	//ContributingURL url to contribution guide
	ContributingURL string `yaml:"contributing_url"`
}

// FirstTimeContributor whether the contributor is new to the repo
func (p *PullRequest) FirstTimeContributor() bool {
	return p.AuthorAssociation == "NONE"
}
