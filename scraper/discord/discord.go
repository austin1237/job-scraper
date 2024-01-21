package discord

import (
	"bytes"
	"encoding/json"
	"net/http"
	"scraper/interest"
)

func SendJobsToDiscord(jobs []interest.Job, webhookURL string) error {
	var message bytes.Buffer
	message.WriteString("```")
	for _, job := range jobs {
		message.WriteString(job.Link)
		message.WriteString(", ")
		message.WriteString(job.Name)
		message.WriteString("\n")
	}
	message.WriteString("```")

	payload := map[string]string{
		"content": message.String(),
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
