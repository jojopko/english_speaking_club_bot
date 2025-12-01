package ai

import (
	"encoding/json"
	"fmt"
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

	payloadBytes, err := json.Marshal(payloadData)
	if err != nil {
		return "", err
	}

	url := c.baseUrl + "api/v1/networks/gpt-4o-mini"

	req, err := c.httpClient.Post(url, "application/json", strings.NewReader(string(payloadBytes)))
	if err != nil {
		return "", err
	}
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	defer req.Body.Close()

	if req.StatusCode != http.StatusOK {
		return "", fmt.Errorf("AI provider returned non-200 status: %d", req.StatusCode)
	}

	var response GenApiResponse
	if err := json.NewDecoder(req.Body).Decode(&response); err != nil {
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
