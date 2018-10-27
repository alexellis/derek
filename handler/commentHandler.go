// Copyright (c) Derek Author(s) 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handler

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/alexellis/derek/auth"
	"github.com/alexellis/derek/config"
	"github.com/alexellis/derek/factory"
	"github.com/alexellis/derek/types"
	"github.com/google/go-github/github"
)

const (
	openConstant            string = "open"
	ClosedConstant          string = "closed"
	closeConstant           string = "close"
	reopenConstant          string = "reopen"
	lockConstant            string = "Lock"
	unlockConstant          string = "Unlock"
	setTitleConstant        string = "SetTitle"
	assignConstant          string = "Assign"
	unassignConstant        string = "Unassign"
	removeLabelConstant     string = "RemoveLabel"
	addLabelConstant        string = "AddLabel"
	setMilestoneConstant    string = "SetMilestone"
	removeMilestoneConstant string = "RemoveMilestone"

	commandTriggerDefault string = "Derek "
	commandTriggerSlash   string = "/"
)

func makeClient(installation int, config config.Config) (*github.Client, context.Context) {
	ctx := context.Background()

	token := os.Getenv("personal_access_token")
	if len(token) == 0 {

		newToken, tokenErr := auth.MakeAccessTokenForInstallation(config.ApplicationID, installation, config.PrivateKey)
		if tokenErr != nil {
			log.Fatalln(tokenErr.Error())
		}

		token = newToken
	}

	client := factory.MakeClient(ctx, token, config)

	return client, ctx
}

func HandleComment(req types.IssueCommentOuter, config config.Config) {

	var feedback string
	var err error

	command := parse(req.Comment.Body, getCommandTrigger())

	switch command.Type {

	case addLabelConstant, removeLabelConstant:

		feedback, err = manageLabel(req, command.Type, command.Value, config)
		break

	case assignConstant, unassignConstant:

		feedback, err = manageAssignment(req, command.Type, command.Value, config)
		break

	case closeConstant, reopenConstant:

		feedback, err = manageState(req, command.Type, config)
		break

	case setTitleConstant:

		feedback, err = manageTitle(req, command.Type, command.Value, config)
		break

	case lockConstant, unlockConstant:

		feedback, err = manageLocking(req, command.Type, config)
		break

	case setMilestoneConstant, removeMilestoneConstant:

		feedback, err = updateMilestone(req, command.Type, command.Value, config)
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

func manageLabel(req types.IssueCommentOuter, cmdType string, labelValue string, config config.Config) (string, error) {

	var buffer bytes.Buffer
	labelAction := strings.Replace(strings.ToLower(cmdType), "label", "", 1)

	buffer.WriteString(fmt.Sprintf("%s wants to %s label of '%s' on issue #%d \n", req.Comment.User.Login, labelAction, labelValue, req.Issue.Number))

	found := findLabel(req.Issue.Labels, labelValue)

	if !validAction(found, cmdType, addLabelConstant, removeLabelConstant) {
		buffer.WriteString(fmt.Sprintf("Request to %s label of '%s' on issue #%d was unnecessary.", labelAction, labelValue, req.Issue.Number))
		return buffer.String(), nil
	}

	client, ctx := makeClient(req.Installation.ID, config)

	var err error

	if cmdType == addLabelConstant {
		_, _, err = client.Issues.AddLabelsToIssue(ctx, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number, []string{labelValue})
	} else {
		if isDcoLabel(labelValue) {
			buffer.WriteString(fmt.Sprintf("%s the request is not allowed.Label `%s` can be removed by owner or by signing out the commit.", req.Repository.Owner.Login, labelValue))
			return buffer.String(), nil
		} else {
			_, err = client.Issues.RemoveLabelForIssue(ctx, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number, labelValue)
		}
	}

	if err != nil {
		return buffer.String(), err
	}

	buffer.WriteString(fmt.Sprintf("Request to %s label of '%s' on issue #%d was successfully completed.", labelAction, labelValue, req.Issue.Number))
	return buffer.String(), nil
}

func manageTitle(req types.IssueCommentOuter, cmdType string, cmdValue string, config config.Config) (string, error) {

	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s wants to set the title of issue #%d\n", req.Comment.User.Login, req.Issue.Number))

	newTitle := cmdValue

	if newTitle == req.Issue.Title || len(newTitle) == 0 {
		buffer.WriteString(fmt.Sprintf("Setting the title of #%d by %s was unsuccessful as the new title was empty or unchanged.\n", req.Issue.Number, req.Comment.User.Login))
		return buffer.String(), nil
	}

	client, ctx := makeClient(req.Installation.ID, config)

	input := &github.IssueRequest{Title: &newTitle}

	_, _, err := client.Issues.Edit(ctx, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number, input)
	if err != nil {
		return buffer.String(), err
	}

	buffer.WriteString(fmt.Sprintf("Request to set the title of issue #%d by %s was successful.\n", req.Issue.Number, req.Comment.User.Login))
	return buffer.String(), nil
}

func manageAssignment(req types.IssueCommentOuter, cmdType string, cmdValue string, config config.Config) (string, error) {

	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("%s wants to %s user '%s' from issue #%d\n", req.Comment.User.Login, strings.ToLower(cmdType), cmdValue, req.Issue.Number))

	client, ctx := makeClient(req.Installation.ID, config)

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

func manageState(req types.IssueCommentOuter, cmdType string, config config.Config) (string, error) {

	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("%s wants to %s issue #%d\n", req.Comment.User.Login, cmdType, req.Issue.Number))

	newState, validTransition := checkTransition(cmdType, req.Issue.State)

	if !validTransition {
		buffer.WriteString(fmt.Sprintf("Request to %s issue #%d by %s was invalid.\n", cmdType, req.Issue.Number, req.Comment.User.Login))
		return buffer.String(), nil
	}

	client, ctx := makeClient(req.Installation.ID, config)
	input := &github.IssueRequest{State: &newState}

	_, _, err := client.Issues.Edit(ctx, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number, input)
	if err != nil {
		return buffer.String(), err
	}

	buffer.WriteString(fmt.Sprintf("Request to %s issue #%d by %s was successful.\n", cmdType, req.Issue.Number, req.Comment.User.Login))
	return buffer.String(), nil

}

func manageLocking(req types.IssueCommentOuter, cmdType string, config config.Config) (string, error) {

	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("%s wants to %s issue #%d\n", req.Comment.User.Login, strings.ToLower(cmdType), req.Issue.Number))

	if !validAction(req.Issue.Locked, cmdType, lockConstant, unlockConstant) {

		buffer.WriteString(fmt.Sprintf("Issue #%d is already %sed\n", req.Issue.Number, strings.ToLower(cmdType)))

		return buffer.String(), nil
	}

	client, ctx := makeClient(req.Installation.ID, config)

	var err error

	if cmdType == lockConstant {
		_, err = client.Issues.Lock(ctx, req.Repository.Owner.Login, req.Repository.Name,
			req.Issue.Number, &github.LockIssueOptions{})
	} else {
		_, err = client.Issues.Unlock(ctx, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number)
	}

	if err != nil {
		return buffer.String(), err
	}

	buffer.WriteString(fmt.Sprintf("Request to %s issue #%d by %s was successful.\n", strings.ToLower(cmdType), req.Issue.Number, req.Comment.User.Login))
	return buffer.String(), nil
}

