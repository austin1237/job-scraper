package discord

import (
	"bytes"
	"encoding/json"
	"jobNotifier/job"
	"net/http"
)

func generateMessages(jobs []job.Job) []string {
	var messages []string
	var message bytes.Buffer
	message.WriteString("```")

	for _, job := range jobs {
		newLine := job.Link + ", " + job.Company + "\n"
		// Discord has a 2000 character limit for messages
		if message.Len()+len(newLine)+3 >= 2000 { // +3 for the ending "```"
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

	for _, message := range messages {
		go func(message string) {
			payload := map[string]string{
				"content": message,
			}

			jsonPayload, err := json.Marshal(payload)
			if err != nil {
				errorChannel <- err
				return
			}

			resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonPayload))
			if err != nil {
				errorChannel <- err
				return
			}
			defer resp.Body.Close()
			errorChannel <- nil
		}(message)
	}

	var errors []error
	for i := 0; i < len(messages); i++ {
		if err := <-errorChannel; err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}
