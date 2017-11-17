package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/alexellis/derek/auth"
	"github.com/alexellis/derek/types"
)

const dcoCheck = "dco_check"
const comments = "comments"

func hmacValidation() bool {
	val := os.Getenv("validate_hmac")
	return len(val) > 0 && (val == "1" || val == "true")
}

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

	switch eventType {
	case "pull_request":
		req := types.PullRequestOuter{}
		if err := json.Unmarshal(bytesIn, &req); err != nil {
			log.Fatalf("Cannot parse input %s", err.Error())
		}

		customer, err := auth.IsCustomer(req.Repository)
		if err != nil {
			log.Fatalf("Unable to verify customer: %s/%s", req.Repository.Owner.Login, req.Repository.Name)
		} else if !customer {
			log.Fatalf("No customer found for: %s/%s", req.Repository.Owner.Login, req.Repository.Name)
		}

		derekConfig, err := getConfig(req.Repository.Owner.Login, req.Repository.Name)
		if err != nil {
			log.Fatalf("Unable to access maintainers file at: %s/%s", req.Repository.Owner.Login, req.Repository.Name)
		}

		if enabledFeature(dcoCheck, derekConfig) {
			handlePullRequest(req)
		}
		break

	case "issue_comment":
		req := types.IssueCommentOuter{}
		if err := json.Unmarshal(bytesIn, &req); err != nil {
			log.Fatalf("Cannot parse input %s", err.Error())
		}

		customer, err := auth.IsCustomer(req.Repository)
		if err != nil {
			log.Fatalf("Unable to verify customer: %s/%s", req.Repository.Owner.Login, req.Repository.Name)
		} else if !customer {
			log.Fatalf("No customer found for: %s/%s", req.Repository.Owner.Login, req.Repository.Name)
		}

		derekConfig, err := getConfig(req.Repository.Owner.Login, req.Repository.Name)
		if err != nil {
			log.Fatalf("Unable to access maintainers file at: %s/%s", req.Repository.Owner.Login, req.Repository.Name)
		}

		if permittedUserFeature(comments, derekConfig, req.Comment.User.Login) {
			handleComment(req)
		}
		break
	default:
		log.Fatalln("X_Github_Event want: ['pull_request', 'issue_comment'], got: " + eventType)
	}
}
