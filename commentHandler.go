package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/alexellis/derek/auth"
	"github.com/alexellis/derek/types"
	"github.com/google/go-github/github"
)

const open = "open"
const closed = "closed"

func makeClient(installation int) (*github.Client, context.Context) {
	ctx := context.Background()

	token := os.Getenv("access_token")
	if len(token) == 0 {

		applicationID := os.Getenv("application")
		privateKeyPath := os.Getenv("private_key")

		newToken, tokenErr := auth.MakeAccessTokenForInstallation(applicationID, installation, privateKeyPath)
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

		newState, validTransition := checkTransition(command.Type, req.Issue.State)

		if !validTransition {
			fmt.Printf("Request to %s issue #%d by %s was invalid.\n", command.Type, req.Issue.Number, req.Comment.User.Login)
			return
		}

		client, ctx := makeClient(req.Installation.ID)
		input := &github.IssueRequest{State: &newState}

		_, _, err := client.Issues.Edit(ctx, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number, input)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("Request to %s issue #%d by %s was successful.\n", command.Type, req.Issue.Number, req.Comment.User.Login)

		break

	case "SetTitle":

		fmt.Printf("%s wants to set the title of issue #%d\n", req.Comment.User.Login, req.Issue.Number)

		newTitle := command.Value

		if newTitle == req.Issue.Title {
			fmt.Printf("Setting the title of #%d by %s was unsuccessful as the new title was empty or unchanged.\n", req.Issue.Number, req.Comment.User.Login)
			return
		}

		client, ctx := makeClient(req.Installation.ID)

		input := &github.IssueRequest{Title: &newTitle}

		_, _, err := client.Issues.Edit(ctx, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number, input)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("Request to set the title of issue #%d by %s was successful.\n", req.Issue.Number, req.Comment.User.Login)

		break

	case "Lock":
		fmt.Printf("%s wants to lock issue #%d\n", req.Comment.User.Login, req.Issue.Number)

		if req.Issue.Locked {
			fmt.Printf("Issue #%d is already locked.\n", req.Issue.Number)
			return
		}

		client, ctx := makeClient(req.Installation.ID)

		_, err := client.Issues.Lock(ctx, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("Request to lock issue #%d by %s was successful.\n", req.Issue.Number, req.Comment.User.Login)

		break

	case "Unlock":
		fmt.Printf("%s wants to unlock issue #%d\n", req.Comment.User.Login, req.Issue.Number)

		if !req.Issue.Locked {
			fmt.Printf("Issue #%d is already unlocked\n", req.Issue.Number)
			return
		}

		client, ctx := makeClient(req.Installation.ID)

		_, err := client.Issues.Unlock(ctx, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("Request to unlock issue #%d by %s was successful.\n", req.Issue.Number, req.Comment.User.Login)

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
		"Derek set title: ":    "SetTitle",
		"Derek lock":           "Lock",
		"Derek unlock":         "Unlock",
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
	return (len(body) > len(trigger) && body[0:len(trigger)] == trigger) ||
		(body == trigger && !strings.HasSuffix(trigger, ": "))
}

func checkTransition(requestedAction string, currentState string) (string, bool) {

	desiredState := ""
	validTransition := false

	if requestedAction == "close" && currentState != closed {
		desiredState = closed
		validTransition = true
	} else if requestedAction == "reopen" && currentState != open {
		desiredState = open
		validTransition = true
	}

	return desiredState, validTransition
}
