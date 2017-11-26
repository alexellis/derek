package types

type Repository struct {
	Owner Owner  `json:"owner"`
	Name  string `json:"name"`
}

type Owner struct {
	Login string `json:"login"`
	Type  string `json:"type"`
}

type PullRequest struct {
	Number int `json:"number"`
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
	Action      string      `json:"action"`
	InstallationRequest
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
	Labels []IssueLabel `json:"labels"`
	Number int          `json:"number"`
	Title  string       `json:"title"`
	Locked bool         `json:"locked"`
	State  string       `json:"state"`
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

type DerekConfig struct {

	// Features can be turned on/off if needed.
	Features []string

	// Users who are enrolled to make use of Derek
	Maintainers []string

	// Curators is an alias for Maintainers and is only used if the Maintainers list is empty.
	Curators []string
}
