package ai

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type MentorRepository interface {
	AskToFixMistake(question string, extraQuestion string) (string, error)
}

type AIMentorRepository struct {
	apiToken string
}

const AIEndpoint = "https://api.gen-api.ru/api/v1/networks/gpt-4o-mini"

func NewAIMentorRepository(apiToken string) AIMentorRepository {
	return AIMentorRepository{
		apiToken: apiToken,
	}
}

type aiRequestPayload struct {
	IsSync   bool      `json:"is_sync"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string        `json:"role"`
	Content []ContentItem `json:"content"`
}

type ContentItem struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func (r *AIMentorRepository) AskToFixMistake(question string, extraQuestion string) (string, error) {
	systemPrompt := "You are a helpful assistant that helps people to learn English. Use plain text formating.\n"
	payloadData := aiRequestPayload{
		IsSync: true,
		Messages: []Message{
			{
				Role: "system",
				Content: []ContentItem{
					{
						Type: "text",
						Text: systemPrompt,
					},
				},
			},
			{
				Role: "user",
				Content: []ContentItem{
					{
						Type: "text",
						Text: question,
					},
				},
			},
		},
	}

	payloadBytes, err := json.Marshal(payloadData)
	if err != nil {
		return "", err
	}
	payload := strings.NewReader(string(payloadBytes))

	client := &http.Client{}
	req, err := http.NewRequest("POST", AIEndpoint, payload)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.apiToken))
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var result struct {
		RequestId int64   `json:"request_id"`
		Model     string  `json:"model"`
		Cost      float64 `json:"cost"`
		Response  []struct {
			Index   int `json:"index"`
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
		} `json:"Response"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Response) == 0 {
		return "", fmt.Errorf("empty response from AI")
	}

	return result.Response[0].Message.Content, nil
}
