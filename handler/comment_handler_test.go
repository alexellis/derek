// Copyright (c) Derek Author(s) 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handler

import (
	"os"
	"strings"
	"testing"

	"github.com/alexellis/derek/types"
	"github.com/google/go-github/github"
)

func Test_fixCommitsMessage(t *testing.T) {
	t.Log(fixCommits)
}

func Test_getCommandTrigger(t *testing.T) {

	triggers := getCommandTriggers()

	want := []string{"Derek ", "/"}
	found := 0

	for _, v := range want {
		for _, a := range triggers {
			if a == v {
				found++
			}
		}
	}

	if found != len(want) {
		t.Errorf("Wanted to find %d triggers, but got %d", len(want), found)
	}
}

func Test_Parsing_OpenClose(t *testing.T) {

	var actionOptions = []struct {
		title          string
		body           string
		expectedAction string
	}{
		{
			title:          "Correct reopen command",
			body:           "reopen",
			expectedAction: "reopen",
		},
		{ //this case replaces Test_Parsing_Close
			title:          "Correct close command",
			body:           "close",
			expectedAction: "close",
		},
		{
			title:          "invalid command",
			body:           "dance",
			expectedAction: "",
		},
		{
			title:          "Longer reopen command",
			body:           "reopen: ",
			expectedAction: "reopen",
		},
		{
			title:          "Longer close command",
			body:           "close: ",
			expectedAction: "close",
		},
	}

	for _, test := range actionOptions {
		t.Run(test.title, func(t *testing.T) {

			for _, trigger := range getCommandTriggers() {
				action := parse(trigger+test.body, getCommandTriggers())
				if action.Type != test.expectedAction {
					t.Errorf("Action - want: %s, got %s", test.expectedAction, action.Type)
				}
			}
		})
	}
}

func Test_Parsing_Labels(t *testing.T) {

	var labelOptions = []struct {
		title        string
		body         string
		expectedType string
		expectedVal  string
	}{
		{ //this case replaces Test_Parsing_AddLabel
			title:        "Add label of demo",
			body:         "add label: demo",
			expectedType: "AddLabel",
			expectedVal:  "demo",
		},
		{
			title:        "Remove label of demo",
			body:         "remove label: demo",
			expectedType: "RemoveLabel",
			expectedVal:  "demo",
		},
		{
			title:        "Invalid label action",
			body:         "peel label: demo",
			expectedType: "",
			expectedVal:  "",
		},
	}

	for _, test := range labelOptions {
		t.Run(test.title, func(t *testing.T) {

			for _, trigger := range getCommandTriggers() {
				action := parse(trigger+test.body, getCommandTriggers())
				if action.Type != test.expectedType || action.Value != test.expectedVal {
					t.Errorf("Action - wanted: %s, got %s\nLabel - wanted: %s, got %s", test.expectedType, action.Type, test.expectedVal, action.Value)
				}
			}
		})
	}
}

func Test_Parsing_Assignments(t *testing.T) {

	var assignmentOptions = []struct {
		title        string
		body         string
		expectedType string
		expectedVal  string
	}{
		{
			title:        "Assign to burt",
			body:         "assign: burt",
			expectedType: assignConstant,
			expectedVal:  "burt",
		},
		{
			title:        "Unassign burt",
			body:         "unassign: burt",
			expectedType: unassignConstant,
			expectedVal:  "burt",
		},
		{
			title:        "Assign to me",
			body:         "assign: me",
			expectedType: assignConstant,
			expectedVal:  "me",
		},
		{
			title:        "Unassign me",
			body:         "unassign: me",
			expectedType: unassignConstant,
			expectedVal:  "me",
		},
		{
			title:        "Invalid assignment action",
			body:         "consign: burt",
			expectedType: "",
			expectedVal:  "",
		},
		{
			title:        "Unassign blank",
			body:         "unassign: ",
			expectedType: "",
			expectedVal:  "",
		},
	}

	for _, test := range assignmentOptions {
		t.Run(test.title, func(t *testing.T) {

			for _, trigger := range getCommandTriggers() {
				action := parse(trigger+test.body, getCommandTriggers())
				if action.Type != test.expectedType || action.Value != test.expectedVal {
					t.Errorf("Action - wanted: %s, got %s\nMaintainer - wanted: %s, got %s", test.expectedType, action.Type, test.expectedVal, action.Value)
				}
			}
		})
	}
}

