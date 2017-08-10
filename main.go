package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/alexellis/derek/types"
)

func main() {
	bytesIn, _ := ioutil.ReadAll(os.Stdin)
	req := types.PullRequestOuter{}
	if err := json.Unmarshal(bytesIn, &req); err != nil {
		log.Fatalf("Cannot parse input %s", err.Error())
		return
	}

	handle(req)
}
