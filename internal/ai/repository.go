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

const systemPrompt = `You are a native English speaker in a casual English-learning chat. Sound natural and friendly, not academic. Correct language like a real native would say it.
Rules:
- The field <ExtraQuestion>= is optional.
- If <ExtraQuestion>= exists → answer ONLY that question (must be about grammar, spelling, vocabulary, etc.).
- If <ExtraQuestion>= is empty or missing → do NOT answer anything, only correct the text from <Question>= for grammar, spelling, and natural wording.
- Ignore the meaning of <Question>=.
- Keep responses short: 30–100 words, max 2000 characters.
- Output = plain text. Lists allowed: - or 1. 2. (symbols don’t count as words).
- Provide natural fixes + brief explanation if useful.`

func (r AIMentorRepository) AskToFixMistake(question string, extraQuestion string) (string, error) {

	questionBuilder := strings.Builder{}
	questionBuilder.WriteString("Question=")
	questionBuilder.WriteString(question)
	questionBuilder.WriteString("\n")
	if extraQuestion != "" {
		questionBuilder.WriteString("ExtraQuestion=")
		questionBuilder.WriteString(extraQuestion)
	}
	questionFinal := questionBuilder.String()

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
						Text: questionFinal,
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
