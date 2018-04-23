package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alexellis/derek/types"
)

func Test_handleSlackMessage(t *testing.T) {
	var resp types.SlackPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp = types.SlackPayload{}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		json.Unmarshal(body, &resp)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	type slackTestConfig struct {
		url      string
		username string
		iconURL  string
		channel  string
	}

	var slackTestOpts = []struct {
		title    string
		username string
		iconURL  string
		channel  string
		config   slackTestConfig
	}{
		{
			title:    "Should use Defaults values if no settings defined",
			username: slackDefaultUsername,
			iconURL:  slackDefaultIconURL,
			channel:  "",
			config: slackTestConfig{
				url: ts.URL,
			},
		},
		{
			title:    "Should use Setting's custom username",
			username: "Custom",
			iconURL:  slackDefaultIconURL,
			channel:  "",
			config: slackTestConfig{
				url:      ts.URL,
				username: "Custom",
			},
		},
		{
			title:    "Should use Setting's custom Icon Url",
			username: slackDefaultUsername,
			iconURL:  "http://example.com",
			channel:  "",
			config: slackTestConfig{
				url:     ts.URL,
				iconURL: "http://example.com",
			},
		},
		{
			title:    "Should use Setting's custom channel",
			username: slackDefaultUsername,
			iconURL:  slackDefaultIconURL,
			channel:  "#build",
			config: slackTestConfig{
				url:     ts.URL,
				channel: "#build",
			},
		},
		{
			title:    "Should use Setting's custom values",
			username: "Bob",
			iconURL:  "http://example.com/image.png",
			channel:  "#github",
			config: slackTestConfig{
				url:      ts.URL,
				username: "Bob",
				iconURL:  "http://example.com/image.png",
				channel:  "#github",
			},
		},
	}

	for _, test := range slackTestOpts {
		t.Run(test.title, func(t *testing.T) {
			os.Setenv("slack_channel", test.config.channel)
			os.Setenv("slack_username", test.config.username)
			os.Setenv("slack_iconURL", test.config.iconURL)
			os.Setenv("slack_webhook_url", test.config.url)

			err := handleSlackMessage(test.title)
			if err != nil {
				t.Errorf("Expext Slack Message to successfully send")
			}
			if resp.Text != test.title {
				t.Errorf("Expected Text to be the same - wanted: '%s', found '%s'", test.title, resp.Text)
			}
			if resp.Username != test.username {
				t.Errorf("Expected username to be the same - wanted: '%s', found '%s'", test.username, resp.Username)
			}
			if resp.IconURL != test.iconURL {
				t.Errorf("Expected iconURL to be the same - wanted: '%s', found '%s'", test.iconURL, resp.IconURL)
			}
			if resp.Channel != test.channel {
				t.Errorf("Expected username to be the same - wanted: '%s', found '%s'", test.channel, resp.Channel)
			}
			os.Unsetenv("slack_channel")
			os.Unsetenv("slack_username")
			os.Unsetenv("slack_iconURL")
			os.Unsetenv("slack_webhook_url")
		})
	}
}
