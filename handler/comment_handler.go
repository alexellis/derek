// Copyright (c) Derek Author(s) 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handler

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/alexellis/derek/auth"
	"github.com/alexellis/derek/config"
	"github.com/alexellis/derek/factory"
	"github.com/alexellis/derek/types"
	"github.com/google/go-github/github"
	log "github.com/sirupsen/logrus"
)

const (
	openConstant             string = "open"
	ClosedConstant           string = "closed"
	closeConstant            string = "close"
	reopenConstant           string = "reopen"
	lockConstant             string = "Lock"
	unlockConstant           string = "Unlock"
	setTitleConstant         string = "SetTitle"
	assignConstant           string = "Assign"
	unassignConstant         string = "Unassign"
	removeLabelConstant      string = "RemoveLabel"
	addLabelConstant         string = "AddLabel"
	setMilestoneConstant     string = "SetMilestone"
	removeMilestoneConstant  string = "RemoveMilestone"
	assignReviewerConstant   string = "AssignReviewer"
	unassignReviewerConstant string = "UnassignReviewer"
	messageConstant          string = "message"

	noDCO             string = "no-dco"
	labelLimitDefault int    = 5
	labelLimitEnvVar  string = "multilabel_limit"
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

// HandleComment handles a comment
func HandleComment(req types.IssueCommentOuter, config config.Config, derekConfig *types.DerekRepoConfig) {

	var feedback string
	var err error

	command := parse(req.Comment.Body, getCommandTriggers())

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

	case assignReviewerConstant, unassignReviewerConstant:

		pr := types.PullRequest{
			Number: req.Issue.Number,
		}
		prReq := types.PullRequestOuter{
			Repository:          req.Repository,
			PullRequest:         pr,
			Action:              req.Action,
			InstallationRequest: req.InstallationRequest,
		}
		feedback, err = editReviewers(prReq, command.Type, command.Value, config)
		break

	case messageConstant:

		feedback, err = createMessage(req, command.Type, command.Value, config, derekConfig)
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

func classifyLabels(currentLabels []types.IssueLabel, labelAction string, labelValue string) ([]string, []string) {

	var actionableLabels, unactionableLabels []string

	requestedLabels := strings.Split(labelValue, ",")

	for _, requestedLabel := range requestedLabels {

		requestedLabel = strings.TrimSpace(requestedLabel)

		found := findLabel(currentLabels, requestedLabel)

		if validAction(found, labelAction, addLabelConstant, removeLabelConstant) {
			actionableLabels = append(actionableLabels, requestedLabel)
		} else {
			unactionableLabels = append(unactionableLabels, requestedLabel)
		}

	}
	return actionableLabels, unactionableLabels
}

func manageLabel(req types.IssueCommentOuter, cmdType string, labelValue string, config config.Config) (string, error) {

	var buffer bytes.Buffer
	labelAction := strings.Replace(strings.ToLower(cmdType), "label", "", 1)
	buffer.WriteString(fmt.Sprintf("%s wants to %s label(s) of '%s' on issue #%d.\n", req.Comment.User.Login, labelAction, labelValue, req.Issue.Number))

	actionableLabels, unactionableLabels := classifyLabels(req.Issue.Labels, cmdType, labelValue)

	if len(unactionableLabels) > 0 {
		buffer.WriteString(fmt.Sprintf("Request to %s label(s) of '%s' on issue #%d was unnecessary.\n", labelAction, strings.Join(unactionableLabels, ", "), req.Issue.Number))

		if len(actionableLabels) == 0 {
			buffer.WriteString(fmt.Sprintf("No further valid labels found - no action taken on issue #%d.\n", req.Issue.Number))
			return buffer.String(), nil
		}
	}

	client, ctx := makeClient(req.Installation.ID, config)

	var err error

	maxActionableLabels := getMultiLabelLimit()

	if len(actionableLabels) > maxActionableLabels {
		buffer.WriteString(fmt.Sprintf("Label(s) '%s' on issue #%d were ignored as they fall outside of the configured limit of %d.\n", strings.Join(actionableLabels[maxActionableLabels:], ", "), req.Issue.Number, maxActionableLabels))
		actionableLabels = actionableLabels[:maxActionableLabels]
	}

	if cmdType == addLabelConstant {

		_, _, err = client.Issues.AddLabelsToIssue(ctx, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number, actionableLabels)

		if err != nil {
			return buffer.String(), err
		}

	} else {

		actionedLabels := actionableLabels[:0]

		for _, actionableLabel := range actionableLabels {

			if isDcoLabel(actionableLabel) {

				buffer.WriteString(fmt.Sprintf("The request to remove `%s` by %s was not allowed - label can be removed by owner or by signing off the commit.\n", actionableLabel, req.Repository.Owner.Login))

			} else {

				_, err = client.Issues.RemoveLabelForIssue(ctx, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number, actionableLabel)

				if err != nil {
					return buffer.String(), err
				}

				actionedLabels = append(actionedLabels, actionableLabel)
			}
		}
		actionableLabels = actionedLabels
	}

	buffer.WriteString(fmt.Sprintf("Request to %s label(s) of '%s' on issue #%d was successfully completed.\n", labelAction, strings.Join(actionableLabels, ", "), req.Issue.Number))
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

func editReviewers(req types.PullRequestOuter, cmdType string, cmdValue string, config config.Config) (string, error) {
	var buffer bytes.Buffer

	client, ctx := makeClient(req.Installation.ID, config)

	reviewer := github.ReviewersRequest{Reviewers: []string{cmdValue}}

	var err error

	if cmdType == unassignReviewerConstant {
		_, err = client.PullRequests.RemoveReviewers(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number, reviewer)
	} else {
		_, _, err = client.PullRequests.RequestReviewers(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number, reviewer)
	}

	if err != nil {
		return buffer.String(), err
	}

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

func parse(body string, commandTriggers []string) *types.CommentAction {
	commentAction := types.CommentAction{}

	for _, commandTrigger := range commandTriggers {
		commands := map[string]string{
			commandTrigger + "add label: ":        addLabelConstant,
			commandTrigger + "remove label: ":     removeLabelConstant,
			commandTrigger + "add labels: ":       addLabelConstant,
			commandTrigger + "remove labels: ":    removeLabelConstant,
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
			commandTrigger + "set reviewer: ":     assignReviewerConstant,
			commandTrigger + "clear reviewer: ":   unassignReviewerConstant,
			commandTrigger + "message: ":          messageConstant,
			commandTrigger + "msg: ":              messageConstant,
		}

		for trigger, commandType := range commands {

			if isValidCommand(body, trigger) {
				commentAction.Type = commandType
				commentAction.Value = getCommandValue(body, len(trigger))
				break
			}
		}
	}

	return &commentAction
}

func getCommandValue(commentBody string, triggerLength int) string {

	val := commentBody[triggerLength:]
	val = strings.Split(val, "\n")[0]
	val = strings.Trim(val, " \t.,\n\r")
	return val
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

// removeMilestone sets milestones field to interface{} aka. null since library does not support that
// reference to issue - https://github.com/google/go-github/issues/236
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
	return strings.ToLower(labelValue) == noDCO
}

func getCommandTriggers() []string {
	return []string{"Derek ", "/"}
}

func getMultiLabelLimit() int {

	val, ok := os.LookupEnv(labelLimitEnvVar)
	if ok {
		configuredLimit, err := strconv.Atoi(val)
		if err != nil {
			return labelLimitDefault
		}
		return configuredLimit
	}
	return labelLimitDefault
}

func createMessage(req types.IssueCommentOuter, cmdType, cmdValue string, config config.Config, derekConfig *types.DerekRepoConfig) (string, error) {
	var err error

	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("%s wants to add message of type '%s' on issue #%d \n", req.Comment.User.Login, cmdValue, req.Issue.Number))

	messageValue, err := createIssueComment(derekConfig.Messages, cmdValue)
	if err != nil {
		return buffer.String(), fmt.Errorf("Error while filtering message: %s", err.Error())
	}

	buffer.WriteString(fmt.Sprintf("Message '%s' found.\n", cmdValue))

	client, ctx := makeClient(req.Installation.ID, config)

	_, resp, err := client.Issues.CreateComment(ctx, req.Repository.Owner.Login, req.Repository.Name, req.Issue.Number, messageValue)
	if err != nil {
		return buffer.String(), err
	}
	buffer.WriteString(fmt.Sprintf("Successfully applied message: `%s` status code: %d\n",
		cmdValue,
		resp.StatusCode))

	return buffer.String(), nil
}

func createIssueComment(messages []types.Message, wantedMessage string) (*github.IssueComment, error) {
	for _, message := range messages {
		if message.Name == wantedMessage {
			return &github.IssueComment{
				Body: &message.Value,
			}, nil
		}
	}
	return nil, fmt.Errorf("Message: `%s` is not configured", wantedMessage)
}