func Test_Parsing_Titles(t *testing.T) {

	var titleOptions = []struct {
		title        string
		body         string
		expectedType string
		expectedVal  string
	}{
		{
			title:        "Set Title",
			body:         "set title: This is a really great Title!",
			expectedType: setTitleConstant,
			expectedVal:  "This is a really great Title!",
		},
		{
			title:        "Mis-spelling of title",
			body:         "set titel: This is a really great Title!",
			expectedType: "",
			expectedVal:  "",
		},
		{
			title:        "Empty Title",
			body:         "set title: ",
			expectedType: "", //blank because it should fail isValidCommand
			expectedVal:  "",
		},
		{
			title:        "Empty Title (Double Space)",
			body:         "set title:  ",
			expectedType: setTitleConstant,
			expectedVal:  "",
		},
	}

	for _, test := range titleOptions {
		t.Run(test.title, func(t *testing.T) {

			for _, trigger := range getCommandTriggers() {
				action := parse(trigger+test.body, getCommandTriggers())
				if action.Type != test.expectedType || action.Value != test.expectedVal {
					t.Errorf("\nAction - wanted: %s, got %s\nValue - wanted: %s, got %s", test.expectedType, action.Type, test.expectedVal, action.Value)
				}
			}
		})
	}
}

func Test_assessState(t *testing.T) {

	var stateOptions = []struct {
		title            string
		requestedAction  string
		currentState     string
		expectedNewState string
		expectedBool     bool
	}{
		{
			title:            "Currently Closed and trying to close",
			requestedAction:  closeConstant,
			currentState:     ClosedConstant,
			expectedNewState: "",
			expectedBool:     false,
		},
		{
			title:            "Currently Open and trying to reopen",
			requestedAction:  reopenConstant,
			currentState:     openConstant,
			expectedNewState: "",
			expectedBool:     false,
		},
		{
			title:            "Currently Closed and trying to open",
			requestedAction:  reopenConstant,
			currentState:     ClosedConstant,
			expectedNewState: openConstant,
			expectedBool:     true,
		},
		{
			title:            "Currently Open and trying to close",
			requestedAction:  closeConstant,
			currentState:     openConstant,
			expectedNewState: ClosedConstant,
			expectedBool:     true,
		},
	}

	for _, test := range stateOptions {
		t.Run(test.title, func(t *testing.T) {

			newState, validTransition := checkTransition(test.requestedAction, test.currentState)

			if newState != test.expectedNewState || validTransition != test.expectedBool {
				t.Errorf("\nStates - wanted: %s, got %s\nValidity - wanted: %t, got %t\n", test.expectedNewState, newState, test.expectedBool, validTransition)
			}
		})
	}
}

func Test_validAction(t *testing.T) {

	var stateOptions = []struct {
		title           string
		running         bool
		requestedAction string
		start           string
		stop            string
		expectedBool    bool
	}{
		{
			title:           "Currently unlocked and trying to lock",
			running:         false,
			requestedAction: lockConstant,
			start:           lockConstant,
			stop:            unlockConstant,
			expectedBool:    true,
		},
		{
			title:           "Currently unlocked and trying to unlock",
			running:         false,
			requestedAction: unlockConstant,
			start:           lockConstant,
			stop:            unlockConstant,
			expectedBool:    false,
		},
		{
			title:           "Currently locked and trying to lock",
			running:         true,
			requestedAction: lockConstant,
			start:           lockConstant,
			stop:            unlockConstant,
			expectedBool:    false,
		},
		{
			title:           "Currently locked and trying to unlock",
			running:         true,
			requestedAction: unlockConstant,
			start:           lockConstant,
			stop:            unlockConstant,
			expectedBool:    true,
		},
	}

	for _, test := range stateOptions {
		t.Run(test.title, func(t *testing.T) {

			isValid := validAction(test.running, test.requestedAction, test.start, test.stop)

			if isValid != test.expectedBool {
				t.Errorf("\nActions - wanted: %t, got %t\n", test.expectedBool, isValid)
			}
		})
	}
}

