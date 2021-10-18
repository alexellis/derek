package function

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/google/go-github/github"

	"github.com/alexellis/derek/auth"
	"github.com/alexellis/derek/config"

	"github.com/alexellis/derek/handler"

	"github.com/alexellis/derek/types"
	"github.com/alexellis/hmac"

	handler "github.com/openfaas/templates-sdk/go-http"
)

// Example handle
// Handle a function invocation
// func Handle(req handler.Request) (handler.Response, error) {
// 	var err error

// 	message := fmt.Sprintf("Body: %s", string(req.Body))

// 	return handler.Response{
// 		Body:       []byte(message),
// 		StatusCode: http.StatusOK,
// 	}, err
// }

const (
	dcoCheck              = "dco_check"
	comments              = "comments"
	deleted               = "deleted"
	prDescriptionRequired = "pr_description_required"
	hacktoberfest         = "hacktoberfest"
	noNewbies             = "no_newbies"
	releaseNotes          = "release_notes"
)

func Handle(req handler.Request) (handler.Response, error) {
	validateHmac := hmacValidation()

	requestRaw, _ := ioutil.ReadAll(req.Body)

	xHubSignature := r.Header.Get("X-Hub-Signature")

	if validateHmac && len(xHubSignature) == 0 {
		return handler.Response{
			Body:       []byte("must provide X-Hub-Signature"),
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	config, configErr := config.NewConfig()
	if configErr != nil {
		return handler.Response{
			Body:       []byte(configErr.Error()),
			StatusCode: http.StatusInternalServerError,
		}, configErr
	}

	if validateHmac {
		err := hmac.Validate(requestRaw, xHubSignature, config.SecretKey)
		if err != nil {
			return handler.Response{
				Body:       []byte(err.Error()),
				StatusCode: http.StatusInternalServerError,
			}, err
		}
	}

	eventType := req.Header.Get("X-Github-Event")

	if err := handleEvent(eventType, requestRaw, config); err != nil {
		return handler.Response{
			Body:       []byte(err.Error()),
			StatusCode: http.StatusInternalServerError,
		}, err
	}
}

func handleEvent(eventType string, bytesIn []byte, config config.Config) error {

	switch eventType {
	case "pull_request":
		req := types.PullRequestOuter{}
		if err := json.Unmarshal(bytesIn, &req); err != nil {
			return fmt.Errorf("Cannot parse input %s", err.Error())
		}

		customer, err := auth.IsCustomer(req.Repository.Owner.Login, &http.Client{})
		if err != nil {
			return fmt.Errorf("Unable to verify customer: %s/%s", req.Repository.Owner.Login, req.Repository.Name)
		} else if !customer {
			return fmt.Errorf("No customer found for: %s/%s", req.Repository.Owner.Login, req.Repository.Name)
		}

		log.Printf("Owner: %s, repo: %s, action: %s", req.Repository.Owner.Login, req.Repository.Name, "pull_request")

		var derekConfig *types.DerekRepoConfig
		if req.Repository.Private {
			derekConfig, err = handler.GetPrivateRepoConfig(req.Repository.Owner.Login, req.Repository.Name, req.Repository.DefaultBranch, req.Installation.ID, config)
		} else {
			derekConfig, err = handler.GetRepoConfig(req.Repository.Owner.Login, req.Repository.Name, req.Repository.DefaultBranch)
		}

		if err != nil {
			return fmt.Errorf("Unable to access maintainers file at: %s/%s\nError: %s",
				req.Repository.Owner.Login,
				req.Repository.Name,
				err.Error())
		}

		if req.Action != handler.ClosedConstant && req.PullRequest.State != handler.ClosedConstant {
			contributingURL := getContributingURL(derekConfig.ContributingURL, req.Repository.Owner.Login, req.Repository.Name)

			if handler.EnabledFeature(dcoCheck, derekConfig) {
				log.Printf("Owner: %s, repo: %s, action: %s", req.Repository.Owner.Login, req.Repository.Name, "derek:dco_check")

				handler.HandlePullRequest(req, contributingURL, config)
			}

			if handler.EnabledFeature(prDescriptionRequired, derekConfig) {
				handler.VerifyPullRequestDescription(req, contributingURL, config)
			}

			if handler.EnabledFeature(noNewbies, derekConfig) {
				isSpamPR, _ := handler.HandleFirstTimerPR(req, contributingURL, config)
				if isSpamPR {
					return nil
				}
			}

			if handler.EnabledFeature(hacktoberfest, derekConfig) {
				isSpamPR, _ := handler.HandleHacktoberfestPR(req, contributingURL, config)
				if isSpamPR {
					return nil
				}
			}
		}
		break

	case "issues":

		req := types.IssuesOuter{}
		if err := json.Unmarshal(bytesIn, &req); err != nil {
			return fmt.Errorf("Cannot parse input %s", err.Error())
		}

		if req.Action == "opened" {
			log.Printf("Owner: %s, repo: %s, action: %s", req.Repository.Owner.Login, req.Repository.Name, "issues")

			customer, err := auth.IsCustomer(req.Repository.Owner.Login, &http.Client{})
			if err != nil {
				return fmt.Errorf("Unable to verify customer: %s/%s", req.Repository.Owner.Login, req.Repository.Name)
			} else if !customer {
				return fmt.Errorf("No customer found for: %s/%s", req.Repository.Owner.Login, req.Repository.Name)
			}

			var derekConfig *types.DerekRepoConfig
			if req.Repository.Private {
				derekConfig, err = handler.GetPrivateRepoConfig(req.Repository.Owner.Login, req.Repository.Name, req.Repository.DefaultBranch, req.Installation.ID, config)
			} else {
				derekConfig, err = handler.GetRepoConfig(req.Repository.Owner.Login, req.Repository.Name, req.Repository.DefaultBranch)
			}
			if err != nil {
				return fmt.Errorf("Unable to access maintainers file at: %s/%s\nError: %s",
					req.Repository.Owner.Login,
					req.Repository.Name,
					err.Error())
			}

			if len(derekConfig.RequiredInIssues) > 0 {
				err := handler.CheckIssueTemplateHeadings(req, derekConfig, config)
				if err != nil {
					return err
				}
			}
		}

	case "issue_comment":
		req := types.IssueCommentOuter{}
		if err := json.Unmarshal(bytesIn, &req); err != nil {
			return fmt.Errorf("Cannot parse input %s", err.Error())
		}

		log.Printf("Owner: %s, repo: %s, action: %s", req.Repository.Owner.Login, req.Repository.Name, "issue_comment")

		customer, err := auth.IsCustomer(req.Repository.Owner.Login, &http.Client{})
		if err != nil {
			return fmt.Errorf("Unable to verify customer: %s/%s", req.Repository.Owner.Login, req.Repository.Name)
		} else if !customer {
			return fmt.Errorf("No customer found for: %s/%s", req.Repository.Owner.Login, req.Repository.Name)
		}

		var derekConfig *types.DerekRepoConfig
		if req.Repository.Private {
			derekConfig, err = handler.GetPrivateRepoConfig(req.Repository.Owner.Login, req.Repository.Name, req.Repository.DefaultBranch, req.Installation.ID, config)
		} else {
			derekConfig, err = handler.GetRepoConfig(req.Repository.Owner.Login, req.Repository.Name, req.Repository.DefaultBranch)
		}

		if err != nil {
			return fmt.Errorf("Unable to access maintainers file at: %s/%s\nError: %s",
				req.Repository.Owner.Login,
				req.Repository.Name,
				err.Error())
		}

		if req.Action != deleted {
			if handler.PermittedUserFeature(comments, derekConfig, req.Comment.User.Login) {
				log.Printf("Owner: %s, repo: %s, action: %s", req.Repository.Owner.Login, req.Repository.Name, "derek:handle_comment")

				handler.HandleComment(req, config, derekConfig)
			}
		}

	case "release":
		req := github.ReleaseEvent{}

		if err := json.Unmarshal(bytesIn, &req); err != nil {
			return fmt.Errorf("Cannot parse input %s", err.Error())
		}

		log.Printf("Owner: %s, repo: %s, action: %s", req.Repo.Owner.GetLogin(), req.Repo.GetName(), "release")

		if req.GetAction() == "created" {
			customer, err := auth.IsCustomer(req.Repo.Owner.GetLogin(), &http.Client{})
			if err != nil {
				return fmt.Errorf("unable to verify customer: %s/%s", req.Repo.Owner.GetLogin(), req.Repo.GetName())
			} else if customer == false {
				return fmt.Errorf("no customer found for: %s/%s", req.Repo.Owner.GetLogin(), req.Repo.GetName())
			}

			var derekConfig *types.DerekRepoConfig
			if req.Repo.GetPrivate() {
				derekConfig, err = handler.GetPrivateRepoConfig(req.Repo.Owner.GetLogin(), req.Repo.GetName(), req.Repo.GetDefaultBranch(), int(req.Installation.GetID()), config)
				if err != nil {
					return fmt.Errorf("unable to get private repo config: %s", err)
				}
			} else {
				derekConfig, err = handler.GetRepoConfig(req.Repo.Owner.GetLogin(), req.Repo.GetName(), req.Repo.GetDefaultBranch())
				if err != nil {
					return fmt.Errorf("unable to get repo config: %s", err)
				}
			}

			err = fmt.Errorf(`"release_notes" feature not enabled`)
			if handler.EnabledFeature(releaseNotes, derekConfig) {
				log.Printf("Owner: %s, repo: %s, action: %s", req.Repo.Owner.GetLogin(), req.Repo.GetName(), "derek:handle_release")

				handler := handler.NewReleaseHandler(config, int(req.Installation.GetID()))
				err = handler.Handle(req)
			}
			return err
		}

	default:
		return fmt.Errorf("X_Github_Event want: ['pull_request', 'issue_comment', 'release'], got: " + eventType)
	}

	return nil
}

func getContributingURL(contributingURL, owner, repositoryName string) string {
	if len(contributingURL) == 0 {
		contributingURL = fmt.Sprintf("https://github.com/%s/%s/blob/master/CONTRIBUTING.md", owner, repositoryName)
	}
	return contributingURL
}

func hmacValidation() bool {
	val := os.Getenv("validate_hmac")
	return (val != "false") && (val != "0")
}
