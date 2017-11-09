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

const open = "open"
const closed = "closed"
const maintainersFileEnv = "maintainers_file"
const defaultMaintFile = "MAINTAINERS"

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
			_, res, err := client.Issues.AddLabelsToIssue(ctx, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number, []string{command.Value})
			if err != nil {
				log.Fatalf("%s, limit: %d, remaining: %d", err, res.Limit, res.Remaining)
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
	case "Assign":
		allowed := isMaintainer(req.Comment.User.Login, req.Repository)
		fmt.Printf("%s wants to %s user %s to issue %d - allowed? %t\n", req.Comment.User.Login, command.Type, command.Value, req.Issue.Number, allowed)

		if allowed {
			client, ctx := makeClient()
			assignee := command.Value
			if assignee == "me" {
				assignee = req.Comment.User.Login
			}
			_, _, err := client.Issues.AddAssignees(ctx, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number, []string{assignee})
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Printf("%s assigned successfully or already assigned.\n", command.Value)
		}

		break
	case "Unassign":
		allowed := isMaintainer(req.Comment.User.Login, req.Repository)
		fmt.Printf("%s wants to %s user %s from issue %d - allowed? %t\n", req.Comment.User.Login, command.Type, command.Value, req.Issue.Number, allowed)

		if allowed {
			client, ctx := makeClient()
			assignee := command.Value
			if assignee == "me" {
				assignee = req.Comment.User.Login
			}
			_, _, err := client.Issues.RemoveAssignees(ctx, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number, []string{assignee})
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Printf("%s unassigned successfully or already unassigned.\n", command.Value)
		}

		break
	case "close", "reopen":
		allowed := isMaintainer(req.Comment.User.Login, req.Repository)
		fmt.Printf("%s wants to %s issue #%d - allowed? %t\n", req.Comment.User.Login, command.Type, req.Issue.Number, allowed)

		if allowed {
			client, ctx := makeClient()

			var state string

			if command.Type == "close" {
				state = closed
			} else if command.Type == "reopen" {
				state = open
			}
			input := &github.IssueRequest{State: &state}

			_, _, err := client.Issues.Edit(ctx, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number, input)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Printf("Request to %s issue #%d by %s was successful.\n", command.Type, req.Issue.Number, req.Comment.User.Login)
		}

		break
	default:
		log.Fatalln("Unable to work with comment: " + req.Comment.Body)
		break
	}
}

func parse(body string) *types.CommentAction {
	commentAction := types.CommentAction{}

	commands := map[string]string{
		"Derek add label: ":    "AddLabel",
		"Derek remove label: ": "RemoveLabel",
		"Derek assign: ":       "Assign",
		"Derek unassign: ":     "Unassign",
		"Derek close":          "close",
		"Derek reopen":         "reopen",
	}

	for trigger, commandType := range commands {

		if isValidCommand(body, trigger) {
			val := body[len(trigger):]
			val = strings.Trim(val, " \t.,\n\r")
			commentAction.Type = commandType
			commentAction.Value = val
			break
		}
	}

	return &commentAction
}

func isValidCommand(body string, trigger string) bool {

	return (len(body) > len(trigger) && body[0:len(trigger)] == trigger) || (body == trigger && !strings.HasSuffix(trigger, ": "))

}

func getMaintainers(owner string, repository string) []string {
	client := http.Client{}

	maintainersFile := getEnv(maintainersFileEnv, defaultMaintFile)
	maintainersFile = fmt.Sprintf("https://github.com/%s/%s/raw/master/%s", owner, repository, strings.Trim(maintainersFile, "/"))

	req, _ := http.NewRequest(http.MethodGet, maintainersFile, nil)

	res, resErr := client.Do(req)
	if resErr != nil {
		log.Fatalln(resErr)
	}

	if res.StatusCode != http.StatusOK {
		log.Fatalln(fmt.Sprintf("HTTP Status code: %d while fetching maintainers list (%s)", res.StatusCode, maintainersFile))
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

func getEnv(envVar, assumed string) string {
	if value, exists := os.LookupEnv(envVar); exists {
		return value
	}
	return assumed
}
