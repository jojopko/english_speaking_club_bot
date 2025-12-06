package ai

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

func (c *client) AskQuestion(question string, extraQuestion string) (string, error) {
	messages := []InputMessage{}
	messages = append(messages, InputMessage{
		Role: "system",
		Content: []ContentItem{
			{
				Type: "text",
				Text: c.systemPrompt,
			},
		},
	})
	messages = append(messages, InputMessage{
		Role: "user",
		Content: []ContentItem{
			{
				Type: "text",
				Text: makeQuestion(question, extraQuestion),
			},
		},
	})

	payloadData := GenApiRequest{
		Messages: messages,
		IsSync:   true,
	}

	slog.Debug("ai payload:", "payload", payloadData)

	payloadBytes, err := json.Marshal(payloadData)
	if err != nil {
		return "", err
	}

	url := c.baseUrl + "/api/v1/networks/gpt-4o-mini"

	req, err := http.NewRequest("POST", url, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return "", err
	}
	for key, value := range c.headers {
		req.Header.Add(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("AI provider returned non-200 status: %d", resp.StatusCode)
	}

	var response GenApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	if len(response.Responses) == 0 {
		return "", fmt.Errorf("empty response from AI")
	}

	return response.Responses[0].Message.Content, nil
}

func makeQuestion(question string, extraQuestion string) string {
	questionBuilder := strings.Builder{}
	questionBuilder.WriteString("Question=")
	questionBuilder.WriteString(question)
	questionBuilder.WriteString("\n")
	if extraQuestion != "" {
		questionBuilder.WriteString("ExtraQuestion=")
		questionBuilder.WriteString(extraQuestion)
	}
	return questionBuilder.String()
}
