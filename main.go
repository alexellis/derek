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
		xHubSignature = os.Getenv("Http_X-Hub-Signature")
	}

	if len(xHubSignature) > 0 {
		secretKey := os.Getenv("secret_key")

		err := auth.ValidateHMAC(bytesIn, xHubSignature, secretKey)
		if err != nil {
			log.Fatal(err.Error())
			return
		}
	} else if len(os.Getenv("validate_hmac")) > 0 {
		log.Fatal("must provide X-Hub-Signature")
	}

	req := types.PullRequestOuter{}
	if err := json.Unmarshal(bytesIn, &req); err != nil {
		log.Fatalf("Cannot parse input %s", err.Error())
	}

	if os.Getenv("Http_X-Github-Event") != "pull_request" {
		log.Fatalln("X-Github-Event want: 'pull_request', got: " + os.Getenv("Http_X-Github-Event"))
	}

	handle(req)
}
