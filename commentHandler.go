package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/alexellis/derek/auth"
	"github.com/alexellis/derek/types"
	"github.com/google/go-github/github"
)

const openConstant string = "open"
const closedConstant string = "closed"
const closeConstant string = "close"
const reopenConstant string = "reopen"
const lockConstant string = "Lock"
const unlockConstant string = "Unlock"
const setTitleConstant string = "SetTitle"
const assignConstant string = "Assign"
const unassignConstant string = "Unassign"
const removeLabelConstant string = "RemoveLabel"
const addLabelConstant string = "AddLabel"

func makeClient(installation int) (*github.Client, context.Context) {
	ctx := context.Background()

	token := os.Getenv("access_token")
	if len(token) == 0 {

		applicationID := os.Getenv("application")

		newToken, tokenErr := auth.MakeAccessTokenForInstallation(applicationID, installation)
		if tokenErr != nil {
			log.Fatalln(tokenErr.Error())
		}

		token = newToken
	}

	client := auth.MakeClient(ctx, token)

	return client, ctx
}

func handleComment(req types.IssueCommentOuter) {

	var feedback string
	var err error

	command := parse(req.Comment.Body)

	switch command.Type {

	case addLabelConstant, removeLabelConstant:

		feedback, err = manageLabel(req, command.Type, command.Value)
		break

	case assignConstant, unassignConstant:

		feedback, err = manageAssignment(req, command.Type, command.Value)
		break

	case closeConstant, reopenConstant:

		feedback, err = manageState(req, command.Type)
		break

	case setTitleConstant:

		feedback, err = manageTitle(req, command.Type, command.Value)
		break

	case lockConstant, unlockConstant:

		feedback, err = manageLocking(req, command.Type)
		break

	default:
		feedback = "Unable to work with comment: " + req.Comment.Body
		err = nil
		break
	}

	fmt.Print(feedback)

	if err != nil {
		fmt.Println(err)
	}
}

func findLabel(currentLabels []types.IssueLabel, cmdLabel string) bool {

	for _, label := range currentLabels {
		if strings.EqualFold(label.Name, cmdLabel) {
			return true
		}
	}
	return false
}

func manageLabel(req types.IssueCommentOuter, cmdType string, labelValue string) (string, error) {

	var buffer bytes.Buffer
	labelAction := strings.Replace(strings.ToLower(cmdType), "label", "", 1)

	buffer.WriteString(fmt.Sprintf("%s wants to %s label of '%s' on issue #%d \n", req.Comment.User.Login, labelAction, labelValue, req.Issue.Number))

	found := findLabel(req.Issue.Labels, labelValue)

	if !validAction(found, cmdType, addLabelConstant, removeLabelConstant) {
		buffer.WriteString(fmt.Sprintf("Request to %s label of '%s' on issue #%d was unnecessary.", labelAction, labelValue, req.Issue.Number))
		return buffer.String(), nil
	}

	client, ctx := makeClient(req.Installation.ID)

	var err error

	if cmdType == addLabelConstant {
		_, _, err = client.Issues.AddLabelsToIssue(ctx, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number, []string{labelValue})
	} else {
		_, err = client.Issues.RemoveLabelForIssue(ctx, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number, labelValue)
	}

	if err != nil {
		return buffer.String(), err
	}

	buffer.WriteString(fmt.Sprintf("Request to %s label of '%s' on issue #%d was successfully completed.", labelAction, labelValue, req.Issue.Number))
	return buffer.String(), nil
}

func manageTitle(req types.IssueCommentOuter, cmdType string, cmdValue string) (string, error) {

	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s wants to set the title of issue #%d\n", req.Comment.User.Login, req.Issue.Number))

	newTitle := cmdValue

	if newTitle == req.Issue.Title || len(newTitle) == 0 {
		buffer.WriteString(fmt.Sprintf("Setting the title of #%d by %s was unsuccessful as the new title was empty or unchanged.\n", req.Issue.Number, req.Comment.User.Login))
		return buffer.String(), nil
	}

	client, ctx := makeClient(req.Installation.ID)

	input := &github.IssueRequest{Title: &newTitle}

	_, _, err := client.Issues.Edit(ctx, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number, input)
	if err != nil {
		return buffer.String(), err
	}

	buffer.WriteString(fmt.Sprintf("Request to set the title of issue #%d by %s was successful.\n", req.Issue.Number, req.Comment.User.Login))
	return buffer.String(), nil
}

