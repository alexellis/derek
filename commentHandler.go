package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/alexellis/derek/auth"
	"github.com/alexellis/derek/types"
	"github.com/google/go-github/github"
)

func makeClient() (*github.Client, context.Context) {
	ctx := context.Background()

	token := os.Getenv("access_token")
	if len(token) == 0 {
		newToken, tokenErr := auth.MakeAccessTokenForInstallation(os.Getenv("application"), os.Getenv("installation"), os.Getenv("private_key"))
		if tokenErr != nil {
			log.Fatalln(tokenErr.Error())
		}

		token = newToken
	}

	client := auth.MakeClient(ctx, token)

	return client, ctx
}

func handleComment(req types.IssueCommentOuter) {

	command := parse(req.Comment.Body)
	switch command.Type {
	case "AddLabel":
		allowed := isMaintainer(req.Comment.User.Login, req.Repository)
		fmt.Printf("%s wants to %s of %s to issue %d - allowed? %t\n", req.Comment.User.Login, command.Type, command.Value, req.Issue.Number, allowed)

		found := false
		for _, label := range req.Issue.Labels {
			if label.Name == command.Value {
				found = true
				break
			}
		}

		if found == true {
			fmt.Println("Label already exists.")
			return
		}

		if allowed {
			client, ctx := makeClient()
			_, _, err := client.Issues.AddLabelsToIssue(ctx, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number, []string{command.Value})
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println("Label added successfully or already existed.")
		}
		break

	case "RemoveLabel":
		allowed := isMaintainer(req.Comment.User.Login, req.Repository)
		fmt.Printf("%s wants to %s of %s to issue %d - allowed? %t\n", req.Comment.User.Login, command.Type, command.Value, req.Issue.Number, allowed)

		found := false
		for _, label := range req.Issue.Labels {
			if label.Name == command.Value {
				found = true
				break
			}
		}

		if found == false {
			fmt.Println("Label didn't exist on issue.")
			return
		}

		if allowed {
			client, ctx := makeClient()
			_, err := client.Issues.RemoveLabelForIssue(ctx, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number, command.Value)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println("Label removed successfully or already removed.")
		}

		break

	default:
		log.Fatalln("Unable to work with comment: " + req.Comment.Body)
		break
	}
}

func parse(body string) *types.LabelAction {
	labelAction := types.LabelAction{}
	add := "Derek add label: "
	remove := "Derek remove label: "

	if len(body) > len(add) && body[0:len(add)] == add {
		label := body[len(add):]
		label = strings.Trim(label, " \t.,\n\r")
		labelAction.Type = "AddLabel"
		labelAction.Value = label
	} else if len(body) > len(remove) && body[0:len(remove)] == remove {
		label := body[len(remove):]
		label = strings.Trim(label, " \t.,\n\r")
		labelAction.Type = "RemoveLabel"
		labelAction.Value = label
	}
	return &labelAction
}

func getMaintainers(owner string, repository string) []string {
	client := http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://github.com/%s/%s/raw/master/MAINTAINERS", owner, repository), nil)

	res, resErr := client.Do(req)
	if resErr != nil {
		log.Fatalln(resErr)
	}

	if res.StatusCode != http.StatusOK {
		log.Fatalln("HTTP Status code: %d while fetching maintainers list", res.StatusCode)
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	bytesOut, _ := ioutil.ReadAll(res.Body)
	lines := string(bytesOut)
	return strings.Split(lines, "\n")
}

func isMaintainer(userLogin string, repository types.Repository) bool {
	maintainers := getMaintainers(repository.Owner.Login, repository.Name)
	fmt.Println("UserLogin: "+userLogin+", Maintainers: ", maintainers)
	allow := false
	for _, maintainer := range maintainers {
		if len(maintainer) > 0 && maintainer == userLogin {
			allow = true
			break
		}
	}

	return allow
}
