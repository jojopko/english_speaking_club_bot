package bot

import (
	"fmt"
	"speaking-club-bot/internal/client/telegram"
	"strings"
)

func askHandler(bot *Bot, update *telegram.Update) error {
	if update.Message == nil {
		return nil
	}
	message := update.Message
	question := getExtraQuestion(message.Text)

	answer, err := bot.MentorClient().AskQuestion("", question)
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

func getExtraQuestion(text string) string {
	// "@nickname_bot help me" -> "help me"
	// "/command extra question" -> "extra question"
	commandSlice := strings.SplitN(text, " ", 2)
	var extraQuestion string
	if len(commandSlice) < 2 {
		extraQuestion = ""
	} else {
		extraQuestion = commandSlice[1]
	}
	extraQuestion = strings.TrimSpace(extraQuestion)

	return extraQuestion
}
