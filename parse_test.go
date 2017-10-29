package main

import (
	"testing"
)

func Test_Parsing_Close(t *testing.T) {
	action := parse("Derek close")

	if action.Type != "close" {
		t.Errorf("Action - want: %s, got %s", "close", action.Type)
	}

}

func Test_Parsing_AddLabel(t *testing.T) {
	action := parse("Derek add label: demo")
	want := "AddLabel"
	if action.Type != want {
		t.Errorf("Action - want: %s, got %s", want, action.Type)

	}
	wantValue := "demo"
	if action.Value != wantValue {
		t.Errorf("Action - want: %s, got %s", wantValue, action.Value)
	}
}
