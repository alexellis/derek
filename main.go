package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/alexellis/derek/auth"
	"github.com/alexellis/derek/types"
)

func main() {
	bytesIn, _ := ioutil.ReadAll(os.Stdin)

	xHubSignature := os.Getenv("XHubSignature")
	if len(xHubSignature) == 0 {
		xHubSignature = os.Getenv("Http_X_Hub_Signature")
	}

	if len(xHubSignature) > 0 {
		secretKey := os.Getenv("secret_key")

		err := auth.ValidateHMAC(bytesIn, xHubSignature, secretKey)
		if err != nil {
			log.Fatal(err.Error())
			return
		}
	} else if len(os.Getenv("validate_hmac")) > 0 {
		log.Fatal("must provide X_Hub_Signature")
	}

	eventType := os.Getenv("Http_X_Github_Event")
	switch eventType {
	case "pull_request":
		req := types.PullRequestOuter{}
		if err := json.Unmarshal(bytesIn, &req); err != nil {
			log.Fatalf("Cannot parse input %s", err.Error())
		}
		handlePullRequest(req)
		break
	case "issue_comment":
		req := types.IssueCommentOuter{}
		if err := json.Unmarshal(bytesIn, &req); err != nil {
			log.Fatalf("Cannot parse input %s", err.Error())
		}
		handleComment(req)
		break
	default:
		log.Fatalln("X_Github_Event want: ['pull_request', 'issue_comment'], got: " + eventType)
	}
}
