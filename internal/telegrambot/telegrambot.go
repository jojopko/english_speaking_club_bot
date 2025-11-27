package telegrambot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"speaking-club-bot/internal/ai"
	"strings"
)

const telegramAPIBaseURL = "https://api.telegram.org/bot"

type Update struct {
	UpdateID int      `json:"update_id"`
	Message  *Message `json:"message,omitempty"`
}

type Message struct {
	MessageID int      `json:"message_id"`
	From      *User    `json:"from,omitempty"`
	Chat      Chat     `json:"chat"`
	ReplyTo   *Message `json:"reply_to_message,omitempty"`
	Text      string   `json:"text"`
}

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

func (u User) Name() string {
	if u.Username != "" {
		return u.Username
	}

	s := fmt.Sprintf("%s %s", u.FirstName, u.LastName)
	s = strings.TrimSpace(s)
	if s != "" {
		return s
	}

	return fmt.Sprintf("user#%d", u.ID)
}

type Chat struct {
	ID int64 `json:"id"`
}

type Bot struct {
	tokenAPI    string
	Whitelist   []int64
	StorageRepo BotStorageRepository
	MentorRepo  ai.AIMentorRepository
}

func NewBot(token string, whitelist []int64, storageRepo BotStorageRepository, mentorRepo ai.AIMentorRepository) (*Bot, error) {
	return &Bot{
		tokenAPI:    token,
		Whitelist:   whitelist,
		StorageRepo: storageRepo,
		MentorRepo:  mentorRepo,
	}, nil
}

func (b *Bot) getURL(method string) string {
	return telegramAPIBaseURL + b.tokenAPI + "/" + method
}

func (b *Bot) GetUpdates(offset int) ([]Update, error) {
	url := b.getURL("getUpdates")
	reqData := map[string]any{
		"offset":  offset,
		"timeout": 60,
	}

	body, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("marshaling error: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Ok     bool     `json:"ok"`
		Result []Update `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if !result.Ok {
		return nil, fmt.Errorf("telegram API returned not ok")
	}

	return result.Result, nil
}

func (b *Bot) SendMessage(chatID int64, text string, replyToMessageID int) error {
	url := b.getURL("sendMessage")
	reqData := map[string]any{
		"chat_id": chatID,
		"text":    text,
	}

	if replyToMessageID != 0 {
		reqData["reply_parameters"] = map[string]any{
			"message_id": replyToMessageID,
		}
	}

	body, err := json.Marshal(reqData)
	if err != nil {
		return fmt.Errorf("marshaling error: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result struct {
		Ok bool `json:"ok"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if !result.Ok {
		return fmt.Errorf("telegram API returned not ok")
	}

	return nil
}
