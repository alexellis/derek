package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alexellis/derek/types"
)

const slack = "slack"
const slackDefaultIconURL string = "https://camo.githubusercontent.com/cf0edcdaf482b61b065bde6ce7744f7fc3164d69/68747470733a2f2f7062732e7477696d672e636f6d2f6d656469612f44506f344f7972577341414f6b5f692e706e67"
const slackDefaultUsername string = "Derek"

var isEnabled bool
var settings types.SlackSetting

func setSlackSettings(config *types.DerekConfig) {
	isEnabled = enabledFeature(slack, config)
	settings = config.Slack
}

func handleSlackMessage(text string) error {
	if isEnabled != true {
		return nil
	}

	url := settings.WebhookURL
	if url == "" {
		return fmt.Errorf("Slack Webhook Url not set in DerekConfig")
	}

	// Build with default values
	payload := types.SlackPayload{
		Text:     text,
		Username: slackDefaultUsername,
		IconURL:  slackDefaultIconURL,
	}

	// Set Overrides
	if settings.Channel != "" {
		payload.Channel = settings.Channel
	}
	if settings.Username != "" {
		payload.Username = settings.Username
	}
	if settings.IconURL != "" {
		payload.IconURL = settings.IconURL
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Slack didn’t respond 200 OK: %s", resp.Status)
	}
	defer resp.Body.Close()

	return nil
}