func Test_findLabel(t *testing.T) {

	var stateOptions = []struct {
		title         string
		currentLabels []types.IssueLabel
		cmdLabel      string
		expectedFound bool
	}{
		{
			title: "Label exists lowercase",
			currentLabels: []types.IssueLabel{
				types.IssueLabel{
					Name: "rod",
				},
				types.IssueLabel{
					Name: "jane",
				},
				types.IssueLabel{
					Name: "freddie",
				},
			},
			cmdLabel:      "rod",
			expectedFound: true,
		},
		{
			title: "Label exists case insensitive",
			currentLabels: []types.IssueLabel{
				types.IssueLabel{
					Name: "rod",
				},
				types.IssueLabel{
					Name: "jane",
				},
				types.IssueLabel{
					Name: "freddie",
				},
			},
			cmdLabel:      "Rod",
			expectedFound: true,
		},
		{
			title: "Label doesnt exist lowercase",
			currentLabels: []types.IssueLabel{
				types.IssueLabel{
					Name: "rod",
				},
				types.IssueLabel{
					Name: "jane",
				},
				types.IssueLabel{
					Name: "freddie",
				},
			},
			cmdLabel:      "derek",
			expectedFound: false,
		},
		{
			title: "Label doesnt exist case insensitive",
			currentLabels: []types.IssueLabel{
				types.IssueLabel{
					Name: "rod",
				},
				types.IssueLabel{
					Name: "jane",
				},
				types.IssueLabel{
					Name: "freddie",
				},
			},
			cmdLabel:      "Derek",
			expectedFound: false,
		},
		{
			title:         "no existing labels lowercase",
			currentLabels: nil,
			cmdLabel:      "derek",
			expectedFound: false,
		},
		{title: "Label doesnt exist case insensitive",
			currentLabels: nil,
			cmdLabel:      "Derek",
			expectedFound: false,
		},
	}

	for _, test := range stateOptions {
		t.Run(test.title, func(t *testing.T) {

			labelFound := findLabel(test.currentLabels, test.cmdLabel)

			if labelFound != test.expectedFound {
				t.Errorf("Find Labels(%s) - wanted: %t, got %t\n", test.title, test.expectedFound, labelFound)
			}
		})
	}
}

func Test_Parsing_Milestones(t *testing.T) {

	var milestonesOptions = []struct {
		title        string
		body         string
		expectedType string
		expectedVal  string
	}{
		{
			title:        "Right set milestone",
			body:         "set milestone: demo",
			expectedType: "SetMilestone",
			expectedVal:  "demo",
		},
		{
			title:        "Right remove milestone",
			body:         "remove milestone: demo",
			expectedType: "RemoveMilestone",
			expectedVal:  "demo",
		},
		{
			title:        "Wrong set milestone",
			body:         "you ok label: demo",
			expectedType: "",
			expectedVal:  "",
		},
	}

	for _, test := range milestonesOptions {
		t.Run(test.title, func(t *testing.T) {

			for _, trigger := range getCommandTriggers() {
				action := parse(trigger+test.body, getCommandTriggers())
				if action.Type != test.expectedType || action.Value != test.expectedVal {
					t.Errorf("Action - wanted: %s, got %s\nLabel - wanted: %s, got %s", test.expectedType, action.Type, test.expectedVal, action.Value)
				}
			}
		})
	}
}

func Test_isDcoLabel(t *testing.T) {
	dcoLabel := []struct {
		title        string
		label        string
		expectedBool bool
	}{
		{
			title:        "Counts as no-dco - case insensitivity",
			label:        "NO-DCO",
			expectedBool: true,
		},
		{
			title:        "Normal no-dco case",
			label:        "no-dco",
			expectedBool: true,
		},
		{
			title:        "Counts as no-dco - case insensitivity",
			label:        "No-Dco",
			expectedBool: true,
		},
		{
			title:        "Does not follow no-dco so it counts as normal label",
			label:        "nodco",
			expectedBool: false,
		},
		{
			title:        "Normal label",
			label:        "randomlabel",
			expectedBool: false,
		},
	}

	for _, test := range dcoLabel {
		t.Run(test.label, func(t *testing.T) {
			itsDco := isDcoLabel(test.label)
			if itsDco != test.expectedBool {
				t.Errorf("Wanted `%s` to return: %t but it returned:  %t.", test.label, test.expectedBool, itsDco)
			}
		})
	}
}

