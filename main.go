package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/alexellis/derek/auth"
	"github.com/alexellis/derek/types"
)

const deleted = "deleted"

func main() {
	bytesIn, _ := ioutil.ReadAll(os.Stdin)

	xHubSignature := os.Getenv("Http_X_Hub_Signature")
	if hmacValidation() && len(xHubSignature) == 0 {
		log.Fatal("must provide X_Hub_Signature")
		return
	}

	if len(xHubSignature) > 0 {
		secretKey := os.Getenv("secret_key")

		err := auth.ValidateHMAC(bytesIn, xHubSignature, secretKey)
		if err != nil {
			log.Fatal(err.Error())
			return
		}
	}

	// HMAC Validated or not turned on.
	eventType := os.Getenv("Http_X_Github_Event")

	if err := handleEvent(eventType, bytesIn); err != nil {
		log.Fatal(err)
	}
}

func handleEvent(eventType string, bytesIn []byte) error {

	switch eventType {
	case "pull_request":
		req := types.PullRequestOuter{}
		if err := json.Unmarshal(bytesIn, &req); err != nil {
			return fmt.Errorf("pull_request handler, cannot parse input: %s", err.Error())
		}

		customer, err := auth.IsCustomer(req.Repository)
		if err != nil {
			return fmt.Errorf("Unable to verify customer: %s/%s", req.Repository.Owner.Login, req.Repository.Name)
		} else if customer == false {
			return fmt.Errorf("No customer found for: %s/%s", req.Repository.Owner.Login, req.Repository.Name)
		}

		derekConfig, err := getConfig(req.Repository.Owner.Login, req.Repository.Name)
		if err != nil {
			return fmt.Errorf("Unable to access maintainers file at: %s/%s", req.Repository.Owner.Login, req.Repository.Name)
		}

		if req.Action != closedConstant {
			prFeatures := types.PullRequestFeatures{
				CommitLintingFeature: enabledFeature(types.CommitLintingFeature, derekConfig),
				DCOCheckFeature:      enabledFeature(types.DCOCheckFeature, derekConfig),
			}

			if prFeatures.Enabled() {
				handlePullRequest(req, prFeatures)
			}
		}

		break

	case "issue_comment":
		req := types.IssueCommentOuter{}
		if err := json.Unmarshal(bytesIn, &req); err != nil {
			return fmt.Errorf("Cannot parse input %s", err.Error())
		}

		customer, err := auth.IsCustomer(req.Repository)
		if err != nil {
			return fmt.Errorf("Unable to verify customer: %s/%s", req.Repository.Owner.Login, req.Repository.Name)
		} else if customer == false {
			return fmt.Errorf("No customer found for: %s/%s", req.Repository.Owner.Login, req.Repository.Name)
		}

		derekConfig, err := getConfig(req.Repository.Owner.Login, req.Repository.Name)
		if err != nil {
			return fmt.Errorf("Unable to access maintainers file at: %s/%s", req.Repository.Owner.Login, req.Repository.Name)
		}

		if req.Action != deleted {
			if permittedUserFeature(types.CommentFeature, derekConfig, req.Comment.User.Login) {
				handleComment(req)
			}
		}

		break

	default:
		return fmt.Errorf("X_Github_Event want: ['pull_request', 'issue_comment'], got: " + eventType)
	}

	return nil
}

func hmacValidation() bool {
	val := os.Getenv("validate_hmac")
	return len(val) > 0 && (val == "1" || val == "true")
}
