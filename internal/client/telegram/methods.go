package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func (c *client) GetUpdates(req GetUpdatesRequest) ([]Update, error) {
	url := c.baseUrl + "/getUpdates"

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshaling error: %w", err)
	}

	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	var result TelegramResponse[[]Update]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if !result.Ok {
		return nil, fmt.Errorf("telegram api error (%d): %s", result.ErrorCode, result.Description)
	}

	return result.Result, nil
}

func (c *client) SendMessage(req SendMessageRequest) (Message, error) {
	url := c.baseUrl + "/sendMessage"

	body, err := json.Marshal(req)
	if err != nil {
		return Message{}, fmt.Errorf("marshaling error: %w", err)
	}

	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return Message{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return Message{}, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	var result TelegramResponse[Message]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return Message{}, err
	}

	if !result.Ok {
		return Message{}, fmt.Errorf("telegram api error (%d): %s", result.ErrorCode, result.Description)
	}

	return result.Result, nil
}
