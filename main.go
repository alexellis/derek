// Copyright (c) Derek Author(s) 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/alexellis/derek/auth"
	"github.com/alexellis/hmac"

	"github.com/alexellis/derek/types"
)

const (
	dcoCheck = "dco_check"
	comments = "comments"
	deleted  = "deleted"
)

const derekSecretKeyFile = "derek-secret-key"
const privateKeyFile = "derek-private-key"

func getSecretPath() (string, error) {
	secretPath := os.Getenv("secret_path")

	if len(secretPath) == 0 {
		return "", fmt.Errorf("secret_path not set, this must be /var/openfaas/secrets or /run/secrets")

	}

	return secretPath, nil
}

func hmacValidation() bool {
	val := os.Getenv("validate_hmac")
	return len(val) > 0 && (val == "1" || val == "true")
}

func removeNewLine(secret []byte) []byte {
	stringSecret := string(secret)
	if newLine := strings.Index(stringSecret, "\n"); newLine != -1 {
		secret = secret[:newLine]
	}
	return secret
}

func main() {

	bytesIn, _ := ioutil.ReadAll(os.Stdin)

	xHubSignature := os.Getenv("Http_X_Hub_Signature")

	if hmacValidation() && len(xHubSignature) == 0 {
		log.Fatal("must provide X_Hub_Signature")
		return
	}

	keyPath, pathErr := getSecretPath()
	if pathErr != nil {
		os.Stderr.Write([]byte(pathErr.Error()))
		os.Exit(1)
	}

	if len(xHubSignature) > 0 {
		secretKeyBytes, readErr := ioutil.ReadFile(keyPath + derekSecretKeyFile)

		if readErr != nil {
			msg := fmt.Errorf("unable to read GitHub symmetrical secret: %s, error: %s",
				keyPath+derekSecretKeyFile, readErr)
			os.Stderr.Write([]byte(msg.Error()))
			os.Exit(1)
		}

		secretKeyBytes = removeNewLine(secretKeyBytes)

		err := hmac.Validate(bytesIn, xHubSignature, string(secretKeyBytes))
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
			return fmt.Errorf("Cannot parse input %s", err.Error())
		}

		customer, err := auth.IsCustomer(req.Repository.Owner.Login, &http.Client{})
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
			if enabledFeature(dcoCheck, derekConfig) {
				handlePullRequest(req)
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

		derekConfig, err := getConfig(req.Repository.Owner.Login, req.Repository.Name)
		if err != nil {
			return fmt.Errorf("Unable to access maintainers file at: %s/%s", req.Repository.Owner.Login, req.Repository.Name)
		}

		if req.Action != deleted {
			if permittedUserFeature(comments, derekConfig, req.Comment.User.Login) {
				handleComment(req)
			}
		}
		break
	default:
		return fmt.Errorf("X_Github_Event want: ['pull_request', 'issue_comment'], got: " + eventType)
	}

	return nil
}