func updateMilestone(req types.IssueCommentOuter, cmdType string, cmdValue string, config config.Config) (string, error) {

	milestoneValue := cmdValue
	var buffer bytes.Buffer
	milestoneAction := strings.Replace(strings.ToLower(cmdType), "milestone", "", 1)
	buffer.WriteString(fmt.Sprintf("%s wants to %s milestone of '%s' on issue #%d \n", req.Comment.User.Login, milestoneAction, milestoneValue, req.Issue.Number))

	allMilestones := &github.MilestoneListOptions{}
	var milestoneNumber *int
	var err error

	client, ctx := makeClient(req.Installation.ID, config)
	theMilestones, _, milErr := client.Issues.ListMilestones(ctx, req.Repository.Owner.Login, req.Repository.Name, allMilestones)
	if milErr != nil {
		return buffer.String(), milErr
	}

	switch cmdType {
	case setMilestoneConstant:
		if req.Issue.Milestone.Title == cmdValue {
			buffer.WriteString(fmt.Sprintf("Setting the milestone of #%d by %s was unnecessary.\n", req.Issue.Number, req.Comment.User.Login))
			return buffer.String(), nil
		}
		for _, mil := range theMilestones {
			if mil != nil && *mil.Title == milestoneValue {
				milestoneNumber = mil.Number
				break
			}
		}
		input := &github.IssueRequest{
			Milestone: milestoneNumber,
		}
		_, _, err = client.Issues.Edit(ctx, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number, input)
		if err != nil {
			return buffer.String(), err
		}
	case removeMilestoneConstant:
		if err = removeMilestone(client, ctx, req.Issue.URL); err != nil {
			return buffer.String(), err
		}
	default:
		buffer.WriteString(fmt.Sprintf("Unknown milestone action %q on issue #%d.", milestoneAction, req.Issue.Number))
		return buffer.String(), nil
	}

	buffer.WriteString(fmt.Sprintf("Request to %s milestone of '%s' on issue #%d was successfully completed.", milestoneAction, milestoneValue, req.Issue.Number))
	return buffer.String(), nil
}

func parse(body, commandTrigger string) *types.CommentAction {

	commentAction := types.CommentAction{}

	commands := map[string]string{
		commandTrigger + "add label: ":        addLabelConstant,
		commandTrigger + "remove label: ":     removeLabelConstant,
		commandTrigger + "assign: ":           assignConstant,
		commandTrigger + "unassign: ":         unassignConstant,
		commandTrigger + "close":              closeConstant,
		commandTrigger + "reopen":             reopenConstant,
		commandTrigger + "set title: ":        setTitleConstant,
		commandTrigger + "edit title: ":       setTitleConstant,
		commandTrigger + "lock":               lockConstant,
		commandTrigger + "unlock":             unlockConstant,
		commandTrigger + "set milestone: ":    setMilestoneConstant,
		commandTrigger + "remove milestone: ": removeMilestoneConstant,
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

	if requestedAction == closeConstant && currentState != ClosedConstant {
		return ClosedConstant, true
	} else if requestedAction == reopenConstant && currentState != openConstant {
		return openConstant, true
	}

	return "", false
}

//removeMilestone sets milestones field to interface{} aka. null since library does not support that
//Reference to issue - https://github.com/google/go-github/issues/236
func removeMilestone(client *github.Client, ctx context.Context, URL string) error {
	req, err := client.NewRequest("PATCH", URL, &struct {
		Milestone interface{} `json:"milestone"`
	}{})
	if err != nil {
		return err
	}
	if _, err = client.Do(ctx, req, nil); err != nil {
		return err
	}
	return nil
}

func isDcoLabel(labelValue string) bool {
	return strings.ToLower(labelValue) == "no-dco"
}

func getCommandTrigger() string {
	commandTrigger := commandTriggerDefault
	if os.Getenv("use_slash_trigger") == "true" {
		commandTrigger = commandTriggerSlash
	}
	return commandTrigger
}
