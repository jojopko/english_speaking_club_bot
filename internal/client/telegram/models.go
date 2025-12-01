package telegram

import (
	"fmt"
	"strings"
)

const (
	AllowedUpdateTypeMessage = "message"
)

type TelegramResponse[T any] struct {
	Ok          bool   `json:"ok"`
	Result      T      `json:"result,omitempty"`
	ErrorCode   int    `json:"error_code,omitempty"`
	Description string `json:"description,omitempty"`
}

type GetUpdatesRequest struct {
	Offset         int64    `json:"offset,omitempty"`
	Limit          int64    `json:"limit,omitempty"`
	Timeout        int64    `json:"timeout,omitempty"`
	AllowedUpdates []string `json:"allowed_updates,omitempty"`
}

type SendMessageRequest struct {
	ChatID          int64            `json:"chat_id"`
	Text            string           `json:"text"`
	ReplyParameters *ReplyParameters `json:"reply_parameters,omitempty"`
}

type Update struct {
	UpdateId int64    `json:"update_id"`
	Message  *Message `json:"message,omitempty"`
}

type ReplyParameters struct {
	MessageID int64 `json:"message_id"`
	ChatID    int64 `json:"chat_id,omitempty"`
}

type Message struct {
	MessageID int64    `json:"message_id"`
	From      *User    `json:"from,omitempty"`
	Chat      Chat     `json:"chat"`
	ReplyTo   *Message `json:"reply_to_message,omitempty"`
	Text      string   `json:"text"`
}

type User struct {
	ID        int64  `json:"id"`
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
	ID   int64  `json:"id"`
	Type string `json:"type,omitempty"`
}
