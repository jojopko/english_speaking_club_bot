package bot

import (
	"fmt"
	"speaking-club-bot/internal/client/telegram"
)

const shouldHadReply = "Use this command to reply to the message you need help on"

func helpmeHandle(bot *Bot, update *telegram.Update) error {
	if update.Message == nil {
		return nil
	}
	message := update.Message

	if message.ReplyTo == nil {
		_, err := bot.TelegramClient().SendMessage(telegram.SendMessageRequest{
			ChatID: message.Chat.ID,
			Text:   shouldHadReply,
		})
		return err
	}

	extraQuestion := getExtraQuestion(message.Text)

	s := message.ReplyTo.Text
	answer, err := bot.MentorClient().AskQuestion(s, extraQuestion)
	if err != nil {
		return fmt.Errorf("mentor has error: %w", err)
	}

	_, err = bot.TelegramClient().SendMessage(telegram.SendMessageRequest{
		ChatID: message.Chat.ID,
		Text:   answer,
		ReplyParameters: &telegram.ReplyParameters{
			MessageID: message.MessageID,
		},
	})
	return err
}
