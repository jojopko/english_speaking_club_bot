package telegram

import (
	"net/http"
	"strings"
	"time"
)

type Client interface {
	GetUpdates(req GetUpdatesRequest) ([]Update, error)
	SendMessage(req SendMessageRequest) (Message, error)
}

type client struct {
	baseUrl      string
	headers      map[string]string
	httpClient   *http.Client
}

func New(baseUrl string, token string) Client {
	baseUrl = strings.TrimSuffix(baseUrl, "/")
	baseUrl = baseUrl + "/bot" + token
	return &client{
		baseUrl: baseUrl,
		httpClient: &http.Client{
			Timeout: time.Second * 60,
		},
		headers:      make(map[string]string),
	}
}
