// Copyright (c) Derek Author(s) 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/alexellis/derek/auth"
	"github.com/alexellis/derek/config"

	"github.com/alexellis/derek/handler"

	"github.com/alexellis/derek/types"
	"github.com/alexellis/hmac"
)

const (
	dcoCheck              = "dco_check"
	comments              = "comments"
	deleted               = "deleted"
	prDescriptionRequired = "pr_description_required"
	hacktoberfest         = "hacktoberfest"
)

func main() {
	validateHmac := hmacValidation()

	requestRaw, _ := ioutil.ReadAll(os.Stdin)

	xHubSignature := os.Getenv("Http_X_Hub_Signature")

	if validateHmac && len(xHubSignature) == 0 {
		os.Stderr.Write([]byte("must provide X_Hub_Signature"))
		os.Exit(1)
	}

	config, configErr := config.NewConfig()
	if configErr != nil {
		os.Stderr.Write([]byte(configErr.Error()))
		os.Exit(1)
	}

	if validateHmac {
		err := hmac.Validate(requestRaw, xHubSignature, config.SecretKey)
		if err != nil {
			os.Stderr.Write([]byte(err.Error()))
			os.Exit(1)
		}
	}

	eventType := os.Getenv("Http_X_Github_Event")

	if err := handleEvent(eventType, requestRaw, config); err != nil {
		os.Stderr.Write([]byte(err.Error()))
		os.Exit(1)
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
		} else if customer == false {
			return fmt.Errorf("No customer found for: %s/%s", req.Repository.Owner.Login, req.Repository.Name)
		}

		var derekConfig *types.DerekRepoConfig
		if req.Repository.Private {
			derekConfig, err = handler.GetPrivateRepoConfig(req.Repository.Owner.Login, req.Repository.Name, req.Installation.ID, config)
		} else {
			derekConfig, err = handler.GetRepoConfig(req.Repository.Owner.Login, req.Repository.Name)
		}
		if err != nil {
			return fmt.Errorf("Unable to access maintainers file at: %s/%s\nError: %s",
				req.Repository.Owner.Login,
				req.Repository.Name,
				err.Error())
		}
		if req.Action != handler.ClosedConstant {
			contributingURL := getContributingURL(derekConfig.ContributingURL, req.Repository.Owner.Login, req.Repository.Name)
			if handler.EnabledFeature(hacktoberfest, derekConfig) {
				isSpamPR, _ := handler.HandleHacktoberfestPR(req, contributingURL, config)
				if isSpamPR {
					return nil
				}
			}
			if handler.EnabledFeature(dcoCheck, derekConfig) {
				handler.HandlePullRequest(req, contributingURL, config)
			}
			if handler.EnabledFeature(prDescriptionRequired, derekConfig) {
				handler.VerifyPullRequestDescription(req, contributingURL, config)
			}
		}
		break

	case "issue_comment":
		req := types.IssueCommentOuter{}
		if err := json.Unmarshal(bytesIn, &req); err != nil {
			return fmt.Errorf("Cannot parse input %s", err.Error())
		}

		customer, err := auth.IsCustomer(req.Repository.Owner.Login, &http.Client{})
		if err != nil {
			return fmt.Errorf("Unable to verify customer: %s/%s", req.Repository.Owner.Login, req.Repository.Name)
		} else if customer == false {
			return fmt.Errorf("No customer found for: %s/%s", req.Repository.Owner.Login, req.Repository.Name)
		}

		var derekConfig *types.DerekRepoConfig
		if req.Repository.Private {
			derekConfig, err = handler.GetPrivateRepoConfig(req.Repository.Owner.Login, req.Repository.Name, req.Installation.ID, config)
		} else {
			derekConfig, err = handler.GetRepoConfig(req.Repository.Owner.Login, req.Repository.Name)
		}
		if err != nil {
			return fmt.Errorf("Unable to access maintainers file at: %s/%s\nError: %s",
				req.Repository.Owner.Login,
				req.Repository.Name,
				err.Error())
		}

		if req.Action != deleted {
			if handler.PermittedUserFeature(comments, derekConfig, req.Comment.User.Login) {
				handler.HandleComment(req, config, derekConfig)
			}
		}
		break
	default:
		return fmt.Errorf("X_Github_Event want: ['pull_request', 'issue_comment'], got: " + eventType)
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