func manageAssignment(req types.IssueCommentOuter, cmdType string, cmdValue string) (string, error) {

	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("%s wants to %s user '%s' from issue #%d\n", req.Comment.User.Login, strings.ToLower(cmdType), cmdValue, req.Issue.Number))

	client, ctx := makeClient(req.Installation.ID)

	if cmdValue == "me" {
		cmdValue = req.Comment.User.Login
	}

	var err error

	if cmdType == unassignConstant {
		_, _, err = client.Issues.RemoveAssignees(ctx, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number, []string{cmdValue})
	} else {
		_, _, err = client.Issues.AddAssignees(ctx, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number, []string{cmdValue})
	}

	if err != nil {
		return buffer.String(), err
	}

	buffer.WriteString(fmt.Sprintf("%s %sed successfully or already %sed.\n", cmdValue, strings.ToLower(cmdType), strings.ToLower(cmdType)))
	return buffer.String(), nil
}

func manageState(req types.IssueCommentOuter, cmdType string) (string, error) {

	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("%s wants to %s issue #%d\n", req.Comment.User.Login, cmdType, req.Issue.Number))

	newState, validTransition := checkTransition(cmdType, req.Issue.State)

	if !validTransition {
		buffer.WriteString(fmt.Sprintf("Request to %s issue #%d by %s was invalid.\n", cmdType, req.Issue.Number, req.Comment.User.Login))
		return buffer.String(), nil
	}

	client, ctx := makeClient(req.Installation.ID)
	input := &github.IssueRequest{State: &newState}

	_, _, err := client.Issues.Edit(ctx, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number, input)
	if err != nil {
		return buffer.String(), err
	}

	buffer.WriteString(fmt.Sprintf("Request to %s issue #%d by %s was successful.\n", cmdType, req.Issue.Number, req.Comment.User.Login))
	return buffer.String(), nil

}

func manageLocking(req types.IssueCommentOuter, cmdType string) (string, error) {

	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("%s wants to %s issue #%d\n", req.Comment.User.Login, strings.ToLower(cmdType), req.Issue.Number))

	if !validAction(req.Issue.Locked, cmdType, lockConstant, unlockConstant) {

		buffer.WriteString(fmt.Sprintf("Issue #%d is already %sed\n", req.Issue.Number, strings.ToLower(cmdType)))

		return buffer.String(), nil
	}

	client, ctx := makeClient(req.Installation.ID)

	var err error

	if cmdType == lockConstant {
		_, err = client.Issues.Lock(ctx, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number)
	} else {
		_, err = client.Issues.Unlock(ctx, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number)
	}

	if err != nil {
		return buffer.String(), err
	}

	buffer.WriteString(fmt.Sprintf("Request to %s issue #%d by %s was successful.\n", strings.ToLower(cmdType), req.Issue.Number, req.Comment.User.Login))
	return buffer.String(), nil
}

func parse(body string) *types.CommentAction {

	commentAction := types.CommentAction{}

	commands := map[string]string{
		"Derek add label: ":    addLabelConstant,
		"Derek remove label: ": removeLabelConstant,
		"Derek assign: ":       assignConstant,
		"Derek unassign: ":     unassignConstant,
		"Derek close":          closeConstant,
		"Derek reopen":         reopenConstant,
		"Derek set title: ":    setTitleConstant,
		"Derek edit title: ":   setTitleConstant,
		"Derek lock":           lockConstant,
		"Derek unlock":         unlockConstant,
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

func validAction(running bool, requestedAction string, start string, stop string) bool {

	return !running && requestedAction == start || running && requestedAction == stop

}

func checkTransition(requestedAction string, currentState string) (string, bool) {

	if requestedAction == closeConstant && currentState != closedConstant {
		return closedConstant, true
	} else if requestedAction == reopenConstant && currentState != openConstant {
		return openConstant, true
	}

	return "", false
}
