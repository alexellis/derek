package main

import (
	"testing"
)

var actionOptions = []struct {
	title          string
	body           string
	expectedAction string
}{
	{
		title:          "Correct reopen command",
		body:           "Derek reopen",
		expectedAction: "reopen",
	},
	{ //this case replaces Test_Parsing_Close
		title:          "Correct close command",
		body:           "Derek close",
		expectedAction: "close",
	},
	{
		title:          "invalid command",
		body:           "Derek dance",
		expectedAction: "",
	},
	{
		title:          "Longer reopen command",
		body:           "Derek reopen: ",
		expectedAction: "reopen",
	},
	{
		title:          "Longer close command",
		body:           "Derek close: ",
		expectedAction: "close",
	},
}

func Test_Parsing_OpenClose(t *testing.T) {

	for _, test := range actionOptions {
		t.Run(test.title, func(t *testing.T) {

			action := parse(test.body)

			if action.Type != test.expectedAction {
				t.Errorf("Action - want: %s, got %s", test.expectedAction, action.Type)
			}

		})
	}
}

var labelOptions = []struct {
	title        string
	body         string
	expectedType string
	expectedVal  string
}{
	{ //this case replaces Test_Parsing_AddLabel
		title:        "Add label of demo",
		body:         "Derek add label: demo",
		expectedType: "AddLabel",
		expectedVal:  "demo",
	},
	{
		title:        "Remove label of demo",
		body:         "Derek remove label: demo",
		expectedType: "RemoveLabel",
		expectedVal:  "demo",
	},
	{
		title:        "Invalid label action",
		body:         "Derek peel label: demo",
		expectedType: "",
		expectedVal:  "",
	},
}

func Test_Parsing_Labels(t *testing.T) {

	for _, test := range labelOptions {
		t.Run(test.title, func(t *testing.T) {

			action := parse(test.body)
			if action.Type != test.expectedType || action.Value != test.expectedVal {
				t.Errorf("Action - wanted: %s, got %s\nLabel - wanted: %s, got %s", test.expectedType, action.Type, test.expectedVal, action.Value)
			}
		})
	}
}

var assignmentOptions = []struct {
	title        string
	body         string
	expectedType string
	expectedVal  string
}{
	{
		title:        "Assign to burt",
		body:         "Derek assign: burt",
		expectedType: "Assign",
		expectedVal:  "burt",
	},
	{
		title:        "Unassign burt",
		body:         "Derek unassign: burt",
		expectedType: "Unassign",
		expectedVal:  "burt",
	},
	{
		title:        "Assign to me",
		body:         "Derek assign: me",
		expectedType: "Assign",
		expectedVal:  "me",
	},
	{
		title:        "Unassign me",
		body:         "Derek unassign: me",
		expectedType: "Unassign",
		expectedVal:  "me",
	},
	{
		title:        "Invalid assignment action",
		body:         "Derek consign: burt",
		expectedType: "",
		expectedVal:  "",
	},
	{
		title:        "Unassign blank",
		body:         "Derek unassign: ",
		expectedType: "",
		expectedVal:  "",
	},
}

func Test_Parsing_Assignments(t *testing.T) {

	for _, test := range assignmentOptions {
		t.Run(test.title, func(t *testing.T) {

			action := parse(test.body)
			if action.Type != test.expectedType || action.Value != test.expectedVal {
				t.Errorf("Action - wanted: %s, got %s\nMaintainer - wanted: %s, got %s", test.expectedType, action.Type, test.expectedVal, action.Value)
			}
		})
	}
}