func Test_Parsing_Reviewers(t *testing.T) {

	var reviewersOptions = []struct {
		title        string
		body         string
		expectedType string
		expectedVal  string
	}{
		{
			title:        "Parse request to assign reviewer with valid message",
			body:         "set reviewer: john",
			expectedType: "AssignReviewer",
			expectedVal:  "john",
		},
		{
			title:        "Parse request to unassign reviewer with valid message",
			body:         "clear reviewer: john",
			expectedType: "UnassignReviewer",
			expectedVal:  "john",
		},
		{
			title:        "Parse request to assign reviewer with invalid message",
			body:         "random message: john",
			expectedType: "",
			expectedVal:  "",
		},
	}

	for _, test := range reviewersOptions {
		t.Run(test.title, func(t *testing.T) {

			for _, trigger := range getCommandTriggers() {
				action := parse(trigger+test.body, getCommandTriggers())
				if action.Type != test.expectedType || action.Value != test.expectedVal {
					t.Errorf("Action - wanted: %s, got %s\nResult - wanted: %s, got %s", test.expectedType, action.Type, test.expectedVal, action.Value)
				}
			}
		})
	}
}
func Test_classifyLabels(t *testing.T) {

	var classifyOptions = []struct {
		title                string
		currentLabels        []types.IssueLabel
		cmdType              string
		labelValue           string
		expectedActionable   []string
		expectedUnactionable []string
	}{
		{title: "Add when all Labels exist",
			currentLabels: []types.IssueLabel{
				types.IssueLabel{
					Name: "rod",
				},
				types.IssueLabel{
					Name: "jane",
				},
				types.IssueLabel{
					Name: "freddie",
				},
			},
			cmdType:              addLabelConstant,
			labelValue:           "rod, jane, freddie",
			expectedActionable:   []string{},
			expectedUnactionable: []string{"rod", "jane", "freddie"},
		},
		{title: "Remove when all Labels exist",
			currentLabels: []types.IssueLabel{
				types.IssueLabel{
					Name: "rod",
				},
				types.IssueLabel{
					Name: "jane",
				},
				types.IssueLabel{
					Name: "freddie",
				},
			},
			cmdType:              removeLabelConstant,
			labelValue:           "rod, jane, freddie",
			expectedActionable:   []string{"rod", "jane", "freddie"},
			expectedUnactionable: []string{},
		},
		{title: "Add when no Labels exist",
			currentLabels:        []types.IssueLabel{},
			cmdType:              addLabelConstant,
			labelValue:           "rod, jane, freddie",
			expectedActionable:   []string{"rod", "jane", "freddie"},
			expectedUnactionable: []string{},
		},
		{title: "Remove when no Labels exist",
			currentLabels:        []types.IssueLabel{},
			cmdType:              removeLabelConstant,
			labelValue:           "rod, jane, freddie",
			expectedActionable:   []string{},
			expectedUnactionable: []string{"rod", "jane", "freddie"},
		},
		{title: "Add when subset of Labels exist",
			currentLabels: []types.IssueLabel{
				types.IssueLabel{
					Name: "rod",
				},
				types.IssueLabel{
					Name: "jane",
				},
			},
			cmdType:              addLabelConstant,
			labelValue:           "rod, jane, freddie",
			expectedActionable:   []string{"freddie"},
			expectedUnactionable: []string{"rod", "jane"},
		},
		{title: "Remove when subset of Labels exist",
			currentLabels: []types.IssueLabel{
				types.IssueLabel{
					Name: "rod",
				},
				types.IssueLabel{
					Name: "jane",
				},
			},
			cmdType:              removeLabelConstant,
			labelValue:           "rod, jane, freddie",
			expectedActionable:   []string{"rod", "jane"},
			expectedUnactionable: []string{"freddie"},
		},
		{title: "Add new value to set",
			currentLabels: []types.IssueLabel{
				types.IssueLabel{
					Name: "rod",
				},
				types.IssueLabel{
					Name: "jane",
				},
			},
			cmdType:              addLabelConstant,
			labelValue:           "freddie, burt",
			expectedActionable:   []string{"freddie", "burt"},
			expectedUnactionable: []string{},
		},
		{title: "remove existing values from set",
			currentLabels: []types.IssueLabel{
				types.IssueLabel{
					Name: "rod",
				},
				types.IssueLabel{
					Name: "jane",
				},
				types.IssueLabel{
					Name: "freddie",
				},
				types.IssueLabel{
					Name: "burt",
				},
			},
			cmdType:              removeLabelConstant,
			labelValue:           "rod, jane",
			expectedActionable:   []string{"freddie", "burt"},
			expectedUnactionable: []string{},
		},
	}

	for _, test := range classifyOptions {
		t.Run(test.title, func(t *testing.T) {

			actionableLabels, unactionableLabels := classifyLabels(test.currentLabels, test.cmdType, test.labelValue)

			if len(actionableLabels) != len(test.expectedActionable) || len(unactionableLabels) != len(test.expectedUnactionable) {
				t.Errorf("Label Classification (%s) - wanted: Actionable(%s) Unactionable(%s), got Actionable(%s) Unactionable(%s)\n", test.title, strings.Join(test.expectedActionable, ", "), strings.Join(test.expectedUnactionable, ", "), strings.Join(actionableLabels, ", "), strings.Join(unactionableLabels, ", "))
			}
		})
	}
}

