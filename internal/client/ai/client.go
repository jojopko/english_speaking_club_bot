package ai

import (
	"net/http"
	"strings"
)

type Client interface {
	AskQuestion(question string, extraQuestion string) (string, error)
}

type client struct {
	baseUrl      string
	headers      map[string]string
	httpClient   *http.Client
	systemPrompt string
}

func New(
	baseUrl string,
	apiToken string,
	systemPrompt string,
) Client {
	baseUrl = strings.TrimSuffix(baseUrl, "/")
	headers := map[string]string{
		"Authorization": "Bearer " + apiToken,
		"Content-Type":  "application/json",
	}
	return &client{
		baseUrl: baseUrl,
		httpClient: &http.Client{
			Timeout: http.DefaultClient.Timeout,
		},
		headers:      headers,
		systemPrompt: systemPrompt,
	}
}
