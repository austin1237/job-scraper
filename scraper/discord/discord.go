package discord

import (
	"bytes"
	"encoding/json"
	"net/http"
	"scraper/job"
)

func generateMessages(jobs []job.Job) []string {
	var messages []string
	var message bytes.Buffer
	message.WriteString("```")

	for _, job := range jobs {
		newLine := job.Link + ", " + job.Company + "\n"
		// Discord has a 2000 character limit for messages
		if message.Len()+len(newLine)+3 > 2000 { // +3 for the ending "```"
			message.WriteString("```")
			messages = append(messages, message.String())
			message.Reset()
			message.WriteString("```")
		}
		message.WriteString(newLine)
	}

	if message.Len() > 0 {
		message.WriteString("```")
		messages = append(messages, message.String())
	}

	return messages
}

func SendJobsToDiscord(jobs []job.Job, webhookURL string) []error {
	if len(jobs) == 0 {
		return nil
	}
	messages := generateMessages(jobs)
	errorChannel := make(chan error, len(messages))

	go func() {
		for _, message := range messages {
			payload := map[string]string{
				"content": message,
			}

			jsonPayload, err := json.Marshal(payload)
			if err != nil {
				errorChannel <- err
				continue
			}

			resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonPayload))
			if err != nil {
				errorChannel <- err
				continue
			}
			defer resp.Body.Close()
		}
		close(errorChannel)
	}()

	var errors []error
	for err := range errorChannel {
		errors = append(errors, err)
	}

	return errors
}