func Test_getMultiLabelLimit(t *testing.T) {

	var labelLimits = []struct {
		title       string
		envVar      string
		envVal      string
		expectedVal int
	}{
		{
			title:       "No ENV var",
			envVar:      "random",
			envVal:      "10",
			expectedVal: labelLimitDefault,
		},
		{
			title:       "ENV var exists - all valid",
			envVar:      labelLimitEnvVar,
			envVal:      "8",
			expectedVal: 8,
		},
		{
			title:       "ENV var exists but cannot be cast as int",
			envVar:      labelLimitEnvVar,
			envVal:      "fred",
			expectedVal: labelLimitDefault,
		},
	}

	for _, test := range labelLimits {
		t.Run(test.title, func(t *testing.T) {

			os.Setenv(test.envVar, test.envVal)

			maxActionableLabels := getMultiLabelLimit()

			os.Unsetenv(test.envVar)

			if maxActionableLabels != test.expectedVal {
				t.Errorf("multi-label limit wrong value found - wanted: %d, found %d", test.expectedVal, maxActionableLabels)
			}
		})
	}
}

func Test_getCommandValue(t *testing.T) {

	var commandValues = []struct {
		title       string
		commentBody string
		trigger     string
		expectedVal string
	}{
		{
			title:       "Single Label",
			commentBody: "Derek add label: burt",
			trigger:     "Derek add label: ",
			expectedVal: "burt",
		},
		{
			title:       "Single Label trailing spaces",
			commentBody: "Derek add label: burt       ",
			trigger:     "Derek add label: ",
			expectedVal: "burt",
		},
		{
			title:       "Single Label trailing dots",
			commentBody: "Derek add label: burt........",
			trigger:     "Derek add label: ",
			expectedVal: "burt",
		},
		{
			title:       "Single Label trailing commas",
			commentBody: "Derek add label: burt,,,,,,,,,,,",
			trigger:     "Derek add label: ",
			expectedVal: "burt",
		},
		{
			title:       "Single Label trailing mixure",
			commentBody: "Derek add label: burt,,. , ,,	,.,",
			trigger:     "Derek add label: ",
			expectedVal: "burt",
		},
		{
			title:       "Multiple Labels",
			commentBody: "Derek add label: burt, and, ernie",
			trigger:     "Derek add label: ",
			expectedVal: "burt, and, ernie",
		},
		{
			title: "Multi-line Labels",
			commentBody: `Derek add label: burt
											, and
											, ernie`,
			trigger:     "Derek add label: ",
			expectedVal: "burt",
		},
		{
			title: "Multi-line Labels with a trailing comma",
			commentBody: `Derek add label: burt,
											 and,
											 ernie`,
			trigger:     "Derek add label: ",
			expectedVal: "burt",
		},
	}

	for _, test := range commandValues {
		t.Run(test.title, func(t *testing.T) {

			val := getCommandValue(test.commentBody, len(test.trigger))

			if val != test.expectedVal {
				t.Errorf("command value error - wanted: %s, found %s", test.expectedVal, val)
			}
		})
	}
}

func Test_createIssueComment(t *testing.T) {
	desiredBody := "The content of the message with slack info"
	tests := []struct {
		title             string
		message           []types.Message
		wantedMessage     string
		wantedErr         bool
		desiredGitHubBody *github.IssueComment
	}{
		{
			title: "Desired message is found.",
			message: []types.Message{
				{Name: "slack", Value: "The content of the message with slack info"},
				{Name: "docs", Value: "The content of the message with docs info"},
			},
			wantedMessage: "slack",
			wantedErr:     false,
			desiredGitHubBody: &github.IssueComment{
				Body: &desiredBody,
			},
		},
		{
			title: "User tried to set non existing message",
			message: []types.Message{
				{Name: "slack", Value: "The content of the message with slack info"},
				{Name: "docs", Value: "The content of the message with docs info"},
			},
			wantedMessage:     "dco",
			wantedErr:         true,
			desiredGitHubBody: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			gitHubComment, err := createIssueComment(test.message, test.wantedMessage)
			if gitHubComment != nil && test.desiredGitHubBody != nil {
				if *gitHubComment.Body != *test.desiredGitHubBody.Body {
					t.Errorf("Expected body to contain: %s got :%s",
						*test.desiredGitHubBody.Body,
						*gitHubComment.Body)
				}
			}
			if err != nil && test.wantedErr == false {
				t.Errorf("Unexpected error: %s",
					err.Error())
			}
		})
	}
}
