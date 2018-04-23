package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/alexellis/derek/types"
)

const slack = "slack"
const slackDefaultIconURL string = "https://camo.githubusercontent.com/cf0edcdaf482b61b065bde6ce7744f7fc3164d69/68747470733a2f2f7062732e7477696d672e636f6d2f6d656469612f44506f344f7972577341414f6b5f692e706e67"
const slackDefaultUsername string = "Derek"

func handleSlackMessage(text string) error {
	url := os.Getenv("slack_webhook_url")

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
	if channel := os.Getenv("slack_channel"); channel != "" {
		payload.Channel = channel
	}
	if username := os.Getenv("slack_username"); username != "" {
		payload.Username = username
	}
	if iconURL := os.Getenv("slack_iconURL"); iconURL != "" {
		payload.IconURL = iconURL
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
