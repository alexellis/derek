package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/google/go-github/github"
	"github.com/alexellis/derek/auth"
	"github.com/alexellis/derek/types"
)

const open = "open"
const closed = "closed"

func makeClient(installation int) (*github.Client, context.Context) {
	ctx := context.Background()

	token := os.Getenv("access_token")
	if len(token) == 0 {
		newToken, tokenErr := auth.MakeAccessTokenForInstallation(os.Getenv("application"), installation, os.Getenv("private_key"))
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

		fmt.Printf("%s wants to %s of %s to issue %d\n", req.Comment.User.Login, command.Type, command.Value, req.Issue.Number)

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

		client, ctx := makeClient(req.Installation.ID)
		_, res, err := client.Issues.AddLabelsToIssue(ctx, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number, []string{command.Value})
		if err != nil {
			log.Fatalf("%s, limit: %d, remaining: %d", err, res.Limit, res.Remaining)
		}

		fmt.Println("Label added successfully or already existed.")

		break

	case "RemoveLabel":

		fmt.Printf("%s wants to %s of %s to issue %d \n", req.Comment.User.Login, command.Type, command.Value, req.Issue.Number)

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

		client, ctx := makeClient(req.Installation.ID)
		_, err := client.Issues.RemoveLabelForIssue(ctx, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number, command.Value)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Label removed successfully or already removed.")

		break
	case "Assign":

		fmt.Printf("%s wants to %s user %s to issue %d\n", req.Comment.User.Login, command.Type, command.Value, req.Issue.Number)

		client, ctx := makeClient(req.Installation.ID)
		assignee := command.Value
		if assignee == "me" {
			assignee = req.Comment.User.Login
		}
		_, _, err := client.Issues.AddAssignees(ctx, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number, []string{assignee})
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%s assigned successfully or already assigned.\n", command.Value)

		break
	case "Unassign":

		fmt.Printf("%s wants to %s user %s from issue %d\n", req.Comment.User.Login, command.Type, command.Value, req.Issue.Number)

		client, ctx := makeClient(req.Installation.ID)
		assignee := command.Value
		if assignee == "me" {
			assignee = req.Comment.User.Login
		}
		_, _, err := client.Issues.RemoveAssignees(ctx, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number, []string{assignee})
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%s unassigned successfully or already unassigned.\n", command.Value)

		break
	case "close", "reopen":
		fmt.Printf("%s wants to %s issue #%d\n", req.Comment.User.Login, command.Type, req.Issue.Number)

		client, ctx := makeClient(req.Installation.ID)

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
